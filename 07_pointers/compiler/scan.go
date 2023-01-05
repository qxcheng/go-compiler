package compiler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

// Error
var tokenError error = errors.New("not supported token")

// DFA的状态
type StateType int

const (
	START StateType = iota
	INCOMMENT
	INSTRING
	INNUM
	INID
	INEQ  // ==
	ININC // ++
	INLE  // <=
	INGE  // >=
	INNE  // !=
	DONE
)

// Scanner
type Scanner struct {
	file     *os.File // 源文件
	buf      *bufio.Reader
	linebuf  string // 当前行
	ch       int    // 当前字符
	linesize int    // 当前行的长度
	linepos  int    // 下一个待读取字符在当前行的位置
	err      error
	trace    map[int]bool
}

func NewScanner(file *os.File) *Scanner {
	s := Scanner{
		file:     file,
		buf:      bufio.NewReader(file),
		linesize: 0,
		linepos:  0,
	}
	s.next()
	return &s
}

func (s *Scanner) error(err error) {
	fmt.Println("Scan Error>> Line%d: %s\n, Position%d: %v\n", GLineno, s.linebuf, s.linepos, s.linebuf[s.linepos])
	panic(err)
}

// next 获取当前行的下一个非空字符，当前行无字符时读取新行
func (s *Scanner) next() {
	if !(s.linepos < s.linesize) {
		GLineno++

		var err error
		s.linebuf, err = s.buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				s.ch = -1
				return
			} else {
				s.error(err)
				return
			}
		}

		s.linesize = len(s.linebuf)
		s.linepos = 0
		s.ch = int(s.linebuf[s.linepos])
		s.linepos++
		return
	} else {
		s.ch = int(s.linebuf[s.linepos])
		s.linepos++
		return
	}
}

// 回退一个字符
func (s *Scanner) unget() {
	if s.ch != -1 {
		s.linepos--
	}
}

// 查看下一个字符
func (s *Scanner) prev() int {
	if s.linepos < s.linesize {
		return int(s.linebuf[s.linepos])
	} else {
		return -1
	}
}

func isalpha(c int) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isdigit(c int) bool {
	return '0' <= c && c <= '9'
}

// GetToken
//
//	@Description: 返回源文件的下一个记号
//	@receiver s
//	@return token 记号
//	@return lit 记号的值（如标识符名、数字）
func (s *Scanner) GetToken() (token Token, lit string) {
	var state StateType = START
	var save bool
	lit = ""

	for state != DONE {
		c := s.ch
		save = true

		switch state {
		case START:
			if isdigit(c) {
				state = INNUM
			} else if isalpha(c) {
				state = INID
			} else if c == ' ' || c == '\t' || c == '\n' {
				save = false
			} else {
				state = DONE
				switch c {
				case -1:
					save = false
					token = ENDFILE
				case '+':
					if s.prev() == '+' {
						state = ININC
					} else {
						token = ADD
					}
				case '-':
					token = SUB
				case '*':
					token = MUL
				case '%':
					token = REM
				case '/':
					if s.prev() == '/' {
						save = false
						state = INCOMMENT
					} else {
						token = QUO
					}
				case '=':
					if s.prev() == '=' {
						state = INEQ
					} else {
						token = ASSIGN
					}
				case '<':
					if s.prev() == '=' {
						state = INLE
					} else {
						token = LT
					}
				case '>':
					if s.prev() == '=' {
						state = INGE
					} else {
						token = GT
					}
				case '!':
					if s.prev() == '=' {
						state = INNE
					} else {
						token = NOT
					}
				case '"':
					state = INSTRING
				case '(':
					token = LPAREN
				case ')':
					token = RPAREN
				case '{':
					token = LBRACE
				case '}':
					token = RBRACE
				case ';':
					token = SEMI
				case '.':
					token = PERIOD
				case ',':
					token = COMMA
				default:
					token = ERROR
				}
			}
		case INCOMMENT:
			save = false
			if c == -1 {
				state = DONE
				token = ENDFILE
			} else if c == '\n' {
				state = START
			}
		case INSTRING:
			if c == '"' {
				state = DONE
				token = STRING
			}
		case INNUM:
			if !isdigit(c) {
				s.unget()
				save = false
				state = DONE
				token = NUM
			}
		case INID:
			if !isalpha(c) {
				s.unget()
				save = false
				state = DONE
				token = ID
			}
		case INEQ:
			state = DONE
			token = EQ
		case INGE:
			state = DONE
			token = GE
		case INLE:
			state = DONE
			token = LE
		case INNE:
			state = DONE
			token = NE
		case ININC:
			state = DONE
			token = INC
		case DONE:
		default:
			state = DONE
			token = ERROR
			s.error(tokenError)
		}

		if save {
			lit += string(c)
		}
		if state == DONE {
			if token == ID {
				// 如果是关键字
				if _token, ok := lit2token[lit]; ok {
					token = _token
				}
			}
		}
		s.next()
	}
	if GTraceScan {
		if s.trace == nil {
			s.trace = make(map[int]bool)
		}
		if s.trace[GLineno] {
			fmt.Printf("\tToken: %-8s, Lit: %s\n", tokens[token], lit)
		} else {
			fmt.Printf("Line%d: %s", GLineno, s.linebuf)
			fmt.Printf("\tToken: %-8s, Lit: %s\n", tokens[token], lit)
			s.trace[GLineno] = true
		}
	}
	return
}
