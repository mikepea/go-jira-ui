package jiraui

import (
	ui "github.com/gizak/termui"
	"strings"
)

type PasswordInputBox struct {
	BaseInputBox
}

func (p *PasswordInputBox) Update() {
	ls := p.uiList
	ls.Items = strings.Split(string(p.text), "\n")
	ui.Render(ls)
}

func (p *PasswordInputBox) Create() {
	ls := ui.NewList()
	p.uiList = ls
	var strs []string
	ls.Items = strs
	ls.ItemFgColor = ui.ColorGreen
	ls.BorderLabel = "Enter Password:"
	ls.BorderFg = ui.ColorRed
	ls.Height = 3
	ls.Width = 30
	ls.X = ui.TermWidth()/2 - ls.Width/2
	ls.Y = ui.TermHeight()/2 - ls.Height/2
	p.Update()
}
