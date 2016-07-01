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

func (p *SortOrderPage) IsPopulated() bool {
	if len(p.cachedResults) > 0 {
		return true
	} else {
		return false
	}
}

func (p *SortOrderPage) itemizeResults() []string {
	items := make([]string, len(p.cachedResults))
	for i, v := range p.cachedResults {
		items[i] = fmt.Sprintf("%s", v.Name)
	}
	return items
}

func (p *SortOrderPage) PreviousPara() {
	newDisplayLine := 0
	sl := p.uiList.Cursor
	if sl == 0 {
		return
	}
	for i := sl - 1; i > 0; i-- {
		if p.cachedResults[i].JQL == "" {
			newDisplayLine = i
			break
		}
	}
	p.PreviousLine(sl - newDisplayLine)
}

func (p *SortOrderPage) NextPara() {
	newDisplayLine := len(p.cachedResults) - 1
	sl := p.uiList.Cursor
	if sl == newDisplayLine {
		return
	}
	for i := sl + 1; i < len(p.cachedResults); i++ {
		if p.cachedResults[i].JQL == "" {
			newDisplayLine = i
			break
		}
	}
	p.NextLine(newDisplayLine - sl)
}

func (p *SortOrderPage) SelectedSort() Sort {
	return p.cachedResults[p.uiList.Cursor]
}

func (p *SortOrderPage) SelectItem() {
	if p.SelectedSort().JQL == "" {
		return
	}
	q := new(TicketListPage)
	// pop old page, we're going to 'replace' it
	var oldTicketListPage Navigable
	if len(previousPages) > 0 {
		oldTicketListPage, previousPages = previousPages[len(previousPages)-1], previousPages[:len(previousPages)-1]
		switch a := oldTicketListPage.(type) {
		case *TicketListPage:
			q.ActiveQuery = a.ActiveQuery
			q.ActiveSort = p.SelectedSort()
			currentPage = q
		}
	}
	changePage()
}

func (p *SortOrderPage) Update() {
	ls := p.uiList
	ui.Render(ls)
}

func (p *SortOrderPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]Sort, 0)
	changePage()
	q.Create()
}

func (p *SortOrderPage) Create() {
	ls := NewScrollableList()
	p.uiList = ls
	p.uiList.Cursor = 0
	if len(p.cachedResults) == 0 {
		p.cachedResults = getSorts()
	}
	ls.Items = p.itemizeResults()
	ls.ItemFgColor = ui.ColorGreen
	ls.BorderLabel = "Sort By..."
	ls.BorderFg = ui.ColorRed
	ls.Height = 10
	ls.Width = 50
	ls.X = ui.TermWidth()/2 - ls.Width/2
	ls.Y = ui.TermHeight()/2 - ls.Height/2
	p.Update()
}
