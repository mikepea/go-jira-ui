package jiraui

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type Sort struct {
	Name string
	JQL  string
}

type SortOrderPage struct {
	BaseListPage
	cachedResults []Sort
}

var baseSorts = []Sort{
	Sort{"default", " "},
	Sort{"created, oldest first", "ORDER BY created ASC"},
	Sort{"updated, newest first", "ORDER BY updated DESC"},
	Sort{"updated, oldest first", "ORDER BY updated ASC"},
	Sort{"---", ""}, // no-op line in UI
}

func getSorts() (sorts []Sort) {
	opts := getJiraOpts()
	if q := opts["sorts"]; q != nil {
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
			sorts = append(sorts, Sort{q2["name"], q2["jql"]})
		}
	}
	return append(baseSorts, sorts...)
}

func (p *SortOrderPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
			p.displayLines[i] = fmt.Sprintf("[%s](%s)", v.Name, selected)
		} else {
			p.displayLines[i] = fmt.Sprintf("%s", v.Name)
		}
	}
}

func (p *SortOrderPage) PreviousPara() {
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

func (p *SortOrderPage) NextPara() {
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

func (p *SortOrderPage) BottomOfPage() {
	p.selectedLine = len(p.cachedResults) - 1
	firstLine := p.selectedLine - (p.uiList.Height - 3)
	if firstLine > 0 {
		p.firstDisplayLine = firstLine
	} else {
		p.firstDisplayLine = 0
	}
}

func (p *SortOrderPage) SelectedSort() Sort {
	return p.cachedResults[p.selectedLine]
}

func (p *SortOrderPage) SelectItem() {
	if p.SelectedSort().JQL == "" {
		return
	}
	q := new(TicketListPage)
	q.ActiveQuery = ticketListPage.ActiveQuery
	q.ActiveSort = p.SelectedSort()
	ticketListPage = q
	currentPage = ticketListPage
	changePage()
}

func (p *SortOrderPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *SortOrderPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]Sort, 0)
	q.Create()
	changePage()
}

func (p *SortOrderPage) Create() {
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	if len(p.cachedResults) == 0 {
		p.cachedResults = getSorts()
		p.displayLines = make([]string, len(p.cachedResults))
	}
	ls.ItemFgColor = ui.ColorGreen
	ls.BorderLabel = "Sort By..."
	ls.BorderFg = ui.ColorRed
	ls.Height = 10
	ls.Width = 50
	ls.X = ui.TermWidth()/2 - ls.Width/2
	ls.Y = ui.TermHeight()/2 - ls.Height/2
	p.Update()
}
