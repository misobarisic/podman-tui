package images

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/rs/zerolog/log"
)

func (img *Images) runCommand(cmd string) {

	switch cmd {
	case "diff":
		img.diff()
	case "history":
		img.history()
	case "inspect":
		img.inspect()
	case "prune":
		img.cprune()
	case "rm":
		img.rm()
	case "search/pull":
		img.searchDialog.Display()
	case "tag":
		img.ctag()
	case "untag":
		img.cuntag()
	}

}

func (img *Images) displayError(title string, err error) {
	var message string
	if title != "" {
		message = fmt.Sprintf("%s: %v", title, err)
	} else {
		message = fmt.Sprintf("%v", err)
	}

	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	img.errorDialog.SetText(message)
	img.errorDialog.Display()
}

func (img *Images) diff() {
	if img.selectedID == "" {
		img.displayError("", fmt.Errorf("there is no image to display diff"))
		return
	}
	img.progressDialog.SetTitle("image diff in progress")
	img.progressDialog.Display()
	diff := func() {
		data, err := images.Diff(img.selectedID)
		img.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) DIFF ERROR", img.selectedID)
			img.displayError(title, err)
			return
		}
		img.messageDialog.SetTitle("podman image diff")
		img.messageDialog.SetText(strings.Join(data, "\n"))
		img.messageDialog.Display()
	}
	go diff()
}

func (img *Images) history() {
	if img.selectedID == "" {
		img.displayError("", fmt.Errorf("there is no image to display history"))
		return
	}
	result, err := images.History(img.selectedID)
	if err != nil {
		title := fmt.Sprintf("IMAGE (%s) HISTORY ERROR", img.selectedID)
		img.displayError(title, err)
	}
	img.historyDialog.UpdateResults(result)
	img.historyDialog.Display()

}

func (img *Images) inspect() {
	if img.selectedID == "" {
		img.displayError("", fmt.Errorf("there is no image to display inspect"))
		return
	}
	data, err := images.Inspect(img.selectedID)
	if err != nil {
		title := fmt.Sprintf("IMAGE (%s) INSPECT ERROR", img.selectedID)
		img.displayError(title, err)
		return
	}
	img.messageDialog.SetTitle("podman image inspect")
	img.messageDialog.SetText(data)
	img.messageDialog.Display()
}

func (img *Images) cprune() {
	img.confirmDialog.SetTitle("podman image prune")
	img.confirmData = "prune"
	img.confirmDialog.SetText("Are you sure you want to remove all unused images")
	img.confirmDialog.Display()
}

func (img *Images) prune() {
	img.progressDialog.SetTitle("image purne in progress")
	img.progressDialog.Display()
	prune := func() {
		err := images.Prune()
		img.progressDialog.Hide()
		if err != nil {
			img.displayError("IMAGE PRUNE ERROR", err)
			return
		}
	}
	go prune()
}

func (img *Images) rm() {
	if img.selectedID == "" {
		img.displayError("", fmt.Errorf("there is no image to remove"))
		return
	}
	img.confirmDialog.SetTitle("podman image remove")
	img.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following image ? \n\nimage name : %s\nimage ID   : %s", img.selectedName, img.selectedID)
	img.confirmDialog.SetText(description)
	img.confirmDialog.Display()
}

func (img *Images) remove() {
	img.progressDialog.SetTitle("image remove in progress")
	img.progressDialog.Display()
	remove := func(id string) {
		data, err := images.Remove(id)
		img.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) REMOVE ERROR", img.selectedID)
			img.displayError(title, err)
		} else {
			img.messageDialog.SetTitle("podman image remove")
			img.messageDialog.SetText(strings.Join(data, "\n"))
			img.messageDialog.Display()
		}

	}
	go remove(img.selectedID)
}

func (img *Images) search(term string) {
	img.progressDialog.SetTitle("image search in progress")
	img.progressDialog.Display()
	search := func(term string) {
		result, err := images.Search(term)
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) SEARCH ERROR", img.selectedID)
			img.displayError(title, err)
		}
		img.searchDialog.UpdateResults(result)
		img.progressDialog.Hide()
	}
	go search(term)
}

func (img *Images) ctag() {
	if img.selectedID == "" {
		img.displayError("", fmt.Errorf("there is no image to tag"))
		return
	}
	img.cmdInputDialog.SetTitle("podman image tag")
	description := fmt.Sprintf("[white::]image name : [black::]%s[white::]\nimage ID   : [black::]%s", img.selectedName, img.selectedID)
	img.cmdInputDialog.SetDescription(description)
	img.cmdInputDialog.SetSelectButtonLabel("tag")
	img.cmdInputDialog.SetLabel("target name")
	img.cmdInputDialog.SetSelectedFunc(func() {
		img.tag(img.cmdInputDialog.GetInputText())
		img.cmdInputDialog.Hide()
	})
	img.cmdInputDialog.Display()
}

func (img *Images) tag(tag string) {
	if err := images.Tag(img.selectedID, tag); err != nil {
		title := fmt.Sprintf("IMAGE (%s) TAG ERROR", img.selectedID)
		img.displayError(title, err)
	}
}

func (img *Images) cuntag() {
	if img.selectedID == "" {
		img.displayError("", fmt.Errorf("there is no image to untag"))
		return
	}
	img.cmdInputDialog.SetTitle("podman image untag")
	img.cmdInputDialog.SetDescription("")
	img.cmdInputDialog.SetSelectButtonLabel("untag")
	img.cmdInputDialog.SetLabel("image")
	img.cmdInputDialog.SetInputText(img.selectedName)
	img.cmdInputDialog.SetSelectedFunc(func() {
		img.untag(img.cmdInputDialog.GetInputText())
		img.cmdInputDialog.Hide()
	})
	img.cmdInputDialog.Display()
}

func (img *Images) untag(id string) {
	if err := images.Untag(id); err != nil {
		title := fmt.Sprintf("IMAGE (%s) UNTAG ERROR", img.selectedID)
		img.displayError(title, err)
	}
}

func (img *Images) pull(image string) {
	img.progressDialog.SetTitle("image pull in progress")

	img.progressDialog.Display()
	pull := func(name string) {
		err := images.Pull(name)
		if err != nil {
			title := fmt.Sprintf("IMAGE (%s) PULL ERROR", img.selectedID)
			img.displayError(title, err)
		}
		img.progressDialog.Hide()
	}
	go pull(image)
}
