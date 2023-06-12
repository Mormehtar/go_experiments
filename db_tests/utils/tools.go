package utils

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func CreteDB(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, name))
	return err
}

func MakeTemplate(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf(`ALTER DATABASE "%s" WITH ALLOW_CONNECTIONS FALSE IS_TEMPLATE TRUE;`, name))
	return err
}

func DropTemplate(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf(`ALTER DATABASE "%s" WITH ALLOW_CONNECTIONS TRUE IS_TEMPLATE FALSE;`, name))
	if err != nil {
		return err
	}
	return DropDB(db, name)
}

func DropDB(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf(`DROP DATABASE "%s";`, name))
	return err
}

func CloneDB(db *sql.DB, dbUrl, templateDBName string) (*sql.DB, string, error) {
	newName := fmt.Sprintf("%s_%d", templateDBName, uuid.New().ID())
	_, err := db.Exec(fmt.Sprintf(`CREATE DATABASE "%s" TEMPLATE "%s";`, newName, templateDBName))
	if err != nil {
		return nil, "", err
	}
	url := ChangeDBName(dbUrl, newName)
	testDB, err := sql.Open("pgx", url)
	if err != nil {
		return nil, "", err
	}
	return testDB, newName, nil
}

func DropClone(db, cloneDB *sql.DB, cloneName string) error {
	err := cloneDB.Close()
	if err != nil {
		return err
	}
	return DropDB(db, cloneName)
}

func ConfigureGoose() {
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	// Silence goose logs
	goose.SetLogger(goose.NopLogger())
}

func ChangeDBName(dsn string, newName string) string {
	if dsn == "" {
		return fmt.Sprintf("database=%s", newName)
	}
	index := strings.Index(dsn, "database=")
	if index == -1 {
		return fmt.Sprintf("%s database=%s", dsn, newName)
	} else {
		lastIndex := strings.Index(dsn[index:], " ")
		if lastIndex == -1 {
			lastIndex = len(dsn) - index
		}
		return fmt.Sprintf("%sdatabase=%s%s", dsn[:index], newName, dsn[index+lastIndex:])
	}
}
