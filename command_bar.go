package jiraui

import (
	ui "github.com/gizak/termui"
)

type CommandBar struct {
	uiList *ui.List
	EditBox
	commandType         byte
	commandHistoryIndex int
	searchHistoryIndex  int
	commandHistory      []string
	searchHistory       []string
}

func (p *CommandBar) resetSearchIndex() {
	p.searchHistoryIndex = len(p.searchHistory) - 1
}

func (p *CommandBar) resetCommandIndex() {
	p.commandHistoryIndex = len(p.commandHistory) - 1
}

func addCommandIfNotSameAsLast(new string, history *[]string) {
	log.Noticef("addCommandIfNotSameAsLast: got %s", new)
	l := len(*history)
	if l > 0 && new == (*history)[l-1] {
		return
	} else {
		log.Noticef("addCommandIfNotSameAsLast: Adding %s", new)
		*history = append(*history, new)
	}
}

func (p *CommandBar) Submit() {
	if obj, ok := currentPage.(CommandBoxer); ok {
		obj.SetCommandMode(false)
		obj.ExecuteCommand()
		if len(p.text) > 1 {
			ct := p.text[0]
			cb := string(p.text[1:])
			switch {
			case ct == ':':
				addCommandIfNotSameAsLast(cb, &p.commandHistory)
				p.resetCommandIndex()
			case (ct == '/' || ct == '?'):
				addCommandIfNotSameAsLast(cb, &p.searchHistory)
				p.resetSearchIndex()
			}
		}
		p.text = []byte("")
	}
	// currentPage may have changed
	if obj, ok := currentPage.(CommandBoxer); ok {
		obj.Update()
	}
}

func (p *CommandBar) PreviousCommand() {
	if obj, ok := currentPage.(CommandBoxer); ok {
		ct := p.commandType
		switch {
		case (ct == ':'):
			if len(p.commandHistory) == 0 {
				return
			}
			p.text = []byte(string(p.commandType) + p.commandHistory[p.commandHistoryIndex])
			if p.commandHistoryIndex > 0 {
				p.commandHistoryIndex = p.commandHistoryIndex - 1
			} else {
				p.resetCommandIndex()
			}
		case (ct == '/' || ct == '?'):
			if len(p.searchHistory) == 0 {
				return
			}
			p.text = []byte(string(p.commandType) + p.searchHistory[p.searchHistoryIndex])
			if p.searchHistoryIndex > 0 {
				p.searchHistoryIndex = p.searchHistoryIndex - 1
			} else {
				p.resetSearchIndex()
			}
		}
		p.MoveCursorToEnd()
		obj.Update()
	}
}

func (p *CommandBar) NextCommand() {
	if obj, ok := currentPage.(CommandBoxer); ok {
		ct := p.commandType
		switch {
		case (ct == ':'):
			if len(p.commandHistory) == 0 {
				return
			}
			p.text = []byte(string(p.commandType) + p.commandHistory[p.commandHistoryIndex])
			if p.commandHistoryIndex < len(p.commandHistory)-1 {
				p.commandHistoryIndex = p.commandHistoryIndex + 1
			} else {
				p.resetCommandIndex()
			}
		case (ct == '/' || ct == '?'):
			if len(p.searchHistory) == 0 {
				return
			}
			p.text = []byte(string(p.commandType) + p.searchHistory[p.searchHistoryIndex])
			if p.searchHistoryIndex < len(p.commandHistory)-1 {
				p.searchHistoryIndex = p.searchHistoryIndex + 1
			} else {
				p.resetSearchIndex()
			}
		}
		p.MoveCursorToEnd()
		obj.Update()
	}
}

func (p *CommandBar) Reset() {
	p.text = []byte(``)
	p.line_voffset = 0
	p.cursor_boffset = 0
	p.cursor_voffset = 0
	p.cursor_coffset = 0
}

func (p *CommandBar) Update() {
	if len(p.text) == 0 {
		if obj, ok := currentPage.(CommandBoxer); ok {
			obj.SetCommandMode(false)
		}
	} else {
		p.commandType = p.text[0]
	}
	ls := p.uiList
	strs := []string{string(p.text)}
	ls.Items = strs
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
