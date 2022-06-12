package data

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	logger "github.com/alecthomas/log4go"
	"github.com/chai2010/gettext-go"
	_ "github.com/mattn/go-sqlite3"
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
	for i, p := range results {
		projectInfos[i] = p.(ProjectInfo)
	}
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
	params := []interface{}{wikiName}

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

func ValidateSubProject(projectName string, form url.Values) map[string]string {
	logger.Debug("======ValidateSubProject: [%s]\n", form.Encode())

	messages := make(map[string]string)

	name := form.Get("name")
	logger.Debug("name: %s", name)
	if len(strings.TrimSpace(name)) == 0 {
		messages["name"] = gettext.Gettext("it will required. please describe.")
	}
	size := form.Get("upload_max_size")
	logger.Debug("upload_max_size: %s", size)
	if len(strings.TrimSpace(size)) == 0 {
		messages["upload_max_size"] = gettext.Gettext("it will required. please describe.")
	} else if _, err := strconv.Atoi(strings.TrimSpace(size)); err != nil {
		messages["upload_max_size"] = gettext.Gettext("it must be number.")
	}
	return messages
}

func RegisterSubProject(projectName string, form url.Values) error {
	logger.Debug("RegisterSubProject")
	err := tran(projectName, func(db *sql.DB, tx *sql.Tx) error {
		logger.Debug("RegisterSubProject name: %s", form["name"])
		result, err := db.Exec(
			"update setting set value = ? where name = 'project_name'",
			form.Get("name"))
		if err != nil {
			logger.Error("RegisterSubProject err: %s", err)
			logger.Error("RegisterSubProject name update failed: %s", form["name"])
			panic("update setting failed")
		}
		if count, _ := result.RowsAffected(); count != 1 {
			logger.Error("RegisterSubProject name update failed: %s", form["name"])
			panic("update setting failed. afeected rows is not 1.")
		}

		result, err = db.Exec(
			"update setting set value = ? where name = 'upload_max_size'",
			form.Get("upload_max_size"))
		if err != nil {
			logger.Error("RegisterSubProject size update failed: %s", form["upload_max_size"])
			panic("update setting failed")
		}
		if count, _ := result.RowsAffected(); count != 1 {
			logger.Error("RegisterSubProject size update failed: %s", form["upload_max_size"])
			panic("update setting failed. afeected rows is not 1.")
		}
		return nil
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return nil
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
