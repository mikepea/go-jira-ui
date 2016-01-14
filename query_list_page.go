package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type Query struct {
	Name string
	JQL  string
}

type QueryPage struct {
	BaseListPage
	cachedResults []Query
}

var baseQueries = []Query{
	Query{"My Assigned Tickets", "assignee = currentUser() AND resolution = Unresolved"},
	Query{"My Reported Tickets", "reporter = currentUser() AND resolution = Unresolved"},
	Query{"My Watched Tickets", "watcher = currentUser() AND resolution = Unresolved"},
	Query{"My Voted Tickets", "voter = currentUser() AND resolution = Unresolved"},
	Query{"---", ""}, // no-op line in UI
}

func getQueries() (queries []Query) {
	opts := getJiraOpts()
	if q := opts["queries"]; q != nil {
		qList := q.([]interface{})
		for _, v := range qList {
			q1 := v.(map[interface{}]interface{})
			q2 := make(map[string]string)
			for k, v := range q1 {
				switch k := k.(type) {
				case string:
					switch v := v.(type) {
					case string:
						q2[k] = v
					}
				}
			}
			queries = append(queries, Query{q2["name"], q2["jql"]})
		}
	}
	return append(baseQueries, queries...)
}

func (p *QueryPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
			p.displayLines[i] = fmt.Sprintf("[%-30s -- %s](%s)", v.Name, v.JQL, selected)
		} else {
			p.displayLines[i] = fmt.Sprintf("%-30s -- %s", v.Name, v.JQL)
		}
	}
}

func (p *QueryPage) PreviousPara() {
	newDisplayLine := 0
	if p.selectedLine == 0 {
		return
	}
	for i := p.selectedLine - 1; i > 0; i-- {
		if p.cachedResults[i].JQL == "" {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(p.selectedLine - newDisplayLine)
}

func (p *QueryPage) NextPara() {
	newDisplayLine := len(p.cachedResults) - 1
	if p.selectedLine == newDisplayLine {
		return
	}
	for i := p.selectedLine + 1; i < len(p.cachedResults); i++ {
		if p.cachedResults[i].JQL == "" {
			newDisplayLine = i
			break
		}
	}
	p.NextLine(newDisplayLine - p.selectedLine)
}

func (p *QueryPage) BottomOfPage() {
	p.selectedLine = len(p.cachedResults) - 1
	firstLine := p.selectedLine - (p.uiList.Height - 3)
	if firstLine > 0 {
		p.firstDisplayLine = firstLine
	} else {
		p.firstDisplayLine = 0
	}
}

func (p *QueryPage) SelectedQuery() Query {
	return p.cachedResults[p.selectedLine]
}

func (p *QueryPage) SelectItem() {
	if p.SelectedQuery().JQL == "" {
		return
	}
	q := new(TicketListPage)
	q.ActiveQuery = p.SelectedQuery()
	ticketListPage = q
	currentPage = ticketListPage
	changePage()
}

func (p *QueryPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *QueryPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]Query, 0)
	q.Create()
	changePage()
}

func (p *QueryPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	p.cachedResults = getQueries()
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Queries"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
