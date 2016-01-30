package jiraui

import (
	"fmt"
	"regexp"

	ui "github.com/gizak/termui"
)

type Search struct {
	command     string
	directionUp bool
	re          *regexp.Regexp
}

type BaseListPage struct {
	selectedLine     int
	uiList           *ui.List
	displayLines     []string
	cachedResults    []string
	firstDisplayLine int
	isPopulated      bool
	ActiveSearch     Search
}

func (p *BaseListPage) SetSearch(searchCommand string) {
	if len(searchCommand) < 2 {
		// must be '/a' minimum
		return
	}
	direction := []byte(searchCommand)[0]
	regex := "(?i)" + string([]byte(searchCommand)[1:])
	s := new(Search)
	s.command = searchCommand
	if direction == '?' {
		s.directionUp = true
	} else if direction == '/' {
		s.directionUp = false
	} else {
		// bad command
		return
	}
	if re, err := regexp.Compile(regex); err != nil {
		// bad regex
		return
	} else {
		s.re = re
		p.ActiveSearch = *s
	}
}

func (p *BaseListPage) IsPopulated() bool {
	if len(p.cachedResults) > 0 || p.isPopulated {
		return true
	} else {
		return false
	}
}

func (p *BaseListPage) FixFirstDisplayLine(n int) {
	if p.selectedLine < p.firstDisplayLine {
		p.firstDisplayLine = p.selectedLine
	} else if p.selectedLine > p.lastDisplayedLine() {
		p.firstDisplayLine = p.selectedLine - (p.PageLines() - 1)
	}
}

func (p *BaseListPage) PreviousLine(n int) {
	p.selectedLine = p.selectedLine - n
	if p.selectedLine < 0 {
		p.selectedLine = 0
	}
	p.FixFirstDisplayLine(n)
}

func (p *BaseListPage) NextLine(n int) {
	if p.selectedLine < len(p.displayLines)-n {
		p.selectedLine = p.selectedLine + n
	} else {
		p.selectedLine = len(p.displayLines) - 1
	}
	p.FixFirstDisplayLine(n)
}

func (p *BaseListPage) PreviousPara() {
	p.PreviousLine(5)
}

func (p *BaseListPage) NextPara() {
	p.NextLine(5)
}

func (p *BaseListPage) PreviousPage() {
	p.PreviousLine(p.PageLines())
}

func (p *BaseListPage) NextPage() {
	p.NextLine(p.PageLines())
}

func (p *BaseListPage) PageLines() int {
	return p.uiList.Height - 2
}

func (p *BaseListPage) TopOfPage() {
	p.selectedLine = 0
	p.firstDisplayLine = 0
}

func (p *BaseListPage) BottomOfPage() {
	p.selectedLine = len(p.cachedResults) - 1
	firstLine := p.selectedLine - (p.uiList.Height - 3)
	if firstLine > 0 {
		p.firstDisplayLine = firstLine
	} else {
		p.firstDisplayLine = 0
	}
}

func (p *BaseListPage) lastDisplayedLine() int {
	return lastLineDisplayed(p.uiList, p.firstDisplayLine, 3)
}

func (p *BaseListPage) SetSelectedLine(line int) {
	if line > 0 && line < len(p.cachedResults) {
		p.selectedLine = line
		p.FixFirstDisplayLine(0)
	}
}

func (p *BaseListPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
			if v == "" {
				v = " "
			} else if ok, _ := regexp.MatchString(`\[.+\]\((fg|bg)-[a-z]{1,6}\)`, v); ok {
				r := regexp.MustCompile(`\[(.*?)\]\((fg|bg)-[a-z]{1,6}\)`)
				v = r.ReplaceAllString(v, `$1`)
			}
			p.displayLines[i] = fmt.Sprintf("[%s](%s)", v, selected)
		} else {
			p.displayLines[i] = v
		}
	}
}

func (p *BaseListPage) Id() string {
	return fmt.Sprintf("BaseListPage(%p)", p)
}

func (p *BaseListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *BaseListPage) Refresh() {
	pDeref := &p
	q := *pDeref
	q.cachedResults = make([]string, 0)
	changePage()
	q.Create()
}

func (p *BaseListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.cachedResults = make([]string, 0)
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Updating, please wait"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
