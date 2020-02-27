-- +migrate Up
CREATE TABLE "organization" ("id" integer unique, "name" text, "description" text, "createdat" datetime, "updatedat" datetime);

-- +migrate Down
DROP TABLE if exists "organization";