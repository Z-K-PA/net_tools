// Code generation
// !!! Do not edit it.
// !!! Use code gen tool to generate.

package sample

import (
    "github.com/Z-K-PA/net_tools/binary_serialization"
)

//序列化例子
type SampleHandler struct {
    *binary_serialization.BinaryHandler
}

//反序列化handler,读取字节流到对象中
func NewReadSampleHandler(data []byte) (*SampleHandler, error) {
    binHandler, err := binary_serialization.NewReadBinaryHandler(data)
    if err != nil {
        return nil, err
    }else{
        return &SampleHandler{
            BinaryHandler: binHandler,
        }, nil
    }
}

//序列化handler,将对象转化成字节流
func NewWriteSampleHandler(data []byte) *SampleHandler {
    return &SampleHandler{
        BinaryHandler: binary_serialization.NewWriteBinaryHandler(data),
    }
}

//读取Sample1
func (p *SampleHandler) ReadSample1() (ret Sample1, err error) {
    ret.Field1, err = p.ReadByteArray()
    if err != nil {
        return
    }
    ret.Field2, err = p.ReadString()
    if err != nil {
        return
    }
    ret.Field3, err = p.ReadFloat64()
    if err != nil {
        return
    }
    return
}

//写入Sample1
func (p *SampleHandler) WriteSample1(ret Sample1) (err error) {
    err = p.WriteByteArray(ret.Field1)
    if err != nil {
        return
    }
    err = p.WriteString(ret.Field2)
    if err != nil {
        return
    }
    err = p.WriteFloat64(ret.Field3)
    if err != nil {
        return
    }
    return
}

//读取Sample1数组
func (p *SampleHandler) ReadSample1Array() (ret []Sample1, err error) {
    var size uint32

    //读长度
    size, err = p.ReadArrayLen()
    if err != nil {
        return
    }
    //读内容
    ret = make([]Sample1, size)
    for i := uint32(0); i < size; i++ {
        ret[i], err = p.ReadSample1()
        if err != nil {
            return
        }
    }
    return
}

//写入Sample1数组
func (p *SampleHandler) WriteSample1Array(v []Sample1) (err error) {
    //写长度
    var size int
    if v == nil{
        size = 0
    }else{
        size = len(v)
    }
    err = p.WriteArrayLen(size)
    if err != nil {
        return
    }

    //写内容
    for i := 0; i < size; i++ {
        err = p.WriteSample1(v[i])
        if err != nil {
            return
        }
    }
    return
}

//读取Sample2
func (p *SampleHandler) ReadSample2() (ret Sample2, err error) {
    ret.Id, err = p.ReadInt32()
    if err != nil {
        return
    }
    ret.Sample1List, err = p.ReadSample1Array()
    if err != nil {
        return
    }
    return
}

//写入Sample2
func (p *SampleHandler) WriteSample2(ret Sample2) (err error) {
    err = p.WriteInt32(ret.Id)
    if err != nil {
        return
    }
    err = p.WriteSample1Array(ret.Sample1List)
    if err != nil {
        return
    }
    return
}

//读取Sample2数组
func (p *SampleHandler) ReadSample2Array() (ret []Sample2, err error) {
    var size uint32

    //读长度
    size, err = p.ReadArrayLen()
    if err != nil {
        return
    }
    //读内容
    ret = make([]Sample2, size)
    for i := uint32(0); i < size; i++ {
        ret[i], err = p.ReadSample2()
        if err != nil {
            return
        }
    }
    return
}

//写入Sample2数组
func (p *SampleHandler) WriteSample2Array(v []Sample2) (err error) {
    //写长度
    var size int
    if v == nil{
        size = 0
    }else{
        size = len(v)
    }
    err = p.WriteArrayLen(size)
    if err != nil {
        return
    }

    //写内容
    for i := 0; i < size; i++ {
        err = p.WriteSample2(v[i])
        if err != nil {
            return
        }
    }
    return
}
