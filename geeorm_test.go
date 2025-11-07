package geeorm

import (
    "errors"
    "geeorm/log"
    "geeorm/session"
    "testing"
)

func OpenDB(t *testing.T) *Engine {
    t.Helper()
    e, err := NewEngine("sqlite3", "gee.db")
    if err != nil {
        t.Fatal(err)
    }
    return e
}

func TestEngine(t *testing.T) {
    e := OpenDB(t)
    defer e.Close()
}

type User struct {
    Name string `geeorm:"PRIMARY KEY"`
    Age  int
}

func testTransactionRollback(t *testing.T) {
    e := OpenDB(t)
    defer e.Close()
    s := e.NewSession()
    _ = s.Model(&User{}).DropTable()
    _, err := e.Transaction(func(s *session.Session) (any, error) {
        _ = s.Model(&User{}).CreateTable()
        _, _ = s.Insert(&User{"Tom", 12})
        log.Error("假设这里遇到错误了")
        return nil, errors.New("error")
    })
    if s.HasTable() {
        t.Fatal("table should not exist")
    }
    if err == nil {
        t.Fatal("failed to rollback", s.HasTable())
    }
}

func testTransactionCommit(t *testing.T) {
    e := OpenDB(t)
    defer e.Close()
    s := e.NewSession()
    _ = s.Model(&User{}).DropTable()
    _, err := e.Transaction(func(s *session.Session) (any, error) {
        s.Model(&User{}).CreateTable()
        s.Insert(&User{"Tom", 12})
        return nil, nil
    })
    if err != nil || !s.HasTable() {
        t.Fatal("failed to commit")
    }
    var user User
    _ = s.First(&user)

    if user.Name != "Tom" || user.Age != 12 {
        t.Fatal("failed to commit")
    }
}

func TestEngine_Transaction(t *testing.T) {
    t.Run("Rollback", func(t *testing.T) { testTransactionRollback(t) })
    t.Run("Commit", func(t *testing.T) { testTransactionCommit(t) })
}
