package main

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
)

var new_type []string

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
}
func map_dec(dec []string) []string {
	for _, val := range dec {
		if strings.Contains(val, "const") {

		}
	}
	return dec
}
func map_const() {

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
