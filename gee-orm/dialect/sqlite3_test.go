package dialect

import (
	"reflect"
	"testing"
)

func TestSqlite3_DataTypeOf(t *testing.T) {
	d, ok := GetDialect("sqlite3")
	if !ok {
		t.Fatal("cannot get dialect sqlite3")
	}
	cases := []struct {
		Value any
		Type  string
	}{
		{"Tom", "text"},
		{123, "integer"},
		{1.2, "real"},
		{[]int{1, 2, 3}, "blob"},
	}
	for _, c := range cases {
		if typ := d.DataTypeOf(reflect.ValueOf(c.Value)); typ != c.Type {
			t.Fatalf("expected %s, got %s", c.Value, typ)
		}
	}
}
