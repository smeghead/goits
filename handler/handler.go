package handler

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	logger "github.com/alecthomas/log4go"
	"github.com/chai2010/gettext-go"
	"github.com/knieriem/markdown"
	"github.com/smeghead/goits/data"
)

type Route struct {
	pattern *regexp.Regexp
	h       func(http.ResponseWriter, *http.Request, []string)
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
	t, _ := template.New("layout.tmpl").
		Funcs(getFuncs()).
		ParseFiles("template/layout.tmpl", fmt.Sprintf("template/%s.tmpl", templateName))
	t.Execute(w, params)
}

func TmplProject(w http.ResponseWriter, templateName string, params map[string]interface{}) {
	t, err := template.New("layout.tmpl").
		Funcs(getFuncs()).
		ParseFiles("template/project/layout.tmpl", fmt.Sprintf("template/project/%s.tmpl", templateName))
	logger.Debug(err)
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
		"json": func(o interface{}) string {
			str, err := json.Marshal(o)
			if err != nil {
				return ""
			}
			return string(str)
		},
		"_": func(messageId string) string {
			return gettext.Gettext(messageId)
		},
		"_f": func(messageId string, args ...interface{}) string {
			return fmt.Sprintf(gettext.Gettext(messageId), args...)
		},
		"defferelementwith": func(element data.Element, messages []data.Message, messageIndex int) bool {
			logger.Debug(element)
			if messageIndex == 0 {
				return element.StrVal != ""
			}
			for _, e := range messages[messageIndex-1].Elements {
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
			if _, ok := params[fieldName]; ok {
				return params.Get(fieldName)
			} else {
				for _, e := range elements {
					if e.ElementType.Id == elementType.Id {
						return e.StrVal
					}
				}
			}
			return ""
		},
		"getvalues": func(elementType data.ElementType, params url.Values, elements []data.Element) []string {
			fieldName := fmt.Sprintf("field%d", elementType.Id)
			if _, ok := params[fieldName]; ok {
				return params[fieldName]
			} else {
				for _, e := range elements {
					if e.ElementType.Id == elementType.Id {
						return strings.Split(e.StrVal, "\t")
					}
				}
			}
			return []string{""}
		},
		"markdown": func(content interface{}) template.HTML {
			str := html.EscapeString(content.(string))
			p := markdown.NewParser(&markdown.Extensions{Smart: true})
			buf := bytes.NewBufferString("")
			w := bufio.NewWriter(buf)
			p.Markdown(bufio.NewReader(bytes.NewBufferString(str)), markdown.ToHTML(w))
			w.Flush()
			return template.HTML(buf.String())
		},
	}
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
