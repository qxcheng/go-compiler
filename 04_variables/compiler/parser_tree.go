package compiler

import (
    "fmt"
    "strings"
)

/* 语法分析相关 */
// 节点类型
type NodeKind int
const (
    // 声明类型
    VarK NodeKind = iota
    FuncK

    // 语句类型
    IfK
    ForK
    AssignK
    PrintK

    // 表达式类型
    OpK
    ConstK
    IdK
)

// 语法树
type ASTNode struct {
    child []*ASTNode   // 子节点
    sibling *ASTNode   // 兄弟节点

    nodeKind NodeKind
    token Token
    intval int     // 数字
    litval string  // 标识符名
}

func NewASTNode(nodeKind NodeKind) *ASTNode {
    var childLen int
    switch nodeKind {
    case OpK:
        childLen = 2
    case ConstK:
        childLen = 0
    case PrintK, VarK, AssignK:
        childLen = 1
    }

    if childLen == 0 {
        return &ASTNode{
            nodeKind:nodeKind,
        }
    } else {
        return &ASTNode{
            child: make([]*ASTNode, childLen),
            nodeKind:nodeKind,
        }
    }
}

func (t *ASTNode) printTree(level int) {
    tab := strings.Repeat(" ", level)
    switch t.nodeKind {
    case OpK:
        fmt.Printf("%sOp: %s\n", tab, tokens[t.token])
    case ConstK:
        fmt.Printf("%sConst: %d\n", tab, t.intval)
    case IdK:
        fmt.Printf("%sId: %d\n", tab, t.litval)
    case PrintK:
        fmt.Printf("%sStmt: %d\n", tab, tokens[t.token])
    }
    for _, child := range t.child {
        if child != nil {
            child.printTree(level+4)
        }
    }
    if t.sibling != nil {
        t.sibling.printTree(0)
    }
}