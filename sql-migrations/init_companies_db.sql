-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS companies (
                                      name varchar PRIMARY KEY,
                                     code varchar NOT NULL,
                                      country varchar NOT NULL,
                                      website varchar NOT NULL,
                                      phone varchar NOT NULL
);
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE IF EXISTS companies;
-- +migrate StatementEnd