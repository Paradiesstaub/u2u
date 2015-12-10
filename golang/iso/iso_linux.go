package iso

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/Paradiesstaub/u2u/golang/root"
)

// TODO implement the ISO writer in pure go.
//
// * How to run the writing part as root?
// * How to unmount a mounted device?

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
	// call dd
	e := w.executor
	cmd := fmt.Sprintf("dd bs=32M if=%s of=%s && sync", iso, device)
	msg := fmt.Sprintf("Write ISO to USB Stick (%s).", device)
	e.Exec(cmd, msg)
	// call sync
	if err := exec.Command("sync").Run(); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("DONE - ISO written to %s\n", device)
}

func NewWriter() Writer {
	executor := root.NewExecutor()
	return NewDummyWriterLinux(executor)
}
