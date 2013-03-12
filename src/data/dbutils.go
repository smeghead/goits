package data

import (
    logger "code.google.com/p/log4go"
    "github.com/gosexy/gettext"
    "os"
    "fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "reflect"
)

func getDatabaseId(projectName string) (int, error) {
    manageDbFile := "./db/1.db"
    if _, err := os.Stat(manageDbFile); err != nil {
        if os.IsNotExist(err) {
            logger.Warn("create top table.")
            CreateTopTables()
        }
    }

    if projectName == "manage" {
        return 1, nil
    }

    db, err := sql.Open("sqlite3", manageDbFile)
    if err != nil {
        logger.Error(err)
        return -1, err
    }
    defer db.Close()

    stmt, err := db.Prepare("select id from project_info where name = ?")
    if err != nil {
        logger.Error(err)
        return -1, err
    }
    defer stmt.Close()

    var id int
    err = stmt.QueryRow(projectName).Scan(&id)
    if err != nil {
        logger.Error(err)
        return -1, err
    }
    return id, nil
}

func query(projectName string, statement string, params []interface{}, callback func(rows *sql.Rows) interface{}) ([]interface{}, error) {
    databaseId, err := getDatabaseId(projectName)
    if err != nil {
        logger.Error(err)
        return nil, err
    }
    projectDbFile := fmt.Sprintf("./db/%d.db", databaseId)
    if _, err := os.Stat(projectDbFile); err != nil {
        if os.IsNotExist(err) {
            logger.Warn("create project table.")
            CreateProjectTables(projectName, databaseId)
        }
    }
    logger.Warn("created project table.")
    db, err := sql.Open("sqlite3", projectDbFile)
    if err != nil {
        logger.Error(err)
        return nil, err
    }
    defer db.Close()

    logger.Debug("sql: %s", statement)
    stmt, err := db.Prepare(statement)
    if err != nil {
        logger.Error(err)
        return nil, err
    }
    defer stmt.Close()

    //execute `rows, err := stmt.Query(arg1, arg2, ...)` by reflect
    values := []reflect.Value{}
    for _, p := range params {
        logger.Debug("[%s]", p)
        values = append(values, reflect.ValueOf(p))
    }
    returnValues := reflect.ValueOf(stmt).MethodByName("Query").Call(values)
    rows := returnValues[0].Interface().(*sql.Rows)
    if !returnValues[1].IsNil() {
        err = returnValues[1].Interface().(error)
    }
    if err != nil {
        logger.Error(err)
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
        logger.Error(err)
        return nil, err
    }
    var v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29, v30 string
    count := len(columns)
    values := []string{}
    switch count {
    case 0:
    case 1:
        rows.Scan(&v1)
        values = append(values, v1)
    case 2:
        rows.Scan(&v1, &v2)
        values = append(values, v1, v2)
    case 3:
        rows.Scan(&v1, &v2, &v3)
        values = append(values, v1, v2, v3)
    case 4:
        rows.Scan(&v1, &v2, &v3, &v4)
        values = append(values, v1, v2, v3, v4)
    case 5:
        rows.Scan(&v1, &v2, &v3, &v4, &v5)
        values = append(values, v1, v2, v3, v4, v5)
    case 6:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6)
        values = append(values, v1, v2, v3, v4, v5, v6)
    case 7:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7)
        values = append(values, v1, v2, v3, v4, v5, v6, v7)
    case 8:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8)
    case 9:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9)
    case 10:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10)
    case 11:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11)
    case 12:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12)
    case 13:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13)
    case 14:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14)
    case 15:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15)
    case 16:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16)
    case 17:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17)
    case 18:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18)
    case 19:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19)
    case 20:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20)
    case 21:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21)
    case 22:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22)
    case 23:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23)
    case 24:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24)
    case 25:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25)
    case 26:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26)
    case 27:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27)
    case 28:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27, &v28)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28)
    case 29:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27, &v28, &v29)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29)
    case 30:
        rows.Scan(&v1, &v2, &v3, &v4, &v5, &v6, &v7, &v8, &v9, &v10, &v11, &v12, &v13, &v14, &v15, &v16, &v17, &v18, &v19, &v20, &v21, &v22, &v23, &v24, &v25, &v26, &v27, &v28, &v29, &v30)
        values = append(values, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22, v23, v24, v25, v26, v27, v28, v29, v30)
    default:
        panic("too many dynamic columns")
    }
    return values, nil
}

func tran(projectName string, callback func(tx *sql.Tx) error) error {
    databaseId, err := getDatabaseId(projectName)
    if err != nil {
        logger.Error(err)
        return err
    }
    db, err := sql.Open("sqlite3", fmt.Sprintf("./db/%d.db", databaseId))
    if err != nil {
        logger.Error(err)
        return err
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        logger.Error(err)
        return err
    }
    err = callback(tx)
    if err != nil {
        logger.Error(err)
        tx.Rollback()
        return err
    }
    tx.Commit()
    return nil
}

func exec(tx *sql.Tx, statement string, params []interface{}) sql.Result {
    result, err := tx.Exec(statement, params)
    if err != nil {
        logger.Error(err)
        tx.Rollback()
        panic(err)
    }
    return result
}

func CreateTopTables() {
    logger.Debug("create_top_tables")

    db, err := sql.Open("sqlite3", "./db/1.db")
    if err != nil {
        logger.Error(err)
        panic(err)
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        logger.Error(err)
        panic(err)
    }
    db.Exec("create table setting( " +
        " name text primary key, " +
        " value text " +
        ");")
    db.Exec("create table setting( " +
            " name text primary key, " +
            " value text " +
            ");")
    db.Exec("insert into setting(name, value)" +
            "values ('project_name', ?);", gettext.Gettext("main project"))
    db.Exec("insert into setting(name, value)" +
            "values ('home_description', ?);", gettext.Gettext("home"))
    db.Exec("insert into setting(name, value)" +
            "values ('home_url', '');")
    db.Exec("insert into setting(name, value)" +
            "values ('locale', ?);", "ja_JP")
    db.Exec("create table project_info ( " +
            " id integer not null primary key, " +
            " name text, " +
            " deleted integer default 0, " +
            " sort integer " +
            ");")
    db.Exec("insert into project_info(id, name, deleted, sort) " +
            "values (1, 'top', 0, 1);")
    db.Exec("insert into project_info(id, name, deleted, sort) " +
            "values (2, 'Sand Box', 0, 1);")
    db.Exec("create table wiki( " +
            " id integer not null primary key, " +
            " name text, " +
            " content text, " +
            " registerdate text " +
            ");")
    db.Exec("insert into wiki(id, name, content, registerdate) " +
            "values (NULL, 'top', ?, current_timestamp);", gettext.Gettext("[wiki syntax example]"))
logger.Debug("wiki content: %s", gettext.Gettext("[wiki syntax example]"))
    logger.Debug("create_top_tables commit")
    tx.Commit()
    logger.Debug("create_top_tables commited")
}

func CreateProjectTables(projectName string, id int) {
    logger.Debug("CreateProjectTables")

    db, err := sql.Open("sqlite3", fmt.Sprintf("./db/%d.db", id))
    if err != nil {
        logger.Error(err)
        panic(err)
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        logger.Error(err)
        panic(err)
    }
    db.Exec("create table setting( " +
            " name text primary key, " +
            " value text " +
            ");")
    db.Exec(
            "insert into setting(name, value)" +
            "values ('project_name', ?);",
            projectName)
    db.Exec("insert into setting(name, value)" +
            "values ('home_description', ?);",
            gettext.Gettext("home"))
    db.Exec("insert into setting(name, value)" +
            "values ('home_url', '');")
    db.Exec("insert into setting(name, value)" +
            "values ('upload_max_size', 512);")
    db.Exec("create table setting_file( " +
            " name text primary key, " +
            " file_name text, " +
            " size integer, " +
            " mime_type text, " +
            " content blob " +
            ");")
    textContent := gettext.Gettext("[style sheet default value]")
    db.Exec(
            "insert into setting_file(name, file_name, size, mime_type, content)" +
            "values ('user.css', 'user.css', ?, 'text/css', ?);",
            len(textContent), textContent)
    db.Exec(
            "insert into setting_file " +
            "values(" +
            " 'top_image', " +
            " 'title.jpg', " +
            " 4067, " +
            " 'image/jpeg', " +
            "X'89504E470D0A1A0A0000000D49484452000000A0000000460806000000" +
            "4BD29FB20000000467414D410000B18F0BFC610500000F9A49444154785E" +
            "ED5CFF8F5C5515C73F41A2123562D4A8C12FC12F5123311AF55720424389" +
            "504C301231242424184BF956B02DA845428D5424B4A5142A15C1408B9416" +
            "E9376A4384022B94D252DA6DB7DD76BADBDD999D37EFEB5CCF9BD99979EF" +
            "9C73EF3DF7CD76BB0BAFC90DCCCCBDE77ECEE7F3B9E7DEFB76673FD46C36" +
            "D559E5BF928133C5406AC0B2951C9C290F9C75A6262EE72D4DDFDA7D55F9" +
            "CF9981E69BBF508DB517A8E0CD23CE63CB0179064A038A1DF11F15AEF9AE" +
            "6A645A30501A504C9FA6A3680B8EB77F4735566BDAF69D1F9033E44E1520" +
            "0E8281C31F90DC1D8E0B87EE9AF4CA5D2A16DC2F8C068CB77E5B35564ADB" +
            "12D184AD985B67A369C180888BE0F5D280B9B3FCC12519BFC8FCA031E0A0" +
            "0AFFF12DE53DE8D2AE56E129CD4A39B55635B2B1FE3D4B0D88F868EC1E2C" +
            "2BE064954B765F8DFCB25845452B60B4F99BCAFB8B63DBFC92460C30F3DF" +
            "512C6D5F87522F486E6A6FDA2FA90071D278B53460CA31EF97A2061C7D4C" +
            "35FEFC0DE5A1D6F82F253BDAD4EB17BCAB330F18F06F28DE269D5967B801" +
            "059C4CADE967321F80EDDDC5C4273DDF1435E0FE45CA5BF6B55CE3CCD723" +
            "1A0CF6D82243B94D3FFF7A3EE6C6596AC065F93CCCBCCC70F3F4B583C06E" +
            "80B8C09EF196993CD1E3869C019397AF52F57BCFCFB4AB5438DA0F9960C0" +
            "35D978F0FFFFDA310BCF4E3B949FE3E57CE5BDFCFEDA8253EDED39511EF2" +
            "7EE9685DD480BBE6A9FAD2AFE69ABFAF800147D6286FE957482C1CBBF7FA" +
            "4E7D15DD77A739CE0699A113925B674E58240F7358B3988078948FB70B0C" +
            "A8CDD3900FAA3E7A5C3ADE35584C554DCBE13C158E34550E83914F3A37AF" +
            "A92C7F7A0B7E07C4FEDD97F36DE51A95B896EC54181CC7F89A071CAD4758" +
            "4C31D69B8D98ECBA12E576A70AC97BF9F9FC773A2600E29DF2E9C4B9B225" +
            "B0E96CC8E132DF202996D662E034E2F494E4A1D59CCE9D7244752A6AC026" +
            "4CB0E43C555FF225D4EE105DABBB249C0403B27170DCCE6B3E7EB2F30A06" +
            "8B2E06BCFFB4DE84FDC5D2F162C092E1D0DFAB3721C565E39A62F1765203" +
            "3AE7DBC1FB90A9E0F4E6CECE193D8D3D63CBA1CD07FB1C30D901A2FF1602" +
            "72ED29D976D7AC80017531D8F735809DE39CA774621BF36230E5E300F14E" +
            "F960FEAE5061853721C565138F62F176200316E0ADADB760EEBF5283464F" +
            "E17C6D710C066C3DDB79F20B6A62E117F5EDC9EDC28BC4A00A96A358E2B1" +
            "29C8EDCA5F78B90A58F198D829664DFC64C7E59A7C50FCCA23CA5BFE083A" +
            "76A4382827FEDB8CA9DE5EC8CF4362B6C7525C0B2DBB0DC5820DC8EB87E3" +
            "F23911330B8E5F743E5B0E1603B688D90682DD0AA41B9AB7CD76130493DC" +
            "8F623C2135AFECF2133D81E2DF8FCD3329349B8F8CA8D642403C9873A7FD" +
            "531EFD3D3427CAB30D930D0B37B72E26D357C39FE91C4B34B8D59683C080" +
            "ED09C1407FFABCAA2D30B5DB8DCF01C9F875536BC0E65BB7237C3C9E78EB" +
            "E5240FEF2D99C953037A8883FA56CBE2834A5AC7BC31B9535C263EDBBB82" +
            "190BFDBC66E03C5A87B5B5CD4F392B1A43F4DB301DE7478F7F4ED57EA36B" +
            "B7694C784805F7A1318F173360BC65AE61FEEC1C735570829244C7EB3073" +
            "A6045151EEF52D872CC71026F7FB56935FDA70C765C3423FAF1938A7BABA" +
            "F0D2E6AA680C2703B68D08A4DEFB5955BBF133B4ADDDC608C2F467FBE92A" +
            "D136E5DDA8998FC3D07AEF32DE802F5E8630BB104D71D45FB4191084598B" +
            "B1D33963675C362C9C46BA5C197EEFA58BC4F62346499E5C8C02066C1B25" +
            "7E618EAADDF069D4E6A8E03836129071CFB9F97E8F7246650C78FC6155BF" +
            "018D2573620CEDD7DEFF980A4830DFE2F06809844258EA2F080CF828C64F" +
            "E7A45CDA70D9B1F0FAE0B8344E4B53A93E99CB4924C8734A0D98060BD79C" +
            "ABAAD7E79B37C018F0F7A8DF1A8901811C141BCF657A4D71C0A2D93C07E1" +
            "BD4585821B5E9B388A6762B3CD80B0F870EED7D339DD7149B0707353BD28" +
            "875C11B19F93A91764DC16AE802D51066E51D5EB3E956B139BB02840C4DD" +
            "F93ED5D55BAD8F70E24D9792D8D5EB2E55C130430683833520892923A96D" +
            "C0ADCAB3E68AB1D131D5BB57D133A0332E2916667E94835D3FBBF95AC568" +
            "35D2F83A19B7FD197078959AB8F693AA9A6913CF33065C9CEF535D653760" +
            "B80A8D594C85EB96F4376ECE6148F1786F305BF0F360EA1CDE9B1D2A2088" +
            "69CD15CDC9E0A2FC406576C6E58285F6CD73D0E399C3663BFB753E277A5D" +
            "2BE3566340007DCD02AB38F17397A8EA359FC835EF752CFC41152CFA78BE" +
            "DFA295965FDF67C6ACD09B568603842678ED39F6044839C9E731F1DC4143" +
            "25A7FDABD75CC2567029FE3616861BD080C5F2FA02A20FD6ABF5DAAA87BD" +
            "0A862B90C602FF687F14D722E4E710B0D57891F27D4C7D81B03B3A9FF7FE" +
            "4B081B5EA9261EEA992C7C888EF15EA34470FD52DC5C5F8AD9D1805D4E3A" +
            "D80C86227D614C26BF5C65790D8C2289CDF66B63A1068405908BC9639556" +
            "385B3FAA838C5BA6021E54FE6DE7A8F19FB9B7DAB37C45889FFD892CDE83" +
            "5BBA15453C468393C34263CA486A93BF45790538E9F1689AABDFD8E72892" +
            "EFEE0532CEB339DD66DB99F495307C10FB45C62D35E0EE9BD4F8BC8F3AB7" +
            "DA06D3760404CFFB983DE6033D03B604978CD161CDC59A7C74B4E16284E1" +
            "26EB31A3B7F2FBC173B1F28F99B7B1986073D380F07F6C85AA15E4CFAC25" +
            "9F47F800D657C62D530181E89F7E448D3B34EF55FB19A1F9EA7C7BCCE559" +
            "0342CCA340A208C77CF8E5018CF922E51FCDE38AD75F8430CC7733A0080B" +
            "C6219F2324396874587011E1A5B61E1700D8C916B8E998D55CA469E61116" +
            "C52ECBDBFCBDE0672E54E373CFD634D904F8EC102ED7C53B5B79AF688CFC" +
            "0A9897C57161CE64243632744CF271C901166617436F9C369FF92B44DF93" +
            "A667ABEC3C88AB6E3EB44FED1964C024514D68F1D3260DF55A8C3BE2A73C" +
            "C8B8EDEF318CF821AEA04296B1ACCF466D1781FCE76D03CADA8BAA3EE76C" +
            "3536E7C399365F05D3A04969C06920D9CD3853B958910175B9A6463DB654" +
            "051B7F9C69F7B87F0DA30097A5010B9076E60C55C09CFB7FA582FD47ECD5" +
            "15FAF9CFFE28D3FE50F008E186B134E0FBD98047972AFF991F765B882E65" +
            "ED857444455B7A7DBAFD5FD96537ED1470571A700A489C9115B1BA4E05FF" +
            "FC81F20BB55FAAA8EA56C98A7250FE7DC07EFFC0DD8C1C0F556DF3F795FF" +
            "64B1160E4D5F524C053CA1A2B777ABD0A1456382D532B6D71EF3F0897CD9" +
            "6F1C1262D9ABA2C318F3808A1B08D7F1018461AFC341DB9D973687F23912" +
            "928346877707082FD1F109B4654EA8F8CD3F2A7FDDF75443DC7EAD8249DD" +
            "459A66760F8A5D963735E0F87E15EE03A11C5B54A91BCE0C27552C893774" +
            "32134338461737176BD288953D28AFFD0E06EC07CF1E15FB96454AB0B969" +
            "40F8F7A10A66B809365ED0FAB3C26CDBB836C78B594B3E8F6408E39571CB" +
            "54C0BA8A0FEF51D181026D4463C2917764F18E8DF40C281DA3C3C9612131" +
            "DF7330E0888A8B70D21D639AABDFD8A015CEB7FA9E8CF36C4E878F16BE78" +
            "24C7B05F64DCF29790D1032A1ADC37D9067991727D4C7D3D950C753ECFFC" +
            "77D4CB271B0CABE8C4A9EE7BC9093A26AED1D5C7F54BB1737D9B04B32637" +
            "F662724AC55D4E3AD80EA824602A02CB0D8CC9E4973BB4D706337C1B62B3" +
            "FD26FB633E9B18AF06EB145DC2A80E326E35B7E071150F0DD9ABC3181037" +
            "042B2DD3E23A16A4A192E17C9F68B8625969CC9893E3FA31221C808BF413" +
            "E4D81528E504E53AD630E441FB4743200A675829FE1616861BC0157358EA" +
            "43441FAC57EBB5550FFB193F39893496F807F2E9EF314C085BC7F060BE55" +
            "7D248AAF920AEA335AB596FA64148DA964B667BC6ABD618223F118D2AA43" +
            "A8DFB07D9175E7AAC242B2E58AE66470C5841F18E38CCB050BD317E7D179" +
            "CD61135648A2D7B08CDBFE0CD8A8A8B87234DF6A013257A09211D467AC66" +
            "3560B306A6C2B12B9054C8188BC191E01B704A2489597130600D16922D57" +
            "8C8D1933D23B6674B761675C522C4C3FC2A96B4E9A4BC8188A539171DB97" +
            "0193DA09158FE65B426E7B60C031D4AF66BA317712ACAB04C5C673995E53" +
            "1C10B70E0B261773D4C180144F52C78B0D8BC3E43ECACCE98C4B82859B9B" +
            "EA453904E344F62D173F78A65E90715BDC800D98607C0435782FC6E04395" +
            "D450BF097401D195F91856309903CFA979CD9DB50866C3B9926082CB14C6" +
            "D208ED957C02E363E674C625C0C2EA83E766E2A4394AF5C97224C993D1B9" +
            "800123A82463602AA679F8FC979A91E9CFF6D3AD3A3843727319DF03A2C9" +
            "4280F8FE388A25380A74496370F8909BED8CE4619E98399D71D9B0701AE9" +
            "726562D5253B143EEF0AF2ECDB803E54244FD774A0818C061AE3DBB62E8D" +
            "198309C3FCD939E0A7029C01C97817A20133CE3DB01990C9BDC1DC9C9D71" +
            "D9B0309F9B3827BABAF032A955C118820A18838152E14DCDB4A532E38382" +
            "06D4559B08B6921C3E0D9E10CE4E288FA6F8BC938A8A3808213753058CA1" +
            "BAE0315CEE0497ED8862C3C27C6EE2BCB500B2CD363F53200AC6301B3004" +
            "611B209AA9D9446882483E8A1108CE4EB6AD2DFB7980E2FBDC51004863F3" +
            "313DCBCB120D98310FC6DC99FE309E353CC165C364C3C2CDAD8BC9F4D5F1" +
            "67D2046BC0557AA72D3825C537B4506A22306080E288C7A6068079004733" +
            "E1B66526768A59173FAD949A9C72F113B841928AD1C681C7B35B7D0C0F8B" +
            "B979745588E0121810C78F503566F5C371F99C121C4B520CC87CB61CDA7A" +
            "F215D02094565C0C3215D16460F299A66A39C701B372E7BF149F292FABB9" +
            "346239E4C82F220E97868B2EC70C166C9A02BCB5F512CCCD1E23F0E2B4C5" +
            "D11A100ECE2CA9B280DD331110C0C7D155214DFC48534D74C287868B415F" +
            "B174BC58768A499CDA45D15A1838471BD7144B93AB5AAEF97638359ED17B" +
            "73E7E62415D09683CE80E9C139004272ADC0A52135208983E3665FEBCE6D" +
            "A631E83393F93A42E730C19CA948069C3DE300F14EF9F4E26A2B5FA7A211" +
            "0C36F12816D680697C564F09A73ACD99B9D31D27C4316D39E80CC808625C" +
            "BDBAF3C154195042A2CD7856A1E1FC44084C09CD92A8115D9BA74C80D68E" +
            "713A0DD8C9DD60C4D602C96230F2295D88B2FCE9199033207B0190FEB886" +
            "11576A18C9E177DAFA38549D69C324D540D80FB4D756D26E4EA7DB80717B" +
            "EBCC35E3AD283598698B6E1B301FCFF6005748D8B40ADD265ECECB4CCC61" +
            "AA30315C606E84CF7A69054CE0864582695646BA623A7D7537CFF4396069" +
            "40FB8FECA675314D9111B962D5F58EECDEC07F2B0E02ABC8B1C5B1E6AB54" +
            "4D1A4BDB77FABE8DE53E13E4873981A349F90F1860FD128AA8D17C2D9331" +
            "8DC4903A3DE0171272E29506148933AB3AC1652CBF40850634FE2C33829B" +
            "4C286DB292DB4C6346B3F40C88B988E16FAACCC6ADF374614E9F08743992" +
            "F941F0CB08E915DD60C25969A62267205834A501ED0BAE6BC2A934E0E95A" +
            "31B32A2E5CA65AD53BD3CA0A6837A4456359059C55462952DD1CC7A4B7BF" +
            "D488A5014B03966730C7C533C38A4959016798201FB405551AB03460DFDB" +
            "683F8BA6346069C0D280FDACA0726C79063CA32BA834E0EC36E0FF01BBDD" +
            "E06F371676E90000000049454E44AE426082')")
    db.Exec(
            "create table element_type(" +
            " id integer not null primary key, " +
            " type integer, " +
            " ticket_property integer not null default 0, " +
            " reply_property integer not null default 0, " +
            " required integer not null default 0, " +
            " name text, " +
            " description text, " +
            " default_value text, " +
            " auto_add_item integer not null default 0, " +
            " display_in_list integer not null default 0, " +
            " sort integer, " +
            " deleted integer not null default 0 " +
            ");")
    db.Exec(
            "create index index_element_type_0 on element_type (id, type, display_in_list, sort)")
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (1, ?, 1, 0, 1, ?, ?, 0, '', 1, 1);", 
            ELEM_TYPE_TEXT, gettext.Gettext("title"), gettext.Gettext("please input title that means content correctly."))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (2, ?, 0, 0, 1, ?, ?, 0, '', 1, 2);",
            ELEM_TYPE_TEXT, gettext.Gettext("registerer"), gettext.Gettext("please input your name."))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (3, ?, 1, 0, 1, ?, ?, 0, ?, 1, 3);", 
            ELEM_TYPE_LIST_SINGLE_RADIO, gettext.Gettext("status"), gettext.Gettext("please select status."), gettext.Gettext("new"))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (4, ?, 1, 0, 0, ?, ?, 1, '', 1, 4);", 
            ELEM_TYPE_LIST_MULTI, gettext.Gettext("category"), gettext.Gettext("please select category."))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (5, ?, 1, 0, 0, ?, ?, 0, '', 1, 5);", 
            ELEM_TYPE_LIST_SINGLE, gettext.Gettext("priority"), gettext.Gettext("please select priority."))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (6, ?, 1, 0, 1, ?, ?, 0, '', 0, 6);",
            ELEM_TYPE_TEXTAREA, gettext.Gettext("detail"), gettext.Gettext("please describe the detail."))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (7, ?, 1, 0, 0, ?, ?, 0, '', 0, 7);", 
            ELEM_TYPE_TEXTAREA, gettext.Gettext("reproduction procedure"), gettext.Gettext("please describe the condition and the procedure to reproduce the problem."))
    db.Exec(
            "insert into element_type(id, type, ticket_property, reply_property, required, name, description, auto_add_item, default_value, display_in_list, sort) " +
            "values (8, ?, 0, 1, 0, ?, ?, 0, '', 0, 8);", 
            ELEM_TYPE_TEXTAREA, gettext.Gettext("comment"), gettext.Gettext("please describe the comment."))
    db.Exec(
            "create table list_item( " +
            " id integer not null primary key, " +
            " element_type_id integer not null default 0, " +
            " name text, " +
            " close integer not null default 0, " +
            " sort integer " +
            ");")
    db.Exec(
            "create index index_list_item_0 on list_item (id, sort)")
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (1, 3, ?, 0, 1);",
            gettext.Gettext("new"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (2, 3, ?, 0, 2);",
            gettext.Gettext("accepted"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (3, 3, ?, 0, 3);",
            gettext.Gettext("fixed"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (4, 3, ?, 0, 4);",
            gettext.Gettext("reservation"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (5, 3, ?, 1, 5);",
            gettext.Gettext("complete"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (6, 3, ?, 1, 6);",
            gettext.Gettext("wontfix"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (7, 3, ?, 1, 7);",
            gettext.Gettext("specification"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (8, 3, ?, 1, 8);",
            gettext.Gettext("deleted"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (9, 4, ?, 0, 1);",
            gettext.Gettext("page"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (10, 4, ?, 0, 2);",
            gettext.Gettext("batch process"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (11, 4, ?, 0, 3);",
            gettext.Gettext("document"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (12, 5, ?, 0, 1);",
            gettext.Gettext("emergency"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (13, 5, ?, 0, 2);",
            gettext.Gettext("high"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (14, 5, ?, 0, 3);",
            gettext.Gettext("middle"))
    db.Exec("insert into list_item(id, element_type_id, name, close, sort) values (15, 5, ?, 0, 4);",
            gettext.Gettext("low"))
    db.Exec(
            "create table ticket(" +
            " id integer not null primary key, " +
            " original_message_id integer not null default 0, " +
            " last_message_id integer not null default 0, " +
            " registerdate text not null, " +
            " closed integer not null default 0" +
            ");")
    db.Exec(
            "create index index_ticket_0 on ticket (id, last_message_id, last_message_id, closed)")
    db.Exec(
            "create table message(" +
            " id integer not null primary key, " +
            " ticket_id integer not null, " +
            " registerdate text, " +
            " field1 text not null default '' , " +
            " field2 text not null default '' , " +
            " field3 text not null default '' , " +
            " field4 text not null default '' , " +
            " field5 text not null default '' , " +
            " field6 text not null default '' , " +
            " field7 text not null default '' , " +
            " field8 text not null default ''  " +
            ");")
    db.Exec(
            "create table element_file(" +
            " id integer not null primary key, " +
            " message_id integer not null, " +
            " element_type_id integer not null, " +
            " filename text, " +
            " size integer, " +
            " mime_type text, " +
            " content blob, " +
            " deleted integer " +
            ");")
    db.Exec(
            "create table wiki( " +
            " id integer not null primary key, " +
            " name text, " +
            " content text, " +
            " registerdate text " +
            ");")
    db.Exec(
            "insert into wiki(id, name, content, registerdate) values (NULL, 'top', ?, current_timestamp);",
            gettext.Gettext("[wiki syntax example]"))
    db.Exec(
            "insert into wiki(id, name, content, registerdate) values (NULL, 'help', ?, current_timestamp);",
            gettext.Gettext("[help content]"))
    db.Exec(
            "insert into wiki(id, name, content, registerdate) values (NULL, 'adminhelp', ?, current_timestamp);",
            gettext.Gettext("[admin help content]"))
    tx.Commit()
}
/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
