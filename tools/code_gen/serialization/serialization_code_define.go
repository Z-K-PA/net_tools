package serialization

import (
	"strings"
	"text/template"
)

type Field struct {
	//属性名
	Name string `yaml:"name" json:"name"`
	//属性类型
	TypeDefine string `yaml:"typeDefine" json:"typeDefine"`
	//注释
	Comment string `yaml:"comment" json:"comment"`
}

type Object struct {
	//消息编号
	Cmd uint16 `yaml:"cmd" json:"cmd"`
	//消息版本
	Version uint16 `yaml:"ver" json:"ver"`
	//结构体名称
	Name string `yaml:"name" json:"name"`
	//结构体注释
	Comment string `yaml:"comment" json:"comment"`
	//结构体中的属性列表
	Fields []Field `yaml:"fields" json:"fields"`
}

type Package struct {
	//包名
	Package string `yaml:"package" json:"package"`
	//序列化handler的名称
	Name string `yaml:"name" json:"name"`
	//注释
	Comment string `yaml:"comment" json:"comment"`
	//定义的结构体列表
	Objects []Object `yaml:"objects" json:"objects"`

	//序列化最大长度
	DataMaxLen int `yaml:"dataMaxLen" json:"dataMaxLen"`
	//支持的字符串长度
	StringMaxLen int `yaml:"stringMaxLen" json:"stringMaxLen"`
	//支持的数组最大长度
	ArrayMaxLen int `yaml:"arrayMaxLen" json:"arrayMaxLen"`
	//扩大容量时额外多分配的字节数
	ExtendExtraSize int `yaml:"extendExtraSize" json:"extendExtraSize"`
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
