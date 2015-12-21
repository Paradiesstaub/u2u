package usb

import "errors"

// About golang build constrains:
// https://golang.org/pkg/go/build/#hdr-Build_Constraints

type Devicer interface {
	Vendor() string
	Model() string
	Path() string
	Sectors() int
	BlockSize() int
	// IsSameDevice checks if the path point still to the same device.
	IsSameDevice() bool
	// Unmount the device partitions if they are mounted.
	Unmount() error
}

var (
	ErrPathEmpty      = errors.New("No device path found.")
	ErrSectorEmpty    = errors.New("Sectors has 0 size.")
	ErrBlocksizeEmpty = errors.New("Block size has 0 size.")
	ErrNotWritable    = errors.New("Device is not writable.")
	ErrNotRemovable   = errors.New("Device is not removable.")
	ErrNoUUID         = errors.New("Device has no UUID.")
)

// SizeToHuman formats the device size to a human friendly string, e.g. 1,1 GiB.
func SizeToHuman(d Devicer) string {
	size := d.Sectors() * d.BlockSize()
	return ByteSize(size).ToHuman()
}
