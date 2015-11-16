package main

import (
	"os"
	"strings"

	"github.com/Paradiesstaub/u2u/golang/iso"
)

// Controler handels the input events of the main view.
type Controler struct {
	view   *View
	writer iso.Writer
	model  *Model
}

// Model object of main.
type Model struct {
	dropdownList []string
}

const defaultDropdownEntry = "No USB drive found"

// NewControler creates a new controler for the main view.
func NewControler(w iso.Writer, v View) *Controler {
	m := &Model{
		dropdownList: []string{defaultDropdownEntry},
	}
	// setup dropw-down
	device := "/dev/sdc" // TODO remove
	m.dropdownList = []string{device}
	v.SetDropdownItems(m.dropdownList)
	return &Controler{
		view:   &v,
		writer: w,
		model:  m,
	}
}

// CheckShowRunButton checks if the run button should be displayed in QML.
func (c *Controler) CheckShowRunButton(iso string) bool {
	if len(iso) == 0 {
		return false
	}
	// test for .iso postfix
	if !strings.HasSuffix(strings.ToLower(iso), ".iso") {
		return false
	}
	// check if default entry is still used
	dd := c.model.dropdownList
	if len(dd) == 1 && dd[0] == defaultDropdownEntry {
		return false
	}
	return true
}

// CreateUsb handels the usb creation.
func (c *Controler) CreateUsb(iso, device string) {
	c.writer.Write(iso, device)
}

// Quit terminates the program.
func (c *Controler) Quit() {
	os.Exit(0)
}
