package geeorm

import (
    "database/sql"
    "geeorm/dialect"
    "geeorm/log"
    "geeorm/session"

    _ "github.com/mattn/go-sqlite3"
)

type Engine struct {
    db      *sql.DB
    dialect dialect.Dialect
}

func NewEngine(driver, source string) (*Engine, error) {
    db, err := sql.Open(driver, source)
    if err != nil {
        log.Error("failed to open database")
        return nil, err
    }
    if err = db.Ping(); err != nil {
        log.Error("failed to ping database")
        return nil, err
    }
    dial, ok := dialect.GetDialect(driver)
    if !ok {
        log.Error("dialect %s is not supported", driver)
        return nil, err
    }
    e := &Engine{db: db, dialect: dial}
    log.Info("connect database success")
    return e, nil
}

func (e *Engine) Close() {
    err := e.db.Close()
    if err != nil {
        log.Error("failed to close database")
        return
    }
    log.Info("close database success")
    return
}

func (e *Engine) NewSession() *session.Session {
    return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (any, error)

func (e *Engine) Transaction(f TxFunc) (result any, err error) {
    s := e.NewSession()
    if err := s.Begin(); err != nil {
        return nil, err
    }
    defer func() {
        if p := recover(); p != nil {
            _ = s.Rollback()
            panic(p)
        } else if err != nil {
            _ = s.Rollback()
        } else {
            err = s.Commit()
        }
    }()
    return f(s)
}
