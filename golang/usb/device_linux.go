package usb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type linuxDevice struct {
	vendor    string
	model     string
	path      string
	sectors   int
	blockSize int
	writable  bool
	removable bool
	// uuid string // TODO
}

// List returns a slice of writable and removable USB devices.
func List() []Devicer {
	arr := FindDevices()
	if len(arr) == 0 {
		return make([]Devicer, 0)
	}
	var devices []Devicer
	for _, p := range arr {
		d, err := info(p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := d.check(); err == nil {
			devices = append(devices, Devicer(&d))
		}
	}
	return devices
}

// FindDevices lists all USB device paths in Linux.
// Returns an empty slice if there was an error.
func FindDevices() []string {
	files, err := ioutil.ReadDir("/dev/disk/by-id")
	if err != nil {
		return make([]string, 0)
	}
	var arr []string
	var name string
	var previousName string
	for _, f := range files {
		// name without partion suffix
		// example device name: usb-TrekStor_USB_Stick_QU_AA04013000007552-0:0
		name = strings.Split(f.Name(), ":")[0]
		// filter non-usb drives
		if !strings.HasPrefix(name, "usb-") {
			continue
		}
		// ignore device partitions, e.g: /dev/sdc1
		if name == previousName {
			continue
		}
		// get path of device for the given device name
		e := fmt.Sprintf("/dev/disk/by-id/%s", f.Name())
		out := pathByDeviceID(e)
		if len(out) == 0 {
			return make([]string, 0)
		}
		arr = append(arr, out)
		previousName = name
	}
	return arr
}

// info returns a Device model for the device path p.
func info(p string) (linuxDevice, error) {
	b := filepath.Base(p)
	vendor := readFile(fmt.Sprintf("/sys/block/%s/device/vendor", b))
	model := readFile(fmt.Sprintf("/sys/block/%s/device/model", b))
	ro := readFile(fmt.Sprintf("/sys/block/%s/ro", b))
	removable := readFile(fmt.Sprintf("/sys/block/%s/removable", b))
	bSize, err := intFromFile(
		fmt.Sprintf("/sys/block/%s/queue/physical_block_size", b))
	if err != nil {
		return linuxDevice{}, fmt.Errorf("E: Could not convert '%v' to int.", bSize)
	}
	sectors, err := intFromFile(fmt.Sprintf("/sys/block/%s/size", b))
	if err != nil {
		return linuxDevice{}, fmt.Errorf("E: Could not convert '%v' to int.", sectors)
	}
	return linuxDevice{
		vendor:    vendor,
		model:     model,
		path:      p,
		sectors:   sectors,
		blockSize: bSize,
		writable:  ro == "0",
		removable: removable == "1",
	}, nil
}

// readFile reads a file with path p and trims space.
func readFile(p string) string {
	f, err := os.Open(p)
	if err != nil {
		fmt.Println("E: Couldn't open file:", p)
		return ""
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("E: Couldn't read file:", p)
		return ""
	}
	out := string(b)
	out = strings.TrimSpace(out)
	return out
}

func intFromFile(p string) (int, error) {
	str := readFile(p)
	size, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("E: Couldn't convert '%s' to int.", str)
	}
	return size, nil
}

// pathByDeviceID returns the device path for the device entry e,
// or in case of an error an empty string.
func pathByDeviceID(e string) string {
	// "/dev/disk/by-id/usb-TrekStor_USB_Stick_QU_AA04013000007552-0:0"
	s, err := os.Readlink(e)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return fmt.Sprintf("/dev/%s", filepath.Base(s))
}

// check if device is readable & writable.
func (d *linuxDevice) check() error {
	switch {
	case len(d.path) == 0:
		return ErrPathEmpty
	case d.sectors <= 0:
		return ErrSectorEmpty
	case d.blockSize <= 0:
		return ErrBlocksizeEmpty
	case !d.writable:
		return ErrNotWritable
	case !d.removable:
		return ErrNotRemovable
		// case len(d.UUID) == 0: // TODO implement UUID on Linux and uc check
		// 	return ErrNoUUID
	}
	return nil
}

func (d *linuxDevice) Vendor() string     { return d.vendor }
func (d *linuxDevice) Model() string      { return d.model }
func (d *linuxDevice) Path() string       { return d.path }
func (d *linuxDevice) Sectors() int       { return d.sectors }
func (d *linuxDevice) BlockSize() int     { return d.blockSize }
func (d *linuxDevice) IsSameDevice() bool { return true } // TODO
func (d *linuxDevice) Unmount() error     { return nil }  // TODO

func (d *linuxDevice) Size() int {
	return d.sectors * d.blockSize
}
