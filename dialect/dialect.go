package dialect

import (
    "reflect"
)

type Dialect interface {
    DataTypeOf(typ reflect.Value) string
    TableExistSQL(tableName string) (string, []any)
}

var dialectsMap = map[string]Dialect{}

func RegisterDialect(name string, dialect Dialect) {
    dialectsMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
    d, ok := dialectsMap[name]
    return d, ok
}
