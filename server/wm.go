package server

import (
	"fmt"
	"os/exec"
	"strings"
)

func NewBspwm() *Bspwm {
	return &Bspwm{oldIds: make([]string, 0)}
}

type Bspwm struct {
	nextIdNum uint
	oldIds    []string
}

func (this *Bspwm) RootDesktop() (id string) {
	names := this.exec("query", "--monitor", "focused", "--desktops")
	splits := strings.SplitN(string(names), "\n", 2)
	return splits[0]
}

func (this *Bspwm) AddDesktop() (id string) {
	if len(this.oldIds) == 0 {
		id = fmt.Sprintf("metadesk-desktop%d", this.nextIdNum)
		this.nextIdNum++
	} else {
		id = this.oldIds[len(this.oldIds)-1]
		this.oldIds = this.oldIds[:len(this.oldIds)-1]
	}

	this.exec("monitor", "--add-desktops", id)
	return id
}

func (this *Bspwm) RemoveDesktop(id string) {
	this.oldIds = append(this.oldIds, id)
	this.exec("monitor", "--remove-desktops", id)
}

func (this *Bspwm) FocusDesktop(id string) {
	this.exec("desktop", "--focus", id)
}

func (this *Bspwm) IsDesktopOccupied(id string) bool {
	return len(this.exec("query", "--desktop", id, "--windows")) > 0
}

func (this *Bspwm) ClaimFocusedWindow(id string) {
	this.exec("window", "focused", "-d", id)
}

func (this *Bspwm) exec(args ...string) []byte {
	output, _ := exec.Command("bspc", args...).Output()

	return output
}
