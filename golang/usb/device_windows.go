package usb

// List returns a slice of writable USB devices.
// Or an empty slice if no device could be found or an error occurred.
func List() []Device {
	// TODO
	return make([]Device, 0)
}

// FindDevices lists all USB device paths.
// Returns an empty slice if there was an error.
func FindDevices() []string {
	// TODO
	return make([]string, 0)
}
