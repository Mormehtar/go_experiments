package utils

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Transactional struct {
	baseDB     *sql.DB
	testDB     *sql.DB
	baseDBUrl  string
	testDBUrl  string
	testDBName string
}

func NewTransactional(url string) *Transactional {
	transactional := &Transactional{baseDBUrl: url}
	return transactional
}

func (transactional *Transactional) Init(migrationsPath string) error {
	var err error

	transactional.baseDB, err = sql.Open("pgx", transactional.baseDBUrl)
	if err != nil {
		return err
	}

	transactional.testDBName = fmt.Sprintf("test_db_%d", uuid.New().ID())
	if CreteDB(transactional.baseDB, transactional.testDBName) != nil {
		return err
	}
	transactional.testDBUrl = ChangeDBName(transactional.baseDBUrl, transactional.testDBName)

	transactional.testDB, err = sql.Open("pgx", transactional.testDBUrl)
	if err != nil {
		return err
	}

	ConfigureGoose()
	return goose.Up(transactional.testDB, migrationsPath)
}

func (transactional *Transactional) Close() error {
	if err := transactional.testDB.Close(); err != nil {
		return err
	}
	if err := DropDB(transactional.baseDB, transactional.testDBName); err != nil {
		return err
	}
	return transactional.baseDB.Close()
}

func (transactional *Transactional) GetTestDB() (*sql.DB, error) {
	_, err := transactional.testDB.Exec("START TRANSACTION")
	if err != nil {
		return nil, err
	}
	return transactional.testDB, nil
}

func (transactional *Transactional) ReturnTestDB() error {
	_, err := transactional.testDB.Exec("ROLLBACK")
	return err
}
