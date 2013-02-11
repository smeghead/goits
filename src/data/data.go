package data

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
//	"os"
)

type Project struct {
	Id int
	Name string
	Sort int
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

func exec(databaseName string, callback func(dbHandle *sql.DB) interface{}) interface{} {
	databaseId, err := getDatabaseId(databaseName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("./db/%d.db", databaseId))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer db.Close()

	return callback(db)
}

func GetProjects() []Project {
	projects := exec("manage", func(db *sql.DB) interface{} {
		fmt.Println("hoge")
		stmt, err := db.Prepare("select id, name, sort from project_info where deleted = 0 order by sort")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer stmt.Close()

		rows, err := stmt.Query()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer rows.Close()

		projects := []Project{}
		for rows.Next() {
			var id int
			var name string
			var sort int
			rows.Scan(&id, &name, &sort)
			projects = append(projects, Project{id, name, sort})
		}
		rows.Close()

		return projects
	})
	return projects.([]Project)
}

//package main
//
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/mattn/go-sqlite3"
//	"os"
//)
//
//func main() {
//	os.Remove("./foo.db")
//
//	db, err := sql.Open("sqlite3", "./foo.db")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer db.Close()
//
//	sqls := []string{
//		"create table foo (id integer not null primary key, name text)",
//		"delete from foo",
//	}
//	for _, sql := range sqls {
//		_, err = db.Exec(sql)
//		if err != nil {
//			fmt.Printf("%q: %s\n", err, sql)
//			return
//		}
//	}
//
//	tx, err := db.Begin()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer stmt.Close()
//	for i := 0; i < 100; i++ {
//		_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}
//	tx.Commit()
//
//	rows, err := db.Query("select id, name from foo")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var id int
//		var name string
//		rows.Scan(&id, &name)
//		fmt.Println(id, name)
//	}
//	rows.Close()
//
//	stmt, err = db.Prepare("select name from foo where id = ?")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer stmt.Close()
//	var name string
//	err = stmt.QueryRow("3").Scan(&name)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(name)
//
//	_, err = db.Exec("delete from foo")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	rows, err = db.Query("select id, name from foo")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var id int
//		var name string
//		rows.Scan(&id, &name)
//		fmt.Println(id, name)
//	}
//	rows.Close()
//
//}
