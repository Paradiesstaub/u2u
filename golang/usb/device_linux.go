package usb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// List returns a slice of writable and removable USB devices.
// Or an empty slice if no device could be found or an error occurred.
func List() []Device {
	arr := FindDevices()
	if len(arr) == 0 {
		return make([]Device, 0)
	}
	var devices []Device
	for _, p := range arr {
		d, err := info(p)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !d.Writable {
			continue
		}
		if !d.Removable {
			continue
		}
		devices = append(devices, d)
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
func info(p string) (Device, error) {
	b := filepath.Base(p)
	vendor := readFile(fmt.Sprintf("/sys/block/%s/device/vendor", b))
	model := readFile(fmt.Sprintf("/sys/block/%s/device/model", b))
	ro := readFile(fmt.Sprintf("/sys/block/%s/ro", b))
	removable := readFile(fmt.Sprintf("/sys/block/%s/removable", b))
	bSize, err := intFromFile(
		fmt.Sprintf("/sys/block/%s/queue/physical_block_size", b))
	if err != nil {
		return Device{}, fmt.Errorf("E: Couldn't convert '%v' to int.", bSize)
	}
	sectors, err := intFromFile(fmt.Sprintf("/sys/block/%s/size", b))
	if err != nil {
		return Device{}, fmt.Errorf("E: Couldn't convert '%v' to int.", sectors)
	}
	return Device{
		Vendor:    vendor,
		Model:     model,
		Path:      p,
		Sectors:   sectors,
		BlockSize: bSize,
		Writable:  ro == "0",
		Removable: removable == "1",
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
