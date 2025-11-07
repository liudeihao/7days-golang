package session

import (
	"database/sql"
	"geeorm/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB  *sql.DB
	Dialect dialect.Dialect
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "../gee.db")
	Dialect, _ = dialect.GetDialect("sqlite3")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, Dialect)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?);", "Tom", "Pete").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expected 2 rows inserted, got", count, err)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	err := row.Scan(&count)
	if err != nil {
		t.Fatal("failed to query db", err)
	}
	if count != 0 {
		t.Fatal("expected 0 rows, got", count)
	}
}
