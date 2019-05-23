package fast_rpc_cluster

import "github.com/pineal-niwan/busybox/fast_rpc"

//集群Gossip消息头
type GMsgHead struct {
	//Id -- 消息id
	Id int64
	//内容大小
	Size uint32
}

//GMsg
type GMsg struct {
	//头部
	GMsgHead
	//内容
	fast_rpc.IMsg
}

