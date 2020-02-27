package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rubenv/sql-migrate"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
	"github.com/Masterminds/squirrel"

	"github.com/lyondhill/hello/sql/migrations"

)

type ID uint64
type Status string
type User struct {
	ID       ID     `json:"id,omitempty"`
	Name     string `json:"name"`
	OAuthID  string `json:"oauthID,omitempty"`
	Status   Status `json:"status"`
	Password string
}

type UserType string
type MappingType uint8
type ResourceType string

type UserResourceMapping struct {
	UserID       ID           `json:"userID"`
	UserType     UserType     `json:"userType"`
	MappingType  MappingType  `json:"mappingType"`
	ResourceType ResourceType `json:"resourceType"`
	ResourceID   ID           `json:"resourceID"`
}

type CRUDLog struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Organization struct {
	ID          ID     `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CRUDLog
}

type BucketType int

type Bucket struct {
	ID                  ID            `json:"id,omitempty"`
	OrgID               ID            `json:"orgID,omitempty"`
	Type                BucketType    `json:"type"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	RetentionPolicyName string        `json:"rp,omitempty"` // This to support v1 sources
	RetentionPeriod     time.Duration `json:"retentionPeriod"`
	CRUDLog
}

func main() {
	s, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	
	theMigration(s)

	rows, err := s.Query("SELECT * FROM sqlite_master")
	printTable(rows)
	
	theModlWay(s)
	theSqlxWay(s)
	theSquirlWay(s)	
}

func printTable(rows *sql.Rows) {
	cols, _ := rows.Columns()
	n := len(cols)

	for i := 0; i < n; i++ {
		fmt.Print(cols[i], "\t")
	}
	fmt.Println()

	var fields []interface{}
	for i := 0; i < n; i++ {
		fields = append(fields, new(string))
	}
	for rows.Next() {
		rows.Scan(fields...)
		for i := 0; i < n; i++ {
			fmt.Print(*(fields[i].(*string)), "\t")
		}
		fmt.Println()
	}
}

func theMigration(s *sql.DB) {
	mgs := &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      "",
	}

	n, err := migrate.Exec(s, "sqlite3", mgs, migrate.Up)
	if err != nil {
    	panic(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func theModlWay(s *sql.DB) {
	dbMap := modl.NewDbMap(s, modl.SqliteDialect{})
	dbMap.AddTable(User{})
	dbMap.AddTable(UserResourceMapping{})
	dbMap.AddTable(Organization{})
	dbMap.AddTable(Bucket{})

	err := dbMap.Insert(&User{Name: "modlDude1", ID: 1, OAuthID: "data", Status: "active", Password: "1234"})
	if err != nil {
		panic(err)
	}	
	err = dbMap.Insert(&User{Name: "modlDude2", ID: 2, OAuthID: "data", Status: "active", Password: "1234"})
	if err != nil {
		panic(err)
	}

	users := []User{}
	err = dbMap.Select(&users, "select * from user")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nmodlUsers:\n %+v\n\n", users)
}

func theSqlxWay(s *sql.DB) {
	sx := sqlx.NewDb(s, "sqlite3")
	sx.MustExec("INSERT INTO user (name, id, oauthid, status, password) VALUES (\"sqlxdude1\", 3, \"data\", \"active\", \"1234\")")
	sx.MustExec("INSERT INTO user (name, id, oauthid, status, password) VALUES (\"sqlxdude2\", 4, \"data\", \"active\", \"1234\")")

	users := []User{}
	err := sx.Select(&users, "select * from user")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nsqlxUsers: \n%+v\n\n", users)
}

func theSquirlWay(s *sql.DB) {
	insert := squirrel.Insert("user").Columns("name", "id", "oauthid", "status", "password").
	Values("sqrldude1", 5, "data", "active", "1234").
	Values("sqrldude2", 6, "data", "active", "1234")

	_, err := insert.RunWith(s).Exec()
	if err != nil {
		panic(err)
	}
	// use sqlx
	sx := sqlx.NewDb(s, "sqlite3")

	q := squirrel.Select("*").From("user")
	ugly, err := q.RunWith(sx).Query()
	if err != nil {
		panic(err)
	}
	printTable(ugly)

	newQ, _, _ := q.ToSql()

	users := []User{}
	err = sx.Select(&users, newQ)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nsqrlxUsers: \n%+v\n\n", users)
}