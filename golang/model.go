package main

import (
	"fmt"

	"github.com/Paradiesstaub/u2u/golang/usb"
)

const defaultDropdownEntry = "No USB drive found"

// Model of the main view.
type Model struct {
	devices []usb.Devicer
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
	var s string
	for k, d := range m.devices {
		vSize := len(d.Vendor())
		mSize := len(d.Model())
		switch {
		case vSize != 0 && mSize != 0:
			s = fmt.Sprintf("%s %s (%s)", d.Vendor(), d.Model(), usb.SizeToHuman(d))
		case vSize == 0 && mSize != 0:
			s = fmt.Sprintf("%s (%s)", d.Model(), usb.SizeToHuman(d))
		case vSize != 0 && mSize == 0:
			s = fmt.Sprintf("%s (%s)", d.Model(), usb.SizeToHuman(d))
		default:
			s = fmt.Sprintf("(%s)", usb.SizeToHuman(d))
		}
		arr[k] = s
	}
	return arr
}
