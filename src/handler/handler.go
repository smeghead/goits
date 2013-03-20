package handler

import (
    logger "code.google.com/p/log4go"
    "github.com/gosexy/gettext"
    "os"
    "fmt"
    "net/http"
    "net/url"
    "html/template"
    "regexp"
    "math"
    "../data"
)

type Route struct {
    pattern *regexp.Regexp
    h func(http.ResponseWriter, *http.Request, []string)
}

var _routes []Route

func RegisterRoute(patternString string, handler func(http.ResponseWriter, *http.Request, []string)) {
    logger.Debug("RegisterRoute: %s", patternString)
    _routes = append(_routes, Route{regexp.MustCompile(patternString), handler})
}

func InitRoutes() {
    RegisterRoutesIndex()
    RegisterRoutesProject()
    // ...
}

func RouteHandler(w http.ResponseWriter, r *http.Request) {
    path := r.RequestURI
    if path == "/favicon.ico" {
        return
    }

    gettext.BindTextdomain("multivac", "locale")
    gettext.Textdomain("multivac")
    os.Setenv("LANGUAGE", "ja_JP.utf8")
    gettext.SetLocale(gettext.LC_ALL, "")

    for _, route := range _routes {
        matches := route.pattern.FindStringSubmatch(r.URL.RequestURI())
        if len(matches) > 0 {
            logger.Debug("------------hit pattern:%s path'%s", route.pattern, path)
            //get submatch
            captures := []string{}
            matchLen := len(matches)
            if matchLen > 1 {
                for i := 1; i < matchLen; i++ {
                    matchString, _ := url.QueryUnescape(matches[i])
                    captures = append(captures, matchString)
                }
            }
            route.h(w, r, captures)
            return
        }
    }
    logger.Error("404: %s", path)
    http.NotFound(w, r)
}

func TmplTop(w http.ResponseWriter, templateName string, params map[string]interface{}) {
    t, _ := template.New("layout_top.tmpl").
        Funcs(getFuncs()).
        ParseFiles("template/layout_top.tmpl", fmt.Sprintf("template/%s.tmpl", templateName))
    t.Execute(w, params)
}

func TmplProject(w http.ResponseWriter, templateName string, params map[string]interface{}) {
    t, _ := template.New("layout_project.tmpl").
        Funcs(getFuncs()).
        ParseFiles("template/layout_project.tmpl", fmt.Sprintf("template/%s.tmpl", templateName))
    t.Execute(w, params)
}

func getFuncs() template.FuncMap {
    return template.FuncMap{
        "eq": func(a, b interface{}) bool {
            return a == b
        },
        "ne": func(a, b interface{}) bool {
            return a != b
        },
        "contains": func(a string, arr []string) bool {
            for _, e := range arr {
                if e == a {
                    return true
                }
            }
            return false
        },
        "odd": func(a int) bool {
            return math.Mod(float64(a), 2) != 0
        },
        "even": func(a int) bool {
            return math.Mod(float64(a), 2) == 0
        },
        "inc": func(a int) int {
            return a + 1
        },
        "_": func(messageId string) string {
            return gettext.Gettext(messageId)
        },
        "_f": func(messageId string, args ...interface{}) string {
            return gettext.Sprintf(gettext.Gettext(messageId), args...)
        },
        "defferelementwith": func(element data.Element, messages []data.Message, messageIndex int) bool {
            logger.Debug(element)
            if messageIndex == 0 {
                return element.StrVal != ""
            }
            for _, e := range messages[messageIndex - 1].Elements {
                if e.ElementType.Id == element.ElementType.Id {
                    return e.StrVal != element.StrVal
                }
            }
            return true
        },
        "haserror": func(errors map[string]string, elementType data.ElementType) bool {
            if len(errors) == 0 {
                return false
            }
            _, exists := errors[fmt.Sprintf("field%d", elementType.Id)]
            return exists
        },
        "geterror": func(errors map[string]string, elementType data.ElementType) string {
            if len(errors) == 0 {
                return ""
            }
            error, _ := errors[fmt.Sprintf("field%d", elementType.Id)]
            return error
        },
        "getvalue": func(elementType data.ElementType, params url.Values, elements []data.Element) string {
            fieldName := fmt.Sprintf("field%d", elementType.Id)
            paramVal := params.Get(fieldName)
            return paramVal
        },
        "getvalues": func(elementType data.ElementType, params url.Values, elements []data.Element) []string {
            fieldName := fmt.Sprintf("field%d", elementType.Id)
            return params[fieldName]
        },
    }
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
