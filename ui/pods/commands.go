package pods

import (
	"fmt"
	"strings"

	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/rs/zerolog/log"
)

func (p *Pods) runCommand(cmd string) {
	switch cmd {
	case "create":
		p.createDialog.Display()
	case "inspect":
		p.inspect()
	case "kill":
		p.kill()
	case "pause":
		p.pause()
	case "prune":
		p.confirmDialog.SetTitle("podman pod prune")
		p.confirmData = "prune"
		p.confirmDialog.SetText("Are you sure you want to remove all stopped pods ?")
		p.confirmDialog.Display()
	case "restart":
		p.restart()
	case "rm":
		p.rm()
	case "start":
		p.start()
	case "stats":
		p.stats()
	case "stop":
		p.stop()
	case "top":
		p.top()
	case "unpause":
		p.unpause()
	}
}

func (p *Pods) displayError(title string, err error) {
	var message string
	if title != "" {
		message = fmt.Sprintf("%s: %v", title, err)
	} else {
		message = fmt.Sprintf("%v", err)
	}

	log.Error().Msgf("%s: %v", strings.ToLower(title), err)
	p.errorDialog.SetText(message)
	p.errorDialog.Display()
}

func (p *Pods) stats() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to perform stats command"))
		return
	}
	podOptions := p.getAllItemsForStats()
	p.statsDialog.SetPodsOptions(podOptions)
	p.statsDialog.Display()
}

func (p *Pods) create() {
	podSpec := p.createDialog.GetPodSpec()
	err := ppods.Create(podSpec)
	if err != nil {
		p.displayError("POD CREATE ERROR", err)
		return
	}
}

func (p *Pods) inspect() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to display inspect"))
		return
	}
	data, err := ppods.Inspect(p.selectedID)
	if err != nil {
		title := fmt.Sprintf("POD (%s) INSPECT ERROR", p.selectedID)
		p.displayError(title, err)
		return
	}
	p.messageDialog.SetTitle("podman pod inspect")
	p.messageDialog.SetText(data)
	p.messageDialog.Display()
}

func (p *Pods) kill() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to kill"))
		return
	}
	p.progressDialog.SetTitle("pod kill in progress")
	p.progressDialog.Display()
	kill := func(id string) {
		err := ppods.Kill(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) KILL ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
	}
	go kill(p.selectedID)
}

func (p *Pods) pause() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to pause"))
		return
	}
	p.progressDialog.SetTitle("pod pause in progress")
	p.progressDialog.Display()
	pause := func(id string) {
		err := ppods.Pause(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) PAUSE ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
	}
	go pause(p.selectedID)
}

func (p *Pods) prune() {
	p.progressDialog.SetTitle("pod purne in progress")
	p.progressDialog.Display()
	unpause := func() {
		errData, err := ppods.Prune()
		p.progressDialog.Hide()
		if err != nil {
			p.displayError("PODS PRUNE ERROR", err)
			return
		}
		if len(errData) > 0 {
			errMessages := fmt.Errorf("%v", errData)
			p.displayError("PODS PRUNE ERROR", errMessages)
		}

	}
	go unpause()
}

func (p *Pods) restart() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to restart"))
		return
	}
	p.progressDialog.SetTitle("pod restart in progress")
	p.progressDialog.Display()
	restart := func(id string) {
		err := ppods.Restart(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) RESTART ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
	}
	go restart(p.selectedID)
}

func (p *Pods) rm() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to remove"))
		return
	}
	p.confirmDialog.SetTitle("podman pod rm")
	p.confirmData = "rm"
	description := fmt.Sprintf("Are you sure you want to remove following pod ? \n\nPOD ID : %s", p.selectedID)
	p.confirmDialog.SetText(description)
	p.confirmDialog.Display()
}

func (p *Pods) remove() {
	p.progressDialog.SetTitle("pod remove in progress")
	p.progressDialog.Display()
	remove := func(id string) {
		errData, err := ppods.Remove(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) REMOVE ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
		if len(errData) > 0 {
			title := fmt.Sprintf("POD (%s) REMOVE ERROR", p.selectedID)
			p.displayError(title, fmt.Errorf("%v", errData))
		}
	}
	go remove(p.selectedID)
}

func (p *Pods) start() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to start"))
		return
	}
	p.progressDialog.SetTitle("pod start in progress")
	p.progressDialog.Display()
	start := func(id string) {
		err := ppods.Start(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) START ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
	}
	go start(p.selectedID)
}

func (p *Pods) stop() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to stop"))
		return
	}
	p.progressDialog.SetTitle("pod stop in progress")
	p.progressDialog.Display()
	stop := func(id string) {
		err := ppods.Stop(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) STOP ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
	}
	go stop(p.selectedID)
}

func (p *Pods) top() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to display top"))
		return
	}
	data, err := ppods.Top(p.selectedID)
	if err != nil {
		title := fmt.Sprintf("POD (%s) TOP ERROR", p.selectedID)
		p.displayError(title, err)
		return
	}
	p.topDialog.UpdateResults(data)
	p.topDialog.Display()
}

func (p *Pods) unpause() {
	if p.selectedID == "" {
		p.displayError("", fmt.Errorf("there is no pod to unpause"))
		return
	}
	p.progressDialog.SetTitle("pod unpause in progress")
	p.progressDialog.Display()
	unpause := func(id string) {
		err := ppods.Unpause(id)
		p.progressDialog.Hide()
		if err != nil {
			title := fmt.Sprintf("POD (%s) UNPAUSE ERROR", p.selectedID)
			p.displayError(title, err)
			return
		}
	}
	go unpause(p.selectedID)
}
