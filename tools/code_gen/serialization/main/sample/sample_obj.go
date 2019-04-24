// Code generation
// !!! Do not edit it.
// !!! Use code gen tool to generate.

package sample

//对象定义1
type Sample1 struct {
	//字节流属性
	Field1 []byte
	//字符串属性
	Field2 string
	//浮点数属性
	Field3 float64
}

//对象定义2
type Sample2 struct {
	//id定义
	Id int32
	//sample1 list
	Sample1List []Sample1
}
