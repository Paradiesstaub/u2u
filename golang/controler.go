package main

import (
	"os"
	"teststuff/isocreator/golang/iso"
)

// Controler handels the input events of the main view.
type Controler struct {
	writer iso.Writer
}

// NewControler creates a new controler for the main view.
func NewControler(w iso.Writer) *Controler {
	return &Controler{
		writer: w,
	}
}

// CreateUsb handels the usb creation.
func (c *Controler) CreateUsb(iso, device string) {
	c.writer.Write(iso, device)
}

// Quit terminates the program.
func (c *Controler) Quit() {
	os.Exit(0)
}
