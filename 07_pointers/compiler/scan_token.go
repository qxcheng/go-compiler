package compiler

/* 词法分析相关 */

type Token int

const (
	// 以下为内部记号
	ENDFILE Token = iota
	ERROR

	// 基础类型
	IDENT // main
	INT
	FLOAT
	CHAR
	STRING

	// 以下为多字符记号
	ID
	NUM

	// 以下为特殊符号
	ASSIGN // =
	EQ     // ==
	LT     // <
	GT     // >
	GE     // >=
	LE     // <=
	NE     // !=

	AND    // &&
	OR     // ||
	NOT    // !

	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %
	INC // ++

	LPAREN // (
	RPAREN // )
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	RBRACE // }
	SEMI   // ;
	COMMA  // ,
	PERIOD // .
	COLON  // :

	// 以下为关键字
	IF
	ELSE
	FOR
	BREAK
	CONTINUE

	PACKAGE
	IMPORT

	VAR
	FUNC
	PRINT
	RETURN
)

var tokens = [...]string{
	"ENDFILE",
	"ERROR",

	// 基础类型
	"IDENT", // main
	"INT",
	"FLOAT",
	"CHAR",
	"STRING",

	// 以下为多字符记号
	"ID",
	"NUM",

	// 以下为特殊符号
	"ASSIGN", // =
	"EQ",     // ==
	"LT",     // <
	"GT",     // >
	"GE",     // >=
	"LE",     // <=
	"NE",     // !=

	"AND",    // &&
	"OR",     // ||
	"NOT",    // !

	"ADD", // +
	"SUB", // -
	"MUL", // *
	"QUO", // /
	"REM", // %
	"INC", // ++

	"LPAREN", // (
	"RPAREN", // )
	"LBRACK", // [
	"RBRACK", // ]
	"LBRACE", // {
	"RBRACE", // }
	"SEMI",   // ;
	"COMMA",  // ,
	"PERIOD", // .
	"COLON",  // :

	// 以下为关键字
	"IF",
	"ELSE",
	"FOR",
	"BREAK",
	"CONTINUE",

	"PACKAGE",
	"IMPORT",

	"VAR",
	"FUNC",
	"PRINT",
	"RETURN",
}

var lit2token = map[string]Token{
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"package":  PACKAGE,
	"import":   IMPORT,
	"var":      VAR,
	"func":     FUNC,
	"main":     IDENT,
	"int":      INT,
	"float":    FLOAT,
	"char":     CHAR,
	"print":    PRINT,
	"return":   RETURN,
}
