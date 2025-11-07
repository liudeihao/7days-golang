package schema

import (
	"geeorm/dialect"
	"testing"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	s := Parse(&User{}, TestDial)
	if s.Name != "User" || len(s.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if s.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}

func TestSchema_RecordValues(t *testing.T) {
	s := Parse(&User{}, TestDial)
	user := &User{
		Name: "Tom",
		Age:  20,
	}
	values := s.RecordValues(user)
	name, age := values[0].(string), values[1].(int)
	if name != "Tom" || age != 20 {
		t.Fatal("failed to get values")
	}
}

type UserTest struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func (u *UserTest) TableName() string {
	return "ns_user_test"
}

func TestSchema_TableName(t *testing.T) {
	schema := Parse(&UserTest{}, TestDial)
	if schema.Name != "ns_user_test" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse UserTest struct")
	}
}
