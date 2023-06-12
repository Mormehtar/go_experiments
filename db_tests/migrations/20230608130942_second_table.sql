-- +goose Up
-- +goose StatementBegin
CREATE TABLE "properties" (
    "id" bigserial NOT NULL PRIMARY KEY,
    "name_id" bigint REFERENCES "names" ("id") NOT NULL,
    "key" varchar(255) NOT NULL,
    "value" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL
);;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "properties";
-- +goose StatementEnd
