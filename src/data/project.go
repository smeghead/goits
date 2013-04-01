package data

import (
    logger "code.google.com/p/log4go"
    _ "github.com/gosexy/gettext"
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "strconv"
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
        return ProjectInfo{id, name, sort, GetProject(name)}
    })
    if err != nil {
        logger.Error(err)
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
        logger.Error(err)
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
        logger.Error("wiki got")
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
        logger.Error(err)
        panic(err)
    }
    return wiki
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
        logger.Error(err)
        panic(err)
    }
    return settingFile
}


/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
