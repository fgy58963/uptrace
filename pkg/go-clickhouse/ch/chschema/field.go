package chschema

import (
	"fmt"
	"reflect"

	"github.com/uptrace/go-clickhouse/ch/internal/tagparser"
)

const (
	customTypeFlag = uint8(1) << iota
)

type Field struct {
	Field reflect.StructField
	Tag   tagparser.Tag
	Type  reflect.Type
	Index []int

	GoName    string // struct field name, e.g. Id
	CHName    string // SQL name, .e.g. id
	Column    Safe   // escaped SQL name, e.g. "id"
	CHType    string
	CHDefault Safe

	NewColumn   NewColumnFunc
	appendValue AppenderFunc

	IsPK    bool
	NotNull bool

	flags uint8
}

func (f *Field) String() string {
	return "field=" + f.GoName
}

func (f *Field) Value(strct reflect.Value) reflect.Value {
	return fieldByIndexAlloc(strct, f.Index)
}

func (f *Field) AppendValue(fmter Formatter, b []byte, strct reflect.Value) []byte {
	fv, ok := fieldByIndex(strct, f.Index)
	if !ok {
		return AppendNull(b)
	}

	if f.appendValue == nil {
		return AppendError(b, fmt.Errorf("ch: AppendValue(unsupported %s)", fv.Type()))
	}
	return f.appendValue(fmter, b, fv)
}

func (f *Field) setFlag(flag uint8) {
	f.flags |= flag
}

func (f *Field) hasFlag(flag uint8) bool {
	return f.flags&flag != 0
}
