/*
program -> declare-sequence
declare-sequence -> declare{;declare]
declare -> var-declare|func-declare

var-declare -> var identifier type
var-type -> int|float|string

func-declare -> func identifier(identifier var-type{, identifier var-type}) (identifier var-type{, identifier var-type}) {
    stmt-sequence
}

stmt-sequence -> statement{;statement]
statement -> if-stmt|for-stmt|assign-stmt|print-stmt

if-stmt -> if exp [stmt-sequence] [else stmt-sequence]
for-stmt -> for assign-stmt;exp;exp [stmt-sequence]
assign-stmt -> identifier = exp
print-stmo -> print exp

exp -> simple-exp[comparison-op simple-exp]
comparison-op -> <|=
simple-exp -> term{addop term}
addop -> +|-
term -> factor{mulop factor}
mulop -> *|/
factor -> (exp)| number | identifier
*/

package compiler

import (
    "fmt"
    "os"
    "strconv"
)

type Parser struct {
    s *Scanner
    curToken Token  // 当前token
    curLit string   // 当前lit
    cacheToken Token  // 向前查看一个token
    cacheLit string
    currentFunc int  // 当前函数的插槽id 用于return语句
}

func NewParser(file *os.File) *Parser {
    s := NewScanner(file)
    p := Parser{
        s: s,
        cacheToken: -1,
        cacheLit: "",
    }
    p.curToken, p.curLit = p.s.GetToken()
    return &p
}

func (p *Parser) error(msg string) {
    fmt.Printf("Parse Error>> Line %d: %s\n, Position %d: %v\n", GLineno, p.s.linebuf, p.s.linepos, p.s.linebuf[p.s.linepos])
    panic(msg)
}

func (p *Parser) match(token Token) {
    if p.curToken == token {
        if p.cacheToken != -1 {
            p.curToken, p.curLit = p.cacheToken, p.cacheLit
            p.cacheToken = -1
            p.cacheLit = ""
        } else {
            p.curToken, p.curLit = p.s.GetToken()
        }
    } else {
        p.error("Error: token not match")
    }
}

func (p *Parser) prev() Token {
    if p.cacheToken != -1 {
        return p.cacheToken
    } else {
        p.cacheToken, p.cacheLit = p.s.GetToken()
        return p.cacheToken
    }
}

// 语法树解析
func (p *Parser) Parse() *ASTNode {
    var t *ASTNode

    t = p.stmt_sequence()
    if GTraceParse{
        t.printTree(0)
    }

    return t
}

// 递归：语句序列
func (p *Parser) stmt_sequence() *ASTNode {
    t := p.statement()  // t指向第一个语句
    var n *ASTNode = t

    for p.curToken != ENDFILE {
        if p.curToken == SEMI {
            p.match(SEMI)
        }
        if p.curToken == RBRACE {
            break  // 匹配到右大括号意味着语句序列的结束
        }
        q := p.statement()
        n.sibling = q
        n = q
    }
    return t
}

// 递归：语句类型
func (p *Parser) statement() *ASTNode {
    var t *ASTNode
    switch p.curToken {
    case PRINT:
        t = p.print_stmt()
    case VAR:
        t = p.var_declaration()
    case FUNC:
        t = p.func_declaration()
    case ID:
        t = p.assign_stmt()
    case IF:
        t = p.if_stmt()
    case FOR:
        t = p.for_stmt()
    case RETURN:
        t = p.return_stmt()
    default:
        return nil
    }
    return t
}

// 添加变量到符号表
func (p *Parser) addglob(token Token, name string) int {
    var i int  // 变量的插槽位置
    switch token {
    case CHAR:
        i = Gsym.Addglob(name, VAR_CHAR)
    case INT:
        i = Gsym.Addglob(name, VAR_INT)
    default:
        p.error("Parse error: unspported vartype")
    }
    return i
}

// 声明: 变量
func (p *Parser) var_declaration() *ASTNode {
    t := NewASTNode(VarK)
    p.match(VAR)
    t.child[0] = NewASTNode(IdK)
    t.child[0].litval = p.curLit
    p.match(ID)
    t.child[0].symbleid = p.addglob(p.curToken, t.child[0].litval)
    //fmt.Println("222: ", t.child[0].symbleid, t.child[0].litval)
    p.match(p.curToken)
    return t
}

// 声明：函数
func (p *Parser) func_declaration() *ASTNode {
    t := NewASTNode(FuncK)
    p.match(FUNC)
    t.token = p.curToken  // ID 或 IDENT(main)
    t.litval = p.curLit   // 函数名
    t.symbleid = Gsym.Addglob(t.litval, VAR_FUNC)  // 添加到符号表
    p.currentFunc = t.symbleid
    p.match(p.curToken)
    p.match(LPAREN)
    // 参数解析，暂支持一个参数
    if p.curToken == ID {
        t.child[0] = NewASTNode(IdK)
        t.child[0].litval = p.curLit
        p.match(ID)
        t.child[0].symbleid = p.addglob(p.curToken, t.child[0].litval)  // 暂时放在符号表
        t.child[0].token = p.curToken  // 暂用于保存形参变量类型
        p.match(p.curToken)
    }
    p.match(RPAREN)
    // 返回值类型解析
    if p.curToken != LBRACE {
        switch p.curToken {
        case CHAR:
            Gsym.Setglob(t.symbleid, VAR_CHAR)
        case INT:
            Gsym.Setglob(t.symbleid, VAR_INT)
        default:
            p.error("not supported return type")
        }
        p.match(p.curToken)
    }
    p.match(LBRACE)
    t.child[1] = p.stmt_sequence()
    p.match(RBRACE)
    return t
}

// 语句：返回语句
func (p *Parser) return_stmt() *ASTNode {
    t := NewASTNode(ReturnK)
    t.symbleid = p.currentFunc
    p.match(RETURN)
    t.child[0] = p.exp()
    //switch p.curToken {
    //case ID:
    //    t.child[0] = NewASTNode(IdK)
    //    t.child[0].litval = p.curLit
    //    t.child[0].symbleid = Gsym.Findglob(t.child[0].litval)
    //    if t.child[0].symbleid == -1 {
    //        p.error("Parse error: return undefined var")
    //    }
    //    p.match(ID)
    //case NUM:
    //    t.child[0] = NewASTNode(ConstK)
    //    t.child[0].intval, _ = strconv.Atoi(p.curLit)
    //    p.match(NUM)
    //}
    return t
}


// 语句：赋值语句
func (p *Parser) assign_stmt() *ASTNode {
    t := NewASTNode(AssignK)
    t.litval = p.curLit
    t.symbleid = Gsym.Findglob(t.litval)
    if t.symbleid == -1 {
        p.error("Parse error: use undefined variable")
    }
    p.match(ID)
    p.match(ASSIGN)
    t.child[0] = p.exp()
    return t
}

// 语句：输出语句
func (p *Parser) print_stmt() *ASTNode {
    t := NewASTNode(PrintK)
    p.match(PRINT)
    t.child[0] = p.exp()
    return t
}

// 语句：条件语句
func (p *Parser) if_stmt() *ASTNode {
    t := NewASTNode(IfK)
    p.match(IF)
    t.child[0] = p.exp()
    p.match(LBRACE)
    t.child[1] = p.stmt_sequence()
    p.match(RBRACE)
    if p.curToken == ELSE {
        p.match(ELSE)
        p.match(LBRACE)
        t.child[2] = p.stmt_sequence()
        p.match(RBRACE)
    }
    return t
}

// 语句：循环语句
func (p *Parser) for_stmt() *ASTNode {
    t := NewASTNode(ForK)
    p.match(FOR)
    t.child[0] = p.exp()
    p.match(LBRACE)
    t.child[1] = p.stmt_sequence()
    p.match(RBRACE)
    return t
}


// 表达式： == < >
func (p *Parser) exp() *ASTNode {
    t := p.simple_exp()
    if p.curToken == EQ || p.curToken == LT || p.curToken == GT || p.curToken == GE || p.curToken == LE || p.curToken == NE {
        n := NewASTNode(OpK)
        n.child[0] = t
        n.token = p.curToken
        t = n
        p.match(p.curToken)
        n.child[1] = p.simple_exp()
    }
    return t
}

// 表达式：+ -
func (p *Parser) simple_exp() *ASTNode {
    t := p.term()
    for p.curToken == ADD || p.curToken == SUB {
        n := NewASTNode(OpK)
        n.child[0] = t
        n.token = p.curToken
        t = n
        p.match(p.curToken)
        n.child[1] = p.term()
    }
    return t
}

// 表达式：* /
func (p *Parser) term() *ASTNode {
    t := p.factor()
    for p.curToken == MUL || p.curToken == QUO {
        n := NewASTNode(OpK)
        n.child[0] = t
        n.token = p.curToken
        t = n
        p.match(p.curToken)
        n.child[1] = p.factor()
    }
    return t
}

// 表达式：exp | 常量 | 变量
func (p *Parser) factor() *ASTNode {
    var t *ASTNode = nil
    switch p.curToken {
    case NUM:
        t = NewASTNode(ConstK)
        t.intval, _ = strconv.Atoi(p.curLit)
        p.match(NUM)
    case ID:
        if p.prev() == LPAREN {
            t = NewASTNode(CallK)
            t.litval = p.curLit  // 函数名
            t.symbleid = Gsym.Findglob(t.litval)
            if t.symbleid == -1 {
                p.error("Parse error: use undefined var")
            }
            p.match(ID)
            p.match(LPAREN)
            switch p.curToken {
            case ID:
                t.child[0] = NewASTNode(IdK)
                t.child[0].litval = p.curLit  // 变量名
                t.child[0].symbleid = Gsym.Findglob(p.curLit)
                if t.symbleid == -1 {
                    p.error("Parse error: use undefined var")
                }
            case NUM:
                t.child[0] = NewASTNode(ConstK)
                t.child[0].intval, _ = strconv.Atoi(p.curLit)
            }
            p.match(p.curToken)
            p.match(RPAREN)
        } else {
            t = NewASTNode(IdK)
            t.litval = p.curLit
            t.symbleid = Gsym.Findglob(t.litval)
            if t.symbleid == -1 {
                p.error("Parse error: use undefined var")
            }
            p.match(ID)
        }
    case LPAREN:
        p.match(LPAREN)
        t = p.exp()
        p.match(RPAREN)
    default:
        p.error("Error: undefined token")
    }
    return t
}


