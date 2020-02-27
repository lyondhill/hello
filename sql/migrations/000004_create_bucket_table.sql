-- +migrate Up
CREATE TABLE "bucket" ("id" integer unique, "orgid" integer, "type" integer, "name" text, "description" text, "retentionpolicyname" text, "retentionperiod" integer, "createdat" datetime, "updatedat" datetime);

-- +migrate Down
DROP TABLE if exists "bucket";