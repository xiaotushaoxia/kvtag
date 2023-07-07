package kvtag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type objTest struct {
	Name  string `filter:"colum:name;opr:like;pattern:%%?%%"`
	ID    []int  `filter:"colum:id;opr:in"`
	Type  string `filter:"colum:type;opr:in;split:in\\;t;sep:,"`
	Start string `filter:"colum:created_at"`
	Desc  string `filter:"colum:desc;opr:in"`
	KK    string
}

func TestUsage(t *testing.T) {
	tagInfo := ParserTag(&objTest{}, "filter", ";")
	//tagInfo := ParserTag(objTest{}, "filter", ";")  和上一句一样的
	fmt.Println(tagInfo.GetFieldTagByName("ID")) // &{1 ID map[colum:id opr:in]} true
	fmt.Println(tagInfo.GetFieldTagByIndex(0))   //&{0 Name map[colum:name opr:like pattern:%%?%%]} true

	bb := &objTest{
		Name:  "1232",
		ID:    []int{1, 2},
		Type:  "223232,55",
		Start: "xzxzx",
		Desc:  "c.ss.s",
		KK:    "ccc",
	}

	ti := ParserTag(bb, "filter", ";")
	v := reflect.Indirect(reflect.ValueOf(bb))

	for i := 0; i < v.NumField(); i++ {
		ft, hasFieldTag := ti.GetFieldTagByIndex(i)
		if !hasFieldTag {
			continue
		}
		fv := v.Field(i)
		fmt.Println(i, fv, ft.FieldName, ft.FieldIndex, ft.TagSetting)
	}
	// 输出
	//0 1232 Name 0 map[colum:name opr:like pattern:%%?%%]
	//1 [1 2] ID 1 map[colum:id opr:in]
	//2 223232,55 Type 2 map[colum:type opr:in sep:, split:in;t]
	//3 xzxzx Start 3 map[colum:created_at]
	//4 c.ss.s Desc 4 map[colum:desc opr:in]
}

func TestParserTag(t *testing.T) {
	type args struct {
		t       any
		tagName string
		sep     string
	}

	test1Want := `[{"FieldIndex":0,"FieldName":"Name","TagSetting":{"colum":"name","opr":"like","pattern":"%%?%%"}},{"FieldIndex":1,"FieldName":"ID","TagSetting":{"colum":"id","opr":"in"}},{"FieldIndex":2,"FieldName":"Type","TagSetting":{"colum":"type","opr":"in","sep":",","split":"in;t"}},{"FieldIndex":3,"FieldName":"Start","TagSetting":{"colum":"created_at"}},{"FieldIndex":4,"FieldName":"Desc","TagSetting":{"colum":"desc","opr":"in"}}]`
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", want: test1Want, args: args{t: objTest{}, tagName: "filter", sep: ";"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parserTagAndConvert2Json(tt.args.t, tt.args.tagName, tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nget :%v\nwant:%v", got, tt.want)
			}
		})
	}

	// test1Want格式化后的结果
	test1Want = `[
		{
			"FieldIndex": 0,
			"FieldName": "Name",
			"TagSetting": {
				"colum": "name",
				"opr": "like",
				"pattern": "%%?%%"
			}
		},
		{
			"FieldIndex": 1,
			"FieldName": "ID",
			"TagSetting": {
				"colum": "id",
				"opr": "in"
			}
		},
		{
			"FieldIndex": 2,
			"FieldName": "Type",
			"TagSetting": {
				"colum": "type",
				"opr": "in",
				"sep": ",",
				"split": "in;t"
			}
		},
		{
			"FieldIndex": 3,
			"FieldName": "Start",
			"TagSetting": {
				"colum": "created_at"
			}
		},
		{
			"FieldIndex": 4,
			"FieldName": "Desc",
			"TagSetting": {
				"colum": "desc",
				"opr": "in"
			}
		}
	]
`
}

func parserTagAndConvert2Json(t any, tagName string, sep string) string {
	aa := ParserTag(t, tagName, sep).FieldTags()
	marshal, err := json.Marshal(aa)
	if err != nil {
		return ""
	}
	return string(marshal)
}
