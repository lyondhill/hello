package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"gopkg.in/gorp.v2"

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
	
	theGorpWay(s)
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

func theGorpWay(s *sql.DB) {
	dbMap := gorp.DbMap{Db: s, Dialect: gorp.SqliteDialect{}}
	dbMap.AddTable(User{})
	dbMap.AddTable(UserResourceMapping{})
	dbMap.AddTable(Organization{})
	dbMap.AddTable(Bucket{})
	
	gtx, err := dbMap.Begin()
	if err != nil {
		panic(err)
	}

	err = gtx.Insert(&User{Name: "gorpDude1", ID: 1, OAuthID: "data", Status: "active", Password: "1234"})
	if err != nil {
		panic(err)
	}
	err = gtx.Insert(&User{Name: "gorpDude2", ID: 2, OAuthID: "data", Status: "active", Password: "1234"})
	if err != nil {
		panic(err)
	}
	err = gtx.Insert(&Organization{ID: 3, Name: "otherorg", Description: "what", CRUDLog:CRUDLog{CreatedAt: time.Now()}})
	if err != nil {
		panic(err)
	}
	
	err = gtx.Commit()
	if err != nil {
		panic(err)
	}

	users := []User{}
	_, err = dbMap.Select(&users, "select * from user")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\ngorpUsers:\n %+v\n\n", users)

	orgs := []Organization{}
	_, err = dbMap.Select(&orgs, "select * from organization")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\ngorpOrgs:\n %+v\n\n", orgs)
}

func theModlWay(s *sql.DB) {
	dbMap := modl.NewDbMap(s, modl.SqliteDialect{})
	dbMap.AddTable(User{})
	dbMap.AddTable(UserResourceMapping{})
	dbMap.AddTable(Organization{})
	dbMap.AddTable(Bucket{})

	err := dbMap.Insert(&User{Name: "modlDude1", ID: 3, OAuthID: "data", Status: "active", Password: "1234"})
	if err != nil {
		panic(err)
	}	
	err = dbMap.Insert(&User{Name: "modlDude2", ID: 4, OAuthID: "data", Status: "active", Password: "1234"})
	if err != nil {
		panic(err)
	}

	// tried using the org but modl doesn't like embedded structs. it cant figure out how to encode the embedded fields (it tried setting CRUDLog as a "text" type)

	users := []User{}
	err = dbMap.Select(&users, "select * from user")
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nmodlUsers:\n %+v\n\n", users)
}

func theSqlxWay(s *sql.DB) {
	sx := sqlx.NewDb(s, "sqlite3")
	tx, err := sx.Beginx()
	if err != nil {
		panic(err)
	}
	tx.Exec("INSERT INTO user (name, id, oauthid, status, password) VALUES (\"sqlxdude1\", 5, \"data\", \"active\", \"1234\")")
	tx.Exec("INSERT INTO user (name, id, oauthid, status, password) VALUES (\"sqlxdude2\", 6, \"data\", \"active\", \"1234\")")
	tx.Exec("INSERT INTO organization (id, name, description, createdat) VALUES (0, \"org\", \"desc\", '2020-02-25T13:21:21')")
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	users := []User{}
	err = sx.Select(&users, "select * from user")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nsqlxUsers: \n%+v\n\n", users)

	org := Organization{}
	err = sx.Get(&org, "select id, name, description, createdat from organization limit 1")
	if err != nil {
		panic(err)
	}	
	fmt.Printf("\nsqlxOrgs: \n%+v\n\n", org)
}

func theSquirlWay(s *sql.DB) {
	insert := squirrel.Insert("user").Columns("name", "id", "oauthid", "status", "password").
	Values("sqrldude1", 7, "data", "active", "1234").
	Values("sqrldude2", 8, "data", "active", "1234")
	// _, err := insert.RunWith(s).Exec()
	// if err != nil {
	// 	panic(err)
	// }
	
	orgIn := Organization{ID: 1, Name: "otherorg", Description: "what", CRUDLog:CRUDLog{CreatedAt: time.Now()}}
	orgInsert := squirrel.Insert("organization").Columns("id", "name", "description", "createdat").
	Values(orgIn.ID, orgIn.Name, orgIn.Description, orgIn.CreatedAt)

	i1, args1, _ := insert.ToSql()
	i2, args2, _ := orgInsert.ToSql()
	// use sqlx
	sx := sqlx.NewDb(s, "sqlite3")
	
	tx, _ := sx.Begin()
	tx.Exec(i1, args1...)
	_, err := tx.Exec(i2, args2...)
	if err != nil { panic(err)}
	err = tx.Commit()
	if err != nil { panic(err)}

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

	orgs := []Organization{}
	err = sx.Select(&orgs, "select id, name, description, createdat from organization")
	if err != nil {
		panic(err)
	}	
	fmt.Printf("\nsqrlxOrgs: \n%+v\n\n", orgs)	
}