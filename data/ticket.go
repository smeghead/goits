package data

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	logger "code.google.com/p/log4go"
	"github.com/gosexy/gettext"
	_ "github.com/mattn/go-sqlite3"
)

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
				logger.Error(err)
				panic(err)
			}
			listItems = make([]ListItem, len(results))
			for i, p := range results {
				listItems[i] = p.(ListItem)
			}
		}
		return ElementType{id, element_type, ticket_property, reply_property, required, name,
			description, auto_add_item, default_value, display_in_list, sort, listItems}
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	elementTypes := make([]ElementType, len(results))
	for i, p := range results {
		elementTypes[i] = p.(ElementType)
	}
	return elementTypes
}

func createColumnsExp(elementTypes []ElementType, table_name string) string {
	columns := []string{}
	for _, elem := range elementTypes {
		columns = append(columns, fmt.Sprintf("%s.field%d", table_name, elem.Id))
	}
	return strings.Join(columns, ", ")
}

func createColumnsLikeExp(elementTypes []ElementType, table_name string, keywords []string) string {
	columns := []string{}
	for _, s := range keywords {
		logger.Debug(s)
		for _, elem := range elementTypes {
			columns = append(columns, fmt.Sprintf("%s.field%d like '%%' || ? || '%%' ", table_name, elem.Id))
		}
	}
	return strings.Join(columns, ", ")
}

func getElementFile(projectName string, messageId int, elementType ElementType) ElementFile {
	elementFile := ElementFile{}
	statement :=
		"select id, filename, mime_type, size, deleted from element_file as ef " +
			"where ef.message_id = ? and ef.element_type_id = ?"
	params := []interface{}{messageId, elementType.Id}

	_, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
		logger.Debug("elementFile got")
		var id int
		var filename string
		var mime_type string
		var size int
		var deleted bool
		rows.Scan(&id, &filename, &mime_type, &size, &deleted)
		elementFile.Id = id
		elementFile.Filename = filename
		elementFile.MimeType = mime_type
		elementFile.Size = size
		elementFile.Deleted = deleted
		return nil
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return elementFile
}
func getLastMessage(projectName string, ticketId int, elementTypes []ElementType, forList bool) Message {
	logger.Debug("getLastMessage ticket id", ticketId)
	statement := fmt.Sprintf(
		"select last_m.id, t.id, org_m.field%d, %s "+
			"  , substr(t.registerdate, 1, 16), substr(last_m.registerdate, 1, 16),  "+
			"  julianday(current_date) - julianday(date(last_m.registerdate)) as passed_date "+
			"from ticket as t "+
			"inner join message as last_m on last_m.id = t.last_message_id "+
			"inner join message as org_m on org_m.id = t.original_message_id "+
			"where t.id = ?", ELEM_ID_SENDER, createColumnsExp(elementTypes, "last_m"))
	params := []interface{}{ticketId}

	var messageId int
	elements := []Element{}
	_, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
		logger.Debug("each ticket got")
		pockets, _ := scanDynamicRows(rows)

		logger.Debug(strings.Join(pockets, ","))

		i := 0
		messageId, _ = strconv.Atoi(pockets[i])
		i++
		if forList {
			/* ID */
			elements = append(elements, Element{ELEMENT_TYPE_ID, pockets[i], ElementFile{}})
			i++
			/* 初回投稿者 */
			elements = append(elements, Element{ELEMENT_TYPE_ORG_SENDER, pockets[i], ElementFile{}})
			i++
		} else {
			i += 2
		}
		/* 動的カラム */
		for _, elmType := range elementTypes {
			elementFile := ElementFile{}
			if elmType.Id == ELEM_TYPE_UPLOADFILE {
				elementFile = getElementFile(projectName, messageId, elmType)
			}
			elements = append(elements, Element{elmType, pockets[i], elementFile})
			i++
		}
		if forList {
			/* 投稿日時 */
			elements = append(elements, Element{ELEMENT_TYPE_REGISTERDATE, pockets[i], ElementFile{}})
			i++
			/* 最終更新日時 */
			elements = append(elements, Element{ELEMENT_TYPE_LASTREGISTERDATE, pockets[i], ElementFile{}})
			i++
			/* 最終更新日時からの経過日数 */
			elements = append(elements, Element{ELEMENT_TYPE_LASTREGISTREDATE_PASSED, pockets[i], ElementFile{}})
			i++
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return Message{messageId, elements, ""}
}

func GetNewestTickets(projectName string, limit int) []Ticket {
	elementTypes := GetElementTypes(projectName)

	statement := fmt.Sprintf("select t.id "+
		"from ticket as t "+
		"inner join message as m on m.id = t.last_message_id "+
		"order by m.registerdate desc "+
		"limit %d ", limit)
	params := []interface{}{}

	results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
		logger.Debug("newest tickets got")
		var id int
		rows.Scan(&id)
		lastMessage := getLastMessage(projectName, id, elementTypes, false)
		return NewTicketWithoutMessages(id, lastMessage)
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}

	tickets := make([]Ticket, len(results))
	for i, p := range results {
		tickets[i] = p.(Ticket)
	}
	return tickets
}

func GetStates(projectName string, notClose bool) []State {
	condition := ""
	if notClose {
		condition = "where l.close = 0 "
	}
	statement := fmt.Sprintf(
		"select m.field%d as name, count(t.id) as count "+
			"from ticket as t "+
			"inner join message as m "+
			" on m.id = t.last_message_id "+
			"inner join list_item as l "+
			" on l.element_type_id = %d and l.name = m.field%d "+
			"%s "+
			"group by m.field%d "+
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
		logger.Error(err)
		panic(err)
	}
	states := make([]State, len(results))
	for i, p := range results {
		states[i] = p.(State)
	}
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

func getSearchSqlString(projectName string, conditions []Condition, sort Condition, keywords []string, elementTypes []ElementType) string {
	sqlString := fmt.Sprintf("select t.id as id, m.field%d as state "+
		"from ticket as t "+
		"inner join message as m on m.id = t.last_message_id "+
		"inner join message as org_m on org_m.id = t.original_message_id ",
		ELEM_ID_STATUS)

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
			logger.Debug("#######getSearchSqlString: %d", cond.ElementTypeId)
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
					cond.ValidValue()))
			}
		}
		sqlString += strings.Join(conds, " and ")
	}

	if len(keywords) > 0 {
		sqlString += fmt.Sprintf(" and (%s) ",
			createColumnsLikeExp(elementTypes, "m_all", keywords))
	}

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
	conditions := []Condition{Condition{ELEM_ID_STATUS, CONDITION_TYPE_NORMAL, status, ""}}
	sort := Condition{}

	elementTypes := GetElementTypes(projectName)

	statement := getSearchSqlString(projectName, conditions, sort, keywords, elementTypes)
	statement += " limit ? "

	params := []interface{}{status, LIST_COUNT_PER_LIST_PAGE}

	results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
		var id int
		var state string
		rows.Scan(&id, &state)
		logger.Debug("ticket id:", id)
		lastMessage := getLastMessage(projectName, id, elementTypes, true)
		logger.Debug("elements count:", len(lastMessage.Elements))
		return NewTicketWithoutMessages(id, lastMessage)
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	tickets := make([]Ticket, len(results))
	for i, p := range results {
		tickets[i] = p.(Ticket)
	}

	sums := []int{}
	//TODO 数値項目の合計
	//sums := set_tickets_number_sum(db, conditions, NULL, keywords_a, result);

	return SearchResult{len(tickets), 0, tickets, nil, sums}
}

func getMessages(projectName string, ticketId int, elementTypes []ElementType) []Message {
	statement := fmt.Sprintf(
		"select m.id, m.registerdate, %s from message as m where m.ticket_id = ? order by m.registerdate",
		createColumnsExp(elementTypes, "m"))

	params := []interface{}{ticketId}

	results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
		logger.Debug("each ticket got")
		pockets, _ := scanDynamicRows(rows)

		elements := []Element{}
		logger.Debug(strings.Join(pockets, ","))

		i := 0
		messageId, _ := strconv.Atoi(pockets[i])
		i++
		registerDate := pockets[i]
		i++

		/* 動的カラム */
		for _, elmType := range elementTypes {
			elements = append(elements, Element{elmType, pockets[i], ElementFile{}})
			i++
		}

		return Message{messageId, elements, registerDate}
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	messages := make([]Message, len(results))
	for i, p := range results {
		messages[i] = p.(Message)
	}

	return messages
}

func GetTicket(projectName string, ticketId int, elementTypes []ElementType) Ticket {
	lastMessage := getLastMessage(projectName, ticketId, elementTypes, false)
	messages := getMessages(projectName, ticketId, elementTypes)

	return NewTicket(ticketId, lastMessage, messages)
}

func getCookieValue(cookies []*http.Cookie, name string) string {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value
		}
	}
	return ""
}

func createCondition(name string, elementType ElementType, conditionType int, form url.Values, cookies []*http.Cookie) Condition {

	logger.Debug("======field: %s [%s]", name, form.Get(name))
	paramValue := form.Get(name)
	cookieValue := getCookieValue(cookies, name)
	return Condition{elementType.Id, conditionType, paramValue, cookieValue}
}

func createConditions(form url.Values, cookies []*http.Cookie, elementTypes []ElementType) []Condition {
	conditions := []Condition{}

	for _, elementType := range elementTypes {
		logger.Debug("======et: %s\n", elementType.Id)
		switch elementType.Type {
		case ELEM_TYPE_DATE:
			conditions = append(conditions,
				createCondition(
					fmt.Sprintf("field%d_from", elementType.Id),
					elementType,
					CONDITION_TYPE_DATE_FROM,
					form,
					cookies))
			conditions = append(conditions,
				createCondition(
					fmt.Sprintf("field%d_to", elementType.Id),
					elementType,
					CONDITION_TYPE_DATE_TO,
					form,
					cookies))
		default:
			conditions = append(conditions,
				createCondition(
					fmt.Sprintf("field%d", elementType.Id),
					elementType,
					CONDITION_TYPE_NORMAL,
					form,
					cookies))
		}
	}

	/* register date */
	conditions = append(conditions,
		createCondition(
			"registerdate.from",
			ELEMENT_TYPE_REGISTERDATE,
			CONDITION_TYPE_DATE_FROM,
			form,
			cookies))
	conditions = append(conditions,
		createCondition(
			"registerdate.to",
			ELEMENT_TYPE_REGISTERDATE,
			CONDITION_TYPE_DATE_TO,
			form,
			cookies))
	/* modify date */
	conditions = append(conditions,
		createCondition(
			"lastregisterdate.from",
			ELEMENT_TYPE_LASTREGISTERDATE,
			CONDITION_TYPE_DATE_FROM,
			form,
			cookies))
	conditions = append(conditions,
		createCondition(
			"lastregisterdate.to",
			ELEMENT_TYPE_LASTREGISTERDATE,
			CONDITION_TYPE_DATE_TO,
			form,
			cookies))
	/* passed days from last update. */
	conditions = append(conditions,
		createCondition(
			"lastregisterdate.days",
			ELEMENT_TYPE_LASTREGISTERDATE,
			CONDITION_TYPE_DAYS,
			form,
			cookies))
	return conditions
}

func createSortCondition(form url.Values) Condition {
	elementTypeId, _ := strconv.Atoi(form.Get("sort"))
	sort := Condition{}
	if elementTypeId > 0 {
		sort.ElementTypeId = elementTypeId
	}
	elementTypeId, _ = strconv.Atoi(form.Get("rsort"))
	if elementTypeId > 0 {
		sort.ElementTypeId = elementTypeId
		sort.Value = "reverse"
	}
	return sort
}

func setConditions(conditions []Condition, keywords []string, elementTypes []ElementType) []interface{} {
	params := []interface{}{}
	for _, cond := range conditions {
		if !validCondition(cond) {
			continue
		}
		if cond.ConditionType == CONDITION_TYPE_DAYS {
			continue /* プレースフォルダが無いためスルーする */
		}

		logger.Debug("#######setConditions: %d\n", cond.ElementTypeId)
		params = append(params, cond.ValidValue())
	}

	if len(keywords) > 0 {
		for _, keyword := range keywords {
			for i, _ := range elementTypes {
				logger.Debug("%d", i)
				params = append(params, keyword)
			}
		}
	}
	return params
}

func getKeywords(form url.Values) []string {
	q := form.Get("q")
	if len(q) == 0 {
		return []string{}
	}
	return strings.Split(q, " ")
}

func SearchTickets(projectName string, form url.Values, cookies []*http.Cookie, elementTypes []ElementType) SearchResult {
	keywords := getKeywords(form)
	logger.Debug("=============keywords", keywords)
	logger.Debug("======field2: [%s]\n", form.Encode())
	conditions := createConditions(form, cookies, elementTypes)
	sort := createSortCondition(form)
	page, _ := strconv.Atoi(form.Get("p"))

	statement := getSearchSqlString(projectName, conditions, sort, keywords, elementTypes)
	statement += " limit ? offset ? "

	params := setConditions(conditions, keywords, elementTypes)

	params = append(params, strconv.Itoa(LIST_COUNT_PER_SEARCH_PAGE))
	params = append(params, strconv.Itoa(page*LIST_COUNT_PER_SEARCH_PAGE))

	/* 1ページ分のticket_idを取得する。 */
	results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
		var id int
		var state string
		rows.Scan(&id, &state)
		logger.Debug("ticket id:", id)
		lastMessage := getLastMessage(projectName, id, elementTypes, true)
		logger.Debug("elements count:", len(lastMessage.Elements))
		return NewTicketWithoutMessages(id, lastMessage)
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	tickets := make([]Ticket, len(results))
	for i, p := range results {
		tickets[i] = p.(Ticket)
	}

	/* hit件数を取得する。 */
	hitCount := 0
	var states []State
	{
		statement := fmt.Sprintf(
			"select count(res.id), res.state from (%s) as res  "+
				"inner join list_item as l on l.element_type_id = %d and l.name = res.state "+
				"group by res.state "+
				"order by l.sort",
			getSearchSqlString(projectName, conditions, sort, keywords, elementTypes),
			ELEM_ID_STATUS)
		params := setConditions(conditions, keywords, elementTypes)

		results, err := query(projectName, statement, params, func(rows *sql.Rows) interface{} {
			var count int
			var state string
			rows.Scan(&count, &state)
			hitCount += count
			return State{0, state, count}
		})
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		states = make([]State, len(results))
		for i, p := range results {
			states[i] = p.(State)
		}
	}

	logger.Debug("count: %d\n", len(states))
	sums := []int{}
	//TODO 数値項目の合計
	//sums := set_tickets_number_sum(db, conditions, NULL, keywords_a, result);

	return SearchResult{hitCount, 0, tickets, states, sums}
}

func ValidateTicket(projectName string, form url.Values, elementTypes []ElementType) map[string]string {
	logger.Debug("======field2: [%s]\n", form.Encode())

	messages := make(map[string]string)
	for _, e := range elementTypes {
		fieldName := fmt.Sprintf("field%d", e.Id)
		formVal := form.Get(fieldName)
		logger.Debug("field%d: %s", e.Id, formVal)
		if e.Required && len(strings.TrimSpace(formVal)) == 0 {
			messages[fieldName] = gettext.Gettext("it will required. please describe.")
		}
	}
	return messages
}

func createMessageInsertSql(ticket_id int, elementTypes []ElementType, form url.Values) (string, []interface{}) {
	cols := []string{}
	phs := []string{}
	params := []interface{}{ticket_id, GetLocalTime()}
	for _, e := range elementTypes {
		cols = append(cols, fmt.Sprintf("field%d", e.Id))
		phs = append(phs, "?")

		fieldName := fmt.Sprintf("field%d", e.Id)
		value := ""
		if len(e.ListItems) > 0 {
			value = strings.Join(form[fieldName], "\t")
		} else {
			value = form.Get(fieldName)
		}
		params = append(params, value)
	}

	return fmt.Sprintf(
		"insert into message(ticket_id, registerdate, %s) values (?, ?, %s)",
		strings.Join(cols, ", "),
		strings.Join(phs, ", ")), params
}

func registerOrUpdateTicket(projectName string, id int, form url.Values, elementTypes []ElementType) int {
	mode := "update"
	if id == -1 {
		mode = "create"
	}

	err := tran(projectName, func(db *sql.DB, tx *sql.Tx) error {
		if mode == "create" {
			/* 新規の場合は、ticketテーブルにレコードを挿入する。 */
			result, err := db.Exec("insert into ticket(id, registerdate, closed) values (NULL, ?, 0)", GetLocalTime())
			if err != nil {
				logger.Error(err)
				panic(err)
			}
			id64, _ := result.LastInsertId()
			id = int(id64)
		}
		/* クローズの状態に変更されたかどうかを判定する。 */
		doClose := 0
	ET:
		for _, e := range elementTypes {
			if len(e.ListItems) == 0 {
				break
			}
			fieldName := fmt.Sprintf("field%d", e.Id)
			vs, _ := form[fieldName]
			for _, v := range vs {
				for _, item := range e.ListItems {
					if v == item.Name && item.Close {
						doClose = 1
						break ET
					}
				}
			}
		}
		/* クローズ状態に変更されていた場合は、closedに1を設定する。 */
		result, err := db.Exec("update ticket set closed = ? where id = ?", doClose, id)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		if count, _ := result.RowsAffected(); count != 1 {
			logger.Error("update ticket failed to close.")
			panic("update ticket failed to close.")
		}
		statement, params := createMessageInsertSql(id, elementTypes, form)
		stmt, err := db.Prepare(statement)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		defer stmt.Close()

		logger.Debug("sql", statement)
		values := []reflect.Value{}
		for _, p := range params {
			logger.Debug("[%s]", p)
			values = append(values, reflect.ValueOf(p))
		}
		returnValues := reflect.ValueOf(stmt).MethodByName("Exec").Call(values)
		result2, _ := returnValues[0].Interface().(sql.Result)
		if !returnValues[1].IsNil() {
			err = returnValues[1].Interface().(error)
		}
		if err != nil {
			logger.Error(err)
			panic(err)
		}
		message_id, _ := result2.LastInsertId()

		/* message_id を更新する。 */
		if mode == "create" {
			db.Exec("update ticket set original_message_id = ?, last_message_id = ? where id = ?",
				message_id, message_id, id)
		} else {
			db.Exec("update ticket set last_message_id = ? where id = ?",
				message_id, id)
		}

		//        elements = ticket->elements;
		//        /* register attachment file. */
		//        foreach (it, elements) {
		//            Element* e = it->element;
		//            if (e->is_file) {
		//                int size;
		//                char filename[DEFAULT_LENGTH];
		//                char mime_type[DEFAULT_LENGTH];
		//                char* fname;
		//                char* ctype;
		//                ElementFile* content_a;
		//                fname = get_upload_filename(e->element_type_id, filename);
		//                size = get_upload_size(e->element_type_id);
		//                ctype = get_upload_content_type(e->element_type_id, mime_type);
		//                content_a = get_upload_content(e->element_type_id);
		//                if (exec_query(
		//                            db,
		//                            "insert into element_file("
		//                            " id, message_id, element_type_id, filename, size, mime_type, content"
		//                            ") values (NULL, ?, ?, ?, ?, ?, ?) ",
		//                            COLUMN_TYPE_INT, message_id,
		//                            COLUMN_TYPE_INT, e->element_type_id,
		//                            COLUMN_TYPE_TEXT, fname,
		//                            COLUMN_TYPE_INT, size,
		//                            COLUMN_TYPE_TEXT, mime_type,
		//                            COLUMN_TYPE_BLOB_ELEMENT_FILE, content_a,
		//                            COLUMN_TYPE_END) == 0)
		//                    die("insert failed.");
		//                element_file_free(content_a);
		//            }
		//        }
		return nil
	})
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	return id
}

func RegisterTicket(projectName string, form url.Values, elementTypes []ElementType) int {
	return registerOrUpdateTicket(projectName, -1, form, elementTypes)
}
func ReplyTicket(projectName string, ticketId int, form url.Values, elementTypes []ElementType) int {
	return registerOrUpdateTicket(projectName, ticketId, form, elementTypes)
}

/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
