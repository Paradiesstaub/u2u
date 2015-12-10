package iso

import "fmt"

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
