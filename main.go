package main

// 用于一些测试操作

import (
	"fmt"
	"redis/sds"
)

func main() {
	sds_test := &sds.SDS{}
	sds.Newsds()
	//首先往里面存数据
	test_string := "helloredis"
	for i := range test_string {
		sds_test.Insert_value(test_string[i])
	}

	fmt.Print("长度为：", sds_test.Len_string())
	fmt.Print("\n长度为：", sds_test.Len)
	fmt.Println()
	fmt.Print("容量为：", sds_test.Capsds())
	fmt.Println()
	fmt.Print("内容为：", sds_test.Bufvalue())
	fmt.Println()
	new_string := "helloredis"
	fmt.Println(sds_test.Sdscmp([]byte(sds_test.Bufvalue()), []byte(new_string)))
}
