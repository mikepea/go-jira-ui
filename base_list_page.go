package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type BaseListPage struct {
	selectedLine     int
	uiList           *ui.List
	displayLines     []string
	cachedResults    []string
	firstDisplayLine int
}

func (p *BaseListPage) PreviousLine(n int) {
	p.selectedLine = p.selectedLine - n
	if p.selectedLine < 0 {
		p.selectedLine = 0
	}
	if p.selectedLine < p.firstDisplayLine {
		p.firstDisplayLine = p.selectedLine
	}
}

func (p *BaseListPage) NextLine(n int) {
	if p.selectedLine < len(p.displayLines)-n {
		p.selectedLine = p.selectedLine + n
	} else {
		p.selectedLine = len(p.displayLines) - 1
	}
	if p.selectedLine > p.lastDisplayedLine() {
		p.firstDisplayLine = p.firstDisplayLine + n
	}
}

func (p *BaseListPage) PreviousPage() {
	p.PreviousLine(p.uiList.Height - 2)
}

func (p *BaseListPage) NextPage() {
	p.NextLine(p.uiList.Height - 2)
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

func (p *BaseListPage) markActiveLine() {
	for i, v := range p.cachedResults {
		selected := ""
		if i == p.selectedLine {
			selected = "fg-white,bg-blue"
			p.displayLines[i] = fmt.Sprintf("[%s](%s)", v, selected)
		} else {
			p.displayLines[i] = v
		}
	}
}

func (p *BaseListPage) Id() string {
	return "X"
}

func (p *BaseListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines[p.firstDisplayLine:]
	ui.Render(ls)
}

func (p *BaseListPage) Create() {
	ui.Clear()
	ls := ui.NewList()
	p.uiList = ls
	p.selectedLine = 0
	p.firstDisplayLine = 0
	p.cachedResults = make([]string, 0)
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = "Base List -- Do Not Use"
	ls.Height = ui.TermHeight()
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.Update()
}
