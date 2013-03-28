package handler

import (
    logger "code.google.com/p/log4go"
    "fmt"
    "net/http"
    "../data"
    "strconv"
)

func RegisterRoutesProject() {
    RegisterRoute("^/([^/]+)$", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        logger.Debug("project: %s", projectName)

        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)
        params["wiki"] = data.GetWiki(projectName, "top")
        params["newestTickets"] = data.GetNewestTickets(projectName, 10)
        params["states"] = data.GetStates(projectName, false)

        TmplProject(w, "index", params)
    })

    RegisterRoute("^/([^/]+)/list", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        logger.Debug("project list: %s", projectName)

        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)
        params["newestTickets"] = data.GetNewestTickets(projectName, 10)
        params["states"] = data.GetStates(projectName, false)

        notClosedStates := data.GetStates(projectName, true)
        ticketsByStatus := []interface{}{}
        for _, s := range notClosedStates {
            statusResult := map[string]interface{}{}
            statusResult["State"] = s
            statusResult["SearchResult"] = data.GetTicketsByStatus(projectName, s.Name)
            ticketsByStatus = append(ticketsByStatus, statusResult)
        }
        params["notClosedStates"] = notClosedStates
        params["ticketsByStatus"] = ticketsByStatus
        params["elementTypes"] = data.GetElementTypes(projectName)
        TmplProject(w, "list", params)
    })

    RegisterRoute("^/([^/]+)/search", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        logger.Debug("project search %s", projectName)

        r.ParseForm()
        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)
        params["newestTickets"] = data.GetNewestTickets(projectName, 10)
        params["states"] = data.GetStates(projectName, false)

        elementTypes := data.GetElementTypes(projectName)
        params["searchResult"] = data.SearchTickets(projectName, r.Form, r.Cookies(), elementTypes)
        params["elementTypes"] = elementTypes
        TmplProject(w, "search", params)
    })

    RegisterRoute("^/([^/]+)/setting_file/([^/]+)", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        name := captures[1]

        file := data.GetSettingFile(projectName, name)

        w.Header().Set("Content-Type", file.MimeType)
        w.Header().Set("Content-Length", strconv.Itoa(file.Size))
        w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", 60 * 24 * 1))

        w.Write([]byte(file.Content))
    })

    RegisterRoute("^/([^/]+)/ticket/(\\d+)", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        ticketId, _ := strconv.Atoi(captures[1])
        logger.Debug("ticket %d", ticketId)

        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)
        params["newestTickets"] = data.GetNewestTickets(projectName, 10)
        params["states"] = data.GetStates(projectName, false)

        elementTypes := data.GetElementTypes(projectName)
        params["ticket"] = data.GetTicket(projectName, ticketId, elementTypes)
        params["elementTypes"] = elementTypes
        TmplProject(w, "ticket", params)
    })

    RegisterRoute("^/([^/]+)/register$", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        logger.Debug("project: %s", projectName)

        project := data.GetProject(projectName)
        elementTypes := data.GetElementTypes(projectName)
        params := make(map[string]interface{})

        if r.Method == "POST" {
            logger.Debug("post")
            r.ParseMultipartForm(int64(project.UploadMaxSize))

            errors := data.ValidateTicket(projectName, r.Form, elementTypes)
            if len(errors) == 0 {
                data.RegisterTicket(projectName, r.Form, elementTypes)
                http.Redirect(w, r, fmt.Sprintf("/%s/list", projectName), 302)
                return
            }
            logger.Warn("validate failed.")
            params["params"] = r.Form
            params["errors"] = errors
        }
        params["topProject"] = data.GetProject("manage")
        params["project"] = project
        params["newestTickets"] = data.GetNewestTickets(projectName, 10)
        params["states"] = data.GetStates(projectName, false)
        params["elementTypes"] = elementTypes

        TmplProject(w, "register", params)
    })
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
