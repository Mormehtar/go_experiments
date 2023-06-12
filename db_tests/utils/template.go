package utils

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

type Template struct {
	baseDBURL      string
	baseDB         *sql.DB
	templateDBName string
}

type TestDB struct {
	db     *sql.DB
	dbName string
}

func (db *TestDB) GetDb() *sql.DB {
	return db.db
}

func NewTemplate(baseDBURL string) *Template {
	return &Template{baseDBURL: baseDBURL}
}

func (template *Template) Init(migrationsPath string) error {
	var err error

	template.baseDB, err = sql.Open("pgx", template.baseDBURL)
	if err != nil {
		return err
	}

	template.templateDBName = fmt.Sprintf("test_db_%d", uuid.New().ID())
	if CreteDB(template.baseDB, template.templateDBName) != nil {
		return err
	}
	templateDBURL := ChangeDBName(template.baseDBURL, template.templateDBName)

	templateDB, err := sql.Open("pgx", templateDBURL)
	if err != nil {
		return err
	}

	ConfigureGoose()

	if err := goose.Up(templateDB, migrationsPath); err != nil {
		return err
	}
	if err := templateDB.Close(); err != nil {
		return err
	}

	return MakeTemplate(template.baseDB, template.templateDBName)
}

func (template *Template) Close() error {
	if err := DropTemplate(template.baseDB, template.templateDBName); err != nil {
		return err
	}
	return template.baseDB.Close()
}

func (template *Template) GetTestDB() (*TestDB, error) {
	testDB, testDBName, err := CloneDB(template.baseDB, template.baseDBURL, template.templateDBName)
	if err != nil {
		return nil, err
	}
	return &TestDB{db: testDB, dbName: testDBName}, nil
}

func (template *Template) DropTestDB(db *TestDB) error {
	return DropClone(template.baseDB, db.db, db.dbName)
}
