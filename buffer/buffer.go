package buffer

const (
	BUF_RESIZE_EXTRA = 128 //缓存不够时重新分配内存时会先多预留一些字节
)

//将一个byte切片添加到一个byte切片尾部,如果前一个切片内存不够,需要重新分配内存
func AppendBytes(x, y []byte) []byte {
	xL := len(x)
	yL := len(y)
	xC := cap(x)
	if xC >= xL+yL {
		x = x[:xL+yL]
		copy(x[xL:], y)
		return x
	} else {
		z := make([]byte, xL+yL, xL+yL+BUF_RESIZE_EXTRA)
		copy(z, x)
		copy(z[xL:], y)
		return z
	}
}

//扩容切片
//长度足够的话直接返回
//长度不够,但容量足够,resize长度
//容量不够,重新分配内存
func BytesExtends(x []byte, size int) []byte {
	xL := len(x)
	xC := cap(x)

	if xL >= size {
		//不需要扩容
		return x
	}

	if xC >= size {
		//需要扩容，但不需要重新分配内存
		return x[:size]
	} else {
		z := make([]byte, size, size+BUF_RESIZE_EXTRA)
		copy(z, x)
		return z
	}
}
