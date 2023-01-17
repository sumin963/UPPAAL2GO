package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	. "github.com/dave/jennifer/jen"
)

var new_type []string

type sort_data struct {
	edge  *etree.Element
	clock int
}
type transition_e_prime struct {
	source string
	target string
	guard  string
}
type tada_transition struct {
	source string
	target string
	action string
	update string
}

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("train-gate.xml"); err != nil {
		panic(err)
	}
	f := NewFile("main")

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

	fmt.Printf("%#v", f)
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
			var transtion []*etree.Element

			for _, t := range e.FindElements("transition") {
				transtion = append(transtion, t)
				e.RemoveChild(t)
			}

			_isinvariant := false
			for _, val := range ta_loc_info[tem_num] {
				if len(val) > 1 {
					newloc := e.CreateElement("location")
					newloc.CreateAttr("id", "exp")
					newloc.CreateAttr("x", "0")
					newloc.CreateAttr("y", "0")
					_isinvariant = true
					break
				}
			}
			time_edge := make([][]transition_e_prime, 0)
			tada_edge := make([][]tada_transition, 0)
			for loc_num, l := range e.FindElements("location") {
				Edge_prime := share_source_loc(transtion, l)
				Edge_prime, Edge_double_prime := Edge_Sort(Edge_prime, clock, tem_num)
				time_flow_loc := make([]string, 0)
				time_flow_edge := make([]transition_e_prime, 0)
				tada_flow_edge := make([]tada_transition, 0)

				_id := l.Attr[0].Value
				_guard := ""
				for num, edge := range Edge_prime {
					for _, label := range edge.FindElements("label") {
						if label.Attr[0].Value == "guard" {
							_guard = label.Text()
							_guard = transformEdge(_guard)
						}
					}
					_edge_element := transition_e_prime{_id, _id + "p", _guard}
					time_flow_edge = append(time_flow_edge, _edge_element)
					_id = _id + "p"
					time_flow_loc = append(time_flow_loc, _id)

					if len(Edge_prime)-1 == num && _isinvariant {
						_inv := ta_loc_info[tem_num][loc_num][1]
						_inv = transformEdge_inv(_inv)
						_edge_element := transition_e_prime{_id, "exp", _inv}
						time_flow_edge = append(time_flow_edge, _edge_element)
					} else if true { //guard가 at일때 추가

					}

				}
				for _, val := range time_flow_loc {
					newloc := e.CreateElement("location")
					newloc.CreateAttr("id", val)
					newloc.CreateAttr("x", "0")
					newloc.CreateAttr("y", "0")
				}
				for i, val := range time_flow_edge {
					for _, edge := range Edge_prime {
						_action := ""
						_update := ""
						_target := ""
						if t_target := edge.SelectElement("target"); t_target != nil {
							_target = t_target.Attr[0].Value
						}
						for _, l := range edge.FindElements("label") {

							if l.Attr[0].Value == "select" {
							} else if l.Attr[0].Value == "synchronisation" {
								_action = l.Text()
							} else if l.Attr[0].Value == "assignment" {
								_update = l.Text()
							}
						}
						if i == 0 && ispossible(val, edge, clock, tem_num, "false") {
							_edge_element := tada_transition{val.source, _target, _action, _update}
							tada_flow_edge = append(tada_flow_edge, _edge_element)
						} else if i != 0 && ispossible(val, edge, clock, tem_num, get_form(time_flow_edge[i-1].guard)) {
							_edge_element := tada_transition{val.source, _target, _action, _update}
							tada_flow_edge = append(tada_flow_edge, _edge_element)
						}
					}
					for _, edge := range Edge_double_prime {
						_action := ""
						_update := ""
						_target := ""
						if t_target := edge.SelectElement("target"); t_target != nil {
							_target = t_target.Attr[0].Value
						}
						for _, l := range edge.FindElements("label") {

							if l.Attr[0].Value == "select" {
							} else if l.Attr[0].Value == "synchronisation" {
								_action = l.Text()
							} else if l.Attr[0].Value == "assignment" {
								_update = l.Text()
							}
						}
						_edge_element := tada_transition{val.source, _target, _action, _update}
						tada_flow_edge = append(tada_flow_edge, _edge_element)
					}

				} //timeflow가 생성이 안될때 더블프라임엣지 생서하는 루프필요
				if len(time_flow_edge) == 0 {
					for _, edge := range Edge_double_prime {
						_action := ""
						_update := ""
						_target := ""
						_source := ""
						if t_source := edge.SelectElement("source"); t_source != nil {
							_source = t_source.Attr[0].Value
						}
						if t_target := edge.SelectElement("target"); t_target != nil {
							_target = t_target.Attr[0].Value
						}
						for _, l := range edge.FindElements("label") {

							if l.Attr[0].Value == "select" {
							} else if l.Attr[0].Value == "synchronisation" {
								_action = l.Text()
							} else if l.Attr[0].Value == "assignment" {
								_update = l.Text()
							}
						}
						_edge_element := tada_transition{_source, _target, _action, _update}
						tada_flow_edge = append(tada_flow_edge, _edge_element)
					}
				}

				time_edge = append(time_edge, time_flow_edge)
				tada_edge = append(tada_edge, tada_flow_edge)
			}

			if l_init := e.SelectElement("init"); l_init != nil {
				_init_loc_id := l_init.Attr[0].Value
				e.RemoveChild(l_init)
				etreeloc := e.CreateElement("init")
				etreeloc.CreateAttr("ref", _init_loc_id)
			}

			for _, j := range time_edge {
				for _, val := range j {
					_tran := e.CreateElement("transition")
					_tran_source := _tran.CreateElement("source")
					_tran_source.CreateAttr("ref", val.source)
					_tran_target := _tran.CreateElement("target")
					_tran_target.CreateAttr("ref", val.target)
					_tran_guard := _tran.CreateElement("label")
					_tran_guard.CreateAttr("kind", "guard")
					_tran_guard.CreateAttr("x", "0")
					_tran_guard.CreateAttr("y", "0")
					_tran_guard.CreateText(val.guard)
				}
			}
			for _, j := range tada_edge {
				for _, val := range j {
					_tran := e.CreateElement("transition")
					_tran_source := _tran.CreateElement("source")
					_tran_source.CreateAttr("ref", val.source)
					_tran_target := _tran.CreateElement("target")
					_tran_target.CreateAttr("ref", val.target)
					if val.action != "" {
						_tran_label := _tran.CreateElement("label")
						_tran_label.CreateAttr("kind", "synchronisation")
						_tran_label.CreateAttr("x", "0")
						_tran_label.CreateAttr("y", "0")
						_tran_label.CreateText(val.action)
					}
					if val.update != "" {
						_tran_label := _tran.CreateElement("label")
						_tran_label.CreateAttr("kind", "assignment")
						_tran_label.CreateAttr("x", "0")
						_tran_label.CreateAttr("y", "0")
						_tran_label.CreateText(val.update)
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
	new_doc.WriteTo(os.Stdout)
}

func ispossible(val transition_e_prime, edge *etree.Element, clock [][]string, tem_num int, before string) bool {
	result := false
	for _, label := range edge.FindElements("label") {
		if label.Attr[0].Value == "guard" {
			_guard := label.Text()
			if strings.Contains(_guard, "<=") {
				_guard_element := del_black(_guard, "<=")
				_timeflow_element := del_black(val.guard, "==")
				if (find_int_element(clock, _guard_element, tem_num) > find_int_element_timeflow(clock, _timeflow_element, tem_num)) || ((find_int_element(clock, _guard_element, tem_num) == find_int_element_timeflow(clock, _timeflow_element, tem_num)) && (get_form(val.guard) == "[x==n]")) {
					return true
				}
			} else if strings.Contains(_guard, "<") {
				_guard_element := del_black(_guard, "<")
				_timeflow_element := del_black(val.guard, "==")
				if (find_int_element(clock, _guard_element, tem_num) > find_int_element_timeflow(clock, _timeflow_element, tem_num)) || ((find_int_element(clock, _guard_element, tem_num) == find_int_element_timeflow(clock, _timeflow_element, tem_num)) && (get_form(val.guard) == "(x==n]")) {
					return true
				}
			} else if strings.Contains(_guard, "==") {
				_guard_element := del_black(_guard, "==")
				_timeflow_element := del_black(val.guard, "==")
				if (find_int_element(clock, _guard_element, tem_num) == find_int_element_timeflow(clock, _timeflow_element, tem_num)) && (get_form(val.guard) == "[x==n)") {
					return true
				}
			} else if strings.Contains(_guard, ">=") {
				_guard_element := del_black(_guard, ">=")
				_timeflow_element := del_black(val.guard, "==")
				if (find_int_element(clock, _guard_element, tem_num) < find_int_element_timeflow(clock, _timeflow_element, tem_num)) || ((find_int_element(clock, _guard_element, tem_num) == find_int_element_timeflow(clock, _timeflow_element, tem_num)) && (before == "(x==n]")) {
					return true
				}
			} else if strings.Contains(_guard, ">") {
				_guard_element := del_black(_guard, ">")
				_timeflow_element := del_black(val.guard, "==")
				if (find_int_element(clock, _guard_element, tem_num) < find_int_element_timeflow(clock, _timeflow_element, tem_num)) || ((find_int_element(clock, _guard_element, tem_num) == find_int_element_timeflow(clock, _timeflow_element, tem_num)) && (before == "(x==n]")) {
					return true
				}
			}
		}
	}
	return result
}

func get_form(guard string) string {
	result := ""
	_start := guard[:1]
	_end := guard[len(guard)-1:]
	result = _start + "x==n" + _end
	return result
}
func transformEdge_inv(guard string) string {
	if strings.Contains(guard, "<=") {
		_guard_element := del_black(guard, "<=")
		guard = "[" + _guard_element[0] + "==" + _guard_element[1] + ")"
	} else if strings.Contains(guard, "<") { //inv < 일때 수정
		_guard_element := del_black(guard, "<")
		guard = "(" + _guard_element[0] + "==" + _guard_element[1] + "]"
	}

	return guard
}
func transformEdge(guard string) string {
	if strings.Contains(guard, "<=") {
		_guard_element := del_black(guard, "<=")
		guard = "[" + _guard_element[0] + "==" + _guard_element[1] + "]"
	} else if strings.Contains(guard, "<") {
		_guard_element := del_black(guard, "<")
		guard = "(" + _guard_element[0] + "==" + _guard_element[1] + "]"
	} else if strings.Contains(guard, "==") {
		_guard_element := del_black(guard, "==")
		guard = "(" + _guard_element[0] + "==" + _guard_element[1] + "]"
	} else if strings.Contains(guard, ">=") {
		_guard_element := del_black(guard, ">=")
		guard = "[" + _guard_element[0] + "==" + _guard_element[1] + "]"
	} else if strings.Contains(guard, ">") {
		_guard_element := del_black(guard, ">")
		guard = "[" + _guard_element[0] + "==" + _guard_element[1] + ")"
	}

	return guard
}
func del_black(guard string, operater string) []string {
	var _guard_element []string
	_slice := strings.Split(guard, operater)
	for _, str := range _slice {
		_guard_element = append(_guard_element, strings.Trim(str, " "))
	}
	return _guard_element
}
func Edge_Sort(edge []*etree.Element, clock [][]string, tem_num int) ([]*etree.Element, []*etree.Element) {
	Edge_prime := make([]*etree.Element, 0)
	Edge_double_prime := make([]*etree.Element, 0)
	edge_no_clock := true
	time_flow := make([]sort_data, 0)
	for _, i := range edge {
		for _, j := range i.FindElements("label") {
			if j.Attr[0].Value == "guard" {
				if strings.Contains(j.Text(), "<=") {
					_slice := strings.Split(j.Text(), "<=")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					if find_clock_element(clock, _guard_element, tem_num) == 1 {
						if find_int_element(clock, _guard_element, tem_num) != 0 {
							time_flow = append(time_flow, sort_data{i, find_int_element(clock, _guard_element, tem_num)})
							edge_no_clock = false

						}
					}
				} else if strings.Contains(j.Text(), "<") {
					_slice := strings.Split(j.Text(), "<")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					if find_clock_element(clock, _guard_element, tem_num) == 1 {
						if find_int_element(clock, _guard_element, tem_num) != 0 {
							time_flow = append(time_flow, sort_data{i, find_int_element(clock, _guard_element, tem_num)})
							edge_no_clock = false

						}
					}
				} else if strings.Contains(j.Text(), "==") {
					_slice := strings.Split(j.Text(), "==")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					if find_clock_element(clock, _guard_element, tem_num) == 1 {
						if find_int_element(clock, _guard_element, tem_num) != 0 {
							time_flow = append(time_flow, sort_data{i, find_int_element(clock, _guard_element, tem_num)})
							edge_no_clock = false

						}
					}
				} else if strings.Contains(j.Text(), ">=") {
					_slice := strings.Split(j.Text(), ">=")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					if find_clock_element(clock, _guard_element, tem_num) == 1 {
						if find_int_element(clock, _guard_element, tem_num) != 0 {
							time_flow = append(time_flow, sort_data{i, find_int_element(clock, _guard_element, tem_num)})
							edge_no_clock = false

						}
					}
				} else if strings.Contains(j.Text(), ">") {
					_slice := strings.Split(j.Text(), ">")
					var _guard_element []string
					for _, str := range _slice {
						_guard_element = append(_guard_element, strings.Trim(str, " "))
					}
					if find_clock_element(clock, _guard_element, tem_num) == 1 {
						if find_int_element(clock, _guard_element, tem_num) != 0 {
							time_flow = append(time_flow, sort_data{i, find_int_element(clock, _guard_element, tem_num)})
							edge_no_clock = false
						}
					}
				}
			}
		}
		if edge_no_clock == true {
			Edge_double_prime = append(Edge_double_prime, i)
		}
		edge_no_clock = true
	}
	sort.Slice(time_flow, func(i, j int) bool {
		return time_flow[i].clock < time_flow[j].clock
	})
	for i := 0; i < len(time_flow); i++ {
		Edge_prime = append(Edge_prime, time_flow[i].edge)
	}
	return Edge_prime, Edge_double_prime
}

func share_source_loc(edge []*etree.Element, loc *etree.Element) []*etree.Element {
	edge_return := make([]*etree.Element, 0)
	Edge_prime := make([]*etree.Element, 0)
	for _, i := range edge {
		for _, j := range i.FindElements("source") {
			if j.Attr[0].Value == loc.Attr[0].Value {
				Edge_prime = append(Edge_prime, i)
			} else {
				edge_return = append(edge_return, i)
			}
		}
	}
	return Edge_prime
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
func find_int_element_timeflow(clock [][]string, _guard_element []string, tem_num int) int {
	val := 0

	for i := 0; i < len(_guard_element); i++ {
		if i == 0 {
			_guard_element[i] = _guard_element[i][1:]
		} else {
			_guard_element[i] = _guard_element[i][:(len(_guard_element[i]) - 1)]
		}
		for j := 0; j < len(clock[tem_num]); j++ {
			if _guard_element[i] != clock[tem_num][j] {
				val, _ = strconv.Atoi(_guard_element[i])
				return val
			}
		}
	}
	return -1
}
