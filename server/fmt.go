package server

import (
	"bytes"
	"fmt"
	"strconv"
)

type Formatter interface {
	Format(d *Desktop) string
}

func GetFormatter(name string) Formatter {
	switch name {
	case "tree":
		return new(TreeFormatter)
	case "lemonbar":
		return new(LemonbarFormatter)
	default:
		return nil
	}
}

type TreeFormatter struct{}

func (this *TreeFormatter) Format(d *Desktop) string {
	return this.format(d, 0)
}

func (this *TreeFormatter) format(d *Desktop, indentLevel int) string {
	var buffer bytes.Buffer

	// generate indent string
	indentRunes := make([]rune, indentLevel)
	for i := range indentRunes {
		indentRunes[i] = '\t'
	}
	indent := string(indentRunes)

	// print attributes
	for _, attr := range d.AttrList() {
		buffer.WriteString(indent)
		buffer.WriteString(fmt.Sprintf("\"%s\": %s\n", attr[0], attr[1]))
	}

	// print children
	if d.NumChildren() > 0 {
		buffer.WriteString(indent)
		buffer.WriteString(fmt.Sprintf("Focused child: %d\n", d.FocusedChild()))

		for i := 0; i < d.NumChildren(); i++ {
			buffer.WriteString(fmt.Sprintf("%sChild %d\n", indent, i))
			buffer.WriteString(this.format(d.ChildAt(i), indentLevel+1))
		}
	}

	str := buffer.String()
	if indentLevel == 0 && len(str) > 0 {
		return str[:len(str)-1]
	} else {
		return str
	}
}

type LemonbarFormatter struct{}

func (this *LemonbarFormatter) Format(d *Desktop) string {
	focusColor := "#525252"
	padding := "    "

	var buffer bytes.Buffer

	for i := 0; i < d.NumChildren(); i++ {
		if i == d.FocusedChild() {
			buffer.WriteString("%{B")
			buffer.WriteString(focusColor)
			buffer.WriteString("}")
		}

		buffer.WriteString("%{A:")
		buffer.WriteString(strconv.Itoa(i))
		buffer.WriteString(":}")

		buffer.WriteString(padding)
		buffer.WriteString(d.ChildAt(i).Attr("name"))
		buffer.WriteString(padding)

		buffer.WriteString("%{A}")

		if i == d.FocusedChild() {
			buffer.WriteString("%{B-}")
		}
	}

	return buffer.String()
}
