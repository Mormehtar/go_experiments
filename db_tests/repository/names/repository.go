package names

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Mormehtar/go_experiments/db_tests/repository/interfaces"
)

type Names struct {
	db *sql.DB
}

var _ interfaces.IName = (*Names)(nil)

func New(db *sql.DB) *Names {
	return &Names{db: db}
}

func (n *Names) Create(name string) (*interfaces.Name, error) {
	query := `
        INSERT INTO "names" ("name", "created_at", "updated_at")
        VALUES ($1, now(), now())
        RETURNING "id", "name", "created_at", "updated_at";
    `

	rows, err := n.db.Query(query, name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Message == "duplicate key value violates unique constraint \"names_name_key\"" {
			return nil, interfaces.NameIsUsedAlready
		}
		return nil, err
	}

	defer rows.Close()

	row := &interfaces.Name{}

	rows.Next()
	err = rows.Scan(&row.Id, &row.Name, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (n *Names) Update(id int64, name string) (*interfaces.Name, error) {
	query := `
        UPDATE "names"
        SET "name" = $1, "updated_at" = now()
        RETURNING "id", "name", "created_at", "updated_at";
    `

	rows, err := n.db.Query(query, name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Message == "duplicate key value violates unique constraint \"names_name_key\"" {
			return nil, interfaces.NameIsUsedAlready
		}
		return nil, err
	}

	defer rows.Close()

	row := &interfaces.Name{}

	rows.Next()
	err = rows.Scan(&row.Id, &row.Name, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (n *Names) Delete(id int64) error {
	query := `DELETE FROM "names" WHERE "id" = $1; `

	result, err := n.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return interfaces.NameIsNotFound
	}

	return nil
}

func (n *Names) Get(name string) (*interfaces.Name, error) {
	query := `SELECT "id", "name", "created_at", "updated_at" FROM "names" WHERE "name" = $1;`

	rows, err := n.db.Query(query, name)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	row := &interfaces.Name{}

	hasData := rows.Next()

	if !hasData {
		return nil, interfaces.NameIsNotFound
	}

	err = rows.Scan(&row.Id, &row.Name, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return row, nil
}
