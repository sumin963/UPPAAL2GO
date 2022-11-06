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
	//var pre_str []string
	/*var rmstr []string
	for _, str := range dec_silce {
		splitstr := strings.Split(str, " ")
		rmstr = removespace(splitstr)
	}*/
	for _, str := range dec_silce {
		var var_const bool = false
		//fmt.Println(str)
		splitstr := strings.Split(str, " ")
		if splitstr[0] == "chan" {
			rmstr := removespace(splitstr)

			map_chan(rmstr)
		}
		if splitstr[0] == "urgent" {
			rmstr := removespace(splitstr)
			rmstr = rmstr[1:]
			map_chan(rmstr)
		}
		if splitstr[0] == "typedef" {
			rmstr := removespace(splitstr)
			map_typedef(rmstr)
		}
		if splitstr[0] == "const" {
			rmstr := removespace(splitstr)
			rmstr = rmstr[1:]
			var_const = true
			if rmstr[0] == "int" {
				map_int(rmstr, var_const)
			}

		}
		if splitstr[0] == "system" {
			rmstr := removespace(splitstr)
			map_system(rmstr)
		}
		if splitstr[0] == "clock" {
			rmstr := removespace(splitstr)
			map_clock(rmstr)
		}
		if splitstr[0] == "void" {
			rmstr := removespace(splitstr)
			map_func(rmstr)
		}
		var_const = false
	}
}
func map_func(str []string) {

}
func map_clock(str []string) {
	str = str[1:]
	for _, value := range str {
		if strings.Contains(value, ",") {
			index := strings.Index(value, ",")
			j := Id(value[:index]+"_now").Op(":=").Qual("time", "Now").Call()
			fmt.Printf("%#v\n", j)
			j = Id(value[:index]).Op(":=").Qual("time", "Since").Call(Id(value[:index] + "_now"))
			fmt.Printf("%#v\n", j)
		} else if strings.Contains(value, ";") {
			index := strings.Index(value, ";")
			j := Id(value[:index]+"_now").Op(":=").Qual("time", "Now").Call()
			fmt.Printf("%#v\n", j)
			j = Id(value[:index]).Op(":=").Qual("time", "Since").Call(Id(value[:index] + "_now"))
			fmt.Printf("%#v\n", j)
		}
	}
}

func map_typedef(str []string) { //int일때만 고려 다른 타입은 추가
	str = str[1:]
	str_join := strings.Join(str, "")
	index1 := strings.Index(str[0], "[")
	index2 := strings.Index(str[0], ",")
	index3 := strings.Index(str[0], "]")

	j := Type().Id(str_join[index3+1:len(str_join)-1]).Struct(
		Id(str_join[index3+1:len(str_join)-1]).Int(),
		Id("min").Int(),
		Id("max").Int(),
	)
	fmt.Printf("%#v\n", j)
	j = Func().Params(Id("t").Id("*"+str_join[index3+1:len(str_join)-1])).Id("new_"+str_join[index3+1:len(str_join)-1]).Params(Id("a"),
		Id("b").Int()).Block( //리턴 타입 정의 *id_t
		Qual("t", "min").Op("=").Id("a"),
		Qual("t", "max").Op("=").Id("b"),
	)

	fmt.Printf("%#v\n", j)
	j = Id(str_join[index3+1:len(str_join)-1]).Op(":=").Id("new_"+str_join[index3+1:len(str_join)-1]).Call(Id(str_join[index1+1:index2]), Id(str_join[index2+1:index3])) // 수정필요
	fmt.Printf("%#v\n", j)
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
func map_chan(str []string) {
	str = str[1:]
	for _, value := range str {
		if strings.Count(value, "[") == 1 {
			index := strings.Index(value, "[")
			fmt.Println(value[:index], ":= make([]chan chan_t,", string(value[index+1]), ")")
		} else if strings.Count(value, "[") == 2 { //2차원 배열 선언 추가
		} else {
			index := strings.Index(value, "[")
			fmt.Println(value[:index], ":= make(chan chan_t")
		}
	}
	//Id(222).Op(":=").Make(Call(Index().Chan().Struct(), Lit(222)))
}

func map_int(str []string, var_const bool) {
	str = str[1:]
	output := strings.Join(str, "")
	index := strings.Index(output, "=")
	if var_const {
		fmt.Println("const var", output[:index], "int =", string(output[index+1]))

	} else {
		fmt.Println("var", output[:index], "int =", string(output[index+1]))

	}
}

func map_system(str []string) {
	str = str[1:]
	output := strings.Join(str, "")
	index := strings.Index(output, ",")
	num := strings.Count(output, ",")

	j := Func().Id("main").Params().BlockFunc(func(g *Group) {
		for i := 0; i <= num; i++ {
			if i == num {
				g.Go().Id(output[:len(output)-1]).Call()
			} else {
				g.Go().Id(output[:index]).Call()
				output = output[index+1:]
				index = strings.Index(output, ",")
			}
		}
	})
	fmt.Printf("%#v\n", j)

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
