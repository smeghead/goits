package handler

import (
    "fmt"
    "net/http"
    "./data"
)

func RegisterRoutesProject() {
    RegisterRoute("^/([^/]+)$", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        fmt.Println("project", projectName)

        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)
        params["wiki"] = data.GetWiki(projectName, "top")
        params["newestTickets"] = data.GetNewestTickets(projectName, 10)

        TmplProject(w, "project", params)
    })

    RegisterRoute("^/([^/]+)/list", func(w http.ResponseWriter, r *http.Request, captures []string) {
        projectName := captures[0]
        fmt.Println("project list", projectName)

        params := make(map[string]interface{})
        params["topProject"] = data.GetProject("manage")
        params["project"] = data.GetProject(projectName)

        TmplProject(w, "project_list", params)

//void list_action()
//{
//    List* element_types_a;
//    List* conditions_a = NULL;
//    Project* project_a = project_new();
//    List* states_a;
//    Iterator* it;
//    char** multi;
//    Database* db_a;
//
//    db_a = db_init(g_project_code);
//    project_a = db_get_project(db_a, project_a);
//    output_header(project_a, _("ticket list by status"), NULL, NAVI_LIST);
//    /* retrieve messages. */
//    if ((cgiFormStringMultiple("message", &multi)) != cgiFormNotFound) {
//        int i = 0;
//        o("<div class=\"complete_message\">");
//        while (multi[i]) {
//            if (strstr(multi[i], "[ERROR]") != NULL) {
//                o("<div class=\"error\">");
//            } else {
//                o("<div>");
//            }
//            h(multi[i]);
//            o("</div>");
//            i++;
//        }
//        o("</div>\n");
//    }
//    cgiStringArrayFree(multi);
//    list_alloc(element_types_a, ElementType, element_type_new, element_type_free);
//    element_types_a = db_get_element_types_4_list(db_a, NULL, element_types_a);
//    o("<h2>"); h(string_rawstr(project_a->name)); o(" - %s</h2>\n", _("ticket list by status"));
//    project_free(project_a);
//    list_alloc(states_a, State, state_new, state_free);
//    states_a = db_get_states(db_a, states_a);
//    output_states(states_a, true);
//    list_free(states_a);
//    fflush(cgiOut);
//    o("<div id=\"ticket_list\">\n"
//      "<h3>%s</h3>\n"
//      "<div class=\"description\">%s\n", _("ticket list by status"), _("display unclosed tickets."));
//    list_alloc(states_a, State, state_new, state_free);
//    states_a = db_get_states_has_not_close(db_a, states_a);
//    foreach (it, states_a) {
//        State* s = it->element;
//        o("\t\t\t<a href=\"#state%d""\" title=\"%s\">", s->id, _("link to part of this page"));
//        hs(s->name);
//        o("</a>&nbsp;\n");
//    }
//    o("</div>\n");
//    foreach (it, states_a) {
//        State* s = it->element;
//        SearchResult* result_a = search_result_new();
//
//        /* search */
//        list_alloc(conditions_a, Condition, condition_new, condition_free);
//        result_a = db_get_tickets_by_status(db_a, string_rawstr(s->name), result_a);
//        list_free(conditions_a);
//
//        o("<a id=\"state%d\" name=\"state%d\"></a>\n", s->id, s->id);
//        o("<div>\n");
//        o("<h4 class=\"status\">");hs(s->name);o("&nbsp;(%d件)&nbsp;<a href=\"#top\">↑</a></h4>\n", s->count);
//        if (result_a->hit_count == LIST_COUNT_PER_LIST_PAGE) {
//            o("\t\t<div class=\"information\">%s%d%s<a href=\"%s/%s/search?field%d=", _("new"), 
//                    result_a->hit_count, _("tickets is desplayed."),
//                    cgiScriptName,
//                    g_project_code_4_url,
//                    ELEM_ID_STATUS);
//            us(s->name);
//            o("\">%s", _("display tickets that status is ")); hs(s->name); o("%s</a></div>\n", _(".(display status)"));
//        }
//        output_ticket_table_status_index(db_a, result_a, element_types_a);
//        if (result_a->hit_count == LIST_COUNT_PER_LIST_PAGE) {
//            o("\t\t<div class=\"information\">%s<a href=\"%s/%s/search?field%d=", _("see more..."), 
//                    cgiScriptName,
//                    g_project_code_4_url,
//                    ELEM_ID_STATUS);
//            us(s->name);
//            o("\">%s", _("display tickets that status is ")); hs(s->name); o("%s</a></div>\n", _(".(display status)"));
//        }
//        search_result_free(result_a);
//        o("</div>\n");
//        fflush(cgiOut);
//    }
//    list_free(states_a);
//    list_free(element_types_a);
//    o("</div>\n");
//    output_footer();
//    db_finish(db_a);
//}
    })

}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
