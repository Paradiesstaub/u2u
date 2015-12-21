package main

import (
	"log"
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

func (c *Controler) deviceByIndex(index int) usb.Devicer {
	return c.model.devices[index]
}

// CreateUsb handels the usb creation.
func (c *Controler) CreateUsb(iso string, dropdownIndex int) {
	device := c.deviceByIndex(dropdownIndex)
	// check that hardware pointed to by the device path hasn't changed.
	if !device.IsSameDevice() {
		log.Fatal("E: CreateUsb - Not same device!") // TODO show error to user
	}
	// unmount all device partitions
	if err := device.Unmount(); err != nil {
		log.Fatal("E: ", err.Error()) // TODO show error to user
	}
	c.writer.Write(iso, device.Path())
}

// Quit terminates the program.
func (c *Controler) Quit() {
	os.Exit(0)
}
