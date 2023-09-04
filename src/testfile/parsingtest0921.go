package main

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/dave/jennifer/jen"
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
	var sys string
	var tem_dec []string
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

			loc_slice := make([][]string, 10)
			i := 0
			for _, l := range e.FindElements("location") {
				loc_slice[i] = make([]string, 3)
				if l.Attr[0].Key == "id" {
					fmt.Printf("  template location id : %s\n", l.Attr[0].Value)
					loc_slice[i][0] = l.Attr[0].Value
				}
				if l_name := l.SelectElement("name"); l_name != nil {
					fmt.Printf("  template location name : %s\n", l_name.Text())
					loc_slice[i][1] = l_name.Text()
				}
				if l_label := l.SelectElement("label"); l_label != nil {
					fmt.Printf("  template location label value: %s\n", l_label.Text())
					fmt.Printf("  template location label kind: %s\n", l_label.Attr[0].Value) //전체 attr볼수있게 수정
					loc_slice[i][2] = l_label.Text()
				}
				i = i + 1
			}

			for _, t := range e.FindElements("transition/*") {
				if t.Tag == "source" {
					fmt.Printf("  template transition source: %s\n", t.Attr[0].Value)
				}
				if t.Tag == "target" {
					fmt.Printf("  template transition target: %s\n", t.Attr[0].Value)
				}
				if t.Tag == "label" && t.Attr[0].Key == "kind" {
					fmt.Printf("  template transition label kind: %s %s\n", t.Attr[0].Value, t.Text())
				}

			}
			//template mapping
			fmt.Println(loc_slice)
			var mapping *jen.Statement

			for i, _ := range loc_slice {
				if len(loc_slice[i]) == 0 {
					break
				}
				if loc_slice[i][1] == "0" {
					mapping = Id(loc_slice[i][0]).Op(":")
				} else {
					mapping = Id(loc_slice[i][1]).Op(":")
				}
			}
			f.Func().Id(e.SelectElement("name").Text()).Params().Block(
				mapping,
			)
		}
		if e.Tag == "system" {
			fmt.Println(e.Text())
			sys = e.Text()

		}
		fmt.Println("\n")
	}
	fmt.Printf("%#v", f)

	string_mapping(dec)
	string_mapping(sys)
	for _, value := range tem_dec {
		string_mapping(value)
	}
}

func string_mapping(init_string string) {
	dec_silce := strings.Split(init_string, "\n") // declaration mapping
	var rmstr [][]string
	for _, str := range dec_silce {
		splitstr := strings.Split(str, " ")
		if len(remove_space(splitstr)) == 0 {
			continue
		} else {
			rmstr = append(rmstr, remove_space(splitstr))
		}
	}
	for _, str_r := range rmstr {
		for _, str := range str_r {
			if strings.Count("/*", str) > 0 {
				fmt.Println(str)
			}
		}
	}

	fmt.Println(rmstr)
}

func remove_space(str []string) []string {
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
1. jennifer툴을 사용하여 함수를 매핑할지
2. train-gate에서 train 템플릿 파라미터


진행해야할 사항
1. 함수 매핑
2. 로케이션, 트랜지션 매핑
3. typedef
*/
