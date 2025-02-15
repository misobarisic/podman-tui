package images

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (img *Images) Draw(screen tcell.Screen) {
	img.refresh()
	img.Box.DrawForSubclass(screen, img)
	img.Box.SetBorder(false)
	x, y, width, height := img.GetInnerRect()
	img.table.SetRect(x, y, width, height)
	img.table.SetBorder(true)

	img.table.Draw(screen)
	x, y, width, height = img.table.GetInnerRect()
	// error dialog
	if img.errorDialog.IsDisplay() {
		img.errorDialog.SetRect(x, y, width, height)
		img.errorDialog.Draw(screen)
		return
	}
	// command dialog dialog
	if img.cmdDialog.IsDisplay() {
		img.cmdDialog.SetRect(x, y, width, height)
		img.cmdDialog.Draw(screen)
		return
	}
	// command input dialog
	if img.cmdInputDialog.IsDisplay() {
		img.cmdInputDialog.SetRect(x, y, width, height)
		img.cmdInputDialog.Draw(screen)
		return
	}
	// message dialog
	if img.messageDialog.IsDisplay() {
		img.messageDialog.SetRect(x, y, width, height+1)
		img.messageDialog.Draw(screen)
		return
	}
	// confirm dialog
	if img.confirmDialog.IsDisplay() {
		img.confirmDialog.SetRect(x, y, width, height)
		img.confirmDialog.Draw(screen)
		return
	}

	// search dialog
	if img.searchDialog.IsDisplay() {
		img.searchDialog.SetRect(x, y, width, height)
		img.searchDialog.Draw(screen)
	}
	// progress dialog
	if img.progressDialog.IsDisplay() {
		img.progressDialog.SetRect(x, y, width, height)
		img.progressDialog.Draw(screen)
	}
	// history dialog
	if img.historyDialog.IsDisplay() {
		img.historyDialog.SetRect(x, y, width, height)
		img.historyDialog.Draw(screen)
		return
	}

}
