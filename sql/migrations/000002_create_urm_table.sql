-- +migrate Up
CREATE TABLE "userresourcemapping" ("userid" integer, "usertype" text, "mappingtype" text, "resourcetype" text, "resourceid" integer);
CREATE UNIQUE INDEX urm_ids ON userresourcemapping (userid, resourceid)
CREATE UNIQUE INDEX urm_userid ON userresourcemapping (userid)
CREATE UNIQUE INDEX urm_resourceid ON userresourcemapping (resourceid)

-- +migrate Down
DROP TABLE if exists "userresourcemapping";