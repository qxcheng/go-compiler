package compiler

import (
    "fmt"
    "os"
)

type Cgen struct {
    tree     *ASTNode   // 语法树
    outfile  *os.File   // 汇编结果
    reglist  []string   // 寄存器列表(64位)
    breglist []string   // 寄存器列表(低8位)
    freereg  []bool     // 寄存器对应的状态
    label    int        // 标签id
}

func NewCgen(tree *ASTNode, outfile *os.File) *Cgen {
    return &Cgen{
        tree:    tree,
        outfile: outfile,
        reglist: []string{"%r8", "%r9", "%r10", "%r11", "%r12", "%r13", "%r14", "%r15"},
        breglist: []string{"%r8b", "%r9b", "%r10b", "%r11b", "%r12b", "%r13b", "%r14b", "%r15b"},
        freereg: []bool{true, true, true, true, true, true, true, true},
        label: 0,
    }
}

func (c *Cgen) GenAST() {
    c.cgpreamble()
    c.genAST(c.tree)
}

func (c *Cgen) genAST(tree *ASTNode) {
    if tree != nil {
        switch tree.nodeKind {
        case PrintK, IfK, VarK, AssignK, ForK, FuncK, ReturnK:
            c.genStmt(tree)
        case OpK, ConstK, IdK, CallK:
            c.genExp(tree)
        default:
            c.error("ERROR: not supported nodekind")
        }
        c.genAST(tree.sibling)
    }
}

func (c *Cgen) genStmt(tree *ASTNode) {
    switch tree.nodeKind {
    case PrintK:
        reg := c.genExp(tree.child[0])
        c.cgprintint(reg)
    case VarK:
        //fmt.Println("111: ", tree.child[0].symbleid)
        c.cgglobsym(tree.child[0].symbleid)
    case AssignK:
        reg := c.genExp(tree.child[0])
        c.cgstoreglob(reg, tree.symbleid)
    case IfK:
        var Lfalse, Lend int
        Lfalse = c.genLabel()  // else分支的标签
        if tree.child[2] != nil {
            Lend = c.genLabel()  // if语句尾的标签
        }

        c.genIfExp(tree.child[0], Lfalse)  // 判断结果为false跳转到else标签
        c.freeall_registers()
        c.genAST(tree.child[1])  // if分支语句
        c.freeall_registers()
        if tree.child[2] != nil {
            c.cgjump(Lend)      // 跳过else分支的语句
        }

        c.cglabel(Lfalse)  // 生成else分支的标签
        if tree.child[2] != nil {
            c.genAST(tree.child[2])
            c.freeall_registers()
            c.cglabel(Lend)
        }
    case ForK:
        Lstart := c.genLabel()
        Lend := c.genLabel()
        c.cglabel(Lstart)
        c.genIfExp(tree.child[0], Lend)
        c.freeall_registers()
        c.genAST(tree.child[1])
        c.freeall_registers()
        c.cgjump(Lstart)
        c.cglabel(Lend)
    case FuncK:
        Lend := c.genLabel()
        c.cgfuncpreamble(tree.litval)
        if tree.child[0] != nil {
            c.cgglobsym(tree.child[0].symbleid)
            _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%%rdi, %s(%%rip)\n", Gsym.symbles[tree.child[0].symbleid].Name)
        }
        c.genAST(tree.child[1])
        c.freeall_registers()
        c.cglabel(Lend)
        Gsym.Setglob(tree.symbleid, Lend)
        c.cgfuncpostamble()
    case ReturnK:
        reg := c.genExp(tree.child[0])
        c.cgreturn(reg, tree.symbleid)
    default:
        c.error("Error: not supported statement")
    }
}


func (c *Cgen) genExp(tree *ASTNode) int {
    var leftreg, rightreg int

    if len(tree.child) == 1 {
        leftreg = c.genExp(tree.child[0])  // 函数调用
    } else if len(tree.child) == 2 {
        //fmt.Println("111", tree.child[0], tree.child[1])
        leftreg = c.genExp(tree.child[0])
        rightreg = c.genExp(tree.child[1])
    }

    switch tree.nodeKind {
    case OpK:
        switch tree.token {
        case ADD:
            return c.cgadd(leftreg, rightreg)
        case SUB:
            return c.cgsub(leftreg, rightreg)
        case MUL:
            return c.cgmul(leftreg, rightreg)
        case QUO:
            return c.cgdiv(leftreg, rightreg)
        case EQ:
            return c.cgcompare_and_set(leftreg, rightreg, EQ)
        case GT:
            return c.cgcompare_and_set(leftreg, rightreg, GT)
        case LT:
            return c.cgcompare_and_set(leftreg, rightreg, LT)
        case LE:
            return c.cgcompare_and_set(leftreg, rightreg, LE)
        case GE:
            return c.cgcompare_and_set(leftreg, rightreg, GE)
        case NE:
            return c.cgcompare_and_set(leftreg, rightreg, NE)
        default:
            return -1
        }
    case ConstK:
        return c.cgloadint(tree.intval)
    case IdK:
        return c.cgloadglob(tree.symbleid)
    case CallK:
        return c.cgcall(leftreg, tree.symbleid)
    default:
        return -1
    }
}

func (c *Cgen) genIfExp(tree *ASTNode, label int) int {
    var leftreg, rightreg int

    if len(tree.child) == 2 {
        //fmt.Println("111", tree.child[0], tree.child[1])
        leftreg = c.genIfExp(tree.child[0], -1)  // 不支持多个比较运算符
        rightreg = c.genIfExp(tree.child[1], -1)
    }

    switch tree.nodeKind {
    case OpK:
        switch tree.token {
        case ADD:
            return c.cgadd(leftreg, rightreg)
        case SUB:
            return c.cgsub(leftreg, rightreg)
        case MUL:
            return c.cgmul(leftreg, rightreg)
        case QUO:
            return c.cgdiv(leftreg, rightreg)
        case EQ:
            return c.cgcompare_and_jump(leftreg, rightreg, EQ, label)
        case GT:
            return c.cgcompare_and_jump(leftreg, rightreg, GT, label)
        case LT:
            return c.cgcompare_and_jump(leftreg, rightreg, LT, label)
        case LE:
            return c.cgcompare_and_jump(leftreg, rightreg, LE, label)
        case GE:
            return c.cgcompare_and_jump(leftreg, rightreg, GE, label)
        case NE:
            return c.cgcompare_and_jump(leftreg, rightreg, NE, label)
        default:
            return -1
        }
    case ConstK:
        return c.cgloadint(tree.intval)
    case IdK:
        return c.cgloadglob(tree.symbleid)
    default:
        return -1
    }
}

func (c *Cgen) genLabel() int {
    c.label++
    return c.label-1
}


/**********************************************/
/**************** 生成汇编语句 *****************/
/**********************************************/

func (c *Cgen) error(msg string) {
    panic(msg)
}

// 设置所有寄存器为可用状态
func (c *Cgen) freeall_registers() {
    for i, _ := range c.freereg {
        c.freereg[i] = true
    }
    //fmt.Println("After Freeall: ", c.freereg)
}

// 分配一个空闲的寄存器
func (c *Cgen) alloc_register() int {
    for i := 0; i < len(c.reglist); i++ {
        if c.freereg[i] {
            c.freereg[i] = false
            return i
        }
    }
    c.error("Error: Out of registers!")
    return 0
}

// 释放一个使用状态的寄存器
func (c *Cgen) free_register(reg int) {
    if c.freereg[reg] != false {
        c.error(fmt.Sprintf("Error: trying to free register %d\n", reg))
    }
    c.freereg[reg] = true
}

// 汇编头
func (c *Cgen) cgpreamble() {
    c.freeall_registers()
    _, _ = c.outfile.WriteString(`    .text
.LC0:
    .string "%d\n"
printint:
	pushq   %rbp
	movq    %rsp, %rbp
	subq    $16, %rsp
	movl    %edi, -4(%rbp)
	movl    -4(%rbp), %eax
	movl    %eax, %esi
	leaq	.LC0(%rip), %rdi
	movl	$0, %eax
	call	printf@PLT
	nop
	leave
	ret
`)
}

// 函数头
func (c *Cgen) cgfuncpreamble(name string) {
    _, _ = fmt.Fprintf(c.outfile, "\n\t.text\n" +
        "\t.globl\t%s\n" +
        "\t.type\t%s, @function\n" +
        "%s:\n" +
        "\tpushq\t%%rbp\n" +
        "\tmovq\t%%rsp, %%rbp\n", name, name, name)
}


// 函数尾
func (c *Cgen) cgfuncpostamble() {
    _, _ = c.outfile.WriteString(`    popq	%rbp
	ret
`)
}

// 加载整型
func (c *Cgen) cgloadint(value int) int {
    //fmt.Println("value: ", value)
    r := c.alloc_register()
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t$%d, %s\n", value, c.reglist[r])
    return r
}

// 加法
func (c *Cgen) cgadd(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\taddq\t%s, %s\n", c.reglist[r1], c.reglist[r2])
    c.free_register(r1)
    return r2
}

// 减法
func (c *Cgen) cgsub(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\tsubq\t%s, %s\n", c.reglist[r2], c.reglist[r1])
    c.free_register(r2)
    return r1
}

// 乘法
func (c *Cgen) cgmul(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\timulq\t%s, %s\n", c.reglist[r1], c.reglist[r2])
    c.free_register(r1)
    return r2
}

// 除法
func (c *Cgen) cgdiv(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s,%%rax\n", c.reglist[r1])
    _, _ = fmt.Fprintf(c.outfile, "\tcqo\n")
    _, _ = fmt.Fprintf(c.outfile, "\tidivq\t%s\n", c.reglist[r2])
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%%rax,%s\n", c.reglist[r1])
    c.free_register(r2)
    return r1
}

// 打印
func (c *Cgen) cgprintint(r int) {
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s, %%rdi\n", c.reglist[r])
    _, _ = fmt.Fprintf(c.outfile, "\tcall\tprintint\n")
    c.free_register(r)
}

// 加载变量
func (c *Cgen) cgloadglob(id int) int {
    r := c.alloc_register()
    switch Gsym.symbles[id].Vartype {
    case VAR_CHAR:
        _, _ = fmt.Fprintf(c.outfile, "\tmovzbq\t%s(%%rip), %s\n", Gsym.symbles[id].Name, c.reglist[r])
    case VAR_INT:
        _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s(%%rip), %s\n", Gsym.symbles[id].Name, c.reglist[r])
    default:
        c.error("Error: unspported vartype")
    }
    return r
}

// 变量赋值
func (c *Cgen) cgstoreglob(r int, id int) int {
    switch Gsym.symbles[id].Vartype {
    case VAR_CHAR:
        _, _ = fmt.Fprintf(c.outfile, "\tmovb\t%s, %s(%%rip)\n", c.breglist[r], Gsym.symbles[id].Name)
    case VAR_INT:
        _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s, %s(%%rip)\n", c.reglist[r], Gsym.symbles[id].Name)
    default:
        c.error("Error: unspported vartype")
    }
    return r
}

// 创建变量
func (c *Cgen) cgglobsym(id int) {
    switch Gsym.symbles[id].Vartype {
    case VAR_CHAR:
        _, _ = fmt.Fprintf(c.outfile, "\t.comm\t%s,1,1\n", Gsym.symbles[id].Name)
    case VAR_INT:
        _, _ = fmt.Fprintf(c.outfile, "\t.comm\t%s,8,8\n", Gsym.symbles[id].Name)
    default:
        c.error("Error: unspported vartype")
    }
}

// 寄存器变量的类型转换
func (c *Cgen) cgwiden(r int, oldtype int, newtype int) int {
    return r
}

// 比较并设置
var cmpdict = map[Token]string{
    EQ: "sete",
    GT: "setg",
    LT: "setl",
    LE: "setle",
    GE: "setge",
    NE: "setne",
}

func (c *Cgen) cgcompare_and_set(r1 int, r2 int, how Token) int {
    set, ok := cmpdict[how]
    if !ok {
        panic("Error: unspported compare token")
    }
    _, _ = fmt.Fprintf(c.outfile, "\tcmpq\t%s, %s\n", c.reglist[r2], c.reglist[r1])
    _, _ = fmt.Fprintf(c.outfile, "\t%s\t%s\n", set, c.breglist[r2])
    _, _ = fmt.Fprintf(c.outfile, "\tmovzbq\t%s, %s\n", c.breglist[r2], c.reglist[r2])
    c.free_register(r1)
    return r2
}

// 生成一个标签
func (c *Cgen) cglabel(l int) {
    _, _ = fmt.Fprintf(c.outfile, "L%d:\n", l)
}

// 跳转到一个标签
func (c *Cgen) cgjump(l int) {
    _, _ = fmt.Fprintf(c.outfile, "\tjmp\tL%d\n", l)
}

// 比较并在false时跳转
var jumpdict = map[Token]string{
    EQ: "jne",
    GT: "jle",
    LT: "jge",
    GE: "jl",
    LE: "jg",
    NE: "je",
}

func (c *Cgen) cgcompare_and_jump(r1 int, r2 int, how Token, label int) int {
    set, ok := jumpdict[how]
    if !ok {
        panic("Error: unspported jump token")
    }
    _, _ = fmt.Fprintf(c.outfile, "\tcmpq\t%s, %s\n", c.reglist[r2], c.reglist[r1])
    _, _ = fmt.Fprintf(c.outfile, "\t%s\tL%d\n", set, label)
    c.freeall_registers()
    return -1
}

// 函数调用
func (c *Cgen) cgcall(r int, id int) int {
    outr := c.alloc_register()
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s, %%rdi\n", c.reglist[r])
    _, _ = fmt.Fprintf(c.outfile, "\tcall\t%s\n", Gsym.symbles[id].Name)
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%%rax, %s\n", c.reglist[outr])
    c.free_register(r)
    return outr
}

// 函数返回一个值
func (c *Cgen) cgreturn(r int, id int) {
    switch Gsym.symbles[id].ReturnType {
    case VAR_CHAR:
        _, _ = fmt.Fprintf(c.outfile, "\tmovzbl\t%s, %%eax\n", c.breglist[r])
    case VAR_INT:
        _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s, %%rax\n", c.reglist[r])
    default:
        c.error("Error: undefined return type")
    }
    c.cgjump(Gsym.symbles[id].EndLabel)
}