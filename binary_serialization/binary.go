package binary_serialization

import (
	"errors"
	"fmt"
	"math"
)

var (
	//数据越界
	ErrOverflow = errors.New("binary handler data overflow")

	//字符串长度越界
	ErrStringOverflow = errors.New("binary handler string overflow")

	//数组越界
	ErrArrayOverflow = errors.New("binary handler array overflow")
)

//序列化结构体定义
type BinaryHandler struct {
	pos  int    //当前buffer的指针
	data []byte //二进制流的内容切片
}

//新建对象
func NewBinaryHandler(data []byte) *BinaryHandler {
	if data == nil {
		return &BinaryHandler{
			data: make([]byte, 0, BUF_SIZE_INIT),
		}
	} else {
		return &BinaryHandler{
			data: data,
		}
	}
}

//获取buffer
func (bh *BinaryHandler) Data() []byte {
	return bh.data
}

//检查当前位置加上一个偏移后是否越界
func (bh *BinaryHandler) checkPos(offset uint32) error {
	if bh.pos > DATA_MAX_LEN {
		return ErrOverflow
	}
	if bh.pos+int(offset) > len(bh.data) {
		return errors.New(fmt.Sprintf("binary handler overflow, pos: %d offset: %d", bh.pos, offset))
	}
	return nil
}

//在写入数据后检查是否越界
func (bh *BinaryHandler) checkOverflow() error {
	if len(bh.data) > DATA_MAX_LEN {
		return ErrOverflow
	}
	return nil
}

//写入bytes
func (bh *BinaryHandler) appendData(data []byte) error {
	if len(bh.data)+len(data) > DATA_MAX_LEN {
		return ErrOverflow
	}
	bh.data = append(bh.data, data...)
	return nil
}

//读取bool型
// - 用一个byte表示bool
func (bh *BinaryHandler) ReadBool() (ret bool, err error) {
	err = bh.checkPos(1)
	if err != nil {
		return
	}
	r := bh.data[bh.pos]
	bh.pos += 1
	ret = r != 0
	return
}

//写入bool型
// - 用一个byte表示bool
func (bh *BinaryHandler) WriteBool(b bool) error {
	if b {
		bh.data = append(bh.data, byte(1))
	} else {
		bh.data = append(bh.data, byte(0))
	}
	return bh.checkOverflow()
}

//读取byte型
func (bh *BinaryHandler) ReadByte() (ret byte, err error) {
	err = bh.checkPos(1)
	if err != nil {
		return
	}
	ret = bh.data[bh.pos]
	bh.pos += 1
	return
}

//写入byte型
func (bh *BinaryHandler) WriteByte(bt byte) error {
	bh.data = append(bh.data, bt)
	return bh.checkOverflow()
}

//读取int8型
func (bh *BinaryHandler) ReadInt8() (ret int8, err error) {
	err = bh.checkPos(1)
	if err != nil {
		return
	}
	ret = int8(bh.data[bh.pos])
	bh.pos += 1
	return
}

//写入int8型
func (bh *BinaryHandler) WriteInt8(i8 int8) error {
	bh.data = append(bh.data, byte(i8))
	return bh.checkOverflow()
}

//读取uint8型
func (bh *BinaryHandler) ReadUint8() (ret uint8, err error) {
	err = bh.checkPos(1)
	if err != nil {
		return
	}
	ret = bh.data[bh.pos]
	bh.pos += 1
	return
}

//写入uint8型
func (bh *BinaryHandler) WriteUint8(i8 uint8) error {
	bh.data = append(bh.data, i8)
	return bh.checkOverflow()
}

//读取uint16型
func (bh *BinaryHandler) ReadUint16() (ret uint16, err error) {
	err = bh.checkPos(2)
	if err != nil {
		return
	}

	ret = uint16(bh.data[bh.pos]) |
		(uint16(bh.data[bh.pos+1]) << 8)

	bh.pos += 2
	return
}

//写入uint16型
func (bh *BinaryHandler) WriteUint16(v uint16) error {
	bh.data = append(bh.data,
		byte(v),
		byte(v>>8))

	return bh.checkOverflow()
}

//读取int16型
func (bh *BinaryHandler) ReadInt16() (ret int16, err error) {
	err = bh.checkPos(2)
	if err != nil {
		return
	}

	ret = int16(bh.data[bh.pos]) |
		(int16(bh.data[bh.pos+1]) << 8)

	bh.pos += 2
	return
}

//写入int16型
func (bh *BinaryHandler) WriteInt16(v int16) error {
	bh.data = append(bh.data,
		byte(v),
		byte(v>>8))

	return bh.checkOverflow()
}

//读取uint32型
func (bh *BinaryHandler) ReadUint32() (ret uint32, err error) {
	err = bh.checkPos(4)
	if err != nil {
		return
	}

	ret = uint32(bh.data[bh.pos]) |
		(uint32(bh.data[bh.pos+1]) << 8) |
		(uint32(bh.data[bh.pos+2]) << 16) |
		(uint32(bh.data[bh.pos+3]) << 24)

	bh.pos += 4
	return
}

//写入uint32型
func (bh *BinaryHandler) WriteUint32(v uint32) error {
	bh.data = append(bh.data,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24))

	return bh.checkOverflow()
}

//读取int32型
func (bh *BinaryHandler) ReadInt32() (ret int32, err error) {
	err = bh.checkPos(4)
	if err != nil {
		return
	}

	ret = int32(bh.data[bh.pos]) |
		(int32(bh.data[bh.pos+1]) << 8) |
		(int32(bh.data[bh.pos+2]) << 16) |
		(int32(bh.data[bh.pos+3]) << 24)

	bh.pos += 4
	return
}

//写入int32型
func (bh *BinaryHandler) WriteInt32(v int32) error {
	bh.data = append(bh.data,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24))

	return bh.checkOverflow()
}

//读取uint64型
func (bh *BinaryHandler) ReadUint64() (ret uint64, err error) {
	err = bh.checkPos(8)
	if err != nil {
		return
	}

	ret = uint64(bh.data[bh.pos]) |
		(uint64(bh.data[bh.pos+1]) << 8) |
		(uint64(bh.data[bh.pos+2]) << 16) |
		(uint64(bh.data[bh.pos+3]) << 24) |
		(uint64(bh.data[bh.pos+4]) << 32) |
		(uint64(bh.data[bh.pos+5]) << 40) |
		(uint64(bh.data[bh.pos+6]) << 48) |
		(uint64(bh.data[bh.pos+7]) << 56)

	bh.pos += 8
	return
}

//写入uint64型
func (bh *BinaryHandler) WriteUint64(v uint64) error {
	bh.data = append(bh.data,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56))

	return bh.checkOverflow()
}

//读取int64型
func (bh *BinaryHandler) ReadInt64() (ret int64, err error) {
	err = bh.checkPos(8)
	if err != nil {
		return
	}

	ret = int64(bh.data[bh.pos]) |
		(int64(bh.data[bh.pos+1]) << 8) |
		(int64(bh.data[bh.pos+2]) << 16) |
		(int64(bh.data[bh.pos+3]) << 24) |
		(int64(bh.data[bh.pos+4]) << 32) |
		(int64(bh.data[bh.pos+5]) << 40) |
		(int64(bh.data[bh.pos+6]) << 48) |
		(int64(bh.data[bh.pos+7]) << 56)

	bh.pos += 8
	return
}

//写入int64型
func (bh *BinaryHandler) WriteInt64(v int64) error {
	bh.data = append(bh.data,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56))

	return bh.checkOverflow()
}

//读取float32型
func (bh *BinaryHandler) ReadFloat32() (ret float32, err error) {
	var _ret uint32
	_ret, err = bh.ReadUint32()
	if err != nil {
		return
	}
	ret = math.Float32frombits(_ret)
	return
}

//写入float32型
func (bh *BinaryHandler) WriteFloat32(v float32) error {
	var _ret uint32
	_ret = math.Float32bits(v)
	return bh.WriteUint32(_ret)
}

//读取float64型
func (bh *BinaryHandler) ReadFloat64() (ret float64, err error) {
	var _ret uint64
	_ret, err = bh.ReadUint64()
	if err != nil {
		return
	}
	ret = math.Float64frombits(_ret)
	return
}

//写入float64型
func (bh *BinaryHandler) WriteFloat64(v float64) error {
	var _ret uint64
	_ret = math.Float64bits(v)
	return bh.WriteUint64(_ret)
}

//读取string
func (bh *BinaryHandler) ReadString() (ret string, err error) {
	var size uint32

	//读取长度
	size, err = bh.ReadUint32()
	if err != nil {
		return
	}
	//检查字符串是否越界
	if size > STR_MAX_LEN {
		err = ErrStringOverflow
		return
	}

	//检查是否有足够的缓冲区
	err = bh.checkPos(size)
	if err != nil {
		return
	}

	ret = string(bh.data[bh.pos : bh.pos+int(size)])
	bh.pos += int(size)
	return
}

//写入string
func (bh *BinaryHandler) WriteString(s string) (err error) {
	b := []byte(s)
	size := len(b)

	if size > STR_MAX_LEN {
		err = ErrStringOverflow
		return
	}

	err = bh.WriteUint32(uint32(size))
	if err != nil {
		return
	}
	return bh.appendData(b)
}

//读取一个数组长度，并判断其是否越界
func (bh *BinaryHandler) ReadArrayLen() (size uint32, err error) {
	//读长度
	size, err = bh.ReadUint32()
	if err != nil {
		return
	}
	if size > ARRAY_MAX_LEN {
		err = ErrArrayOverflow
	}
	return
}

//写入一个数组长度，并判断其是否越界
func (bh *BinaryHandler) WriteArrayLen(size int) (err error) {
	if size > ARRAY_MAX_LEN {
		err = ErrArrayOverflow
		return
	}
	return bh.WriteUint32(uint32(size))
}

//读取byte数组
func (bh *BinaryHandler) ReadByteArray() (ret []byte, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	err = bh.checkPos(size)
	if err != nil {
		return
	}

	ret = make([]byte, size)
	copy(ret, bh.data[bh.pos:bh.pos+int(size)])
	bh.pos += int(size)
	return
}

//写入byte数组
func (bh *BinaryHandler) WriteByteArray(v []byte) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	err = bh.appendData(v)
	return
}

//读取uint8数组
func (bh *BinaryHandler) ReadUint8Array() (ret []uint8, err error) {
	return bh.ReadByteArray()
}

//写入uint8数组
func (bh *BinaryHandler) WriteUint8Array(v []uint8) (err error) {
	return bh.WriteByteArray(v)
}

//读取int8数组
func (bh *BinaryHandler) ReadInt8Array() (ret []int8, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]int8, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadInt8()
		if err != nil {
			return
		}
	}
	return
}

//写入int8数组
func (bh *BinaryHandler) WriteInt8Array(v []int8) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteInt8(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取bool数组
func (bh *BinaryHandler) ReadBoolArray() (ret []bool, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]bool, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadBool()
		if err != nil {
			return
		}
	}
	return
}

//写入bool数组
func (bh *BinaryHandler) WriteBoolArray(v []bool) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteBool(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取int16数组
func (bh *BinaryHandler) ReadInt16Array() (ret []int16, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]int16, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadInt16()
		if err != nil {
			return
		}
	}
	return
}

//写入int16数组
func (bh *BinaryHandler) WriteInt16Array(v []int16) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteInt16(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取uint16数组
func (bh *BinaryHandler) ReadUint16Array() (ret []uint16, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]uint16, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadUint16()
		if err != nil {
			return
		}
	}
	return
}

//写入uint16数组
func (bh *BinaryHandler) WriteUint16Array(v []uint16) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteUint16(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取int32数组
func (bh *BinaryHandler) ReadInt32Array() (ret []int32, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]int32, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadInt32()
		if err != nil {
			return
		}
	}
	return
}

//写入int32数组
func (bh *BinaryHandler) WriteInt32Array(v []int32) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteInt32(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取uint32数组
func (bh *BinaryHandler) ReadUint32Array() (ret []uint32, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]uint32, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadUint32()
		if err != nil {
			return
		}
	}
	return
}

//写入uint32数组
func (bh *BinaryHandler) WriteUint32Array(v []uint32) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteUint32(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取int64数组
func (bh *BinaryHandler) ReadInt64Array() (ret []int64, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]int64, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadInt64()
		if err != nil {
			return
		}
	}
	return
}

//写入int64数组
func (bh *BinaryHandler) WriteInt64Array(v []int64) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteInt64(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取uint64数组
func (bh *BinaryHandler) ReadUint64Array() (ret []uint64, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]uint64, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadUint64()
		if err != nil {
			return
		}
	}
	return
}

//写入uint64数组
func (bh *BinaryHandler) WriteUint64Array(v []uint64) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteUint64(v[i])
		if err != nil {
			return
		}
	}
	return
}

//读取string数组
func (bh *BinaryHandler) ReadStringArray() (ret []string, err error) {
	var size uint32

	//读长度
	size, err = bh.ReadArrayLen()
	if err != nil {
		return
	}
	//读内容
	ret = make([]string, size)
	for i := uint32(0); i < size; i++ {
		ret[i], err = bh.ReadString()
		if err != nil {
			return
		}
	}
	return
}

//写入string数组
func (bh *BinaryHandler) WriteStringArray(v []string) (err error) {
	//写长度
	var size int
	if v == nil {
		size = 0
	} else {
		size = len(v)
	}
	err = bh.WriteArrayLen(size)
	if err != nil {
		return
	}

	//写内容
	for i := 0; i < size; i++ {
		err = bh.WriteString(v[i])
		if err != nil {
			return
		}
	}
	return
}
