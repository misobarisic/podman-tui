package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/rs/zerolog/log"
)

func (nets *Networks) runCommand(cmd string) {
	switch cmd {
	case "create":
		nets.createDialog.Display()
	case "inspect":
		nets.inspect()
	case "prune":
		nets.cprune()
	case "rm":
		nets.rm()
	}
}

func (nets *Networks) displayError(title string, err error) {
	var message string
	if title != "" {
		message = fmt.Sprintf("%s: %v", title, err)
	} else {
		message = fmt.Sprintf("%v", err)
	}

	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	nets.errorDialog.SetText(message)
	nets.errorDialog.Display()
}

func (nets *Networks) create() {
	createOpts := nets.createDialog.NetworkCreateOptions()
	filename, err := networks.Create(createOpts)
	if err != nil {
		nets.displayError("NETWORK CREATE ERROR", err)
		return
	}
	nets.UpdateData()
	nets.messageDialog.SetTitle("podman network create")
	nets.messageDialog.SetText(filename)
	nets.messageDialog.Display()

}

func (nets *Networks) inspect() {
	if nets.selectedID == "" {
		nets.displayError("", fmt.Errorf("there is no network to display inspect"))
		return
	}
	data, err := networks.Inspect(nets.selectedID)
	if err != nil {
		title := fmt.Sprintf("NETWORK (%s) INSPECT ERROR", nets.selectedID)
		nets.displayError(title, err)
		return
	}
	nets.messageDialog.SetTitle("podman network inspect")
	nets.messageDialog.SetText(data)
	nets.messageDialog.Display()
}

func (nets *Networks) cprune() {
	nets.confirmDialog.SetTitle("podman network prune")
	nets.confirmData = "prune"
	nets.confirmDialog.SetText("Are you sure you want to remove all un used network ?")
	nets.confirmDialog.Display()
}

func (nets *Networks) prune() {
	nets.progressDialog.SetTitle("network purne in progress")
	nets.progressDialog.Display()
	prune := func() {
		err := networks.Prune()
		nets.progressDialog.Hide()
		if err != nil {
			nets.displayError("NETWORK PRUNE ERROR", err)
			return
		}
	}
	go prune()
}

func (nets *Networks) rm() {
	if nets.selectedID == "" {
		nets.displayError("", fmt.Errorf("there is no network to remove"))
		return
	}
	nets.confirmDialog.SetTitle("podman network remove")
	nets.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following network? \n\nNETWORK NAME : %s", nets.selectedID)
	nets.confirmDialog.SetText(description)
	nets.confirmDialog.Display()
}

func (nets *Networks) remove() {
	nets.progressDialog.SetTitle("network remove in progress")
	nets.progressDialog.Display()
	remove := func(id string) {
		err := networks.Remove(id)
		nets.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("NETWORK (%s) REMOVE ERROR", nets.selectedID)
			nets.displayError(title, err)
			return
		}
		nets.UpdateData()
	}
	go remove(nets.selectedID)
}
