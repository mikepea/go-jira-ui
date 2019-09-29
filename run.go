package jiraui

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	ui "gopkg.in/gizak/termui.v2"
)

var exitNow = false

type EditPager interface {
	DeleteRuneBackward()
	InsertRune(r rune)
	Update()
	Create()
}

type TicketCommander interface {
	ActiveTicketId() string
	Refresh()
}

type Searcher interface {
	SetSearch(string)
	Search()
}

type CommandBoxer interface {
	SetCommandMode(bool)
	ExecuteCommand()
	CommandMode() bool
	CommandBar() *CommandBar
	Update()
}

type NextTicketer interface {
	NextTicket()
}

type PrevTicketer interface {
	PrevTicket()
}

type GoBacker interface {
	GoBack()
}

type Refresher interface {
	Refresh()
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
	NextPara()
	PreviousPara()
	NextPage()
	PreviousPage()
	TopOfPage()
	BottomOfPage()
	IsPopulated() bool
	Update()
}

type Navigable interface {
	Create()
	Update()
	Id() string
}

type RankSelector interface {
	MarkItemForRanking()
}

var currentPage Navigable
var previousPages []Navigable

var ticketQueryPage *QueryPage
var helpPage *HelpPage
var labelListPage *LabelListPage
var sortOrderPage *SortOrderPage
var passwordInputBox *PasswordInputBox
var commandBar *CommandBar

func changePage() {
	if currentPage == nil {
		currentPage = new(QueryPage)
	}
	switch currentPage.(type) {
	case *QueryPage:
		log.Debugf("changePage: QueryPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *TicketListPage:
		log.Debugf("changePage: TicketListPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *SortOrderPage:
		log.Debugf("changePage: SortOrderPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *LabelListPage:
		log.Debugf("changePage: LabelListPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *TicketShowPage:
		log.Debugf("changePage: TicketShowPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	case *HelpPage:
		log.Debugf("changePage: HelpPage %s (%p)", currentPage.Id(), currentPage)
		currentPage.Create()
	}
}

const LOG_MODULE = "jiraui"

var (
	log    = logging.MustGetLogger(LOG_MODULE)
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

var cliOpts map[string]interface{}

func Run(app *kingpin.Application) {

	usage := func(ok bool) {
		printer := fmt.Printf
		if !ok {
			printer = func(format string, args ...interface{}) (int, error) {
				return fmt.Fprintf(os.Stderr, format, args...)
			}
			defer func() {
				os.Exit(1)
			}()
		} else {
			defer func() {
				os.Exit(0)
			}()
		}
		output := fmt.Sprintf(`
Usage:
  jira-ui ls <Query Options> 
  jira-ui ISSUE
  jira-ui

General Options:
  -e --endpoint=URI   URI to use for jira
  -l --log=FILE       FILE to use for log (default /dev/null)
  -h --help           Show this usage
  -u --user=USER      Username to use for authenticaion
  -v --verbose        Increase output logging
  --skiplogin         Skip the login check. You must have a valid session token (eg via 'jira login')
  --version           Print version

Ticket View Options:
  -t --template=FILE  Template file to use for viewing tickets
  -m --max_wrap=VAL   Maximum word-wrap width when viewing ticket text (0 disables)

Query Options:
  -q --query=JQL            Jira Query Language expression for the search
  -f --queryfields=FIELDS   Fields that are used in "list" view

`)
		printer(output)
	}

	jiraCommands := map[string]string{
		"list":     "list",
		"ls":       "list",
		"password": "password",
		"passwd":   "password",
	}

	cliOpts = make(map[string]interface{})
	cliOpts["log"] = "/dev/null"
	/*
		setopt := func(name string, value interface{}) {
			cliOpts[name] = value
		}
	*/

	logging.SetLevel(logging.NOTICE, "jira")
	logging.SetLevel(logging.NOTICE, LOG_MODULE)

	/*
		op := optigo.NewDirectAssignParser(map[string]interface{}{
			"h|help": usage,
			"version": func() {
				fmt.Println(fmt.Sprintf("version: %s", VERSION))
				os.Exit(0)
			},
			"v|verbose+": func() {
				logging.SetLevel(logging.GetLevel(LOG_MODULE)+1, LOG_MODULE)
			},
			"l|log=s":         setopt,
			"u|user=s":        setopt,
			"endpoint=s":      setopt,
			"q|query=s":       setopt,
			"f|queryfields=s": setopt,
			"t|template=s":    setopt,
			"m|max_wrap=i":    setopt,
			"skip_login":      setopt,
		})

		if err := op.ProcessAll(os.Args[1:]); err != nil {
			log.Errorf("%s", err)
			usage(false)
		}
		args := op.Args
	*/
	var args []string
	f, err := os.Create(cliOpts["log"].(string))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	backend := logging.NewLogBackend(f, "", 0)
	logging.SetBackend(backend)

	var command string
	if len(args) > 0 {
		if alias, ok := jiraCommands[args[0]]; ok {
			command = alias
			args = args[1:]
		} else {
			command = "view"
			args = args[0:]
		}
	} else {
		command = "toplevel"
	}

	requireArgs := func(count int) {
		if len(args) < count {
			log.Errorf("Not enough arguments. %d required, %d provided", count, len(args))
			usage(false)
		}
	}

	if val, ok := cliOpts["skip_login"]; !ok || !val.(bool) {
		err = ensureLoggedIntoJira()
		if err != nil {
			log.Error("Login failed. Aborting")
			os.Exit(2)
		}
	}

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	registerEventHandlers()

	helpPage = new(HelpPage)
	commandBar = new(CommandBar)

	switch command {
	case "list":
		if query := cliOpts["query"]; query == nil {
			log.Errorf("Must supply a --query option to %q", command)
			os.Exit(1)
		} else {
			p := new(TicketListPage)
			p.ActiveQuery.JQL = query.(string)
			p.ActiveQuery.Name = "adhoc"
			currentPage = p
		}
	case "view":
		requireArgs(1)
		p := new(TicketShowPage)
		p.TicketId = args[0]
		currentPage = p
	case "toplevel":
		currentPage = new(QueryPage)
	case "password":
		currentPage = new(PasswordInputBox)
	default:
		log.Errorf("Unknown command %s", command)
		os.Exit(1)
	}

	for exitNow != true {

		currentPage.Create()
		ui.Loop()

	}

}
