### 二进制版本的golang序列化

支持基本类型有: 
- bool
- byte
- uint8
- int8
- int16
- int32
- uint32
- int64
- uint64
- float32
- float64
- string  

在基本类型基础上支持由基本类型组成的复合类型  
在基本类型基础上支持由基本类型组成的数组类型  
在复合类型基础上支持复合类嵌套以及复合类型数组

#### 数据格式
- 小端顺序
- string - (x x x x - 32位字符串长度) (x x x ... - string转[]byte)  
  字符串 "abc" -- [3 0 0 0 97 98 99] :  
  [3 0 0 0] -- uint32 值为3,表示字符串长度为3  
  [97 98 99] -- 字符串内容 string([]byte{97,98,99}) = "abc"
  
- array - (x x x x - 32位字符串长度) (x x x ... - 数组内容的二进制流) 

#### 使用append来组装[]byte
经过测试，append与copy性能差距不大  
平均测试结果：  
128字节 copy - 7ns append - 8ns  
256字节 copy - 9ns append - 10ns  
512字节 copy - 13ns append - 14ns  
1024字节 copy - 20ns append - 22ns  
4096字节 copy - 51ns append - 53ns  
在有时的测试中,append甚至和copy耗时一样或者更少  
此测试结果证明append内部实现应该与copy类似,应该不是用循环迭代而是采用内存拷贝的方式.  

用append的好处是在组装过程中不用去考虑切片长度的问题.
