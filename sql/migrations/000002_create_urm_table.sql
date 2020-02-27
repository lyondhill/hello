-- +migrate Up
CREATE TABLE "userresourcemapping" ("userid" integer, "usertype" text, "mappingtype" text, "resourcetype" text, "resourceid" integer);

-- +migrate Down
DROP TABLE if exists "userresourcemapping";