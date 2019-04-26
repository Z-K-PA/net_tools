package binary

const (
	DATA_MAX_LEN = 32 * 1024 * 1024 //序列化最大长度为32M

	STR_MAX_LEN = 1024 * 1024 * 16 //支持的字符串长度最大为16M

	ARRAY_MAX_LEN = 1024 * 1024 * 16 //支持的数组最大长度为16M

	BUF_SIZE_INIT = 256 //初始256字节容量

	MinBufferSize = 16
)
