package usb

// About golang build constrains:
// https://golang.org/pkg/go/build/#hdr-Build_Constraints

// Device contains USB device information.
type Device struct {
	Vendor string
	Model  string
	// device path
	Path string
	// numbers of sectors
	Sectors int
	// physical block size in bytes
	BlockSize int
	Writable  bool
	Removable bool
}

// Size returns the dvice size in bytes.
func (d *Device) Size() int {
	return d.Sectors * d.BlockSize
}

// SizeToHuman formats the device size to a human friendly string, e.g. 1,1 GiB.
func (d *Device) SizeToHuman() string {
	return ByteSize(d.Size()).ToHuman()
}
