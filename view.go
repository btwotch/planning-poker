package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type view struct {
	model  *model
	player *player
	app    *tview.Application
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

	flex.AddItem(top, 0, 3, false)
	flex.AddItem(middle, 0, 3, true)
	flex.AddItem(bottom, 0, 3, false)

	return flex
}

func (v *view) tableFlex(flex *tview.Flex) {
	flex.SetDirection(tview.FlexRow)
	flex.SetBorder(true)

	top := tview.NewBox().SetBorder(false)
	bottom := tview.NewBox().SetBorder(false)
	middle := tview.NewTable().SetFixed(1, 0).SetSelectable(false, false).SetBorders(false)

	for i, p := range v.model.getPlayers() {
		label := fmt.Sprintf("%s", p.getName())
		nameCell := tview.NewTableCell(label)
		nameCell.SetAlign(tview.AlignCenter)

		fibCell := tview.NewTableCell(fmt.Sprintf("%d", p.getChoice()))
		middle.SetCell(i, 1, nameCell)
		if p.hasChosen() {
			middle.SetCell(i, 2, fibCell)
		}
	}

	flex.AddItem(top, 0, 3, false)
	flex.AddItem(middle, 0, 3, true)
	flex.AddItem(bottom, 0, 3, false)
}

func (v *view) flex() *tview.Flex {
	flex := tview.NewFlex()

	chooseFibWin := v.chooseFibWin()
	othersChoiceWin := tview.NewFlex()
	v.tableFlex(othersChoiceWin)

	flex.AddItem(chooseFibWin, 0, 3, false)
	flex.AddItem(othersChoiceWin, 0, 7, false)

	v.model.subscribe(func() {
		go v.app.QueueUpdateDraw(func() {
			othersChoiceWin.Clear()

			v.tableFlex(othersChoiceWin)
		})
	})

	return flex
}
