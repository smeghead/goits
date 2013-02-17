package data

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "reflect"
)

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
        return nil, err
    }
    db, err := sql.Open("sqlite3", fmt.Sprintf("./db/%d.db", databaseId))
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    defer db.Close()

    fmt.Println("sql:", statement)
    stmt, err := db.Prepare(statement)
    if err != nil {
        fmt.Println(err)
        return nil, err
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
        return nil, err
    }
    defer rows.Close()

    results := []interface{}{}
    for rows.Next() {
        results = append(results, callback(rows))
    }
    rows.Close()

    return results, nil
}

func scanDynamicRows(rows *sql.Rows) ([]string, error) {
    columns, err := rows.Columns()
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    count := len(columns)
    values := []string{}
    switch count {
    case 0:
    case 1:
        var v1 string
        rows.Scan(&v1)
        values = append(values, v1)
    case 2:
        var v1, v2 string
        rows.Scan(&v1, &v2)
        values = append(values, v1, v2)
    case 3:
        var v1, v2, v3 string
        rows.Scan(&v1, &v2, &v3)
        values = append(values, v1, v2, v3)
    case 4:
        var v1, v2, v3, v4 string
        rows.Scan(&v1, &v2, &v3, &v4)
        values = append(values, v1, v2, v3, v4)
    case 5:
        var v1, v2, v3, v4, v5 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5)
        values = append(values, v1, v2, v3, v4, v5)
    case 6:
        var v1, v2, v3, v4, v5, v6 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6)
        values = append(values, v1, v2, v3, v4, v5, v6)
    case 7:
        var v1, v2, v3, v4, v5, v6, v7 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7)
        values = append(values, v1, v2, v3, v4, v5, v6, v7)
    case 8:
        var v1, v2, v3, v4, v5, v6, v7, v8 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8)
    case 9:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9)
    case 10:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10)
    case 11:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11)
    case 12:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12)
    case 13:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13)
    case 14:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14)
    case 15:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15)
    case 16:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16)
    case 17:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17)
    case 18:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18)
    case 19:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19)
    case 20:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20)
    case 21:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21)
    case 22:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22)
    case 23:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23)
    case 24:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24)
    case 25:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25)
    case 26:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26)
    case 27:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27)
    case 28:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27, &v28)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28)
    case 29:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27, &v28, &v29)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29)
    case 30:
        var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29, v30 string
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27, &v28, &v29, &v30)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29, v30)
    default:
        panic("too many dynamic columns")
    }
    return values, nil
}
