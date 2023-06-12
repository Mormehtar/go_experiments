package properties

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Mormehtar/go_experiments/db_tests/repository/interfaces"
)

type Properties struct {
	db *sql.DB
}

var _ interfaces.IProperty = (*Properties)(nil)

func New(db *sql.DB) *Properties {
	return &Properties{db: db}
}

func (p *Properties) Create(nameId int64, key, value string) (*interfaces.Property, error) {
	query := `
        INSERT INTO "properties" ("name_id", "key", "value", "created_at", "updated_at")
        VALUES ($1, $2, $3, statement_timestamp(), statement_timestamp())
        RETURNING "id", "name_id", "key", "value", "created_at", "updated_at";
    `

	rows, err := p.db.Query(query, nameId, key, value)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "properties_name_id_fkey" {
			return nil, interfaces.NameIsNotFound
		}
		return nil, err
	}

	defer rows.Close()

	row := &interfaces.Property{}

	rows.Next()
	err = rows.Scan(&row.Id, &row.NameId, &row.Key, &row.Value, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (p *Properties) Get(nameId int64, key string) ([]*interfaces.Property, error) {
	query := `
        SELECT "id", "name_id", "key", "value", "created_at", "updated_at"
        FROM "properties"
        WHERE "name_id" = $1 AND "key" = $2;
    `

	rows, err := p.db.Query(query, nameId, key)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*interfaces.Property, 0)

	for rows.Next() {
		row := &interfaces.Property{}
		err = rows.Scan(&row.Id, &row.NameId, &row.Key, &row.Value, &row.CreatedAt, &row.UpdatedAt)

		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}
	return result, nil
}

func (p *Properties) Update(id int64, key, value string) (*interfaces.Property, error) {
	query := `
        UPDATE "properties"
        SET "key" = $2, "value" = $3, "updated_at" = statement_timestamp()
        WHERE "id" = $1
        RETURNING "id", "name_id", "key", "value", "created_at", "updated_at";
    `

	rows, err := p.db.Query(query, id, key, value)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	row := &interfaces.Property{}

	if !rows.Next() {
		return nil, interfaces.PropertyIsNotFound
	}
	err = rows.Scan(&row.Id, &row.NameId, &row.Key, &row.Value, &row.CreatedAt, &row.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (p *Properties) Delete(id int64) error {
	query := `DELETE FROM "properties" WHERE "id" = $1;`

	result, err := p.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return interfaces.PropertyIsNotFound
	}

	return nil
}
