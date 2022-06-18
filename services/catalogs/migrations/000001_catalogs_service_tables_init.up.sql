-- https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md
--https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS products CASCADE;


CREATE TABLE products
(
    product_id  UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    name        VARCHAR(250)  NOT NULL CHECK ( name <> '' ),
    description VARCHAR(5000) NOT NULL CHECK ( description <> '' ),
    price       NUMERIC       NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

