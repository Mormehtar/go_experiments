package utils

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type CleanDB struct {
	dbUrl          string
	testDBName     string
	baseDB, testDB *sql.DB
}

func NewCleanDB(dbUrl string) *CleanDB {
	return &CleanDB{dbUrl: dbUrl}
}

func (cleanDB *CleanDB) Init() error {
	var err error
	cleanDB.baseDB, err = sql.Open("pgx", cleanDB.dbUrl)

	if err != nil {
		return err
	}

	cleanDB.testDBName = fmt.Sprintf("test_db_%d", uuid.New().ID())
	if CreteDB(cleanDB.baseDB, cleanDB.testDBName) != nil {
		panic(err)
	}
	url := ChangeDBName(cleanDB.dbUrl, cleanDB.testDBName)
	cleanDB.testDB, err = sql.Open("pgx", url)
	return err
}

func (cleanDB *CleanDB) GetTestDB() *sql.DB {
	return cleanDB.testDB
}

func (cleanDB *CleanDB) Close() error {
	if err := cleanDB.testDB.Close(); err != nil {
		return err
	}
	if err := DropDB(cleanDB.baseDB, cleanDB.testDBName); err != nil {
		return err
	}
	return cleanDB.baseDB.Close()
}
