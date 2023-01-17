package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	. "github.com/dave/jennifer/jen"
)

var new_type []string

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("train-gate.xml"); err != nil {
		panic(err)
	}
	f := NewFile("main")

	//root := doc.SelectElement("nta")
	var dec string
	var tem_dec []string
	var num_loc []int
	var loc []string          //loc 정보
	var edge []string         //edge 정보
	var global_clock []string //clock 을 템플릿별로 나누어서 가드랑 비교하는게 필요**

	for _, e := range doc.FindElements("./nta/*") {
		fmt.Println(e.Tag)
		if e.Tag == "declaration" { //declaration parsing필요
			fmt.Println(e.Text())
			dec = e.Text()
		}
		if e.Tag == "template" {
			if name := e.SelectElement("name"); name != nil {
				fmt.Printf("  template name : %s\n", name.Text())
			}
			if param := e.SelectElement("parameter"); param != nil {
				fmt.Printf("  template param : %s\n", param.Text()) //parameter parsing필요
			}
			if declaration := e.SelectElement("declaration"); declaration != nil {
				fmt.Printf("  template declaration : %s\n", declaration.Text()) //declaration parsing필요
				tem_dec = append(tem_dec, declaration.Text())
			}
			i := 0
			for _, l := range e.FindElements("location") {
				if l.Attr[0].Key == "id" {
					fmt.Printf("  template location id : %s\n", l.Attr[0].Value)
					loc = append(loc, l.Attr[0].Value)
				}
				if l_name := l.SelectElement("name"); l_name != nil {
					fmt.Printf("  template location name : %s\n", l_name.Text())
				}
				if l_label := l.SelectElement("label"); l_label != nil {
					fmt.Printf("  template location label value: %s\n", l_label.Text())
					fmt.Printf("  template location label kind: %s\n", l_label.Attr[0].Value)
					if l_label.Attr[0].Value == "invariant" {
						lenloc := len(loc)
						loc[lenloc-1] = loc[lenloc-1] + " " + l_label.Text()
					}
				}
				i = i + 1
			}
			num_loc = append(num_loc, i)

			for _, t := range e.FindElements("transition/*") {
				if t.Tag == "source" {
					fmt.Printf("  template transition source: %s\n", t.Attr[0].Value)
					edge = append(edge, t.Attr[0].Value)
				}
				if t.Tag == "target" {
					fmt.Printf("  template transition target: %s\n", t.Attr[0].Value)
					lenedge := len(edge)
					edge[lenedge-1] = edge[lenedge-1] + " " + t.Attr[0].Value
				}
				if t.Tag == "label" && t.Attr[0].Key == "kind" {
					fmt.Printf("  template transition label kind: %s %s\n", t.Attr[0].Value, t.Text())

					if t.Attr[0].Value == "guard" {
						lenedge := len(edge)
						edge[lenedge-1] = edge[lenedge-1] + " " + t.Text()
					}
				}
			}

		}
		if e.Tag == "system" {
			fmt.Println(e.Text())
			//sys = e.Text()

		}
		fmt.Println("\n")
	}
	fmt.Printf("%#v", f)

	loc_edge := make([][]string, len(loc))
	for i := 0; i < len(loc); i++ {
		loc_edge[i] = make([]string, 0)
		loc_edge[i] = append(loc_edge[i], loc[i])
		for j := 0; j < len(edge); j++ {
			a, _ := strconv.Atoi(string(edge[j][2]))
			if a == i {
				loc_edge[i] = append(loc_edge[i], edge[j])
			}
		}
	}
	fmt.Println(loc_edge)

	_clock := clock_extraction(dec)
	global_clock = append(global_clock, _clock...)
	for _, value := range tem_dec {
		_clock = clock_extraction(value)
		global_clock = append(global_clock, _clock...)
	}
	fmt.Println(global_clock)
	//fmt.Println(local_clock)

	//조건 분석
	translator_condition := make([][]int, len(loc))

	for i := 0; i < len(loc); i++ {
		translator_condition[i] = make([]int, 0)

		if len(loc_edge[i][0]) == 3 { //invariant가 없는 경우
			translator_condition[i] = append(translator_condition[i], 0)
		} else { //invariant가 있는 경우 이때 invariant는 무조건 시간과 관련된 변수로 가정 then 체크안하고 바로 1.
			translator_condition[i] = append(translator_condition[i], 1)
		}

		for j := 1; j < len(loc_edge[i]); j++ {
			if len(loc_edge[i][j]) == 7 && len(translator_condition[i]) == 1 { //guard가 없는 경우
				translator_condition[i] = append(translator_condition[i], 0)
			} else { //guard가 있는 경우
				//fmt.Println(loc_edge[i][j][7:])
				//a := regexp.MustCompile(`>=`)
				var _guard_condition int
				if strings.Contains(loc_edge[i][j][7:], "<=") {
					_slice := strings.Split(loc_edge[i][j][7:], "<=")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					_guard_condition = find_clock_element(global_clock, _guard_element)
					if _guard_condition == 0 && len(translator_condition[i]) == 1 {
						translator_condition[i] = append(translator_condition[i], 0)
					} else if _guard_condition == 1 {
						translator_condition[i] = append(translator_condition[i], 2)
					}

				} else if strings.Contains(loc_edge[i][j][7:], ">=") {
					_slice := strings.Split(loc_edge[i][j][7:], ">=")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					_guard_condition = find_clock_element(global_clock, _guard_element)
					if _guard_condition == 0 && len(translator_condition[i]) == 1 {
						translator_condition[i] = append(translator_condition[i], 0)
					} else if _guard_condition == 1 {
						translator_condition[i] = append(translator_condition[i], 5)
					}

				} else if strings.Contains(loc_edge[i][j][7:], "==") {
					_slice := strings.Split(loc_edge[i][j][7:], "==")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					_guard_condition = find_clock_element(global_clock, _guard_element)
					if _guard_condition == 0 && len(translator_condition[i]) == 1 {
						translator_condition[i] = append(translator_condition[i], 0)
					} else if _guard_condition == 1 {
						translator_condition[i] = append(translator_condition[i], 3)
					}

				} else if strings.Contains(loc_edge[i][j][7:], "<") {
					_slice := strings.Split(loc_edge[i][j][7:], "<")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					_guard_condition = find_clock_element(global_clock, _guard_element)
					if _guard_condition == 0 && len(translator_condition[i]) == 1 {
						translator_condition[i] = append(translator_condition[i], 0)
					} else if _guard_condition == 1 {
						translator_condition[i] = append(translator_condition[i], 1)
					}

				} else if strings.Contains(loc_edge[i][j][7:], ">") {
					_slice := strings.Split(loc_edge[i][j][7:], "<")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					_guard_condition = find_clock_element(global_clock, _guard_element)
					if _guard_condition == 0 && len(translator_condition[i]) == 1 {
						translator_condition[i] = append(translator_condition[i], 0)
					} else if _guard_condition == 1 {
						translator_condition[i] = append(translator_condition[i], 4)
					}
				}
			}
		}
	}
	fmt.Println(translator_condition)

	//xml translator
	new_doc := doc
	//new_doc.WriteTo(os.Stdout)
	_loc_num := 0
	_template_num := 0
	root := new_doc.SelectElement("nta")
	for _, e := range new_doc.FindElements("./nta/*") {
		if e.Tag == "template" {
			var _transtion_no_delete []*etree.Element
			var _transtion_delete []*etree.Element
			var _transtion_version []int
			for _, l := range e.FindElements("location") {
				if l_label := l.SelectElement("label"); l_label != nil {
					l.RemoveChild(l_label)
				}
				if translator_condition[_loc_num][1] > 0 { //생성할 loc 정보//&& len(translator_condition[_loc_num]) == 2
					etreeloc := e.CreateElement("location")
					etreeloc.CreateAttr("id", "id"+strconv.Itoa(_loc_num)+"p")
					etreeloc.CreateAttr("x", "0")
					etreeloc.CreateAttr("y", "0")
				}
				_loc_num++
			}
			for _, t := range e.FindElements("transition") {
				_init_len := len(_transtion_no_delete) + len(_transtion_delete)
				for _, i := range t.SelectElements("label") {
					if i.Attr[0].Value == "guard" { //guard가 있는 경우
						if strings.Contains(i.Text(), "<=") {
							_slice := strings.Split(i.Text(), "<=")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(global_clock, _guard_element) == 0 { //clock 요소가 없는 경우
								_transtion_no_delete = append(_transtion_no_delete, t.Copy())
								break
							} else if find_clock_element(global_clock, _guard_element) == 1 {
								_transtion_delete = append(_transtion_delete, t.Copy())
								_transtion_version = append(_transtion_version, 2)
								break
							}
						} else if strings.Contains(i.Text(), ">=") {
							_slice := strings.Split(i.Text(), ">=")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(global_clock, _guard_element) == 0 { //clock 요소가 없는 경우
								_transtion_no_delete = append(_transtion_no_delete, t.Copy())
								break
							} else if find_clock_element(global_clock, _guard_element) == 1 {
								_transtion_delete = append(_transtion_delete, t.Copy())
								_transtion_version = append(_transtion_version, 5)
								break
							}
						} else if strings.Contains(i.Text(), "==") {
							_slice := strings.Split(i.Text(), "==")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(global_clock, _guard_element) == 0 { //clock 요소가 없는 경우
								_transtion_no_delete = append(_transtion_no_delete, t.Copy())
								break
							} else if find_clock_element(global_clock, _guard_element) == 1 {
								_transtion_delete = append(_transtion_delete, t.Copy())
								_transtion_version = append(_transtion_version, 3)
								break
							}
						} else if strings.Contains(i.Text(), "<") {
							_slice := strings.Split(i.Text(), "<")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(global_clock, _guard_element) == 0 { //clock 요소가 없는 경우
								_transtion_no_delete = append(_transtion_no_delete, t.Copy())
								break
							} else if find_clock_element(global_clock, _guard_element) == 1 {
								_transtion_delete = append(_transtion_delete, t.Copy())
								_transtion_version = append(_transtion_version, 1)
								break
							}
						} else if strings.Contains(i.Text(), ">") {
							_slice := strings.Split(i.Text(), ">")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(global_clock, _guard_element) == 0 { //clock 요소가 없는 경우
								_transtion_no_delete = append(_transtion_no_delete, t.Copy())
								break
							} else if find_clock_element(global_clock, _guard_element) == 1 {
								_transtion_delete = append(_transtion_delete, t.Copy())
								_transtion_version = append(_transtion_version, 4)
								break
							}
						}
					}
				}
				if _init_len == len(_transtion_no_delete)+len(_transtion_delete) {
					_transtion_no_delete = append(_transtion_no_delete, t.Copy())
				}
				e.RemoveChild(t)
			}
			fmt.Println(_transtion_delete)
			fmt.Println(_transtion_version)

			for i := 0; i < num_loc[_template_num]; i++ {
				if translator_condition[i][0] == 1 { //exp loc 삽입 수정해야됨
					etreeloc := e.CreateElement("location")
					etreeloc.CreateAttr("id", "exp")
					etreeloc.CreateAttr("x", "0")
					etreeloc.CreateAttr("y", "0")
					break
				}
			}
			if l_init := e.SelectElement("init"); l_init != nil {
				_init_loc_id := l_init.Attr[0].Value
				e.RemoveChild(l_init)
				etreeloc := e.CreateElement("init")
				etreeloc.CreateAttr("ref", _init_loc_id)
			}
			for i := 0; i < len(_transtion_no_delete); i++ { //guard 없는 transtion 삽입
				e.AddChild(_transtion_no_delete[i])
			}
			for i := 0; i < len(_transtion_delete); i++ { //guard 있는 transtion 삽입
				e.AddChild(_transtion_delete[i])
				fmt.Println(len(e.FindElements("transition")))
				switch _transtion_version[i] {
				case 1:

				case 2:
					for q, s := range e.SelectElements("transition") {
						if len(e.FindElements("transition")) == q+1 {
							var _source_loc string
							//var _target_loc string
							var _guard_text string
							for _, a := range s.SelectElements("source") {
								_source_loc = a.Attr[0].Value
							}
							/*
								for _, a := range s.SelectElements("target") {
									_target_loc = a.Attr[0].Value
									a.Attr[0].Value = _source_loc + "p"
								}
							*/
							for _, a := range s.SelectElements("label") {
								if a.Attr[0].Value == "guard" {
									_guard_text = a.Text()
									s.RemoveChildAt(a.Index())
								}
							}
							_element := e.CreateElement("transition")
							_source_element := _element.CreateElement("source")
							_source_element.CreateAttr("ref", _source_loc)
							_target_element := _element.CreateElement("target")
							_target_element.CreateAttr("ref", _source_loc+"p")
							_guard_element := _element.CreateElement("label")
							_guard_element.CreateAttr("kind", "guard")
							_guard_element.CreateAttr("x", "0")
							_guard_element.CreateAttr("y", "0")
							_guard_element.CreateText(_guard_text) //내용 수정//exp 엣지 추가
						}
					}
				case 5:
					for q, s := range e.SelectElements("transition") {
						if len(e.FindElements("transition")) == q+1 {
							var _source_loc string
							//var _target_loc string
							var _guard_text string
							for _, a := range s.SelectElements("source") {
								_source_loc = a.Attr[0].Value
								a.Attr[0].Value = a.Attr[0].Value + "p"
							}

							for _, a := range s.SelectElements("label") {
								if a.Attr[0].Value == "guard" {
									_guard_text = a.Text()
									s.RemoveChildAt(a.Index())
								}
							}
							_element := e.CreateElement("transition")
							_source_element := _element.CreateElement("source")
							_source_element.CreateAttr("ref", _source_loc)
							_target_element := _element.CreateElement("target")
							_target_element.CreateAttr("ref", _source_loc+"p")
							_guard_element := _element.CreateElement("label")
							_guard_element.CreateAttr("kind", "guard")
							_guard_element.CreateAttr("x", "0")
							_guard_element.CreateAttr("y", "0")
							_guard_element.CreateText(_guard_text) //내용 수정

							_exp_element := e.CreateElement("transition")
							_exp_source_element := _exp_element.CreateElement("source")
							_exp_source_element.CreateAttr("ref", _source_loc+"p")
							_exp_target_element := _exp_element.CreateElement("target")
							_exp_target_element.CreateAttr("ref", "exp")
							_exp_guard_element := _exp_element.CreateElement("label")
							_exp_guard_element.CreateAttr("kind", "guard")
							_exp_guard_element.CreateAttr("x", "0")
							_exp_guard_element.CreateAttr("y", "0")
							_exp_guard_element.CreateText(_guard_text) //내용 수정
						}
					}
				}
			}

			_template_num++
		}
	}

	//transition추가 필요

	for _, t := range root.FindElements("queries") {
		root.RemoveChild(t)
	}
	new_doc.Indent(2)

	new_doc.WriteTo(os.Stdout)
	fmt.Println(num_loc)
}

func find_clock_element(clock []string, _guard_element []string) int {
	for i := 0; i < len(_guard_element); i++ {
		for j := 0; j < len(clock); j++ {
			if _guard_element[i] == clock[j] {
				return 1
			}
		}
	}
	return 0
}

func clock_extraction(init_string string) []string {
	dec_silce := strings.Split(init_string, "\n") // declaration mapping
	var clock []string
	for _, str := range dec_silce {
		//fmt.Println(str)
		splitstr := strings.Split(str, " ")
		//fmt.Println(splitstr)
		if splitstr[0] == "clock" {
			rmstr := removespace(splitstr)
			map1 := map_clock(rmstr)
			clock = append(clock, map1...)
			//map_clock(rmstr)
		}
	}
	return clock
}

func map_clock(str []string) []string {
	str = str[1:]
	var clock []string
	for _, value := range str {
		if strings.Contains(value, ",") {
			index := strings.Index(value, ",")
			clock = append(clock, value[:index])
		} else if strings.Contains(value, ";") {
			index := strings.Index(value, ";")
			clock = append(clock, value[:index])
		}
	}
	return clock
}

func removespace(str []string) []string {
	var mapstr []string
	for _, value := range str {
		if value == "//" {
			break
		} else if len(value) >= 1 {
			mapstr = append(mapstr, value)
		}
	}
	return mapstr
}

/*
고려 사항

*/
