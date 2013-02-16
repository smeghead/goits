package handler

import (
    "fmt"
    "net/http"
    "./data"
)

func RegisterRoutesIndex() {
    RegisterRoute("^/$", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("index")
        params := make(map[string]interface{})
        params["project"] = data.GetProject("manage")
        params["wiki"] = data.GetWiki("manage", "top")
        params["projectInfos"] = data.GetProjectInfos()

        Tmpl(w, "index", params)
        fmt.Println("index end")
    })
//void top_top_action()
//{
//    Database* db_a;
//    List* project_infos_a;
//    Iterator* it;
//    Project* top_project_a = project_new();
//
//    list_alloc(project_infos_a, ProjectInfo, project_info_new, project_info_free);
//
//    db_a = db_init(g_project_code);
//    project_infos_a = db_top_get_all_project_infos(db_a, project_infos_a);
//    top_project_a = db_get_project(db_a, top_project_a);
//    top_output_header(_("top page"), top_project_a);
//    o(      "<div id=\"project_list\">\n"
//            "\t<h2 class=\"radius-right\">%s</h2>\n"
//            "\t<ul>\n", _("sub projects list"));
//    foreach (it, project_infos_a) {
//        ProjectInfo* p = it->element;
//        Database* db_project_a;
//        Project* project_a;
//        if (p->id == 1) {
//            /* idが1のプロジェクトはトップなので、表示しない。 */
//            continue;
//        }
//        if (p->deleted) {
//            continue;
//        }
//        project_a = project_new();
//        db_project_a = db_init(string_rawstr(p->code));
//        project_a = db_get_project(db_project_a, project_a);
//        o(      "\t\t\t\t<li><a href=\"%s/", cgiScriptName); u(string_rawstr(p->code)); o("\">"); h(string_rawstr(project_a->name)); o("</a></li>\n");
//        project_free(project_a);
//        db_finish(db_project_a);
//    }
//    list_free(project_infos_a);
//    o(      "\t</ul>\n");
//    o(      "\t<h2 class=\"radius-right\">%s</h2>\n"
//            "\t<form action=\"%s/%s/top_search\" method=\"get\">\n", _("search"), cgiScriptName, g_project_code_4_url);
//    o(      "\t\t<div>%s</div>\n", _("search by keyword for all sub projects."));
//    o(      "\t\t<input type=\"text\" name=\"q\" />\n"
//            "\t\t<input type=\"submit\" value=\"%s\" />\n"
//            "\t</form>\n", _("search"));
//
//    o(      "</div>\n");
//    o(      "<div id=\"dashboard\">\n"
//            "\t<h2 class=\"radius-left\">%s</h2>\n", _("description"));
//    o(      "<div id=\"top_page_edit_link\"><a href=\"%s/%s/top_edit_top\">%s</a></div>\n", cgiScriptName, g_project_code_4_url, _("edit top page"));
//    wiki_out(db_a, "top");
//    o(      "</div>\n");
//    o(      "<br clear=\"all\" />\n");
//    project_free(top_project_a);
//    top_output_footer();
//    db_finish(db_a);
//}
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
