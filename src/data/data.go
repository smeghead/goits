package data

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "strconv"
    "errors"
    "reflect"
)

type ProjectInfo struct {
    Id int
    Name string
    Sort int
}
func (p *ProjectInfo) IsManage() bool {
    return p.Name == "top"
}

type Project struct {
    ProjectName string
    HomeDescription string
    HomeUrl string
    UploadMaxSize int
    Locale string
}

type Wiki struct {
    Id int
    Name string
    Content string
}

func getDatabaseId(databaseName string) (int, error) {
    if databaseName == "manage" {
        return 1, nil
    }

    db, err := sql.Open("sqlite3", "./db/1.db")
    if err != nil {
        fmt.Println(err)
        return -1, err
    }
    defer db.Close()

    stmt, err := db.Prepare("select id from project_info where name = ?")
    if err != nil {
        fmt.Println(err)
        return -1, err
    }
    defer stmt.Close()

    var id int
    err = stmt.QueryRow(databaseName).Scan(&id)
    if err != nil {
        fmt.Println(err)
        return -1, err
    }
    return id, nil
}

func query(databaseName string, statement string, params []interface{}, callback func(rows *sql.Rows) interface{}) ([]interface{}, error) {
    databaseId, err := getDatabaseId(databaseName)
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("failed to retreive databaseId.")
    }
    db, err := sql.Open("sqlite3", fmt.Sprintf("./db/%d.db", databaseId))
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("failed to open database")
    }
    defer db.Close()

    fmt.Println(statement)
    stmt, err := db.Prepare(statement)
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("failed to create statement")
    }
    defer stmt.Close()

    //execute `rows, err := stmt.Query(arg1, arg2, ...)` by reflect
    values := []reflect.Value{}
    for _, p := range params { values = append(values, reflect.ValueOf(p)) }
    returnValues := reflect.ValueOf(stmt).MethodByName("Query").Call(values)
    rows := returnValues[0].Interface().(*sql.Rows)
    if !returnValues[1].IsNil() {
        err = returnValues[1].Interface().(error)
    }
    if err != nil {
        fmt.Println(err)
        return nil, errors.New("failed to execute query")
    }
    defer rows.Close()

    results := []interface{}{}
    for rows.Next() {
        results = append(results, callback(rows))
    }
    rows.Close()

    return results, nil
}

func GetProjectInfos() []ProjectInfo {
    results, err := query(
            "manage",
            "select id, name, sort from project_info where deleted = 0 order by sort",
            []interface{}{},
            func(rows *sql.Rows) interface{} {
        var id int
        var name string
        var sort int
        rows.Scan(&id, &name, &sort)
        return ProjectInfo{id, name, sort}
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    projectInfos := make([]ProjectInfo, len(results))
    for i, p := range results { projectInfos[i] = p.(ProjectInfo) }
    return projectInfos
}

func GetProject(databaseName string) Project {
    project := Project{}
    _, err := query(databaseName, "select name, value from setting", []interface{}{}, func(rows *sql.Rows) interface{} {
        var name string
        var value string
        rows.Scan(&name, &value)
        switch name {
        case "project_name":
            project.ProjectName = value
        case "home_description":
            project.HomeDescription = value
        case "home_url":
            project.HomeUrl = value
        case "upload_max_size":
            i, _ := strconv.Atoi(value)
            project.UploadMaxSize = i
        case "locale":
            project.Locale = value
        }
        return nil
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    return project
}

func GetWiki(databaseName string, wikiName string) Wiki {
    wiki := Wiki{}
    statement := "select w.id, w.name, w.content " +
            "from wiki as w " +
            "where name = ? " +
            "order by w.registerdate desc " +
            "limit 1 "
    params := []interface{} {wikiName}

    _, err := query(databaseName, statement, params, func(rows *sql.Rows) interface{} {
        fmt.Println("wiki got")
        var id int
        var name string
        var content string
        rows.Scan(&id, &name, &content)
        wiki.Id = id
        wiki.Name = name
        wiki.Content = content
        return nil
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    return wiki
}

//package main
//
//import (
//  "database/sql"
//  "fmt"
//  _ "github.com/mattn/go-sqlite3"
//  "os"
//)
//
//func main() {
//  os.Remove("./foo.db")
//
//  db, err := sql.Open("sqlite3", "./foo.db")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  defer db.Close()
//
//  sqls := []string{
//      "create table foo (id integer not null primary key, name text)",
//      "delete from foo",
//  }
//  for _, sql := range sqls {
//      _, err = db.Exec(sql)
//      if err != nil {
//          fmt.Printf("%q: %s\n", err, sql)
//          return
//      }
//  }
//
//  tx, err := db.Begin()
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  defer stmt.Close()
//  for i := 0; i < 100; i++ {
//      _, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
//      if err != nil {
//          fmt.Println(err)
//          return
//      }
//  }
//  tx.Commit()
//
//  rows, err := db.Query("select id, name from foo")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  defer rows.Close()
//  for rows.Next() {
//      var id int
//      var name string
//      rows.Scan(&id, &name)
//      fmt.Println(id, name)
//  }
//  rows.Close()
//
//  stmt, err = db.Prepare("select name from foo where id = ?")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  defer stmt.Close()
//  var name string
//  err = stmt.QueryRow("3").Scan(&name)
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  fmt.Println(name)
//
//  _, err = db.Exec("delete from foo")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//
//  _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//
//  rows, err = db.Query("select id, name from foo")
//  if err != nil {
//      fmt.Println(err)
//      return
//  }
//  defer rows.Close()
//  for rows.Next() {
//      var id int
//      var name string
//      rows.Scan(&id, &name)
//      fmt.Println(id, name)
//  }
//  rows.Close()
//
//}
/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
