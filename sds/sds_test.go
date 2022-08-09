package sds

import (
	_ "strings"
	"testing"
)

func TestSDS_value(t *testing.T) {
	sds_test := &SDS{}
	Newsds()
	value := sds_test.Len_string()
	if value == 0 {
		t.Error(`sds len is ture`)
	}

}

func TestSDS_Insert(t *testing.T) {
	sds_test := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds_test.Insert_value(test_string[i])
	}
	len_value := sds_test.Len_string()
	if len_value == 10 {
		t.Error(`Insert is ok`)
	}

}

func TestType(t *testing.T) {
	sds_test := &SDS{}
	Newsds()
	if sds_test.Type() == "string" {
		t.Error(`return Type is ok`)
	}
}

func TestFreeall(t *testing.T) {
	sds := &SDS{}
	Newsds()
	//首先往里面存数据
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	// 第一种完全释放 即[]、0、0
	sds.Freeallsds()
	if sds.Len_string() == 0 && cap(sds.buf) == 0 {
		t.Error(`slice clean all`)
	}
}

func TestAvail(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	value := sds.Sdsavail()
	if value != 6 {
		t.Error(`return avail error`)
	}

}

func TestCopy(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	sds_copy := sds.Sdsdup()
	if string(sds_copy) != "helloredis" {
		t.Error(`return copy failed`)
	}
}

func TestCat(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	new_string := "hello golang"

	sds.Sdscat(new_string)

	if string(sds.Bufvalue()) != "helloredishello golang" {
		t.Error(`return failed`)
	}
}

func TestCatsds(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	sds1 := &SDS{}
	Newsds()
	new_string := "hello golang"
	for i := range new_string {
		sds1.Insert_value(new_string[i])
	}
	sds.Sdscatsds([]byte(sds1.Bufvalue()))

	if string(sds.Bufvalue()) != "helloredishello golang" {
		t.Error(`failed`)
	}

}

func TestSdscpy(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	new_string := "hello golang"

	sds.Sdscpy(new_string)
	if string(sds.Bufvalue()) != "hello golang" {
		t.Error(`failed`)
	}

}

func TestSdsgrowzero(t *testing.T) {
	sds := &SDS{}
	Newsds()
	n := 100
	sds.Sdsgrowzero(byte(n))
	if sds.free != 100 {
		t.Error(`failed`)
	}

}

func TestSdstrim(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	new_string := "hello"
	sds.Sdstrim(new_string)
	if string(sds.Bufvalue()) != "rdis" {
		t.Error(`failed`)
	}
}

func TestSdscmp(t *testing.T) {
	sds := &SDS{}
	Newsds()
	test_string := "helloredis"
	for i := range test_string {
		sds.Insert_value(test_string[i])
	}
	new_string := "hello golang"

	if !sds.Sdscmp(sds.buf, []byte(new_string)) {
		t.Error(`return ok`)
	}
}
