package compiler

import (
    "fmt"
    "os"
)

type Cgen struct {
    tree    *ASTNode   // 语法树
    outfile *os.File   // 汇编结果
    reglist []string   // 寄存器列表
    freereg []bool     // 寄存器对应的状态
}

func NewCgen(tree *ASTNode, outfile *os.File) *Cgen {
    return &Cgen{
        tree:    tree,
        outfile: outfile,
        reglist: []string{"%r8", "%r9", "%r10", "%r11"},
        freereg: []bool{true, true, true, true},
    }
}

func (c *Cgen) GenAST() {
    c.cgpreamble()
    c.genAST(c.tree)
    c.cgpostamble()
}

func (c *Cgen) genAST(tree *ASTNode) {
    if tree != nil {
        switch tree.nodeKind {
        case PrintK, IfK:
            c.genStmt(tree)
        case OpK, ConstK, IdK:
            c.genExp(tree)
        default:
            c.error()
        }
        c.genAST(tree.sibling)
    }
}

func (c *Cgen) genStmt(tree *ASTNode) {
    switch tree.nodeKind {
    case PrintK:
        reg := c.genExp(tree.child[0])
        c.cgprintint(reg)
    default:
        c.error()
    }
}

func (c *Cgen) genExp(tree *ASTNode) int {
    var leftreg, rightreg int

    if len(tree.child) == 2 {
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
        default:
            return -1
        }
    case ConstK:
        return c.cgload(tree.intval)
    case IdK:
        return -1
    default:
        return -1
    }
}

/**************** 生成汇编语句 *****************/

func (c *Cgen) error() {
    panic(0)
}

// 设置所有寄存器为可用状态
func (c *Cgen) freeall_registers() {
    for i, _ := range c.freereg {
        c.freereg[i] = true
    }
    fmt.Println("After Freeall: ", c.freereg)
}

func (c *Cgen) alloc_register() int {
    for i := 0; i < 4; i++ {
        if c.freereg[i] {
            c.freereg[i] = false
            return i
        }
    }
    fmt.Println("Out of registers!")
    c.error()
    return 0
}

func (c *Cgen) free_register(reg int) {
    if c.freereg[reg] != false {
        fmt.Printf("Error trying to free register %d\n", reg)
        c.error()
    }
    c.freereg[reg] = true
}

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
	
	.globl  main
	.type   main, @function
main:
	pushq   %rbp
	movq	%rsp, %rbp

`)
}

func (c *Cgen) cgpostamble() {
    _, _ = c.outfile.WriteString(`
    movl	$0, %eax
	popq	%rbp
	ret
`)
}

// 寄存器赋值为value
func (c *Cgen) cgload(value int) int {
    //fmt.Println("value: ", value)
    r := c.alloc_register()
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t$%d, %s\n", value, c.reglist[r])
    return r
}

func (c *Cgen) cgadd(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\taddq\t%s, %s\n", c.reglist[r1], c.reglist[r2])
    c.free_register(r1)
    return r2
}

func (c *Cgen) cgsub(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\tsubq\t%s, %s\n", c.reglist[r2], c.reglist[r1])
    c.free_register(r2)
    return r1
}

func (c *Cgen) cgmul(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\timulq\t%s, %s\n", c.reglist[r1], c.reglist[r2])
    c.free_register(r1)
    return r2
}

func (c *Cgen) cgdiv(r1, r2 int) int {
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s,%%rax\n", c.reglist[r1])
    _, _ = fmt.Fprintf(c.outfile, "\tcqo\n")
    _, _ = fmt.Fprintf(c.outfile, "\tidivq\t%s\n", c.reglist[r2])
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%%rax,%s\n", c.reglist[r1])
    c.free_register(r2)
    return r1
}

func (c *Cgen) cgprintint(r int) {
    _, _ = fmt.Fprintf(c.outfile, "\tmovq\t%s, %%rdi\n", c.reglist[r])
    _, _ = fmt.Fprintf(c.outfile, "\tcall\tprintint\n")
    c.free_register(r)
}