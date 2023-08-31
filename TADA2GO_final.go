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

// 입력할 xml파일 경로
var input_xml = "C:\\Users\\jsm96\\gitfolder\\UPPAAL2GO\\av_algo_tada.xml"

var path = "lexer_input.txt"
var dec_path = "cgo_input.txt"

func main() {
	//etree 객체 생성 및 xml 파일 열기
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(input_xml); err != nil { //TADA.xml
		panic(err)
	}

	//TADA의 정보가 담긴 xml 파일의 내용을 etree에 삽입
	//input : doc  //xml파일
	//output : xml_loc - location data, xml_transition - transition data, xml_dec - dec data, xml_tem - tem data
	xml_loc, xml_transition, xml_dec, xml_tem := input_xml_to_slice(doc)
	dec := xml_dec[0]
	sys_dec := xml_dec[1]
	tem_dec := xml_tem[0]
	tem_name := xml_tem[1]
	tem_param := xml_tem[2]

	//xml파일의 declalation data를 lexer에 input하기 위해 전처리 과정
	//declalation의 전처리
	dec_string_comment_del := string_comment_del(dec) // /**/ 제거
	//template declalation의 전처리
	tem_dec_string_comment_del := tem_string_comment_del(tem_dec)
	//system declalation의 전처리
	sys_dec_string_comment_del := string_comment_del(sys_dec)
	//declalation의 전처리
	dec_comment_del := dec_line_comment_del(dec_string_comment_del) // //제거
	tem_dec_comment_del := tem_line_comment_del(tem_dec_string_comment_del)
	//template declalation의 전처리
	sys_dec_comment_del := dec_line_comment_del(sys_dec_string_comment_del)
	//system declalation의 전처리
	//fmt.Println(xml_loc, xml_transition, tem_name, tem_param, dec_comment_del, tem_dec_comment_del, sys_dec_comment_del)

	//전처리한 data를 lexer에 input하기 위해 파일 생성 후 파일에 전처리한 data 입력
	lexer_input_file_generator(tem_name, tem_param, dec_comment_del, tem_dec_comment_del, sys_dec_comment_del)

	//이전에 생성한 파일의 경로를 통해 lexer
	rst_lexer, rst_token := Lexer_TADA()
	//fmt.Println(rst_lexer, rst_token)

	//lexer를 통해 얻은 결과를 입력으로 parer
	syntax, syntax_lex_data := parse_TADA(rst_lexer, rst_token)
	//fmt.Println(syntax, syntax_lex_data)

	//처음 xml에서 받아온 tem_name을 이용하여  이전 함수를 통해 얻은 paser 결과를 후처리
	cgo_dec := after_treatment(tem_name)
	fmt.Println(cgo_dec)

	//xml에서 얻은 data를 통해 transition을 source location을 공유하는 trasition별로 분류하고 guard의 시간의 흐름에따라 정렬
	srt_trans := sort_tada_trans(xml_loc, xml_transition)

	//location별로 분류된 데이터에 guard의 시간의 흐름에따라 정렬
	make_srt_trans := sort_make_tada_trans(xml_loc, xml_transition)

	tada_timepassage := define_tada_timepassage(xml_loc, xml_transition)
	//지금 까지 정리한 data들을 통해 golang 코드 생성
	code_generator(syntax, syntax_lex_data, cgo_dec, make_srt_trans, tem_name, xml_loc, srt_trans, tada_timepassage)
}

func code_generator(syntax [][]Token, syntax_lex_data [][][]string, cgo_dec [][]byte, make_srt_trans [][][]TADA_transition, tem_name []string, tada_loc [][]TADA_loc, srt_trans [][][]TADA_transition, tada_timepassage [][][]string) {

	tem_val, channel_tada, clock_tada, param_tada, sys_tada := map_token_2_c(syntax, syntax_lex_data)
	f := NewFilePathName("/uppaal2go_result.go", "main")
	for _, val := range cgo_dec {
		f.CgoPreamble(string(val))
	}
	f.Func().Id("main").Params().BlockFunc(func(g *Group) {
		g.Id("eps").Op(":=").Qual("time", "Millisecond").Op("*").Lit(10)

		//채널 선언
		for _, val := range channel_tada {
			g.Id(val[1] + "_chan").Op(":=").Do(func(s *Statement) {
				if strings.Contains(val[0], "[") {
					_lbracket := strings.Index(val[0], "[")
					_rbracket := strings.Index(val[0], "]")
					_string := val[0][_lbracket+1 : _rbracket]
					//Lit("C."+_string)
					s.Make(Index().Chan().Bool(), Qual("C", _string)) //용량 설정 - C.N이 아닌 다른 숫자가 들어와도 가능하게끔만들어야함
					g.For(
						Id("i").Op(":=").Range().Id(val[1] + "_chan"),
					).Block(
						Id(val[1] + "_chan").Index(Id("i")).Op("=").Make(Chan().Bool()),
					)
				} else {
					s.Make(Index().Chan().Bool())
				}
			})
		}
		for i, val := range tem_name { //template 선언 Id("id").Int()
			var id_type string
			var param_id []string
			var clock_id []string
			g.Id(val).Op(":=").Func().CallFunc(func(p *Group) { // 파라미터가 여러개일때 처리
				for i, val := range param_tada[i] {
					if i == 0 || i%2 == 0 {
						id_type = val
					} else {
						if id_type == "int" {
							p.Id(val).Int()
							param_id = append(param_id, val)
						} else if id_type == "string" {
							p.Id(val).String()
							param_id = append(param_id, val)
						} else if id_type == "float" {
							p.Id(val).Float64()
							param_id = append(param_id, val)
						} else {
							p.Id(val).Qual("C", id_type)
							param_id = append(param_id, val)
						}
					}
				}
			}).BlockFunc(func(t *Group) {
				//local val 초기화
				if len(tem_val[i]) != 0 {
					t.Id("local_val").Op(":=").Qual("C", val).ValuesFunc(func(l *Group) {
						for _, v := range tem_val[i] {
							if len(v) > 2 {
								l.Id(v[1]).Op(":").Lit(v[2])
							} else { //0으로 초기화
								if strings.Contains(v[0], "[") {
									slice_name, slice_len := slice_make_C(v[0])
									slice_len_val, _ := strconv.Atoi(slice_len)
									//fmt.Println(slice_len)  //slice_len = N+1
									if slice_len_val == 0 { // 인덱스 []사이의 값이 문자일때
										slice_op_index := strings.Index(slice_len, "+")

										// 여기서 문제 현재상황에서 N값에 접근할수 없음.
										if slice_op_index == 0 { // 인덱스 값이 하나의 변수일때
											l.Id(v[1]).Op(":").Id(slice_name).ValuesFunc(func(block *Group) {
												for o := 0; o < slice_len_val; o++ { //수정
													block.Lit(0)
												}
											})
										} else { // 인덱스 값에 +가 포함될때
											l.Id(v[1]).Op(":").Id(slice_name).ValuesFunc(func(block *Group) {
												for o := 0; o < slice_len_val; o++ {
													block.Lit(0)
												}
											})
										}
									} else { // 인덱스 []사이의 값이 숫자일때
										l.Id(v[1]).Op(":").Id(slice_name).ValuesFunc(func(block *Group) {
											for o := 0; o < slice_len_val; o++ {
												block.Lit(0)
											}
										})
									}

								} else {
									l.Id(v[1]).Op(":").Lit(0)
								}
							}
						}
					})
				}
				//clock 선언
				for clock_name_index, clock_name := range clock_tada[i+1] {
					if clock_name_index > 0 {
						t.Id(clock_name+"_now").Op(":=").Qual("time", "Now").Call()
						t.Id(clock_name).Op(":=").Qual("time", "Since").Call(Id(clock_name + "_now"))
						clock_id = append(clock_id, clock_name)
					}
				}
				//time_passage 선언
				// defind_time_passge := defind_time_passge(make_srt_trans[i])
				// for _, element_time_passge := range defind_time_passge {
				// 	t.Var().Id(element_time_passge + "_passage").Index().String()
				// }
				for _, element_time_passge := range tada_timepassage[i] {
					if len(element_time_passge) > 0 {
						t.Var().Id(element_time_passge[0] + "_passage").Index().String()
					}
				}

				//location 선언
				for j, val_loc := range tada_loc[i] {
					t.Id(val_loc.id).Op(":")
					for clock_name_index, clock_name := range clock_tada[i+1] {
						if clock_name_index > 0 {
							t.Id(clock_name).Op("=").Qual("time", "Since").Call(Id(clock_name + "_now"))
							if val_loc.name == "" {
								t.Qual("fmt", "Println").Call(Lit(val), Lit("template"), Lit(val_loc.id), Lit("location"), Lit(clock_name), Lit(":"), Id(clock_name))
							} else {
								t.Qual("fmt", "Println").Call(Lit(val), Lit("template"), Lit(val_loc.name), Lit("location"), Lit(clock_name), Lit(":"), Id(clock_name))
							}
						}
					}
					//location time passage
					time_passge, time_passage_target := make_time_passge(make_srt_trans[i][j])
					if len(time_passge) > 2 {
						t.Id(val_loc.id + "_passage").Op("=").Index().String().ValuesFunc(func(q *Group) {
							for _, t_guard := range time_passge {
								q.Lit(t_guard)
							}
						})
						//x 수정
						t.Switch(Id("time_passage").Call(Id(val_loc.id+"_passage"), Id("x"))).BlockFunc(func(q *Group) {
							q.Case(Lit(0)).Block()
							for t_index, _ := range time_passge {
								q.Case(Lit(t_index + 1)).Block(
									Goto().Id(time_passage_target[t_index]),
								)
							}
						})
					}
					//location selcet
					t.Select().BlockFunc(func(s *Group) {
						for _, trans_val := range srt_trans[i][j] {
							transition_case := make_trans(trans_val.selects, trans_val.guard, trans_val.sync, clock_id, param_id)
							make_update := make_update(trans_val.assign, param_id, clock_id)
							for _, t_case := range transition_case {
								s.Case(
									t_case,
								).BlockFunc(func(q *Group) {
									for _, new_update := range make_update {
										q.Add(new_update)
									}
									q.Goto().Id(trans_val.target)
								})
							}
						}
					})

				}
			})
		}
		for _, val := range sys_tada {
			fmt.Println(sys_tada)
			g.Go().Id(val).Call()

		}
		g.Op("<-").Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(20))
	})
	f.Func().Id("when").Params(Id("guard").Bool(), Id("channel").Chan().Bool()).Chan().Bool().BlockFunc(func(g *Group) {
		g.If(
			Op("!").Id("guard"),
		).Block(
			Return(Nil()),
		)
		g.Return(Id("channel"))
	})
	f.Func().Id("when_guard").Params(Id("guard").Bool()).Op("<-").Chan().Qual("time", "Time").BlockFunc(func(g *Group) {
		g.If(
			Op("!").Id("guard"),
		).Block(
			Return(Nil()),
		)
		g.Return(Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(0)))
	})
	f.Func().Id("time_passage").Params(Id("time_passage").Index().String(), Id("ctime").Qual("time", "Duration")).Int().BlockFunc(func(g *Group) {
		g.For(
			List(Id("i"), Id("val")).Op(":=").Range().Id("time_passage"),
		).BlockFunc(func(t *Group) {
			t.If(
				Qual("strings", "Contains").Call(Id("val"), Lit("==")),
			).Block(
				List(Id("num"), Id("_")).Op(":=").Qual("strconv", "Atoi").Call(Id("val").Index(Qual("strings", "Index").Call(Id("val"), Lit("==")).Op("+").Lit(2), Empty())),
				If(
					Qual("time", "Second").Op("*").Qual("time", "Duration").Call(Id("num")).Op(">").Id("ctime"),
				).Block(
					Return(Id("i")),
				),
			).Else().If(
				Qual("strings", "Contains").Call(Id("val"), Lit("<")),
			).Block(
				List(Id("num"), Id("_")).Op(":=").Qual("strconv", "Atoi").Call(Id("val").Index(Qual("strings", "Index").Call(Id("val"), Lit("==")).Op("+").Lit(1), Empty())),
				If(
					Qual("time", "Second").Op("*").Qual("time", "Duration").Call(Id("num")).Op("==").Id("ctime"),
				).Block(
					Return(Id("i")),
				),
			)
		})
		g.Return(Len(Id("time_passage")))
	})

	a := f.Save("uppaal2go_result.go")
	fmt.Printf("%#v", f, a)
}
func slice_make_C(str string) (string, string) {
	l_parenthesis_index := strings.Index(str, "[")
	r_parenthesis_index := strings.Index(str, "]")
	slice_len := str[l_parenthesis_index+1 : r_parenthesis_index]
	slice_name := str[r_parenthesis_index+1:]
	val, _ := strconv.Atoi(slice_len)
	if val == 0 {
		str = "[" + "C." + slice_len + "]" + "C." + slice_name

		return str, slice_len

	} else {
		str = "[" + slice_len + "]" + "C." + slice_name
		return str, slice_len
	}
}
func make_time_passge(transition []TADA_transition) ([]string, []string) {
	var rst []string
	var rst_target []string

	for _, val := range transition {
		if val.guard != "" && val.assign == "" && val.selects == "" && val.sync == "" {
			rst = append(rst, val.guard)
			rst_target = append(rst_target, val.target)
		}
	}
	return rst, rst_target
}

func define_tada_timepassage(loc [][]TADA_loc, transition [][]TADA_transition) [][][]string {
	var rst_loc [][][]string
	for _, tem_loc := range loc {
		_tem_loc := make([][]string, 0)
		for _, info_loc := range tem_loc {
			_info_loc := make([]string, 0)
			if strings.Contains(info_loc.name, "_0") {
				//_info_loc = append(_info_loc, info_loc.id, info_loc.name[:len(info_loc.name)-2])
				_info_loc = append(_info_loc, info_loc.id, info_loc.name[:len(info_loc.name)-2])
			}
			_tem_loc = append(_tem_loc, _info_loc)
		}
		rst_loc = append(rst_loc, _tem_loc)
	}
	return rst_loc

}
func defind_time_passge(transition [][]TADA_transition) []string {
	var rst_source []string
	for _, trnas := range transition {
		for _, val := range trnas {
			if val.guard != "" && val.assign == "" && val.selects == "" && val.sync == "" {
				_source := strings.Trim(val.source, "p")
				_flage := false
				for _, element := range rst_source {
					if element == _source {
						_flage = true
					}
				}
				if _flage == false {
					rst_source = append(rst_source, _source)
				}
			}
		}
	}
	return rst_source
}
func make_update(update string, param_id []string, clock_id []string) []*Statement { // =, ()이게 둘다 있는 구문 추가
	var rst []*Statement
	if strings.Contains(update, "=") { //value 초기화
		index := strings.Index(update, "=")

		_id := update[:index]
		_value := update[index+1:]
		_id = strings.Trim(_id, " ")
		_value = strings.Trim(_value, " ")
		//fmt.Println(_value, "22")
		check_clock(_id, clock_id)
		if check_clock(_id, clock_id) {
			rst = append(rst, Id(_id+"_now").Op("=").Qual("time", "Now").Call())
		} else if check_clock(_id, param_id) {

		} else {
			rst = append(rst, Qual("C", _id).Op("=").Lit(_value))
		}
	} else if strings.Contains(update, "(") { //함수 사용
		lpindex := strings.Index(update, "(")
		rpindex := strings.Index(update, ")")
		_id := update[:lpindex]
		_value := update[lpindex+1 : rpindex]
		_id = strings.Trim(_id, " ")
		_value = strings.Trim(_value, " ")
		_int_or_string, _ := strconv.Atoi(_value)
		if _int_or_string == 0 { //문자열이라면
			//rst = append(rst, Qual("C", _id).Call(Op("&").Id("local_val"), Id("C."+_value)))

		} else {
			//rst = append(rst, Qual("C", _id).Call(Op("&").Id("local_val"), Id(_value)))
		}
		rst = append(rst, Qual("C", _id).Call(Op("&").Id("local_val"), Id(_value)))

	}
	return rst
}
func make_trans(selects string, guard string, sync string, clock []string, param []string) []*Statement {
	var rst []*Statement
	if selects == "" {
		if guard == "" && sync == "" {
			rst = append(rst, Op("<-").Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(0)))
			return rst
		} else if guard != "" && sync == "" { //가드만 있을 때 시간만 있다고 가정해서 시간이 아닌 표현식이 들어왔을때 구분해주는게 필요
			if strings.Contains(guard, "==") {
				_time, _val, _ := guard_preprocessing(guard)
				_value, _ := strconv.Atoi(_val)
				rst = append(rst, Op("<-").Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(_value).Op("-").Id(_time).Op("-").Id("eps")))
				return rst
			} else if strings.Contains(guard, ">") {
				_time, _val, _ := guard_preprocessing(guard)
				_value, _ := strconv.Atoi(_val)
				rst = append(rst, Op("<-").Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(_value).Op("-").Id(_time)))
				return rst
			}
			rst = append(rst, Op("<-").Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(40)))
			return rst
		} else if guard == "" && sync != "" {
			if strings.Contains(sync, "!") {
				sync = strings.Trim(sync, "!")
				rst = append(rst, sync_index(sync, "!", param))
				return rst
			} else if strings.Contains(sync, "?") {
				sync = strings.Trim(sync, "?")
				rst = append(rst, sync_index(sync, "?", param))
				return rst
			}
		} else if guard != "" && sync != "" {
			if strings.Contains(sync, "!") {
				sync = strings.Trim(sync, "!")
				rst = append(rst, when_sync(sync, guard, "!", param))
				return rst
			} else if strings.Contains(sync, "?") {
				sync = strings.Trim(sync, "?")
				rst = append(rst, when_sync(sync, guard, "?", param))
				return rst
			}
		}
	} else { //select
		select_preprocessing(selects)
	}
	rst = append(rst, Op("<-").Qual("time", "After").Call(Qual("time", "Second").Op("*").Lit(40)))
	return rst
}
func guard_preprocessing(_guard string) (string, string, string) {
	if strings.Contains(_guard, "==") {
		index := strings.Index(_guard, "==")
		_value := _guard[:index]
		_type := _guard[index+2:]
		_value = strings.Trim(_value, " ")
		_type = strings.Trim(_type, " ")
		return _value, _type, "=="
	} else if strings.Contains(_guard, ">") {
		index := strings.Index(_guard, ">")
		_value := _guard[:index]
		_type := _guard[index+1:]
		_value = strings.Trim(_value, " ")
		_type = strings.Trim(_type, " ")
		return _value, _type, ">"
	} else if strings.Contains(_guard, ">=") {
		index := strings.Index(_guard, ">=")
		_value := _guard[:index]
		_type := _guard[index+1:]
		_value = strings.Trim(_value, " ")
		_type = strings.Trim(_type, " ")
		return _value, _type, ">="
	} else if strings.Contains(_guard, "<") {
		index := strings.Index(_guard, "<")
		_value := _guard[:index]
		_type := _guard[index+1:]
		_value = strings.Trim(_value, " ")
		_type = strings.Trim(_type, " ")
		return _value, _type, "<"
	} else if strings.Contains(_guard, "<=") {
		index := strings.Index(_guard, "<=")
		_value := _guard[:index]
		_type := _guard[index+1:]
		_value = strings.Trim(_value, " ")
		_type = strings.Trim(_type, " ")
		return _value, _type, "<="
	} else {
		return "", "", ""
	}
}
func when_sync(sync string, guard string, op string, param []string) *Statement {
	var rst *Statement
	_left, _right, _op := guard_preprocessing(guard)
	_guard := guard_postprocessing(_left, _right, param, _op)
	if strings.Contains(sync, "[") {
		index_lbracket := strings.Index(sync, "[")
		index_rbracket := strings.Index(sync, "]")

		id := sync[:index_lbracket]
		num := sync[index_lbracket+1 : index_rbracket]
		id = strings.Trim(id, " ")
		num = strings.Trim(num, " ") // num확인 필요

		trans_num, _ := strconv.Atoi(num)
		rst = Id(id + "_chan").Index(Id(num))

		if num == "0" {
			rst = Id(id + "_chan").Index(Lit(0))
		} else if trans_num != 0 {
			rst = Id(id + "_chan").Index(Lit(trans_num))
		} else if strings.Contains(num, "(") {
			index_lbracket = strings.Index(num, "(")
			//index_rbracket := strings.Index(num, ")")
			func_id := num[:index_lbracket]
			//func_param := num[index_lbracket+1 : index_rbracket] //param이 없다고 가정, 수정 필요
			rst = Id(id + "_chan").Index(Qual("C", func_id).Call(Op("&").Id("local_val")))
		} else {
			_flage := false
			for _, val := range param {
				if num == val {
					_flage = true
				}
			}
			if _flage == true {
				rst = Id(id + "_chan").Index(Id(num))
			} else {
				rst = Id(id + "_chan").Index(Qual("C", num))
			}
		}

		if op == "?" {
			return Op("<-").Id("when").Call(Add(_guard), rst)
		} else {
			return Id("when").Call(Add(_guard), rst).Op("<-").Lit(true)
		}
	} else {
		if op == "?" {
			return Op("<-").Id("when").Call(Add(_guard), Id(sync+"_chan"))
		} else {
			return Id("when").Call(Add(_guard), Id(sync+"_chan")).Op("<-").Lit(true)
		}
	}
}
func sync_index(sync string, op string, param []string) *Statement {
	var rst *Statement
	if strings.Contains(sync, "[") {
		index_lbracket := strings.Index(sync, "[")
		index_rbracket := strings.Index(sync, "]")

		id := sync[:index_lbracket]
		num := sync[index_lbracket+1 : index_rbracket]
		id = strings.Trim(id, " ")
		num = strings.Trim(num, " ")
		trans_num, _ := strconv.Atoi(num)
		rst = Id(id + "_chan").Index(Id(num))
		if num == "0" {
			rst = Id(id + "_chan").Index(Lit(0))
		} else if trans_num != 0 {
			rst = Id(id + "_chan").Index(Lit(trans_num))
		} else if strings.Contains(num, "(") {
			index_lbracket = strings.Index(num, "(")
			//index_rbracket := strings.Index(num, ")")
			func_id := num[:index_lbracket]
			//func_param := num[index_lbracket+1 : index_rbracket] //param이 없다고 가정, 수정 필요
			rst = Id(id + "_chan").Index(Qual("C", func_id).Call(Op("&").Id("local_val")))
		} else {
			_flage := false
			for _, val := range param {
				if num == val {
					_flage = true
				}
			}
			if _flage == true {
				rst = Id(id + "_chan").Index(Id(num))
			} else {
				rst = Id(id + "_chan").Index(Qual("C", num))
			}
		}

		if op == "?" {
			return Op("<-").Add(rst)
		} else {
			return rst.Op("<-").Lit(true)
		}
	} else {
		if op == "?" {
			return Op("<-").Id(sync + "_chan")
		} else {
			return Id(sync + "_chan").Op("<-").Lit(true)
		}
	}
}

func guard_postprocessing(_left string, _right string, param []string, op string) *Statement {
	var rst *Statement
	trans_num, _ := strconv.Atoi(_left)
	if _left == "0" {
		rst = Lit(0)
	} else if trans_num != 0 {
		rst = Lit(trans_num)
	} else if strings.Contains(_left, "(") {
	} else {
		_flage := false
		for _, val := range param {
			if _left == val {
				_flage = true
			}
		}
		if _flage == true {
			rst = Id(_left)
		} else {
			rst = Id("local_val").Dot(_left)
		}
	}
	if op == "==" {
		rst = rst.Op("==")
	} else if op == ">" {
		rst = rst.Op(">")
	} else if op == "<" {
		rst = rst.Op("<")
	} else if op == "<=" {
		rst = rst.Op("<=")
	} else if op == ">=" {
		rst = rst.Op(">=")
	}
	trans_num, _ = strconv.Atoi(_right)
	if _right == "0" {
		rst = rst.Lit(0)
	} else if trans_num != 0 {
		rst = rst.Lit(trans_num)
	} else if strings.Contains(_right, "(") {
	} else {
		_flage := false
		for _, val := range param {
			if _right == val {
				_flage = true
			}
		}
		if _flage == true {
			rst = rst.Id(_right)
		} else {
			rst = rst.Id("local_val").Dot(_left)
		}
	}
	return rst
}
func select_preprocessing(_select string) (string, string) {
	index := strings.Index(_select, ":")
	_value := _select[:index]
	_type := _select[index+1:]
	_value = strings.Trim(_value, " ")
	_type = strings.Trim(_type, " ")
	return _value, _type
}
func check_clock(id string, clock []string) bool {
	for _, val := range clock {
		if val == id {
			return true
		}
	}
	return false
}
func sort_make_tada_trans(loc [][]TADA_loc, trans [][]TADA_transition) [][][]TADA_transition {
	srt_data := make([][][]TADA_transition, 0)
	for i, tem := range loc { //tem
		_srt_data := make([][]TADA_transition, 0)
		for _, val := range tem { //location
			_trnas := make([]TADA_transition, 0)
			if !strings.Contains(val.id, "p") {
				for _, trans_q := range trans[i] {
					_loc := trans_q.source
					_loc = strings.Trim(_loc, "p")
					if _loc == val.id {
						_trnas = append(_trnas, trans_q)
					}
				}
			}

			_srt_data = append(_srt_data, _trnas)
		}
		srt_data = append(srt_data, _srt_data)
	}

	return srt_data
}
func sort_tada_tran(loc [][]TADA_loc, trans [][]TADA_transition) [][][]TADA_transition {
	srt_data := make([][][]TADA_transition, 0)
	for i, tem := range loc {
		_srt_data := make([][]TADA_transition, 0)
		for _, val := range tem {
			_trnas := make([]TADA_transition, 0)
			for _, trans_q := range trans[i] {
				if trans_q.source == val.id {
					_trnas = append(_trnas, trans_q)
				}
			}
			_srt_data = append(_srt_data, _trnas)
		}
		srt_data = append(srt_data, _srt_data)
	}
	return srt_data
}
func sort_tada_trans(loc [][]TADA_loc, trans [][]TADA_transition) [][][]TADA_transition {
	srt_data := make([][][]TADA_transition, 0)
	for i, tem := range loc {
		_srt_data := make([][]TADA_transition, 0)
		for _, val := range tem {
			_trnas := make([]TADA_transition, 0)
			for _, trans_q := range trans[i] {
				if trans_q.source == val.id {
					_trnas = append(_trnas, trans_q)
				}
			}
			_srt_data = append(_srt_data, _trnas)
		}
		srt_data = append(srt_data, _srt_data)
	}

	return srt_data
}
func make_chan(name string, isMap bool) *Statement {
	return Id(name).Op(":=").Do(func(s *Statement) {
		if isMap {
			s.Map(String()).String()
		} else {
			s.Index().String()
		}
	}).Values()
}
func after_treatment(tem_name []string) [][]byte {
	input_file, err := os.Open(dec_path)
	check(err)
	reader := bufio.NewReader(input_file)
	input_file_reader := make([][]byte, 0)
	output_file := make([][]byte, 0)
	for {
		line, _, err := reader.ReadLine()
		input_file_reader = append(input_file_reader, line)
		if err != nil {
			break
		}
	}
	defer input_file.Close()
	rbrace_struct := false
	_index := 0
	_tem_name := ""
	for i, val := range input_file_reader {
		if rbrace_struct == false {
			for j := 0; j < len(tem_name); j++ {
				if strings.Contains(string(val), tem_name[j]) && strings.Contains(string(val), "struct") {
					rbrace_struct = true
					_tem_name = tem_name[j]
				}
			}
		} else {
			if !strings.Contains(string(val), "       ") {
				//output_file = append(output_file, input_file_reader[_index:i])
				for j := _index; j < i; j++ {
					output_file = append(output_file, input_file_reader[j])
				}
				output_file = append(output_file, []byte("}"+_tem_name+"; "))
				rbrace_struct = false
				_index = i
				for j := 0; j < len(tem_name); j++ {
					if strings.Contains(string(val), tem_name[j]) && strings.Contains(string(val), "struct") {
						rbrace_struct = true
						_tem_name = tem_name[j]
					}
				}
			}
		}
	}
	for j := _index; j < len(input_file_reader); j++ {
		output_file = append(output_file, input_file_reader[j])
	}
	return output_file
}

func contains(elems []Token, v Token) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
func contains_string(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
func slice_count(s []Token, a Token) int {
	num := 0
	for _, v := range s {
		if v == a {
			num += 1
		}
	}
	return num
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
func map_token_2_c(parse [][]Token, parse_lexr_data [][][]string) ([][][]string, [][]string, [][]string, [][]string, []string) {
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
	tem_val := make([][][]string, 0)
	tem_chan := make([][]string, 0)
	tem_clock := make([][]string, 0)
	tem_param := make([][]string, 0)
	sys_dec := make([]string, 0)

	for {
		line, _, err := reader.ReadLine()
		input_file_reader = append(input_file_reader, line)
		if err != nil {
			break
		}
	}
	//fmt.Println(input_file_reader)
	defer input_file.Close()
	defer output_file.Close()
	//fmt.Println(parse, "\n", parse_lexr_data, len(parse), len(parse_lexr_data))
	_local := false
	_param := false
	_system := false
	for i, _parse := range parse {
		//fmt.Println(i, _parse)
		if _local && contains(_parse, ASSIGN) { //여기서 부터 시작

		} else if _local && contains(_parse, IDENT) && contains(_parse, RPARENTHESIS) && contains(_parse, RBRACE) {
			//fmt.Println("local _ func:", parse_lexr_data[i])
			_start_line, _ := strconv.Atoi(parse_lexr_data[i][0][0])
			_end_line, _ := strconv.Atoi(parse_lexr_data[i][len(parse_lexr_data[i])-1][0])
			for j := _start_line - 1; j < _end_line; j++ {
				if j == _start_line-1 {
					_lparen := contain_index(parse[i], LPARENTHESIS)
					_index, _ := strconv.Atoi(parse_lexr_data[i][_lparen][1])
					//fmt.Println(string(input_file_reader[j][:_index]) + tem_name[len(tem_name)-1] + " *" + tem_name[len(tem_name)-1] + " " + string(input_file_reader[j][_index:]))
					if string(input_file_reader[j][_index:]) == ")" {
						_, err := output_file.Write([]byte(string(input_file_reader[j][:_index]) + tem_name[len(tem_name)-1] + " *" + tem_name[len(tem_name)-1] + " " + string(input_file_reader[j][_index:]) + "\n"))
						check(err)
					} else {
						_, err := output_file.Write([]byte(string(input_file_reader[j][:_index]) + tem_name[len(tem_name)-1] + " *" + tem_name[len(tem_name)-1] + ", " + string(input_file_reader[j][_index:]) + "\n"))
						check(err)
					}

				} else {
					_mapping_bool := false
					_string := string(input_file_reader[j])
					for k := 0; k < len(tem_val[len(tem_name)-1]); k++ {
						if strings.Contains(_string, strings.Trim(tem_val[len(tem_name)-1][k][1], " ")) {
							_string = strings.ReplaceAll(_string, strings.Trim(tem_val[len(tem_name)-1][k][1], " "), tem_name[len(tem_name)-1]+"->"+tem_val[len(tem_name)-1][k][1])

							_mapping_bool = true
							//fmt.Println(tem_val[len(tem_name)-1][k][1])

						}
					}
					if _mapping_bool == false {
						_, err := output_file.Write([]byte(string(input_file_reader[j]) + "\n"))
						check(err)
					} else {
						_, err := output_file.Write([]byte(_string + "\n"))
						check(err)
					}
				}
			}

		} else if contains(_parse, ASSIGN) { //initializer								const int ctimemax[4] =  {1,1,1,1};
			if parse[i][0] == PREFIX && parse_lexr_data[i][0][2] == "const" { //const int N = 6;		#define N 6
				_ident := contain_index(parse[i], IDENT)
				_int := contain_index(parse[i], INT)
				_type := contain_index(parse[i], TYPEID)
				_state := 0
				for _, val_parse := range parse[i] {
					if val_parse == LBRACKET {
						_state = 1
					}
					if val_parse == LBRACE {
						_state = 2
					}
				}
				if _state == 1 {
					_lbracket := contain_index(parse[i], LBRACKET)
					_rbracket := contain_index(parse[i], RBRACKET)
					_mappintstring := "const" + " " + mapping(parse_lexr_data, input_file_reader, i, _type) + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _lbracket) + mapping(parse_lexr_data, input_file_reader, i, _rbracket) + ";"
					_mappintstring = strings.Trim(_mappintstring, " ")
					_mappintstring = _mappintstring[:len(_mappintstring)-1]
					_, err := output_file.Write([]byte(_mappintstring + "\n"))
					check(err)
				} else if _state == 2 {
					_lbrace := contain_index(parse[i], LBRACE)
					_rbrace := contain_index(parse[i], RBRACE)

					_lbracket := contain_index(parse[i], LBRACKET)
					_rbracket := contain_index(parse[i], RBRACKET)
					_mappintstring := "const" + " " + mapping(parse_lexr_data, input_file_reader, i, _type) + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _lbracket) + mapping(parse_lexr_data, input_file_reader, i, _rbracket) + " = " + mapping(parse_lexr_data, input_file_reader, i, _lbrace) + mapping(parse_lexr_data, input_file_reader, i, _rbrace) + ";" //수정
					_mappintstring = strings.Trim(_mappintstring, " ")
					_mappintstring = _mappintstring[:len(_mappintstring)-1]
					_, err := output_file.Write([]byte(_mappintstring + "\n"))
					check(err)
				} else {
					_mappintstring := "const" + " " + mapping(parse_lexr_data, input_file_reader, i, _type) + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + " = " + mapping(parse_lexr_data, input_file_reader, i, _int) + ";"
					_mappintstring = strings.Trim(_mappintstring, " ")
					_mappintstring = _mappintstring[:len(_mappintstring)-1]
					_, err := output_file.Write([]byte(_mappintstring + "\n"))
					check(err)
					//fmt.Println(mapping(parse_lexr_data, input_file_reader, i, _ident), mapping(parse_lexr_data, input_file_reader, i, _int)) //문제@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
					//fmt.Println("#define" + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + " " + mapping(parse_lexr_data, input_file_reader, i, _int))
				}

			} else {
				if parse[i][0] == TYPEID && parse_lexr_data[i][0][2] == "int" {
					_ident := contain_index(parse[i], IDENT)
					_type := contain_index(parse[i], TYPEID)
					_state := 0
					for _, val_parse := range parse[i] {
						if val_parse == LBRACKET {
							_state = 1
						}
						if val_parse == LBRACE {
							_state = 2
						}
					}
					if _state == 1 {
						_lbracket := contain_index(parse[i], LBRACKET)
						_rbracket := contain_index(parse[i], RBRACKET)
						_mappintstring := "int " + mapping(parse_lexr_data, input_file_reader, i, _type) + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _lbracket) + mapping(parse_lexr_data, input_file_reader, i, _rbracket) + ";"
						_mappintstring = strings.Trim(_mappintstring, " ")
						_mappintstring = _mappintstring[:len(_mappintstring)-1]
						_, err := output_file.Write([]byte(_mappintstring + "\n"))
						check(err)
					} else if _state == 2 {
						_lbrace := contain_index(parse[i], LBRACE)
						_rbrace := contain_index(parse[i], RBRACE)

						_lbracket := contain_index(parse[i], LBRACKET)
						_rbracket := contain_index(parse[i], RBRACKET)
						_mappintstring := "int " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _lbracket) + mapping(parse_lexr_data, input_file_reader, i, _rbracket) + " = " + mapping(parse_lexr_data, input_file_reader, i, _lbrace) + mapping(parse_lexr_data, input_file_reader, i, _rbrace) + ";" //수정
						_mappintstring = strings.Trim(_mappintstring, " ")
						_mappintstring = _mappintstring[:len(_mappintstring)-1]
						_, err := output_file.Write([]byte(_mappintstring + "\n"))
						check(err)
					} else {
						_int := contain_index(parse[i], INT)
						_mappintstring := "int" + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + mapping(parse_lexr_data, input_file_reader, i, _int)
						_mappintstring = strings.Trim(_mappintstring, " ")
						_mappintstring = _mappintstring[:len(_mappintstring)-1]
						_, err := output_file.Write([]byte(_mappintstring + "\n"))
						check(err)
					}

					//fmt.Println(_mappintstring)
					//fmt.Println("#define" + " " + mapping(parse_lexr_data, input_file_reader, i, _ident) + " " + mapping(parse_lexr_data, input_file_reader, i, _int))

				}
			}
		} else if contains(_parse, IDENT) && contains(_parse, RPARENTHESIS) && contains(_parse, RBRACE) { //func 수정필요
			//파라미터에 Local *local
			//local->list 구조체 멤버 접근시
		} else if contains(_parse, DIV) && contains(_parse, PARAM) {
			_param = true
			_local = false
			tem_param = append(tem_param, make([]string, 0))

		} else if contains(_parse, DIV) && contains(_parse, SYSTEM) {
			_system = true
			_param = false
		} else if contains(_parse, DIV) { //바꾸어야 할지도\
			tem_name = append(tem_name, parse_lexr_data[i][2][2])
			tem_val = append(tem_val, make([][]string, 0)) //중요
			tem_clock = append(tem_clock, make([]string, 0))
			_, err := output_file.Write([]byte("typedef struct " + parse_lexr_data[i][2][2] + "{\n"))
			check(err)
			_local = true
		} else if _param {
			for j, _ident := range _parse {
				if _ident != PREFIX && _ident != COMMA {
					tem_param[len(tem_param)-1] = append(tem_param[len(tem_param)-1], parse_lexr_data[i][j][2]) // type value
				}
			}

		} else if _system {
			if contains(_parse, SYSTEM) {

				for j, _ident := range _parse {
					if _ident == IDENT {
						sys_dec = append(sys_dec, parse_lexr_data[i][j][2])
					}
				}
			}
		} else { //dec
			//chan ,clock
			//int[0,6]
			//로컬의 경우 struct로
			if _local {
				if contains(_parse, CLOCK) {
					_clock := make([]string, 0)
					_clock = append(_clock, tem_name[len(tem_name)-1], parse_lexr_data[i][1][2])
					tem_clock = append(tem_clock, _clock)
				} else {
					_lbracket := contain_index(parse[i], LBRACKET)
					if parse[i][_lbracket+1] == COMMA { //[,] check

						_ident_line, _ := strconv.Atoi(parse_lexr_data[i][1][0])
						_rbracket := contain_index(parse[i], RBRACKET)
						_rbracket_stack, _ := strconv.Atoi(parse_lexr_data[i][_rbracket][1])

						_slice := []string{mapping(parse_lexr_data, input_file_reader, i, 0), string(input_file_reader[_ident_line-1][_rbracket_stack : len(input_file_reader[_ident_line-1])-1])}
						tem_val[len(tem_name)-1] = append(tem_val[len(tem_name)-1], _slice)
						_, err := output_file.Write([]byte("        " + mapping(parse_lexr_data, input_file_reader, i, 0) + string(input_file_reader[_ident_line-1][_rbracket_stack:]) + "\n"))
						check(err)
						//fmt.Println("55", mapping(parse_lexr_data, input_file_reader, i, 0), string(input_file_reader[_ident_line-1][_rbracket_stack:]))
					} else {
						if contains(_parse, COMMA) { //[],[] check
						} else {
							_ident_line, _ := strconv.Atoi(parse_lexr_data[i][1][0])
							_ident_stack, _ := strconv.Atoi(parse_lexr_data[i][1][1])
							_rbracket := contain_index(parse[i], RBRACKET)
							_lbracket := contain_index(parse[i], LBRACKET)
							_lbracket_stack, _ := strconv.Atoi(parse_lexr_data[i][_lbracket][1])
							_rbracket_stack, _ := strconv.Atoi(parse_lexr_data[i][_rbracket][1])
							_slice := []string{string(input_file_reader[_ident_line-1][_lbracket_stack-1:_rbracket_stack]) + mapping(parse_lexr_data, input_file_reader, i, 0), string(input_file_reader[_ident_line-1][_ident_stack-1 : _lbracket_stack-1])}
							tem_val[len(tem_name)-1] = append(tem_val[len(tem_name)-1], _slice)
							//fmt.Println("56", mapping(parse_lexr_data, input_file_reader, i, 0), string(input_file_reader[_ident_line-1][_ident_stack-1:_lbracket_stack-1])+string(input_file_reader[_ident_line-1][_lbracket_stack-1:_rbracket_stack]))

							_, err := output_file.Write([]byte("        " + mapping(parse_lexr_data, input_file_reader, i, 0) + string(input_file_reader[_ident_line-1][_ident_stack-1:_lbracket_stack-1]) + string(input_file_reader[_ident_line-1][_lbracket_stack-1:_rbracket_stack]) + ";\n"))
							check(err)

							// _, err = output_file.Write([]byte(string(input_file_reader[_ident_line-1][_lbracket_stack-1:_rbracket_stack]) + mapping(parse_lexr_data, input_file_reader, i, 0) + string(input_file_reader[_ident_line-1][_ident_stack-1:_lbracket_stack-1]) + "\n}" + parse_lexr_data[i][2][2] + ";\n"))
							// check(err)
							// fmt.Println(_output_file_reader[:len(_output_file_reader)-2])
						}
					}
				}
				//id_t list[N+1];
			} else {
				if contains(_parse, CLOCK) {
					_clock := make([]string, 0)
					_clock = append(_clock, "gobal", parse_lexr_data[i][1][2])
					tem_clock = append(tem_clock, _clock)
				} else if contains(_parse, CHANNEL) {
					if _parse[0] == PREFIX {
						_num, _ := strconv.Atoi(parse_lexr_data[i][1][0])
						for j, val := range _parse {
							if val == IDENT {
								_chan := make([]string, 0)
								if strings.Contains(string(input_file_reader[_num-1]), "[") {
									_lbracket_index := strings.Index(string(input_file_reader[_num-1]), "[")
									_rbracket_index := strings.Index(string(input_file_reader[_num-1]), "]")
									_chan = append(_chan, "broadcast ["+string(input_file_reader[_num-1][_lbracket_index+1:_rbracket_index])+"]"+parse_lexr_data[i][1][2])
								} else {
									_chan = append(_chan, "broadcast"+parse_lexr_data[i][1][2])
								}
								_chan = append(_chan, parse_lexr_data[i][j][2])
								tem_chan = append(tem_chan, _chan)
							}
						}
					} else if _parse[0] == CHANNEL {
						_num, _ := strconv.Atoi(parse_lexr_data[i][1][0])
						for j, val := range _parse {
							if val == IDENT {
								_chan := make([]string, 0)
								if strings.Contains(string(input_file_reader[_num-1]), "[") {
									_lbracket_index := strings.Index(string(input_file_reader[_num-1]), "[")
									_rbracket_index := strings.Index(string(input_file_reader[_num-1]), "]")
									_chan = append(_chan, "["+string(input_file_reader[_num-1][_lbracket_index+1:_rbracket_index])+"]"+parse_lexr_data[i][0][2])
								} else {
									_chan = append(_chan, parse_lexr_data[i][0][2])
								}
								_chan = append(_chan, parse_lexr_data[i][j][2])
								tem_chan = append(tem_chan, _chan)
							}
						}
					}
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
	//return chan, clock 추가 해야 할듯
	return tem_val, tem_chan, tem_clock, tem_param, sys_dec
	//return tem_chan, tem_clock
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
					stack_lex[len(stack_token)-3][2] = stack_lex[len(stack_token)-3][2] + stack_lex[len(stack_token)-2][2] + stack_lex[len(stack_token)-1][2]
					stack_token = stack_token[:len(stack_token)-2]
					stack_lex = stack_lex[:len(stack_lex)-2]
				}
			}
		}
	}
	return syntax, syntax_lex_data
}

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
	SYSTEM
	PARAM
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
	SYSTEM:  "SYSTEM",
	PARAM:   "PARAM",
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
				case "param":
					return startPos, PARAM, lit
				case "system":
					return startPos, SYSTEM, lit
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

func Lexer_param(param []string, dec string) ([][]string, []Token) {
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
func lexer_input_file_generator(tem_name []string, tem_param []string, dec_comment_del []string, tem_dec_comment_del [][]string, sys_dec_comment_del []string) {
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
	for i, val := range tem_param {
		_, err := file.Write([]byte("//" + "param_" + tem_name[i] + ";" + "\n"))
		check(err)
		_, err = file.Write([]byte(val + ";" + "\n"))
		check(err)
	}
	_, err = file.Write([]byte("//" + "system_dec;" + "\n"))
	check(err)
	for i, _ := range sys_dec_comment_del {

		_, err = file.Write([]byte(sys_dec_comment_del[i] + "\n"))
		check(err)
	}

	file.Close()
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Token int
type TADA_loc struct {
	id   string
	name string
}
type TADA_transition struct {
	source  string
	target  string
	selects string
	guard   string
	sync    string
	assign  string
}
type TADA_tem struct {
	loc  TADA_loc
	tans TADA_transition
}

func input_xml_to_slice(doc *etree.Document) ([][]TADA_loc, [][]TADA_transition, []string, [][]string) {
	tada_loc := make([][]TADA_loc, 0)
	tada_trans := make([][]TADA_transition, 0)
	var dec string
	var tem_dec []string
	var tem_name []string
	var tem_param []string
	var sys_dec string
	for _, e := range doc.FindElements("./nta/*") {
		if e.Tag == "declaration" {
			dec = e.Text()
		}
		if e.Tag == "template" {
			_tada_loc := make([]TADA_loc, 0)
			_tada_trans := make([]TADA_transition, 0)

			if name := e.SelectElement("name"); name != nil {
				tem_name = append(tem_name, name.Text())
				tem_param = append(tem_param, "")

			}
			if param := e.SelectElement("parameter"); param != nil {

				tem_param[len(tem_param)-1] = param.Text()
			}
			if declaration := e.SelectElement("declaration"); declaration != nil {
				tem_dec = append(tem_dec, declaration.Text())
			}
			for _, l := range e.FindElements("location") {
				var _id string
				var _name string
				if l.Attr[0].Key == "id" {
					_id = l.Attr[0].Value
				}
				if l_name := l.SelectElement("name"); l_name != nil {
					_name = l_name.Text()
				}
				_tada_loc = append(_tada_loc, TADA_loc{_id, _name})
			}
			tada_loc = append(tada_loc, _tada_loc)

			for _, t := range e.FindElements("transition") {
				var _source string
				var _target string
				var _select string
				var _guard string
				var _sync string
				var _assign string

				if t_source := t.SelectElement("source"); t_source != nil {
					_source = t_source.Attr[0].Value
				}
				if t_target := t.SelectElement("target"); t_target != nil {
					_target = t_target.Attr[0].Value
				}
				for _, l := range t.FindElements("label") {
					if l.Attr[0].Value == "select" {
						_select = l.Text()
					} else if l.Attr[0].Value == "guard" {
						_guard = l.Text()
					} else if l.Attr[0].Value == "synchronisation" {
						_sync = l.Text()
					} else if l.Attr[0].Value == "assignment" {
						_assign = l.Text()
					}
				}
				_tada_trans = append(_tada_trans, TADA_transition{_source, _target, _select, _guard, _sync, _assign})
			}
			tada_trans = append(tada_trans, _tada_trans)
		}
		if e.Tag == "system" {
			sys_dec = e.Text()

		}
		//fmt.Println("\n")
	}
	return_val_dec := make([]string, 0)
	return_val_dec = append(return_val_dec, dec, sys_dec)
	return_val_tem := make([][]string, 0)
	return_val_tem = append(return_val_tem, tem_dec, tem_name, tem_param)
	return tada_loc, tada_trans, return_val_dec, return_val_tem
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
	if string_counts == 0 {
		return dec
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
