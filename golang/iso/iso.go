package iso

import (
	"fmt"
	"log"
	"strings"

	"github.com/Paradiesstaub/u2u/golang/root"
)

// Writer handels platform independent writing of an ISO file to USB.
type Writer interface {
	// Writes the iso file to the device.
	Write(iso, device string)
}

// FakeWriter doesn't manipulate hardware, just prints an info.
type FakeWriter struct{}

func (w FakeWriter) Write(iso, device string) {
	fmt.Printf("FakeWriter: %s to %s\n", iso, device)
}

// DummyWriterLinux uses dd to write the ISO to usb.
// ATTENTION: real hardware is manipulated when using this object!
type DummyWriterLinux struct {
	executor root.Executor
}

// NewDummyWriterLinux creates a writer which manipulates hardaware using dd.
func NewDummyWriterLinux(e root.Executor) *DummyWriterLinux {
	return &DummyWriterLinux{executor: e}
}

func (w DummyWriterLinux) Write(iso, device string) {
	if len(iso) == 0 {
		log.Fatalln("Missing path to ISO file!")
	}
	if len(device) == 0 {
		log.Fatalln("No device path, can't write ISO file!")
	}
	if !strings.HasPrefix(device, "/dev/") {
		log.Fatalln("Passed device path has to start with '/dev/'")
	}
	e := w.executor
	cmd := fmt.Sprintf("dd bs=4M if=%s of=%s && sync", iso, device)
	msg := fmt.Sprintf("Write ISO to USB stick (%s).", device)
	e.Exec(cmd, msg)
	fmt.Printf("DONE - ISO written to %s!\n", device)
}
