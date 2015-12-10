package usb

import (
	"errors"
	"fmt"
)

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
	// Unmount the device if it is not already unmounted.
	Unmount() error
	// Size returns the dvice size in bytes.
	Size() int
}

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
	UUID      string // currently only set on Mac
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
	return ByteSize(d.Size()).ToHuman()
}

func (d *Device) String() string {
	return fmt.Sprintf("Vendor: %s\nModel: %s\nPath: %s\nSectors: %d\n"+
		"BlockSize: %d\nWritable: %v\nRemovable: %v\n",
		d.Vendor,
		d.Model,
		d.Path,
		d.Sectors,
		d.BlockSize,
		d.Writable,
		d.Removable)
}

// check if device is readable & writable.
func (d *Device) check() error {
	switch {
	case len(d.Path) == 0:
		return ErrPathEmpty
	case d.Sectors <= 0:
		return ErrSectorEmpty
	case d.BlockSize <= 0:
		return ErrBlocksizeEmpty
	case !d.Writable:
		return ErrNotWritable
	case !d.Removable:
		return ErrNotRemovable
		// case len(d.UUID) == 0: // TODO implement UUID on Linux and uc check
		// 	return ErrNoUUID
	}
	return nil
}
