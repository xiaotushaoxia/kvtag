## kvtag  go结构体tag解析工具

***用于解析形如`filter:"colum:type;opr:in;split:in\\;t;sep:,;test:xxsds;name:xxx"`的tag***

```go
package kvtag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type objTest struct {
	eEObj1 `filter:"inline"`
	Name   string `filter:"colum:name;opr:like;pattern:%%?%%"`
	ID     []int  `filter:"colum:id;opr:in"`
	Type   string `filter:"colum:type;opr:in;split:in\\;t;sep:,"`
	Start  string `filter:"colum:created_at"`
	Desc   string `filter:"colum:desc;opr:in"`
	KK     string
	EE1    string `filter:"colum:ee1;opr:=xx"`
}

type eEObj1 struct {
	eEObj2 `filter:"inline"`
	EE1    string `filter:"colum:ee1;opr:="`
	EE2    string
}

type eEObj2 struct {
	EE1 string `filter:"colum:ee1;opr:in"`
	EE2 string `filter:"colum:ee2;opr:in"`
}

func TestUsage(t *testing.T) {
	tagInfo := ParserTag(&objTest{}, "filter", ";")
	//tagInfo := ParserTag(objTest{}, "filter", ";")  和上一句一样的
	fmt.Println(tagInfo.GetFieldTagByName("ID")) // &{1 ID map[colum:id opr:in]} true
	//fmt.Println(tagInfo.GetFieldTagByIndex(0))   //&{0 Name map[colum:name opr:like pattern:%%?%%]} true

	bb := &objTest{
		eEObj1: eEObj1{
			eEObj2: eEObj2{
				EE1: "ee1",
				EE2: "ee2",
			},
			EE1: "e1",
			EE2: "e2",
		},
		EE1:   "xxx",
		Name:  "1232",
		ID:    []int{1, 2},
		Type:  "223232,55",
		Start: "xzxzx",
		Desc:  "c.ss.s",
		KK:    "ccc",
	}

	v := reflect.Indirect(reflect.ValueOf(bb))

	for i, tag := range tagInfo.FieldTags() {
		fv := v.FieldByName(tag.FieldName)
		fmt.Println(i, fv, tag.FieldName, tag.TagSetting)
	}

	// 输出
	//&{ID map[colum:id opr:in]} true
	//0 1232 Name map[colum:name opr:like pattern:%%?%%]
	//1 [1 2] ID map[colum:id opr:in]
	//2 223232,55 Type map[colum:type opr:in sep:, split:in;t]
	//3 xzxzx Start map[colum:created_at]
	//4 c.ss.s Desc map[colum:desc opr:in]
	//5 xxx EE1 map[colum:ee1 opr:=xx]
	//6 e2 EE2 map[colum:ee2 opr:in]
}

```