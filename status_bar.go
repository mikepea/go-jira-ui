package jiraui

import (
	ui "github.com/gizak/termui"
)

type StatusBar struct {
	uiList *ui.List
	lines  []string
}

func (p *StatusBar) StatusLines() []string {
	return p.lines
}

func (p *StatusBar) Update() {
	ls := p.uiList
	ls.Items = p.StatusLines()
	ui.Render(ls)
}

func (p *StatusBar) Create() {
	ls := ui.NewList()
	p.uiList = ls
	p.lines = append(p.lines, "STATUS BAR... STATUS BAR... STATUS BAR...")
	ls.ItemFgColor = ui.ColorWhite
	ls.ItemBgColor = ui.ColorRed
	ls.Bg = ui.ColorRed
	ls.Border = false
	ls.Height = 1
	ls.Width = ui.TermWidth()
	ls.X = 0
	ls.Y = ui.TermHeight() - 2
	p.Update()
}
