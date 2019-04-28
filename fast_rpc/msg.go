package fast_rpc

import "github.com/pineal-niwan/busybox/binary"

const (
	//消息头长度
	MsgHeadSize = 8
)

//消息接口
type IMsg interface {
	//获取命令行
	GetCmd() uint16
	//获取版本号
	GetVersion() uint16
	//获取code
	GetCode() uint32
	//序列化
	Marshal(buf []byte, option *binary.Option) (count int, out []byte, err error)
	//反序列化
	Unmarshal(buf []byte, option *binary.Option) (err error)
}

//消息头
type MsgHead struct {
	//消息体大小
	Size uint32
	//消息编号
	Cmd uint16
	//消息版本号
	Version uint16
}

//消息头获取code
func (h MsgHead) GetCode() uint32 {
	return uint32(h.Cmd) | (uint32(h.Version) << 16)
}

//反序列化消息头
func UnmarshalMsgHead(buf []byte, option *binary.Option) (head MsgHead, err error) {
	var reader *binary.BinaryHandler

	reader, err = binary.NewReadBinaryHandler(buf, option)
	if err != nil {
		return
	}

	head.Size, err = reader.ReadUint32()
	if err != nil {
		return
	}
	head.Cmd, err = reader.ReadUint16()
	if err != nil {
		return
	}
	head.Version, err = reader.ReadUint16()
	return
}

//序列化消息头
func MarshalMsgHead(writer *binary.BinaryHandler, head MsgHead) (err error) {
	err = writer.WriteUint32(head.Size)
	if err != nil {
		return
	}
	err = writer.WriteUint16(head.Cmd)
	if err != nil {
		return
	}
	err = writer.WriteUint16(head.Version)
	return
}
