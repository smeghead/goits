package data

import (
    "strconv"
    "strings"
)

const (
        LIST_COUNT_PER_LIST_PAGE = 15
        LIST_COUNT_PER_SEARCH_PAGE = 30
)

type ProjectInfo struct {
    Id int
    Name string
    Sort int
}
func (p *ProjectInfo) IsManage() bool {
    return p.Name == "top"
}

type Project struct {
    Url string
    Name string
    HomeDescription string
    HomeUrl string
    UploadMaxSize int
    Locale string
}

type Element struct {
    ElementType ElementType
    StrVal string
    ElementFile ElementFile
}
func (e *Element) GetSelectedItemId() string {
    if e.ElementType.ListItems == nil {
        return ""
    }
    for _, item := range e.ElementType.ListItems {
        if strings.Index(e.StrVal, item.Name) > -1 {
            return strconv.Itoa(item.Id)
        }
    }
    return ""
}
func (e *Element) HasTicketLink() bool {
    return e.ElementType.Id == ELEM_ID_ID || e.ElementType.Id == ELEM_ID_TITLE
}
func (e *Element) IsFile() bool {
    return e.ElementType.Type == ELEM_TYPE_UPLOADFILE
}

type SettingFile struct {
    Name string
    FileName string
    Size int
    MimeType string
    Content string
}

type Ticket struct {
    ProjectId int
    ProjectCode string
    ProjectName string
    Id int
    Title string
    Status string
    LastMessage Message
    Messages []Message
}
func NewTicketWithoutMessages(ticketId int, lastMessage Message) Ticket {
    return NewTicket(ticketId, lastMessage, []Message{})
}
func NewTicket(ticketId int, lastMessage Message, messages []Message) Ticket {
    ticket := Ticket{}

    ticket.Id = ticketId
    ticket.Title = GetElementField(lastMessage.Elements, ELEM_ID_TITLE)
    ticket.Status = GetElementField(lastMessage.Elements, ELEM_ID_STATUS)
    ticket.LastMessage = lastMessage
    ticket.Messages = messages

    return ticket
}

func GetElementField(elements []Element, id int) string {
    for _, element := range elements {
        if element.ElementType.Id == id {
            return element.StrVal
        }
    }
    return ""
}

type Message struct {
    Id int
    Elements []Element
    RegisterDate string
}
func (m *Message) GetSender() string {
    for _, element := range m.Elements {
        if element.ElementType.Id == ELEM_ID_SENDER {
            return element.StrVal
        }
    }
    return ""
}


type ListItem struct {
    Id int
    Name string
    Close bool
    Sort int
}

type ElementType struct {
    Id int
    Type int
    TicketProperty bool
    ReplyProperty bool
    Required bool
    Name string
    Description string
    AutoAddItem bool
    DefaultValue string
    DisplayInList bool
    Sort int
    ListItems []ListItem
}

type ElementFile struct {
    Id int
    Filename string
    Size int
    MimeType string
    Content string
    Deleted bool
}

const (
    ELEM_TYPE_TEXT = iota
    ELEM_TYPE_TEXTAREA
    ELEM_TYPE_CHECKBOX
    ELEM_TYPE_LIST_SINGLE
    ELEM_TYPE_LIST_MULTI
    ELEM_TYPE_UPLOADFILE
    ELEM_TYPE_DATE
    ELEM_TYPE_LIST_SINGLE_RADIO
    ELEM_TYPE_NUM
    /* this values match database value, so, if you add ELEM_TYPE, add list of tail. DBの値と連動しているので、追加する場合は、後に追加する必要がある。*/
)

var ELEMENT_TYPE_ID ElementType = ElementType{ELEM_ID_ID, 0, true, false, true, "ID", "", false, "", true, 0, nil}
var ELEMENT_TYPE_REGISTERDATE ElementType = ElementType{ELEM_ID_REGISTERDATE, 0, true, false, true, "register date", "", false, "", true, 0, nil}
var ELEMENT_TYPE_LASTREGISTERDATE ElementType = ElementType{ELEM_ID_LASTREGISTERDATE, 0, true, false, true, "last register date", "", false, "", true, 0, nil}
var ELEMENT_TYPE_ORG_SENDER ElementType = ElementType{ELEM_ID_ORG_SENDER, 0, true, false, true, "org sender", "", false, "", false, 0, nil}
var ELEMENT_TYPE_LASTREGISTREDATE_PASSED ElementType = ElementType{ELEM_ID_LASTREGISTERDATE_PASSED, 0, true, false, true, "last register date passed", "", false, "", true, 0, nil}

const (
    ELEM_ID_ID = -1
    ELEM_ID_TITLE = 1
    ELEM_ID_SENDER = 2
    ELEM_ID_STATUS = 3
    ELEM_ID_REGISTERDATE = -2
    ELEM_ID_LASTREGISTERDATE = -3
    ELEM_ID_ORG_SENDER = -4
    ELEM_ID_LASTREGISTERDATE_PASSED = -5
)
const BASIC_ELEMENT_MAX = 3

type Condition struct {
    ElementTypeId int
    ConditionType int
    Value string
    CookieValue string
}
func (c *Condition) ValidValue() string {
    if len(c.Value) > 0 {
        return c.Value
    } else {
        return c.CookieValue
    }
    return ""
}

const (
    CONDITION_TYPE_NORMAL = iota
    CONDITION_TYPE_DATE_FROM
    CONDITION_TYPE_DATE_TO
    CONDITION_TYPE_DAYS
)

type State struct {
    Id int
    Name string
    Count int
}

type SearchResult struct {
    HitCount int
    Page int
    Tickets []Ticket
    States []State
    Sums []int
}

type Wiki struct {
    Id int
    Name string
    Content string
}


type DbInfo struct {
    Id int
    FieldCount int
}

type BurndownChartPoint struct {
    All int
    NotClosed int
    Day string
}

type UserRanking struct {
    Name string
    Count int
}
/* vim: set ts=4 sw=4 sts=4 expandtab fenc=utf-8: */
