package sds

//import "fmt"

const sds_type = "string"

//SDS 包含 三个属性 len;free;[]byte
// 在设计时候将Len设置成可以包外显示的
type SDS struct {
	Len uint32 // 计算字符串长度

	free uint32 // 判断当前还剩余多少空间

	buf []byte // 存放value

}

// sds初始化 初始化长度为0 ，容量为10
func Newsds() *SDS {
	//free_value := 10
	return &SDS{
		Len:  0,
		free: 0,
		buf:  make([]byte, 0),
	}
}

//计算字符串长度
func (s *SDS) Len_string() int {
	s.Len = uint32(len(s.buf))
	return int(s.Len)
}

//计算内存大小
func (s *SDS) Capsds() int {
	s.free = uint32(cap(s.buf))
	return int(s.free)
}

//判断 SDS类型 --这边只是简单认为是字符串类型
func (s *SDS) Type() string {
	return sds_type
}

//数值的插入
func (s *SDS) Insert_value(c byte) {
	//在插入值的同时更新长度和容量
	s.buf = append(s.buf, c)
	s.Len = uint32(s.Len_string())
	s.free = uint32(s.Capsds())
}

//清空slice 即清空SDS
//完全清空
func (s *SDS) Freeallsds() {
	s.buf = nil
}

//用于显示 sds.buf
func (s *SDS) Bufvalue() string {
	if s.Len == 0 {
		return " "
	}
	return string(s.buf[:])
}

// 未使用的空间字符数
func (s *SDS) Sdsavail() uint32 {
	return s.free - s.Len
}

//创建一个SDS的副本
func (s *SDS) Sdsdup() []byte {
	sds_copy := make([]byte, s.Len)
	copy(sds_copy, s.buf)
	return sds_copy
}

// 字符串的拼接
func (s *SDS) Sdscat(c string) {
	for i := range c {
		s.buf = append(s.buf, c[i])
	}
}

// 将给定SDS字符串拼接到另一个SDS字符串的末尾
func (s *SDS) Sdscatsds(c []byte) {
	s.buf = append(s.buf, c...)
}

// 将给定字符串复制到SDS中并覆盖原来值
func (s *SDS) Sdscpy(c string) {
	//首先清空再重新添加
	s.Freeallsds()
	for i := range c {
		s.buf = append(s.buf, c[i])
	}
}

// 扩展SDS中buf的大小，其中在使用append时已经包含了扩容功能
func (s *SDS) Sdsgrowzero(N byte) {
	n := int(N)
	buflen := s.Len + 1
	buf1 := make([]byte, buflen, s.free+uint32(n))
	copy(buf1, s.buf)
	s.buf = buf1
	s.Len = uint32(len(s.buf))
	s.free = uint32(cap(s.buf))

}

// 接受一个SDS和一个字符串作为参数，从SDS中移除所有在C字符串中出现过的字符
func (s *SDS) Sdstrim(c string) {
	mapvalue := make(map[byte]bool, s.Len)
	for i := 0; i < len(c); i++ {
		mapvalue[c[i]] = true
	}
	if s.Len == 1 {
		if mapvalue[s.buf[0]] {
			s.Freeallsds()
		}
	}
	//这里可以在原来的buf上进行修改 但是遇到些问题
	buf1 := make([]byte, 0)
	for i := 1; i < len(s.buf); i++ {
		if !mapvalue[s.buf[i]] {
			a := s.buf[i]
			buf1 = append(buf1, a)
		}
	}
	s.buf = nil
	s.buf = buf1
}

// 对比SDS字符串
func (s *SDS) Sdscmp(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
