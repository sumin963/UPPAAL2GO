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

var path = "lexer_input.txt"
var dec_path = "cgo_input.txt"

type Token int

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
	file, err := os.Open(path)
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
	var tem_name []string
	for _, e := range doc.FindElements("./nta/*") {
		if e.Tag == "declaration" {
			dec = e.Text()
		}
		if e.Tag == "template" {
			if name := e.SelectElement("name"); name != nil {
				tem_name = append(tem_name, name.Text())
			}
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
		//fmt.Println("\n")
	}
	//fmt.Println(dec)
	//fmt.Println(tem_dec)
	dec_string_comment_del := string_comment_del(dec) // /**/ 제거
	tem_dec_string_comment_del := tem_string_comment_del(tem_dec)
	dec_comment_del := dec_line_comment_del(dec_string_comment_del) // //제거
	tem_dec_comment_del := tem_line_comment_del(tem_dec_string_comment_del)
	//
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC, // 파일이 없으면 생성,읽기/쓰기, 파일을 연 뒤 내용 삭제
		os.FileMode(0644))                // 파일 권한 666
	check(err)
	for i, _ := range dec_comment_del {
		_, err := file.Write([]byte(dec_comment_del[i] + "\n"))
		check(err)
	}
	for i, val := range tem_dec_comment_del {
		_, err := file.Write([]byte("//" + tem_name[i] + ";" + "\n"))
		check(err)
		for j, _ := range val {
			_, err := file.Write([]byte(tem_dec_comment_del[i][j] + "\n"))
			check(err)
		}
	}
	file.Close()
	rst_lexer, rst_token := Lexer_TADA()
	syntax, syntax_lex_data := parse_TADA(rst_lexer, rst_token)
	map_token_2_c(syntax, syntax_lex_data)
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
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func contains(elems []Token, v Token) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
func map_token_2_c(parse [][]Token, parse_lexr_data [][][]string) {
	input_file, err := os.Open(path)
	check(err)
	output_file, err := os.OpenFile(
		dec_path,
		os.O_CREATE|os.O_RDWR|os.O_TRUNC,
		os.FileMode(0644))
	check(err)
	reader := bufio.NewReader(input_file)
	input_file_reader := make([][]byte, 0)
	tem_name := make([]string, 0)
	//tem_val := make([][]interface{}, 0)
	for {
		line, _, err := reader.ReadLine()
		input_file_reader = append(input_file_reader, line)
		if err != nil {
			break
		}
	}
	fmt.Println(input_file_reader)
	defer input_file.Close()
	defer output_file.Close()
	fmt.Println(parse, "\n", parse_lexr_data, len(parse), len(parse_lexr_data))
	_local := false
	for i, _parse := range parse {
		fmt.Println(i, _parse)
		if _local && contains(_parse, ASSIGN) { //여기서 부터 시작;;;;;;;;;;;;;;;

		} else if _local && contains(_parse, IDENT) && contains(_parse, RPARENTHESIS) && contains(_parse, RBRACE) {
			//fmt.Println("local _ func:", parse_lexr_data[i])
			_start_line, _ := strconv.Atoi(parse_lexr_data[i][0][0])
			_end_line, _ := strconv.Atoi(parse_lexr_data[i][len(parse_lexr_data[i])-1][0])
			for j := _start_line - 1; j < _end_line; j++ {
				if j == _start_line-1 {
					_lparen := contain_index(parse[i], LPARENTHESIS)
					_index, _ := strconv.Atoi(parse_lexr_data[i][_lparen][1])
					fmt.Println(string(input_file_reader[j][:_index]) + tem_name[len(tem_name)-1] + " *" + tem_name[len(tem_name)-1] + " " + string(input_file_reader[j][_index:]))
				} else {
					if true { // 조건추가 ->

					}
					fmt.Println(string(input_file_reader[j]))
				}
			}

		} else if contains(_parse, ASSIGN) { //initializer
			if parse[i][0] == PREFIX && parse_lexr_data[i][0][2] == "const" { //const int N = 6;		#define N 6
				_ident := contain_index(parse[i], IDENT)
				_int := contain_index(parse[i], INT)
				_, err := output_file.Write([]byte("#define" + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _int) + "\n"))
				check(err)
				//fmt.Println("#define" + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + " " + mapping(parse_lexr_data, input_file_reader, i, _int))
			}
		} else if contains(_parse, IDENT) && contains(_parse, RPARENTHESIS) && contains(_parse, RBRACE) { //func 수정필요
			//파라미터에 Local *local
			//local->list 구조체 멤버 접근시
		} else if contains(_parse, DIV) { //바꾸어야 할지도\
			tem_name = append(tem_name, parse_lexr_data[i][2][2])
			//tem_val[len(tem_name)-1] = make([]interface{}, 0)
			_local = true
		} else { //dec
			//chan ,clock
			//int[0,6]
			//로컬의 경우 struct로
			if _local {
				if contains(_parse, CLOCK) {
				} else {
					_lbracket := contain_index(parse[i], LBRACKET)
					if parse[i][_lbracket+1] == COMMA { //[,] check
						fmt.Println("55", _lbracket, parse_lexr_data[i], _parse)
					} else {
						if contains(_parse, COMMA) { //[],[] check
						} else {
							// _ident := contain_index(parse[i], IDENT)
							// _int := contain_index(parse[i], INT)
							// _, err := output_file.Write([]byte("#define" + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _int) + "\n"))
							// check(err)      아래를 주석과 같은 형태로 변환
							_start_line, _ := strconv.Atoi(parse_lexr_data[i][0][0])
							_scd_line, _ := strconv.Atoi(parse_lexr_data[i][1][0])
							fmt.Println("56", string(input_file_reader[_start_line][0]), string(input_file_reader[_start_line][_scd_line]), parse_lexr_data[i])
						}
					}
				}
				//id_t list[N+1];
			} else {
				if contains(_parse, CLOCK) {
				} else if contains(_parse, CHANNEL) {
				} else if contains(_parse, TYPEDEF) { //typedef int[0,6] id_t;     typedef int id_t;
					if contains(_parse, COMMA) {
						_typedef := contain_index(parse[i], TYPEDEF)
						_typeid := contain_index(parse[i], TYPEID)
						_ident := contain_index(parse[i], IDENT)
						_, err := output_file.Write([]byte(mapping(parse_lexr_data, input_file_reader, i, _typedef) + " " + mapping(parse_lexr_data, input_file_reader, i, _typeid) + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + "\n"))
						check(err)
					}

				} else {

				}
			}
		}

	}
}
func contain_index(s []Token, substr Token) int {
	for i, v := range s {
		if v == substr {
			return i
		}
	}
	return 0
}
func mapping(parse_lexr_data [][][]string, input_file_reader [][]byte, i int, j int) string {
	if len(parse_lexr_data[i])-1 > j {
		return string(input_file_reader[ext_start_index(parse_lexr_data, i, j)][ext_scd_index(parse_lexr_data, i, j):ext_thd_index(parse_lexr_data, i, j)])

	}
	return string(input_file_reader[ext_start_index(parse_lexr_data, i, j)][ext_scd_index(parse_lexr_data, i, j):])

}
func ext_start_index(parse_lexr_data [][][]string, i int, j int) int {
	rst, _ := strconv.Atoi(parse_lexr_data[i][j][0])
	return rst - 1
}
func ext_scd_index(parse_lexr_data [][][]string, i int, j int) int {
	rst, _ := strconv.Atoi(parse_lexr_data[i][j][1])
	return rst - 1
}
func ext_thd_index(parse_lexr_data [][][]string, i int, j int) int {

	rst, _ := strconv.Atoi(parse_lexr_data[i][j+1][1]) //i=line, j=token, h 0 = start line, h 1 start string, h 2 token
	return rst - 1

}
func parse_TADA(lexer_data [][]string, token []Token) ([][]Token, [][][]string) {
	stack_b := make([]Token, 0)
	stack_token := make([]Token, 0)
	syntax := make([][]Token, 0)
	stack_lex := make([][]string, 0)
	syntax_lex_data := make([][][]string, 0)
	var _pop_item_b Token
	var _pop_item_token Token
	var _comma_token Token
	var _comma_lex []string
	var _comma bool = false
	for i, _ := range lexer_data {

		switch token[i] {
		case SEMI:

			if len(stack_b) == 0 {
				_syntax := stack_token
				_syntax_lex_data := stack_lex
				stack_token = make([]Token, 0)
				stack_lex = make([][]string, 0)
				syntax = append(syntax, _syntax)
				syntax_lex_data = append(syntax_lex_data, _syntax_lex_data)
			}
		case RBRACE: //}

			_pop_item_b = stack_b[len(stack_b)-1]
			stack_b = stack_b[:len(stack_b)-1]
			if len(stack_b) == 0 {
				for {
					_pop_item_token = stack_token[len(stack_token)-1]
					stack_token = stack_token[:len(stack_token)-1]
					_pop_lex_token := stack_lex[len(stack_lex)-1]
					stack_lex = stack_lex[:len(stack_lex)-1]

					if _pop_item_token == _pop_item_b {
						stack_token = append(stack_token, _pop_item_token, token[i])
						stack_lex = append(stack_lex, _pop_lex_token, lexer_data[i])

						_syntax := stack_token
						_syntax_lex_data := stack_lex
						stack_token = make([]Token, 0)
						stack_lex = make([][]string, 0)
						syntax = append(syntax, _syntax)
						syntax_lex_data = append(syntax_lex_data, _syntax_lex_data)
						break
					}
				}
			} else {
				for {
					_pop_item_token = stack_token[len(stack_token)-1]
					stack_token = stack_token[:len(stack_token)-1]
					stack_lex = stack_lex[:len(stack_lex)-1]
					if _pop_item_token == _pop_item_b {
						break
					}
				}
			}
		case RBRACKET, RPARENTHESIS:
			_pop_item_b = stack_b[len(stack_b)-1]
			stack_b = stack_b[:len(stack_b)-1]
			for {
				_pop_item_token = stack_token[len(stack_token)-1]
				stack_token = stack_token[:len(stack_token)-1]
				_pop_lex_token := stack_lex[len(stack_lex)-1]
				stack_lex = stack_lex[:len(stack_lex)-1]
				if _pop_item_token == COMMA {
					_comma_token = _pop_item_token
					_comma_lex = _pop_lex_token
					_comma = true
				}
				if _pop_item_token == _pop_item_b {
					stack_token = append(stack_token, _pop_item_token)
					stack_lex = append(stack_lex, _pop_lex_token)
					if _comma {
						stack_token = append(stack_token, _comma_token)
						stack_lex = append(stack_lex, _comma_lex)
						_comma = false
					}
					stack_token = append(stack_token, token[i])
					stack_lex = append(stack_lex, lexer_data[i])
					break
				}
			}

		case LBRACE, LBRACKET, LPARENTHESIS: //{
			stack_b = append(stack_b, token[i])
			stack_token = append(stack_token, token[i])
			stack_lex = append(stack_lex, lexer_data[i])

		default:
			stack_token = append(stack_token, token[i])
			stack_lex = append(stack_lex, lexer_data[i])
			if len(stack_token) >= 3 { //ident _ ident 처리
				if stack_token[len(stack_token)-2] == UNDERBAR && stack_token[len(stack_token)-3] == IDENT {
					stack_token = stack_token[:len(stack_token)-2]
					stack_lex = stack_lex[:len(stack_lex)-2]
				}
			}
		}
	}
	return syntax, syntax_lex_data
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
