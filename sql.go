package main

import (
	"time"
	"fmt"
	"database/sql"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type ID uint64
type Status string
type User struct {
	ID      ID     `json:"id,omitempty"`
	Name    string `json:"name"`
	OAuthID string `json:"oauthID,omitempty"`
	Status  Status `json:"status"`
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
	dbMap := modl.NewDbMap(s, modl.SqliteDialect{})
	dbMap.AddTable(User{}).ColMap("id").SetUnique(true)
	tm :=dbMap.AddTable(UserResourceMapping{})
	tm.ColMap("userid").SetUnique(true)
	tm.ColMap("resourceid").SetUnique(true)
	dbMap.AddTable(Organization{}).ColMap("id").SetUnique(true)
	dbMap.AddTable(Bucket{}).ColMap("id").SetUnique(true)

	err = dbMap.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}

	rows, err := dbMap.Dbx.Query("SELECT * FROM sqlite_master")
	printTable(rows)

	dbMap.Insert(&User{Name: "lyon"}, &Organization{Name: "influx"})
	err = dbMap.Insert(&User{Name: "lyon", ID: 2})
	if err != nil {
		panic(err)
	}
	users := []User{}
	err = dbMap.Select(&users, "select * from user")
	if err != nil {
		panic(err)
	}
	fmt.Printf("users: %#v\n", users)
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