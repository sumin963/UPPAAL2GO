package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/etree"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("C:\\Users\\jsm96\\gitfolder\\UPPAAL2GO\\train-gate.xml"); err != nil {
		panic(err)
	}

	var dec string
	var tem_dec []string
	ta_loc_info := make([][][]string, 0)
	ta_tran_info := make([][]string, 0)
	for _, e := range doc.FindElements("./nta/*") {
		fmt.Println(e.Tag)
		if e.Tag == "declaration" {
			dec = e.Text()
		}
		template_loc_info := make([][]string, 0)
		if e.Tag == "template" {
			if declaration := e.SelectElement("declaration"); declaration != nil {
				tem_dec = append(tem_dec, declaration.Text())
			}
			for _, l := range e.FindElements("location") {
				loc_info := make([]string, 0)
				if l.Attr[0].Key == "id" {
					fmt.Printf("  template location id : %s\n", l.Attr[0].Value)
					loc_info = append(loc_info, l.Attr[0].Value)
				}
				if l_label := l.SelectElement("label"); l_label != nil {
					fmt.Printf("  template location label value: %s\n", l_label.Text())
					fmt.Printf("  template location label kind: %s\n", l_label.Attr[0].Value)
					loc_info = append(loc_info, l_label.Text())
				}
				template_loc_info = append(template_loc_info, loc_info)
			}

			for _, t := range e.FindElements("transition") {
				tran_info := make([]string, 6, 6)
				if t_source := t.SelectElement("source"); t_source != nil {
					fmt.Printf("  template transition source: %s\n", t_source.Attr[0].Value)
					tran_info[0] = t_source.Attr[0].Value
				}
				if t_target := t.SelectElement("target"); t_target != nil {
					fmt.Printf("  template transition target: %s\n", t_target.Attr[0].Value)
					tran_info[1] = t_target.Attr[0].Value
				}
				for _, l := range t.FindElements("label") {
					if l.Attr[0].Value == "select" {
						tran_info[2] = l.Text()
					} else if l.Attr[0].Value == "guard" {
						tran_info[3] = l.Text()
					} else if l.Attr[0].Value == "synchronisation" {
						tran_info[4] = l.Text()
					} else if l.Attr[0].Value == "assignment" {
						tran_info[5] = l.Text()
					}
				}
				fmt.Println(tran_info)
				ta_tran_info = append(ta_tran_info, tran_info)
			}
			ta_loc_info = append(ta_loc_info, template_loc_info)
		}
		fmt.Println("\n")
	}
	clock := make([][]string, len(ta_loc_info))
	for i, value := range tem_dec { //declaration과 template declaraion에서 clock 추출
		clock[i] = make([]string, 0)
		_clock := clock_extraction(value)
		clock[i] = append(clock[i], _clock...)
		_clock = clock_extraction(dec)
		clock[i] = append(clock[i], _clock...)
	}

	fmt.Println(ta_loc_info)
	fmt.Println(ta_tran_info)
	fmt.Println(clock)

	new_doc := doc
	tem_num := 0
	for _, e := range new_doc.FindElements("./nta/*") {
		if e.Tag == "template" {
			for _, l := range e.FindElements("location") {
				if l_label := l.SelectElement("label"); l_label != nil {
					l.RemoveChild(l_label)
				}
			}
			var transtion_no_clock []*etree.Element
			possible := make([][]*etree.Element, 0)
			guard_possible := make([][]int, 0)
			guard_possible_condition := make([][]string, 0)
			new_loc := make([][]string, 0)
			for _, t := range e.FindElements("transition") {
				transtion_no_clock = append(transtion_no_clock, t)
				e.RemoveChild(t)
			}
			for _, l := range e.FindElements("location") {
				_possible := make([]*etree.Element, 0)
				_guard_possible := make([]int, 0)
				_guard_possible_condition := make([]string, 0)
				_new_loc := make([]string, 0)
				source := make([][]string, 0)
				for _, value := range ta_tran_info {
					if l.Attr[0].Value == value[0] {
						source = append(source, value)
					}
				}
				id := "p"
				for _, value := range source { //소스로케이션을 공유하는 트랜지션 배열
					_create_loc := 0
					if value[3] != "" { //해당 트랜지션에 가드가 있는지 확인
						if strings.Contains(value[3], "<=") {
							_slice := strings.Split(value[3], "<=")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(clock, _guard_element, tem_num) == 1 { //가드에 시간요소가 있는 경우
								_create := true
								for _, a := range _guard_possible {
									if a == find_int_element(clock, _guard_element, tem_num) {
										_create = false
									}
								}
								if _create == true { //가드가 같은값이 있을때 처리 *****수정필요
									etreeloc := e.CreateElement("location")
									for a := 0; a < _create_loc; a++ {
										id = id + "p"
									}
									etreeloc.CreateAttr("id", l.Attr[0].Value+id)
									etreeloc.CreateAttr("x", "0")
									etreeloc.CreateAttr("y", "0")
									_new_loc = append(_new_loc, l.Attr[0].Value+id)
									_create_loc++
								}
								//value(트랜지션 정보)에 해당하는 값을 transtion_no_clock에서 제거, _transtion_delete 삽입 추가
								for i, val := range transtion_no_clock {
									if val.SelectElement("source").Attr[0].Value == value[0] && val.SelectElement("target").Attr[0].Value == value[1] {
										//value 셀렉트 가드 싱크 업데이트 비교해서 찾아야함
										equal_check := true
										for _, a := range val.FindElements("label") {
											if a.Attr[0].Value == "select" {
												if a.Text() != value[3] {
													equal_check = false
												}
											} else if a.Attr[0].Value == "guard" {
												if a.Text() != value[3] {
													equal_check = false
												}
											} else if a.Attr[0].Value == "synchronisation" {
												if a.Text() != value[3] {
													equal_check = false
												}
											} else if a.Attr[0].Value == "assignment" {
												if a.Text() != value[3] {
													equal_check = false
												}
											}
											if equal_check == true {
												/*
													for _, v := range val.FindElements("label") {
														if v.Attr[0].Value == "guard" {
															_guard := val.CreateElement("label")
															_guard.CreateAttr("kind", "guard")
															_guard.CreateAttr("x", "0")
															_guard.CreateAttr("y", "0")
															new_guard := "[" + _guard_element[0] + "==" + _guard_element[1] + "]"
															_guard.CreateText(new_guard)
															val.RemoveChildAt(v.Index())
														}
													}*/

												_possible = append(_possible, val)
												transtion_no_clock = remove(transtion_no_clock, i)
												_guard_possible = append(_guard_possible, find_int_element(clock, _guard_element, tem_num))
												_append_condition := true

												for i, _ := range _guard_possible_condition {
													if _guard_possible_condition[i] == "["+_guard_element[0]+"=="+_guard_element[1]+"]" {
														_append_condition = false
													}
												}
												if _append_condition == true {
													_guard_possible_condition = append(_guard_possible_condition, "["+_guard_element[0]+"=="+_guard_element[1]+"]")

												}
											}
										}
									}
								}
							}
						} else if strings.Contains(value[3], "<") {
						} else if strings.Contains(value[3], "==") {
						} else if strings.Contains(value[3], ">=") {
							_slice := strings.Split(value[3], ">=")
							var _guard_element []string
							for _, str := range _slice {
								_guard_element = append(_guard_element, strings.Trim(str, " "))
							}
							if find_clock_element(clock, _guard_element, tem_num) == 1 { //가드에 시간요소가 있는 경우
								_create := true
								for _, a := range _guard_possible {
									if a == find_int_element(clock, _guard_element, tem_num) {
										_create = false
									}
								}
								if _create == true { //가드가 같은값이 있을때 처리
									etreeloc := e.CreateElement("location")
									for a := 0; a < _create_loc; a++ {
										id = id + "p"
									}
									etreeloc.CreateAttr("id", l.Attr[0].Value+id)
									etreeloc.CreateAttr("x", "0")
									etreeloc.CreateAttr("y", "0")
									_new_loc = append(_new_loc, l.Attr[0].Value+id)
									_create_loc++
								}
								//value(트랜지션 정보)에 해당하는 값을 transtion_no_clock에서 제거, _transtion_delete 삽입 추가
								for i, val := range transtion_no_clock {
									if val.SelectElement("source").Attr[0].Value == value[0] && val.SelectElement("target").Attr[0].Value == value[1] {
										//value 셀렉트 가드 싱크 업데이트 비교해서 찾아야함
										equal_check := true
										for _, a := range val.FindElements("label") {
											if a.Attr[0].Value == "select" {
												if a.Text() != value[3] {
													equal_check = false
												}
											} else if a.Attr[0].Value == "guard" {
												if a.Text() != value[3] {
													equal_check = false
												}
											} else if a.Attr[0].Value == "synchronisation" {
												if a.Text() != value[3] {
													equal_check = false
												}
											} else if a.Attr[0].Value == "assignment" {
												if a.Text() != value[3] {
													equal_check = false
												}
											}
											if equal_check == true {
												/*
													for _, v := range val.FindElements("label") {
														if v.Attr[0].Value == "guard" {
															_guard := val.CreateElement("label")
															_guard.CreateAttr("kind", "guard")
															_guard.CreateAttr("x", "0")
															_guard.CreateAttr("y", "0")
															new_guard := "[" + _guard_element[0] + "==" + _guard_element[1] + "]"
															_guard.CreateText(new_guard)
															val.RemoveChildAt(v.Index())
														}
													}*/
												_possible = append(_possible, val)
												transtion_no_clock = remove(transtion_no_clock, i)
												_guard_possible = append(_guard_possible, find_int_element(clock, _guard_element, tem_num))
												_append_condition := true
												for i, _ := range _guard_possible_condition {
													if _guard_possible_condition[i] == "["+_guard_element[0]+"=="+_guard_element[1]+"]" {
														_append_condition = false
													}
												}
												if _append_condition == true {
													_guard_possible_condition = append(_guard_possible_condition, "["+_guard_element[0]+"=="+_guard_element[1]+"]")

												}
											}
										}
									}
								}
							}
						} else if strings.Contains(value[3], ">") {
						}
					}
				}
				possible = append(possible, _possible)
				guard_possible = append(guard_possible, _guard_possible)
				new_loc = append(new_loc, _new_loc)
				guard_possible_condition = append(guard_possible_condition, _guard_possible_condition)
			}
			fmt.Println("possible :", possible)
			fmt.Println("guard_possible :", guard_possible)
			fmt.Println("transtion_no_clock :", transtion_no_clock)
			fmt.Println("new_loc :", new_loc)
			fmt.Println("guard_possible_condition :", guard_possible_condition)

			for j := 0; j < len(ta_loc_info[tem_num]); j++ {
				fmt.Println("33", ta_loc_info[tem_num])
				if len(ta_loc_info[tem_num][j]) > 1 {
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
			for i := 0; i < len(transtion_no_clock); i++ { //guard 없는 transtion 삽입
				e.AddChild(transtion_no_clock[i])
			}
			for i, loc_array := range new_loc {
				for j, loc := range loc_array {
					if loc != "" {
						_tran := e.CreateElement("transition")
						_tran_source := _tran.CreateElement("source")
						_tran_source.CreateAttr("ref", loc[:len(loc)-1])
						_tran_target := _tran.CreateElement("target")
						_tran_target.CreateAttr("ref", loc)
						_tran_guard := _tran.CreateElement("label")
						_tran_guard.CreateAttr("kind", "guard")
						_tran_guard.CreateAttr("x", "0")
						_tran_guard.CreateAttr("y", "0")
						_tran_guard.CreateText(guard_possible_condition[i][j])
					}
					if loc != "" && len(loc_array)-1 == j {
						for _, q := range ta_loc_info[tem_num] {
							if len(q) > 1 && q[0] == loc[:3] {
								_tran := e.CreateElement("transition")
								_tran_source := _tran.CreateElement("source")
								_tran_source.CreateAttr("ref", loc)
								_tran_target := _tran.CreateElement("target")
								_tran_target.CreateAttr("ref", "exp")
								_tran_guard := _tran.CreateElement("label")
								_tran_guard.CreateAttr("kind", "guard")
								_tran_guard.CreateAttr("x", "0")
								_tran_guard.CreateAttr("y", "0") //q[4:]

								_slice := strings.Split(q[1], "<=")
								var _inv_element []string
								for _, str := range _slice {
									_inv_element = append(_inv_element, strings.Trim(str, " "))
								}
								if strings.Contains(q[1], "<=") {
									_tran_guard.CreateText("[" + _inv_element[0] + "==" + _inv_element[1] + ")")
								} else if strings.Contains(q[1], "<") {
									_tran_guard.CreateText("(" + _inv_element[0] + "==" + _inv_element[1] + ")")

								}

							}
						}
					}
				}
			}
			for i, loc_array := range new_loc {
				for j, loc := range loc_array {
					if j == 0 {
						for a, k := range possible[i] {
							for _, m := range k.SelectElements("label") {
								if m.Attr[0].Value == "guard" {
									if strings.Contains(m.Text(), "<=") || strings.Contains(m.Text(), "<") {
										k.RemoveChild(m)
										e.AddChild(k)
										possible[i] = remove(possible[i], a)
									}
								}
							}

						}
					}
					for a, k := range possible[i] {
						if t_source := k.SelectElement("source"); t_source != nil {
							t_source.Attr[0].Value = loc
						}
						for _, m := range k.SelectElements("label") {
							if m.Attr[0].Value == "guard" {
								k.RemoveChild(m)
								e.AddChild(k)
								possible[i] = remove(possible[i], a)
							}
						}
					}

				}
			}
			tem_num++
		}
	}
	root := new_doc.SelectElement("nta")
	for _, t := range root.FindElements("queries") {
		root.RemoveChild(t)
	}
	new_doc.Indent(2)
	new_doc.WriteToFile("C:\\Users\\jsm96\\gitfolder\\UPPAAL2GO\\TADA.xml")
}
func remove(slice []*etree.Element, s int) []*etree.Element {
	return append(slice[:s], slice[s+1:]...)
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
func find_clock_element(clock [][]string, _guard_element []string, tem_num int) int {
	for i := 0; i < len(_guard_element); i++ {
		for j := 0; j < len(clock[tem_num]); j++ {
			if _guard_element[i] == clock[tem_num][j] {
				return 1
			}
		}
	}
	return 0
}
func find_int_element(clock [][]string, _guard_element []string, tem_num int) int {
	val := 0
	for i := 0; i < len(_guard_element); i++ {
		for j := 0; j < len(clock[tem_num]); j++ {
			if _guard_element[i] != clock[tem_num][j] {
				val, _ = strconv.Atoi(_guard_element[i])
				return val
			}
		}
	}
	return 0
}

//제약사항 또는 수정사항
//clock은 하나로 가정
//inv는 무조건 clock으로 가정
//엣지 생성조건이 명확하지않음
