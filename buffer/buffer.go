package buffer

//扩容切片
//长度足够的话直接返回
//长度不够,但容量足够,resize长度
//容量不够,重新分配内存
func BytesExtends(x []byte, size int, extraSize int) []byte {
	xC := cap(x)

	if xC >= size {
		//不需要扩容
		return x[:xC]
	} else {
		//需要扩容
		z := make([]byte, size+extraSize)
		copy(z, x)
		return z
	}
}
