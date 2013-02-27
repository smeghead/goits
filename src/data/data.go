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

func GetProject(projectName string) Project {
    project := Project{}
    project.Url = fmt.Sprintf("/%s", url.QueryEscape(projectName))
    _, err := query(projectName, "select name, value from setting", []interface{}{}, func(rows *sql.Rows) interface{} {
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

func GetWiki(projectName string, wikiName string) Wiki {
    wiki := Wiki{}
    statement := "select w.id, w.name, w.content " +
            "from wiki as w " +
            "where name = ? " +
            "order by w.registerdate desc " +
            "limit 1 "
    params := []interface{} {wikiName}

    _, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
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

func GetElementTypes(projectName string) []ElementType {
    statement := "select id, type, ticket_property, reply_property, required, name, " +
                 "  description, display_in_list, sort, default_value, auto_add_item " +
                 "from element_type " +
                 "where deleted = 0 order by sort"
    params := []interface{}{}

    results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
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

        var listItems []ListItem
        switch element_type {
        case ELEM_TYPE_CHECKBOX, ELEM_TYPE_LIST_SINGLE, ELEM_TYPE_LIST_MULTI, ELEM_TYPE_LIST_SINGLE_RADIO:
            statement := "select id, name, close, sort from list_item where element_type_id = ? order by sort"
            params := []interface{}{id}

            results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
                var id int
                var name string
                var close bool
                var sort int
                rows.Scan(&id, &name, &close, &sort)
                return ListItem{id, name, close, sort}
            })
            if err != nil {
                fmt.Println(err)
                panic(err)
            }
            listItems = make([]ListItem, len(results))
            for i, p := range results { listItems[i] = p.(ListItem) }
        }
        return ElementType{id, element_type, ticket_property, reply_property, required, name,
            description, auto_add_item, default_value, display_in_list, sort, listItems}
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

func getLastMessage(projectName string, ticketId int, elementTypes []ElementType, forList bool) Message {
    fmt.Println("getLastMessage ticket id", ticketId)
    statement := fmt.Sprintf(
        "select last_m.id, t.id, org_m.field%d, %s " + 
        "  , substr(t.registerdate, 1, 16), substr(last_m.registerdate, 1, 16),  " +
        "  julianday(current_date) - julianday(date(last_m.registerdate)) as passed_date " +
        "from ticket as t " +
        "inner join message as last_m on last_m.id = t.last_message_id " +
        "inner join message as org_m on org_m.id = t.original_message_id " +
        "where t.id = ?", ELEM_ID_SENDER, createColumnsExp(elementTypes, "last_m"))
    params := []interface{} {ticketId}

    var messageId int
    elements := []Element{}
    _, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
        fmt.Println("each ticket got")
        pockets, _ := scanDynamicRows(rows)

        fmt.Println(strings.Join(pockets, ","))

        i := 0
        messageId, _ = strconv.Atoi(pockets[i])
        i++
        if forList {
            /* ID */
            elements = append(elements, Element{ELEMENT_TYPE_ID, pockets[i], false})
            i++
            /* 初回投稿者 */
            elements = append(elements, Element{ELEMENT_TYPE_ORG_SENDER, pockets[i], false})
            i++
        } else {
            i += 2
        }
        /* 動的カラム */
        for _, elmType := range elementTypes {
            elements = append(elements, Element{elmType, pockets[i], false})
            i++
        }
        if forList {
            /* 投稿日時 */
            elements = append(elements, Element{ELEMENT_TYPE_REGISTERDATE, pockets[i], false})
            i++
            /* 最終更新日時 */
            elements = append(elements, Element{ELEMENT_TYPE_LASTREGISTERDATE, pockets[i], false})
            i++
            /* 最終更新日時からの経過日数 */
            elements = append(elements, Element{ELEMENT_TYPE_LASTREGISTREDATE_PASSED, pockets[i], false})
            i++
        }

        return nil
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    return Message{messageId, elements, ""}
}

func GetNewestTickets(projectName string, limit int) []Ticket {
    elementTypes := GetElementTypes(projectName)

    statement := fmt.Sprintf("select t.id " +
            "from ticket as t " +
            "inner join message as m on m.id = t.last_message_id " +
            "order by m.registerdate desc " +
            "limit %d ", limit)
    params := []interface{} {}

    results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
        fmt.Println("newest tickets got")
        var id int
        rows.Scan(&id)
        lastMessage := getLastMessage(projectName, id, elementTypes, false)
        return NewTicketWithoutMessages(id, lastMessage)
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    tickets := make([]Ticket, len(results))
    for i, p := range results { tickets[i] = p.(Ticket) }
    return tickets
}

func GetStates(projectName string, notClose bool) []State {
    condition := ""
    if notClose {
        condition = "where l.close = 0 "
    }
    statement := fmt.Sprintf(
            "select m.field%d as name, count(t.id) as count " +
            "from ticket as t " +
            "inner join message as m " +
            " on m.id = t.last_message_id " +
            "inner join list_item as l " +
            " on l.element_type_id = %d and l.name = m.m.field%d " +
            "%s " +
            "group by m.field%d " +
            "order by l.sort ", ELEM_ID_STATUS, ELEM_ID_STATUS, ELEM_ID_STATUS, condition, ELEM_ID_STATUS)
    params := []interface{}{}

    results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
        var id int
        var name string
        var count int

        rows.Scan(&name, &count)
        return State{id, name, count}
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    states := make([]State, len(results))
    for i, p := range results { states[i] = p.(State) }
    return states
}

func validConditionSize(conditions []Condition) int {
    size := 0
    for _, c := range conditions {
        if validCondition(c) {
            size++
        }
    }
    return size
}

func validCondition(c Condition) bool {
    return len(c.Value) > 0 || len(c.CookieValue) > 0
}

func getConditionValidValue(c Condition) string {
    if len(c.Value) > 0 {
        return c.Value
    } else {
        return c.CookieValue
    }
    return ""
}

func getSearchSqlString(projectName string, conditions []Condition, sort Condition, keywords []string) string {
    sqlString := "select t.id as id " +
        "from ticket as t " +
        "inner join message as m on m.id = t.last_message_id " +
        "inner join message as org_m on org_m.id = t.original_message_id "

    if len(keywords) > 0 {
        sqlString += "inner join message as m_all on m_all.ticket_id = t.id "
    }

    if validConditionSize(conditions) > 0 || len(keywords) > 0 {
        sqlString += "where "
    }

    if validConditionSize(conditions) > 0 {
        conds := []string{}
        for _, cond := range conditions {
            if cond.ElementTypeId < 0 {
                continue
            }
            if !validCondition(cond) {
                continue
            }
            switch cond.ConditionType {
            case CONDITION_TYPE_DATE_FROM:
                conds = append(conds,
                    fmt.Sprintf(
                        "(length(m.field%d) > 0 and m.field%d >= ?)",
                        cond.ElementTypeId,
                        cond.ElementTypeId))
            case CONDITION_TYPE_DATE_TO:
                conds = append(conds,
                    fmt.Sprintf("(length(m.field%d) > 0 and m.field%d <= ?)",
                        cond.ElementTypeId,
                        cond.ElementTypeId))
            default:
                senderTablePrefix := ""
                if cond.ElementTypeId == ELEM_ID_SENDER {
                    senderTablePrefix = "org_"
                }
                conds = append(conds,
                    fmt.Sprintf("(%sm.field%d like '%%' || ? || '%%')",
                        senderTablePrefix, /* 投稿者は初回投稿者が検索対象になる。 */
                        cond.ElementTypeId))
            }
        }
        sqlString += strings.Join(conds, " and ")

        conds = []string{}
        for _, cond := range conditions {
            if cond.ElementTypeId > 0 {
                continue
            }
            if !validCondition(cond) {
                continue
            }

            name := ""
            switch cond.ElementTypeId {
                case ELEM_ID_REGISTERDATE:
                    name = "org_m.registerdate"
                case ELEM_ID_LASTREGISTERDATE:
                    name = "m.registerdate"
            }

            switch cond.ConditionType {
                case CONDITION_TYPE_DATE_FROM:
                    conds = append(conds, fmt.Sprintf("(%s >= ?)", name))
                case CONDITION_TYPE_DATE_TO:
                    conds = append(conds, fmt.Sprintf("(%s <= ?)", name))
                case CONDITION_TYPE_DAYS:
                    conds = append(conds, fmt.Sprintf("(%s >= datetime(current_timestamp, 'utc', '-%s days', 'localtime'))",
                            name,
                            getConditionValidValue(cond)))
            }
        }
        sqlString += strings.Join(conds, " and ")
    }

//    if len(keywords) > 0 {
//        String* columns_a = string_new(0);
//        List* element_types_a;
//        list_alloc(element_types_a, ElementType, element_type_new, element_type_free);
//        element_types_a = db_get_element_types_all(db, NULL, element_types_a);
//        columns_a = create_columns_like_exp(element_types_a, "m_all", keywords, columns_a);
//        if (valid_condition_size(conditions))
//            string_append(sql_string, " and ");
//        string_appendf(sql_string, "(%s)", string_rawstr(columns_a));
//        string_free(columns_a);
//        list_free(element_types_a);
//    }

    sqlString += " group by t.id order by "

    if sort.ElementTypeId != 0 {
        sort_type := "asc"
        if sort.Value == "reverse" {
            sort_type = "desc"
        }
        switch sort.ElementTypeId {
            case -1:
                sqlString += fmt.Sprintf("t.id %s, ", sort_type)
            case -2:
                sqlString += fmt.Sprintf("t.registerdate %s, ", sort_type)
            case -3:
                sqlString += fmt.Sprintf("m.registerdate %s, ", sort_type)
            case ELEM_ID_SENDER:
                sqlString += fmt.Sprintf("org_m.field%d %s, ", sort.ElementTypeId, sort_type)
            default:
                sqlString += fmt.Sprintf("m.field%d %s, ", sort.ElementTypeId, sort_type)
        }
    }
    sqlString += "t.registerdate desc, t.id desc "

    return sqlString
}

func GetTicketsByStatus(projectName string, status string) SearchResult {
    keywords := []string{}
    conditions := []Condition{Condition{ELEM_ID_STATUS, 0, status, ""}}
    sort := Condition{}

    elementTypes := GetElementTypes(projectName)

    statement := getSearchSqlString(projectName, conditions, sort, keywords)
    statement += " limit ? "

    params := []interface{}{status, LIST_COUNT_PER_LIST_PAGE}

    results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
        var id int
        rows.Scan(&id)
        fmt.Println("ticket id:" , id)
        lastMessage := getLastMessage(projectName, id, elementTypes, true)
        fmt.Println("elements count:" , len(lastMessage.Elements))
        return NewTicketWithoutMessages(id, lastMessage)
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    tickets := make([]Ticket, len(results))
    for i, p := range results { tickets[i] = p.(Ticket) }

    sums := []int{}
    //TODO 数値項目の合計
    //sums := set_tickets_number_sum(db, conditions, NULL, keywords_a, result);

    return SearchResult{len(tickets), 0, tickets, nil, sums}
}

func GetSettingFile(projectName string, name string) SettingFile {
    settingFile := SettingFile{}

    statement := "select name, file_name, size, mime_type, content " +
        "from setting_file " +
        "where name = ? "
    _, err := query(projectName, statement, []interface{}{name}, func(rows *sql.Rows) interface{} {
        var name string
        var filename string
        var size int
        var mime_type string
        var content string
        rows.Scan(&name, &filename, &size, &mime_type, &content)
        settingFile = SettingFile{name, filename, size, mime_type, content}
        return nil
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    return settingFile
}

func getMessages(projectName string, ticketId int, elementTypes []ElementType) []Message {
    statement := fmt.Sprintf(
        "select m.id, m.registerdate, %s from message as m where m.ticket_id = ? order by m.registerdate",
        createColumnsExp(elementTypes, "m"))

    params := []interface{} {ticketId}

    results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
        fmt.Println("each ticket got")
        pockets, _ := scanDynamicRows(rows)

        elements := []Element{}
        fmt.Println(strings.Join(pockets, ","))

        i := 0
        messageId, _ := strconv.Atoi(pockets[i])
        i++
        registerDate := pockets[i]

        /* 動的カラム */
        for _, elmType := range elementTypes {
            elements = append(elements, Element{elmType, pockets[i], false})
            i++
        }

        return Message{messageId, elements, registerDate}
    })
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
    messages := make([]Message, len(results))
    for i, p := range results { messages[i] = p.(Message) }

    return messages
}

func GetTicket(projectName string, ticketId int, elementTypes []ElementType) Ticket {
    lastMessage := getLastMessage(projectName, ticketId, elementTypes, false)
    messages := getMessages(projectName, ticketId, elementTypes)

    return NewTicket(ticketId, lastMessage, messages)
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
