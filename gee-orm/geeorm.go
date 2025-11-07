package geeorm

import (
    "database/sql"
    "fmt"
    "geeorm/dialect"
    "geeorm/log"
    "geeorm/session"
    "strings"

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

func difference(a, b []string) []string {
    var diff []string
    mapB := make(map[string]bool)
    for _, v := range b {
        mapB[v] = true
    }
    for _, v := range a {
        if _, ok := mapB[v]; !ok {
            diff = append(diff, v)
        }
    }
    return diff
}

func (e *Engine) Migrate(value any) error {
    _, err := e.Transaction(func(s *session.Session) (any, error) {
        if !s.Model(value).HasTable() {
            log.Infof("table %s does not exist, creating...", s.RefTable().Name)
            return nil, s.CreateTable()
        }
        table := s.RefTable()
        rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1;", table.Name)).QueryRows()
        cols, _ := rows.Columns()
        addCols := difference(table.FieldNames, cols)
        delCols := difference(cols, table.FieldNames)
        log.Infof("add cols %v, deleted cols %v", addCols, delCols)

        for _, col := range addCols {
            f := table.GetField(col)
            sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
            if _, err := s.Raw(sqlStr).Exec(); err != nil {
                return nil, err
            }
        }
        if len(delCols) == 0 {
            return nil, nil
        }
        tmp := "tmp_" + table.Name
        fieldStr := strings.Join(table.FieldNames, ", ")
        s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fieldStr, table.Name))
        s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
        s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tmp, table.Name))
        _, err := s.Exec()
        return nil, err
    })
    return err
}
