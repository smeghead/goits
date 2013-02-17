package data

import (
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
    ElementTypeId int
    StrVal string
    IsFile bool
}

type SettingFile struct {
    name string
    FileName string
    Size int
    MimeType string
    Content string
}

type Message struct {
    Id int
    Elements []Element
}

type ListItem struct {
    Id int
    ElementTypeId int
    Name string
    Close int
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
}

type ElementFile struct {
    Id int
    ElementTypeId int
    Name string
    Size int
    MimeType string
    Content string
}

const (
    ELEM_TYPE_TEXT = 0
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

const (
    CONDITION_TYPE_NORMAL = 0
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
    Messages []Message
    States []State
    Sums []int
}

type Wiki struct {
    Id int
    Name string
    Content string
}

type Ticket struct {
    ProjectId int
    ProjectCode string
    ProjectName string
    Id int
    Title string
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
