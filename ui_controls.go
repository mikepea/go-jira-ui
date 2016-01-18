package jiraui

import (
	ui "github.com/gizak/termui"
	"os"
)

func registerKeyboardHandlers() {
	ui.Handle("/sys/kbd/", func(ev ui.Event) {
		handleAnyKey(ev)
	})
	ui.Handle("/sys/kbd/q", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleBackKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/C-r", func(ui.Event) {
		handleRefreshKey()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.Close()
		os.Exit(0)
	})
	ui.Handle("/sys/kbd/Q", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			ui.Close()
			os.Exit(0)
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/}", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleParaDownKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/{", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleParaUpKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/j", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleDownKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/k", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleUpKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		handleDownKey()
	})
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		handleUpKey()
	})
	ui.Handle("/sys/kbd/g", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleTopOfPageKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/G", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleBottomOfPageKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/L", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleLabelViewKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/S", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleSortOrderKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/<enter>", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleSelectKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		handlePageDownKey()
	})
	ui.Handle("/sys/kbd/C-f", func(ui.Event) {
		handlePageDownKey()
	})
	ui.Handle("/sys/kbd/C-b", func(ui.Event) {
		handlePageUpKey()
	})
	ui.Handle("/sys/kbd/E", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleEditKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/kbd/C", func(ev ui.Event) {
		if _, ok := currentPage.(PagePager); ok {
			handleCommentKey()
		} else {
			handleAnyKey(ev)
		}
	})
	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		handleResize()
	})
}

func handleLabelViewKey() {
	switch page := currentPage.(type) {
	case *TicketListPage:
		q := new(LabelListPage)
		q.ActiveQuery = page.ActiveQuery
		currentPage = q
		changePage()
	}
	return
}

func handleSortOrderKey() {
	switch currentPage.(type) {
	case *TicketListPage:
		q := new(SortOrderPage)
		currentPage = q
		changePage()
	}
	return
}

func handleRefreshKey() {
	if obj, ok := currentPage.(Refresher); ok {
		obj.Refresh()
	}
}

func handleEditKey() {
	if obj, ok := currentPage.(TicketEditer); ok {
		obj.EditTicket()
	}
}

func handleCommentKey() {
	if obj, ok := currentPage.(TicketCommenter); ok {
		obj.CommentTicket()
	}
}

func handleBackKey() {
	if obj, ok := currentPage.(GoBacker); ok {
		obj.GoBack()
	} else {
		ui.StopLoop()
		exitNow = true
	}
}

func handleResize() {
	changePage()
}

func handleSelectKey() {
	if obj, ok := currentPage.(ItemSelecter); ok {
		obj.SelectItem()
	}
}

func handleTopOfPageKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.TopOfPage()
		obj.Update()
	}
}

func handleBottomOfPageKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.BottomOfPage()
		obj.Update()
	}
}

func handleUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousLine(1)
		obj.Update()
	}
}

func handleDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextLine(1)
		obj.Update()
	}
}

func handlePageUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousPage()
		obj.Update()
	}
}

func handlePageDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextPage()
		obj.Update()
	}
}

func handleParaUpKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.PreviousPara()
		obj.Update()
	}
}

func handleParaDownKey() {
	if obj, ok := currentPage.(PagePager); ok {
		obj.NextPara()
		obj.Update()
	}
}

func handleAnyKey(e ui.Event) {
	if obj, ok := currentPage.(EditPager); ok {
		key := e.Data.(ui.EvtKbd).KeyStr
		var str string
		switch {
		case len(key) == 1:
			str = key
		case key == "<enter>":
			str = "\n"
		case key == "<space>":
			log.Noticef("space!")
			str = ` `
		case key == "<backspace>" || key == "C-8":
			log.Noticef("backspace!")
			obj.DeleteRuneBackward()
			obj.Update()
			return
		default:
			return
		}
		r := decodeTermuiKbdStringToRune(str)
		obj.InsertRune(r)
		obj.Update()
	}
}
