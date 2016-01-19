package jiraui

import (
	ui "github.com/gizak/termui"
)

type CommandBar struct {
	uiList *ui.List
	text   []byte
}

func (p *CommandBar) Update() {
	ls := p.uiList
	ui.Render(ls)
}

func (p *CommandBar) Create() {
	ls := ui.NewList()
	p.uiList = ls
	ls.ItemFgColor = ui.ColorGreen
	ls.Border = false
	ls.Height = 1
	ls.Width = ui.TermWidth()
	ls.X = 0
	ls.Y = ui.TermHeight() - 1
	p.Update()
}
