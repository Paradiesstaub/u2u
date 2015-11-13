package main

import (
	"fmt"

	"github.com/Paradiesstaub/u2u/golang/usb"
)

const defaultDropdownEntry = "No USB drive found"

// Model of the main view.
type Model struct {
	devices []usb.Device

	//dropdownList []string
}

// NewModel creates a new model for the main view.
func NewModel() *Model {
	return &Model{
		devices: usb.List(),
	}
}

// DropwdownList converts []device to a []string.
func (m *Model) DropwdownList() []string {
	if len(m.devices) == 0 {
		return []string{defaultDropdownEntry}
	}
	arr := make([]string, len(m.devices))
	for k, d := range m.devices {
		s := fmt.Sprintf("%s %s (%s)", d.Vendor, d.Model, d.SizeToHuman())
		arr[k] = s
	}
	return arr
}
