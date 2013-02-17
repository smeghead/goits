package data

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "strconv"
    "strings"
    "net/url"
)

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
    project.Url = url.QueryEscape(databaseName)
    _, err := query(databaseName, "select name, value from setting", []interface{}{}, func(rows *sql.Rows) interface{} {
        var name string
        var value string
        rows.Scan(&name, &value)
        switch name {
        case "project_name":
            project.Name = value
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

func getElementTypes(databaseName string) []ElementType {
    statement := "select id, type, ticket_property, reply_property, required, name, " +
                 "  description, display_in_list, sort, default_value, auto_add_item " +
                 "from element_type " +
                 "where deleted = 0 order by sort"
    params := []interface{}{}

    results, err := query(databaseName, statement, params, func(rows *sql.Rows) interface{} {
        var id int
        var element_type int
        var ticket_property bool
        var reply_property bool
        var required bool
        var name string
        var description string
        var display_in_list bool
        var sort int
        var default_value string
        var auto_add_item bool

        rows.Scan(&id, &element_type, &ticket_property, &reply_property, &required, &name,
            &description, &display_in_list, &sort, &default_value, &auto_add_item)
        return ElementType{id, element_type, ticket_property, reply_property, required, name,
            description, auto_add_item, default_value, display_in_list, sort}
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    elementTypes := make([]ElementType, len(results))
    for i, p := range results { elementTypes[i] = p.(ElementType) }
    return elementTypes
}

func createColumnsExp(elementTypes []ElementType, table_name string) string {
    columns := []string{}
    for _, elem := range elementTypes {
        columns = append(columns, fmt.Sprintf("%s.field%d", table_name, elem.Id))
    }
    return strings.Join(columns, ", ")
}

func GetNewestTickets(databaseName string, limit int) []Message {
    statement := fmt.Sprintf("select t.id " +
            "from ticket as t " +
            "inner join message as m on m.id = t.last_message_id " +
            "order by m.registerdate desc " +
            "limit %d ", limit)
    params := []interface{} {}

    results, err := query(databaseName, statement, params, func(rows *sql.Rows) interface{} {
        fmt.Println("newest tickets got")
        var id int
        rows.Scan(&id)
        return Message{id, []Element{}}
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    elementTypes := getElementTypes(databaseName)

    tickets := []Message{}
    for _, t := range results {
        ticket := t.(Message)
        statement := fmt.Sprintf(
            "select t.id, org_m.field%d, %s " + 
            "  , substr(t.registerdate, 1, 16), substr(last_m.registerdate, 1, 16),  " +
            "  julianday(current_date) - julianday(date(last_m.registerdate)) as passed_date " +
            "from ticket as t " +
            "inner join message as last_m on last_m.id = t.last_message_id " +
            "inner join message as org_m on org_m.id = t.original_message_id " +
            "where t.id = ?", ELEM_ID_SENDER, createColumnsExp(elementTypes, "last_m"))
        params := []interface{} {ticket.Id}

        _, err := query(databaseName, statement, params, func(rows *sql.Rows) interface{} {
            fmt.Println("newest each ticket got")
            pockets, _ := scanDynamicRows(rows)

            fmt.Println(strings.Join(pockets, ","))
            elements := []Element{}

            i := 0
            /* ID */
            elements = append(elements, Element{ELEM_ID_ID, pockets[i], false})
            i++
            /* 初回投稿者 */
            elements = append(elements, Element{ELEM_ID_ORG_SENDER, pockets[i], false})
            i++
            /* 動的カラム */
            for _, elmType := range elementTypes {
                elements = append(elements, Element{elmType.Id, pockets[i], false})
                i++
            }
            /* 投稿日時 */
            elements = append(elements, Element{ELEM_ID_REGISTERDATE, pockets[i], false})
            i++
            /* 最終更新日時 */
            elements = append(elements, Element{ELEM_ID_LASTREGISTERDATE, pockets[i], false})
            i++
            /* 最終更新日時からの経過日数 */
            elements = append(elements, Element{ELEM_ID_LASTREGISTERDATE_PASSED, pockets[i], false})
            i++

            ticket.Elements = elements
            return nil
        })
        if err != nil {
            fmt.Println(err)
            panic(err)
        }
        tickets = append(tickets, ticket)
    }
    return tickets

//    strcat(sql, sql_suf);

//List* db_get_last_elements_4_list(Database* db, const int ticket_id, List* elements)
//{
//    char sql[DEFAULT_LENGTH] = "";
//    char sql_suf[DEFAULT_LENGTH] = "";
//    sqlite3_stmt *stmt = NULL;
//    int r;
//    char columns[DEFAULT_LENGTH] = "";
//    List* element_types_a = NULL;
//    Iterator* it;
//
//    list_alloc(element_types_a, ElementType, element_type_new, element_type_free);
//    element_types_a = db_get_element_types_4_list(db, NULL, element_types_a);
//    create_columns_exp(element_types_a, "last_m", columns);
//    sprintf(sql, "select t.id, org_m.field%d ", ELEM_ID_SENDER);
//    strcat(sql, columns);
//    sprintf(sql_suf, 
//            "  , substr(t.registerdate, 1, 16), substr(last_m.registerdate, 1, 16),  "
//            "  julianday(current_date) - julianday(date(last_m.registerdate)) as passed_date "
//            "from ticket as t "
//            "inner join message as last_m on last_m.id = t.last_message_id "
//            "inner join message as org_m on org_m.id = t.original_message_id "
//            "where t.id = ?");
//    strcat(sql, sql_suf);
//    if (sqlite3_prepare(db->handle, sql, strlen(sql), &stmt, NULL) == SQLITE_ERROR) goto error;
//    sqlite3_reset(stmt);
//    sqlite3_bind_int(stmt, 1, ticket_id);
//
//    while (SQLITE_ROW == (r = sqlite3_step(stmt))) {
//        int i = 0;
//        Element* e;
//        /* ID */
//        e = list_new_element(elements);
//        e->element_type_id = ELEM_ID_ID;
//        set_str_val(e, sqlite3_column_text(stmt, i++));
//        list_add(elements, e);
//        /* 初回投稿者 */
//        e = list_new_element(elements);
//        e->element_type_id = ELEM_ID_ORG_SENDER;
//        set_str_val(e, sqlite3_column_text(stmt, i++));
//        list_add(elements, e);
//        foreach (it, element_types_a) {
//            ElementType* et = it->element;
//            e = list_new_element(elements);
//            e->element_type_id = et->id;
//            set_str_val(e, sqlite3_column_text(stmt, i++));
//            list_add(elements, e);
//        }
//        /* 投稿日時 */
//        e = list_new_element(elements);
//        e->element_type_id = ELEM_ID_REGISTERDATE;
//        set_str_val(e, sqlite3_column_text(stmt, i++));
//        list_add(elements, e);
//        /* 最終更新日時 */
//        e = list_new_element(elements);
//        e->element_type_id = ELEM_ID_LASTREGISTERDATE;
//        set_str_val(e, sqlite3_column_text(stmt, i++));
//        list_add(elements, e);
//        /* 最終更新日時からの経過日数 */
//        e = list_new_element(elements);
//        e->element_type_id = ELEM_ID_LASTREGISTERDATE_PASSED;
//        set_str_val(e, sqlite3_column_text(stmt, i++));
//        list_add(elements, e);
//    }
//    if (SQLITE_DONE != r)
//        goto error;
//
//    sqlite3_finalize(stmt);
//    list_free(element_types_a);
//
//    return elements;
//ERROR_LABEL(db->handle)
//}


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
