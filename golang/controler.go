package main

import (
	"os"
	"strings"

	"github.com/Paradiesstaub/u2u/golang/iso"
	"github.com/Paradiesstaub/u2u/golang/usb"
)

// Controler handels the input events of the main view.
type Controler struct {
	view   *View
	writer iso.Writer
	model  *Model
}

// NewControler creates a new controler for the main view.
func NewControler(w iso.Writer, v View, m *Model) *Controler {
	v.SetDropdownItems(m.DropwdownList())
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
	if len(c.model.devices) == 0 {
		return false
	}
	return true
}

func (c *Controler) deviceByIndex(index int) usb.Device {
	return c.model.devices[index]
}

// CreateUsb handels the usb creation.
func (c *Controler) CreateUsb(iso string, dropdownIndex int) {
	d := c.deviceByIndex(dropdownIndex)
	c.writer.Write(iso, d.Path)
}

// Quit terminates the program.
func (c *Controler) Quit() {
	os.Exit(0)
}
