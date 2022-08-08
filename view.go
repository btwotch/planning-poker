package main

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type view struct {
	model  *model
	player *player
	app    *tview.Application
	ssh    ssh.Session
	sync.Mutex
}

func (v *view) chooseFibWin() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true)

	top := tview.NewBox().SetBorder(false)
	bottom := tview.NewBox().SetBorder(false)
	middle := tview.NewTable().SetFixed(1, 0).SetSelectable(false, false).SetBorders(false)

	var checkedCell *tview.TableCell

	for i, num := range []uint8{0, 1, 2, 3, 5, 8, 13, 20, 40, 100, 200} {
		currentNum := num
		label := fmt.Sprintf("   %d   ", num)
		cell := tview.NewTableCell(label)
		cell.SetAlign(tview.AlignCenter)
		cell.SetClickedFunc(func() bool {
			if checkedCell != nil {
				checkedCell.SetTextColor(tcell.ColorWhite)
				checkedCell.SetBackgroundColor(tcell.ColorBlack)
			}
			cell.SetTextColor(tcell.ColorBlack)
			cell.SetBackgroundColor(tcell.ColorWhite)

			checkedCell = cell

			v.player.setChoice(currentNum)
			return true
		})
		middle.SetCell(i, 0, cell)
	}

	flex.AddItem(top, 0, 1, false)
	flex.AddItem(middle, 0, 8, true)
	flex.AddItem(bottom, 0, 1, false)

	return flex
}

func (v *view) toggleFlex(flex *tview.Flex) {
	flex.SetBorder(false)
	flex.AddItem(nil, 0, 2, false)
	button := tview.NewButton("Disclose")
	if v.model.getDisclosed() {
		button.SetLabelColor(tcell.ColorRed)
	} else {
		button.SetLabelColor(tcell.ColorWhite)
	}
	button.SetSelectedFunc(func() {
		go v.app.QueueUpdate(func() {
			v.model.toggleDisclose()
		})
	})
	flex.AddItem(button, 10, 3, false)
	clear := tview.NewButton("Clear")
	clear.SetSelectedFunc(func() {
		v.model.clearChoices()
		v.model.setDisclose(false)
	})
	flex.AddItem(clear, 10, 3, false)
	flex.AddItem(nil, 0, 2, false)
}

func (v *view) tableFlex(flex *tview.Flex) {
	flex.SetDirection(tview.FlexRow)
	flex.SetBorder(true)

	top := tview.NewFlex()
	v.toggleFlex(top)

	bottom := tview.NewBox().SetBorder(false)
	middle := tview.NewTable().SetFixed(1, 0).SetSelectable(false, false).SetBorders(false)

	for i, p := range v.model.getPlayers() {
		label := fmt.Sprintf("%s", p.getName())
		nameCell := tview.NewTableCell(label)
		nameCell.SetAlign(tview.AlignCenter)

		disclosedFibCell := tview.NewTableCell(fmt.Sprintf("[%d]", p.getChoice()))
		censoredFibCell := tview.NewTableCell("[ ]")
		middle.SetCell(i, 1, nameCell)
		if p.getName() == v.player.getName() && p.hasChosen() {
			middle.SetCell(i, 2, disclosedFibCell)
		} else if p.hasChosen() && v.model.getDisclosed() {
			middle.SetCell(i, 2, disclosedFibCell)
		} else if p.hasChosen() && !v.model.getDisclosed() {
			middle.SetCell(i, 2, censoredFibCell)
		}
	}

	flex.AddItem(top, 3, 3, false)
	flex.AddItem(middle, 0, 3, true)
	flex.AddItem(bottom, 0, 3, false)
}

func (v *view) flex() *tview.Flex {
	v.Lock()
	defer v.Unlock()

	flex := tview.NewFlex()

	chooseFibWin := v.chooseFibWin()
	othersChoiceWin := tview.NewFlex()
	v.tableFlex(othersChoiceWin)

	flex.AddItem(chooseFibWin, 0, 3, false)
	flex.AddItem(othersChoiceWin, 0, 7, false)

	v.model.subscribe(func() {
		v.Lock()
		defer v.Unlock()

		go v.app.QueueUpdateDraw(func() {
			v.Lock()
			defer v.Unlock()

			othersChoiceWin.Clear()

			v.tableFlex(othersChoiceWin)
		})
	})

	v.player.subscribe(func() {
		v.Lock()
		defer v.Unlock()

		v.app.EnableMouse(false)
		v.ssh.Exit(0)
		v.app.Stop()
	})

	return flex
}
