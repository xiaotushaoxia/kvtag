package kvtag

import (
	"go/ast"
	"reflect"
	"strings"
	"sync"
)

const embeddedFlag1 = "embedded"

// gopkg.in/yaml use inline as flag
const embeddedFlag2 = "inline"

var tagCache sync.Map
var onceMap sync.Map

type TagInfo struct {
	Type      reflect.Type
	fieldTags []*FieldTag
	nmap      map[string]*FieldTag
}

func (ti *TagInfo) GetFieldTagByName(name string) (*FieldTag, bool) {
	v, ok := ti.nmap[name]
	return v, ok
}

func (ti *TagInfo) FieldTags() []*FieldTag {
	return ti.fieldTags
}

type FieldTag struct {
	FieldName  string
	TagSetting map[string]string
}

func ParserTag(obj any, tagName string, sep string) *TagInfo {
	t := reflect.Indirect(reflect.ValueOf(obj)).Type()

	store, _ := onceMap.LoadOrStore(t, &sync.Once{})
	store.(*sync.Once).Do(func() {
		tagCache.Store(t, parserTag(t, tagName, sep))
	})

	value, ok := tagCache.Load(t)
	if !ok { // should not happen
		return nil
	}
	return value.(*TagInfo)
}

func parserTag(t reflect.Type, tagName string, sep string) *TagInfo {
	v := &TagInfo{Type: t, fieldTags: make([]*FieldTag, 0), nmap: map[string]*FieldTag{}}
	_parserTag(t, tagName, sep, v)
	return v
}

func _parserTag(t reflect.Type, tagName, sep string, result *TagInfo) {
	var embeddedFields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		fieldStruct := t.Field(i)
		s := fieldStruct.Tag.Get(tagName)
		if s == "" {
			continue
		}
		if s == embeddedFlag1 || s == embeddedFlag2 {
			embeddedFields = append(embeddedFields, fieldStruct)
			continue
		}
		if !ast.IsExported(fieldStruct.Name) {
			continue
		}
		_, existed := result.nmap[fieldStruct.Name]
		if existed {
			continue
		}
		ft := &FieldTag{
			FieldName:  fieldStruct.Name,
			TagSetting: parseTagSetting(s, sep),
		}
		result.fieldTags = append(result.fieldTags, ft)
		result.nmap[fieldStruct.Name] = ft
	}
	for _, s := range embeddedFields {
		_parserTag(s.Type, tagName, sep, result)
	}
}

func parseTagSetting(str string, sep string) map[string]string {
	// gorm抄来的解析tag的部分
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(values[0])

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
}
