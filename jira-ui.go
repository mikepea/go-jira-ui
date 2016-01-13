package main

import (
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
	"os"
)

const (
	ticketQuery = iota
	ticketList  = iota
	labelList   = iota
	ticketShow  = iota
)

const (
	default_list_template = `{{ range .issues }}{{ .key | printf "%-20s"}}  {{ dateFormat "2006-01-02" .fields.created }}/{{ dateFormat "2006-01-02T15:04" .fields.updated }}  {{ .fields.summary | printf "%-75s"}} -- labels({{ join "," .fields.labels }})
{{ end }}`
)

var exitNow = false

type GoBacker interface {
	GoBack()
}

type ItemSelecter interface {
	SelectItem()
}

type TicketEditer interface {
	EditTicket()
}

type TicketCommenter interface {
	CommentTicket()
}

type PagePager interface {
	NextLine(int)
	PreviousLine(int)
	NextPage()
	PreviousPage()
	TopOfPage()
	BottomOfPage()
	Update()
}

type Navigable interface {
	Create()
	Update()
	PreviousLine(int)
	NextLine(int)
	PreviousPage()
	NextPage()
	Id() string
}

var currentPage Navigable

var ticketQueryPage *QueryPage
var ticketListPage *TicketListPage
var labelListPage *LabelListPage

func changePage() {
	switch currentPage.(type) {
	case *QueryPage:
		log.Debugf("changePage: QueryPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *TicketListPage:
		log.Debugf("changePage: TicketListPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *LabelListPage:
		log.Debugf("changePage: LabelListPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *TicketShowPage:
		log.Debugf("changePage: TicketShowPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	}
}

var (
	log    = logging.MustGetLogger("jira")
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

func main() {

	var err error
	logging.SetLevel(logging.NOTICE, "")

	err = ensureLoggedIntoJira()
	if err != nil {
		log.Error("Login failed. Aborting")
		os.Exit(2)
	}

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	registerKeyboardHandlers()

	ticketQueryPage = new(QueryPage)
	currentPage = ticketQueryPage

	for exitNow != true {

		currentPage.Create()
		ui.Loop()

	}

}
