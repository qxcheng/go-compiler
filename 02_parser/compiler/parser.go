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
statement -> if-stmt|for-stmt|assign-stmt

if-stmt -> if exp [stmt-sequence] [else stmt-sequence]
for-stmt -> for assign-stmt;exp;exp [stmt-sequence]
assign-stmt -> identifier = exp

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
}

func NewParser(file *os.File) *Parser {
    s := NewScanner(file)
    p := Parser{
        s: s,
    }
    p.curToken, p.curLit = p.s.GetToken()
    return &p
}

func (p *Parser) error() {
    fmt.Println("Parse Error>> Line%d: %s\n, Position%d: %v\n", GLineno, p.s.linebuf, p.s.linepos, p.s.linebuf[p.s.linepos])
    panic(0)
}

func (p *Parser) match(token Token) {
    if p.curToken == token {
        p.curToken, p.curLit = p.s.GetToken()
    } else {
        p.error()
    }
}

// 语法树解析
func (p *Parser) Parse() *ASTNode {
    var t *ASTNode
    for p.curToken != ENDFILE {
        t = p.exp()
    }
    t.printTree(0)
    return t
}

func (p *Parser) exp() *ASTNode {
    t := p.simple_exp()
    if p.curToken == EQ || p.curToken == LT || p.curToken == GT {
        n := NewASTNode(OpK)
        n.child[0] = t
        n.token = p.curToken
        t = n
        p.match(p.curToken)
        n.child[1] = p.simple_exp()
    }
    return t
}

// 递归：简单表达式：+ -
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

// 递归：* /
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

// 递归：因子：exp | 常量 | 变量
func (p *Parser) factor() *ASTNode {
    var t *ASTNode = nil
    switch p.curToken {
    case NUM:
        t = NewASTNode(ConstK)
        t.intval, _ = strconv.Atoi(p.curLit)
        p.match(NUM)
    case ID:
        t = NewASTNode(IdK)
        t.litval = p.curLit
        p.match(ID)
    case LPAREN:
        p.match(LPAREN)
        t = p.exp()
        p.match(RPAREN)
    default:
        p.error()
    }
    return t
}


