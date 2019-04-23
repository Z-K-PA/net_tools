package serialization

import (
	"strings"
	"text/template"
)

type Field struct {
	Name       string `yaml:"name" json:"name"`
	TypeDefine string `yaml:"typeDefine" json:"typeDefine"`
	Comment    string `yaml:"comment" json:"comment"`
}

type Object struct {
	Name    string  `yaml:"name" json:"name"`
	Comment string  `yaml:"comment" json:"comment"`
	Fields  []Field `yaml:"fields" json:"fields"`
}

type Package struct {
	Package string   `yaml:"package" json:"package"`
	Name    string   `yaml:"name" json:"name"`
	Comment string   `yaml:"comment" json:"comment"`
	Objects []Object `yaml:"objects" json:"objects"`
}

var (
	FuncHash = template.FuncMap{
		"addArrayPrefix": func(s string) string {
			if strings.Contains(s, "Array") {
				return "[]" + strings.TrimSuffix(s, "Array")
			}
			return s
		},
		"upperLetter": func(s string) string {
			x := []rune(s)
			if x[0] >= 'a' && x[0] <= 'z' {
				x[0] -= 'a' - 'A'
			}
			return string(x)
		},
	}
)
