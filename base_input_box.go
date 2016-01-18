package jiraui

import (
	ui "github.com/gizak/termui"
	"strings"
)

type BaseInputBox struct {
	eb     *EditBox
	uiList *ui.List
}

func (p *BaseInputBox) Update() {
	ls := p.uiList
	ls.Items = strings.Split(string(p.eb.text), "\n")
	ui.Render(ls)
}

func (p *BaseInputBox) Create() {
	ls := ui.NewList()
	var strs []string
	p.uiList = ls
	p.eb = new(EditBox)
	ls.Items = strs
	ls.ItemFgColor = ui.ColorGreen
	ls.BorderFg = ui.ColorRed
	ls.Height = 1
	ls.Width = 30
	ls.Overflow = "wrap"
	ls.X = ui.TermWidth()/2 - ls.Width/2
	ls.Y = ui.TermHeight()/2 - ls.Height/2
	p.Update()
}
