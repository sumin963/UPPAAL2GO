package main

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("2doors.xml"); err != nil {
		panic(err)
	}
	f := NewFile("main")

	//root := doc.SelectElement("nta")

	for _, e := range doc.FindElements("./nta/*") {
		fmt.Println(e.Tag)
		if e.Tag == "declaration" { //declaration parsing필요
			fmt.Println(e.Text())
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
		}
		fmt.Println("\n")
	}
	fmt.Printf("%#v", f)
}
