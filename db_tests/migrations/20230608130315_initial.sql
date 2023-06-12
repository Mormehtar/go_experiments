-- +goose Up
-- +goose StatementBegin
CREATE TABLE "names" (
    "id" bigserial NOT NULL PRIMARY KEY,
    "name" varchar(255) UNIQUE NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "names";
-- +goose StatementEnd
