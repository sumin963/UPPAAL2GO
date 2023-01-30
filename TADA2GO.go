package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/beevik/etree"
	. "github.com/dave/jennifer/jen"
)

var new_type []string
var path = "/input.test"

type Token int
type Stack []interface{}

const (
	EOF     = iota //프로그램의 끝
	ILLEGAL        //정의되지않은 문자
	IDENT
	INT
	STRING
	BOOL
	SEMI     // ;
	COMMA    //,
	UNDERBAR //_

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =

	//괄호
	LPARENTHESIS //(
	RPARENTHESIS //)
	LBRACE       //{
	RBRACE       //}
	LBRACKET     //[
	RBRACKET     //]

	PREFIX
	CHANNEL
	CLOCK
	TYPEDEF
	TYPEID
	RETURN
)

var tokens = []string{
	EOF:      "EOF",
	ILLEGAL:  "ILLEGAL",
	IDENT:    "IDENT",
	INT:      "INT",
	BOOL:     "BOOL",   //
	STRING:   "STRING", //
	SEMI:     ";",
	COMMA:    ",",
	UNDERBAR: "_",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	LPARENTHESIS: "(",
	RPARENTHESIS: ")",
	LBRACE:       "{",
	RBRACE:       "}",
	LBRACKET:     "[",
	RBRACKET:     "]",

	ASSIGN: "=",

	PREFIX:  "PREPIX",
	CHANNEL: "CHANNEL",
	CLOCK:   "CLOCK",
	TYPEDEF: "TYPEDEF",
	TYPEID:  "TYPEID",
	RETURN:  "RETURN",
}

// IsEmpty - 스택이 비어있는지 확인하는 함수
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push - 스택에 값을 추가하는 함수.
func (s *Stack) Push(data interface{}) {
	*s = append(*s, data) // 스택 끝(top)에 값을 추가함.
	//fmt.Printf("%d pushed to stack\n", data)
}

// Pop - 스택에 값을 제거하고 top위치에 값을 반환하는 함수.
func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		//  fmt.Println("stack is empty")
		return nil
	} else {
		top := len(*s) - 1
		data := (*s)[top] // top 위치에 있는 값을 가져 옴
		*s = (*s)[:top]   // 스택에 마지막 데이터 제거함
		return data

	}
}
func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer { //하나의 토큰을 반환하는 함수
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) Lex() (Position, Token, string) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, ""
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			panic(err)
		}

		// update the column to the position of the newly read in rune
		l.pos.column++

		switch r {
		case '\n':
			l.resetPosition()
		case ';':
			return l.pos, SEMI, ";"
		case ',':
			return l.pos, COMMA, ","
		case '_':
			return l.pos, UNDERBAR, "_"
		case '+':
			return l.pos, ADD, "+"
		case '-':
			return l.pos, SUB, "-"
		case '*':
			return l.pos, MUL, "*"
		case '/':
			return l.pos, DIV, "/"
		case '=':
			return l.pos, ASSIGN, "="
		case '(':
			return l.pos, LPARENTHESIS, "("
		case ')':
			return l.pos, RPARENTHESIS, ")"
		case '{':
			return l.pos, LBRACE, "{"
		case '}':
			return l.pos, RBRACE, "}"
		case '[':
			return l.pos, LBRACKET, "["
		case ']':
			return l.pos, RBRACKET, "]"
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsDigit(r) {
				// backup and let lexInt rescan the beginning of the int
				startPos := l.pos
				l.backup()
				lit := l.lexInt()
				return startPos, INT, lit
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				l.backup()
				lit := l.lexIdent()
				//fmt.Println(lit) //switch
				switch lit {
				case "return":
					return startPos, RETURN, lit
				case "typedef":
					return startPos, TYPEDEF, lit
				case "chan":
					return startPos, CHANNEL, lit
				case "clock":
					return startPos, CLOCK, lit
				case "urgent":
					fallthrough
				case "broadcast":
					fallthrough
				case "const":
					return startPos, PREFIX, lit
				case "void":
					fallthrough
				case "int":
					fallthrough
				case "bool":
					fallthrough
				case "string":
					fallthrough
				case "double":
					return startPos, TYPEID, lit
				default:
					return startPos, IDENT, lit
				}
			} else {
				return l.pos, ILLEGAL, string(r)
			}
		}
	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) backup() {
	if err := l.reader.UnreadRune(); err != nil {
		panic(err)
	}

	l.pos.column--
}

// lexInt scans the input until the end of an integer and then returns the
// literal.
func (l *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the int
				return lit
			}
		}

		l.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the integer
			l.backup()
			return lit
		}
	}
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() string {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit
			}
		}

		l.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return lit
		}
	}
}

func Lexer_TADA() ([][]string, []Token) {
	file, err := os.Open("input.test")
	lexer_data := make([][]string, 0)
	lexer_token_data := make([]Token, 0)
	if err != nil {
		panic(err)
	}

	lexer := NewLexer(file)
	for {
		pos, tok, lit := lexer.Lex()
		if tok == EOF {
			break
		}
		_lexer_data := make([]string, 0)
		_lexer_data = append(_lexer_data, strconv.Itoa(pos.line), strconv.Itoa(pos.column), lit)
		lexer_data = append(lexer_data, _lexer_data)
		lexer_token_data = append(lexer_token_data, tok)
		//fmt.Printf("%d:%d\t%s\t%s\n", pos.line, pos.column, tok, lit)
	}
	return lexer_data, lexer_token_data
}

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("C:\\Users\\jsm96\\gitfolder\\UPPAAL2GO\\TADA.xml"); err != nil {
		panic(err)
	}
	var dec string
	var tem_dec []string
	for _, e := range doc.FindElements("./nta/*") {
		if e.Tag == "declaration" {
			dec = e.Text()
		}
		if e.Tag == "template" {
			if declaration := e.SelectElement("declaration"); declaration != nil {
				tem_dec = append(tem_dec, declaration.Text())

			}
			for _, l := range e.FindElements("location") {
				if l.Attr[0].Key == "id" {
				}
				if l_label := l.SelectElement("label"); l_label != nil {
				}
			}

			for _, t := range e.FindElements("transition") {
				if t_source := t.SelectElement("source"); t_source != nil {
				}
				if t_target := t.SelectElement("target"); t_target != nil {
				}
				for _, l := range t.FindElements("label") {
					if l.Attr[0].Value == "select" {
					} else if l.Attr[0].Value == "guard" {
					} else if l.Attr[0].Value == "synchronisation" {
					} else if l.Attr[0].Value == "assignment" {
					}
				}
			}
		}
		fmt.Println("\n")
	}
	//fmt.Println(dec)
	//fmt.Println(tem_dec)
	dec_string_comment_del := string_comment_del(dec)
	tem_dec_string_comment_del := tem_string_comment_del(tem_dec)
	dec_comment_del := dec_line_comment_del(dec_string_comment_del)
	tem_dec_comment_del := tem_line_comment_del(tem_dec_string_comment_del)
	fmt.Println(dec_comment_del)
	fmt.Println(tem_dec_comment_del)
	//
	file, err := os.Create("hello.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	for i, _ := range dec_comment_del {
		n, err := file.Write([]byte(dec_comment_del[i] + "\n"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(n, "바이트 저장 완료")
	}
	for i, val := range tem_dec_comment_del {
		for j, _ := range val {
			n, err := file.Write([]byte(tem_dec_comment_del[i][j] + "\n"))
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(n, "바이트 저장 완료")
		}
	}
	rst_lexer, rst_token := Lexer_TADA()
	map_token_2_c(parse_TADA(rst_lexer, rst_token), rst_lexer)
	//
	f := NewFile("a")
	for i, _ := range dec_comment_del {
		f.Comment(dec_comment_del[i])
	}
	for i, val := range tem_dec_comment_del {
		for j, _ := range val {
			f.Comment(tem_dec_comment_del[i][j])
		}
	}
	//fmt.Printf("%#v", f)
}
func map_token_2_c(parse [][]interface{}, lexer_data [][]string) {

}
func parse_TADA(lexer_data [][]string, token []Token) [][]interface{} {
	var stack_b Stack
	var stack_token Stack
	var stack_syntax Stack
	syntax := make([][]interface{}, 0)
	for i, _ := range lexer_data {
		switch token[i] {
		case SEMI:
			if stack_b.IsEmpty() {
				//fmt.Println(stack_token)
				_syntax := make([]interface{}, 0)
				for !stack_token.IsEmpty() {
					_token := stack_token.Pop()
					stack_syntax.Push(_token)
				}
				for !stack_syntax.IsEmpty() {
					_token := stack_syntax.Pop()
					_syntax = append(_syntax, _token)
				}
				syntax = append(syntax, _syntax)
			}
		case RBRACE: //}
			_pop_item_b := stack_b.Pop()
			if stack_b.IsEmpty() { //stack_b가 비어있을때
				for {
					_pop_item_token := stack_token.Pop()
					if _pop_item_token == _pop_item_b {
						stack_token.Push(_pop_item_token)
						stack_token.Push(token[i])
						//fmt.Println(stack_token)
						_syntax := make([]interface{}, 0)
						for !stack_token.IsEmpty() {
							_token := stack_token.Pop()
							stack_syntax.Push(_token)
						}
						for !stack_syntax.IsEmpty() {
							_token := stack_syntax.Pop()
							_syntax = append(_syntax, _token)
						}
						syntax = append(syntax, _syntax)
						// for !stack_token.IsEmpty() {
						// 	stack_token.Pop()
						// }
						break
					}
				}
			} else {
				for {
					_pop_item_token := stack_token.Pop()
					if _pop_item_token == _pop_item_b {

						break
					}
				}
			}
		case RBRACKET, RPARENTHESIS:
			_pop_item_b := stack_b.Pop()
			//_comma := false
			for {
				_pop_item_token := stack_token.Pop()
				if _pop_item_token == _pop_item_b {
					stack_token.Push(_pop_item_token)
					// if token[i] == RBRACKET && _comma {
					// 	stack_token.Push(COMMA)
					// 	_comma = false
					// }
					stack_token.Push(token[i])
					break
				}
			}

		case LBRACE, LBRACKET, LPARENTHESIS: //{
			stack_b.Push(token[i])
			stack_token.Push(token[i])
		default:
			stack_token.Push(token[i])
		}
	}
	return syntax
}

func string_comment_del(dec string) string {
	string_counts := strings.Count(dec, "/*")
	dec_comment_del := ""
	i := 0
	for i < string_counts {
		start := strings.Index(dec, "/*")
		end := strings.Index(dec, "*/")
		dec_comment_del += dec[0:start] + dec[end+2:len(dec)]
		i++
	}
	return dec_comment_del
}
func dec_line_comment_del(dec string) []string {
	dec_silce := strings.Split(dec, "\n")
	for i, val := range dec_silce {
		if strings.Index(val, "//") != (-1) {
			dec_silce[i] = dec_silce[i][:strings.Index(val, "//")]
		}
	}
	return dec_silce
}
func tem_line_comment_del(dec []string) [][]string {
	dec_comment_del := make([][]string, 0)
	for _, val := range dec {
		array := dec_line_comment_del(val)
		dec_comment_del = append(dec_comment_del, array)
	}

	return dec_comment_del
}

func tem_string_comment_del(dec []string) []string {
	dec_comment_del := make([]string, len(dec))

	for i := 0; i < len(dec); i++ {
		string_counts := strings.Count(dec[i], "/*")
		if string_counts > 0 {
			dec_comment_del_index := ""
			j := 0
			for j < string_counts {
				start := strings.Index(dec[i], "/*")
				end := strings.Index(dec[i], "*/")
				dec_comment_del_index += dec[i][0:start] + dec[i][end+2:len(dec)]
				j++
			}
			dec_comment_del[i] = dec_comment_del_index
		} else {
			dec_comment_del[i] = dec[i]
		}
	}
	return dec_comment_del
}
