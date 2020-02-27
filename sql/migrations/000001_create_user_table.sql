-- +migrate Up
CREATE TABLE "user" ("id" integer unique, "name" text, "oauthid" text, "status" text, "password" text);

-- +migrate Down
DROP TABLE if exists "user";