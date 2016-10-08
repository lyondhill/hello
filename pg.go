package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
    "os/exec"
    "os"
    "strings"
)

var db *sql.DB

type piper struct {
  leftover string
}

func (p *piper) Write(d []byte) (int, error) {
  str := string(d)
  strArr := strings.Split(str, "\n")
  for index, object := range strArr {
    if index == 0 {
      object = p.leftover+object
    }
    if index + 1 == len(strArr) && !strings.HasSuffix(object, "\n") {
      p.leftover = object
      break
    }

    id := strings.Join(strings.Split(object, "/")[3:], "-")
    // fmt.Printf("id: %s...", id)
    if existsInDatabase(id) {
      // fmt.Printf("ok\n")
    } else {
      fmt.Printf("id: %s\n", id)
      // fmt.Printf("FAILED!\n")
      err := os.Remove(object)
      if err != nil {
        fmt.Println("error:", err.Error())
      }
    }
    // fmt.Println("id: "+object)
  }
  return len(d), nil
}

func main() {
  establishDatabseConnection()

  fmt.Println("First Check for files that dont exist in the database:")

  cmd := exec.Command("find", "/mnt/data/", "-mindepth", "5")
  cmd.Stdout = &piper{}
  cmd.Run()

  fmt.Println("Now Check for Database records that dont exist on the file system:")

  rows, err := db.Query("SELECT id FROM objects")
  checkErr(err)
  
  for rows.Next() {
    var id string
    rows.Scan(&id)
    file := strings.Replace(id, "-", "/", -1)
    if !existsInFileSystem(file) {
      fmt.Println("nofile: "+file)
    }
  }


    // // // Query
    // ids := []string{

    // }
    // for _, _ := range ids {
    // }

    // // Delete
    // stmt, err = db.Prepare("delete from userinfo where uid=$1")
    // checkErr(err)

    // res, err = stmt.Exec(1)
    // checkErr(err)

    // affect, err = res.RowsAffected()
    // checkErr(err)

    // fmt.Println(affect)

    // db.Close()

}

func establishDatabseConnection() {
  if db != nil {
    db.Close()
  }
  var err error
  db, err = sql.Open("postgres", "dbname=mdata user=blobstache sslmode=disable host=192.168.0.100 port=5432")
  if err != nil {
    fmt.Println(err)
  }
}

func existsInDatabase(id string) bool {
  rows, err := db.Query("SELECT * FROM objects WHERE id='"+id+"'")
  if err != nil {
    establishDatabseConnection()
    return true
  }
  
  for rows.Next() {
    return true
  }
  return false
}

func existsInFileSystem(file string) bool {
  f, err := os.Open("/mnt/data/"+file)
  if err != nil {
    return false
  }
  defer f.Close()

  _, err = f.Stat()
  if err != nil {
    return false
  }
  return true
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}