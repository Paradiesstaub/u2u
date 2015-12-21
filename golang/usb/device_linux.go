package usb

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	uuid      string
}

// List returns a slice of writable and removable USB devices.
func List() []Devicer {
	var devices []Devicer
	for _, name := range names() {
		if d, err := device(name); err == nil {
			devices = append(devices, d)
		}
	}
	return devices
}

// names returns a list of all USB device names.
func names() []string {
	files, err := ioutil.ReadDir("/dev/disk/by-id")
	if err != nil {
		return make([]string, 0)
	}
	var names []string
	var previousName string
	for _, f := range files {
		// remove partion suffix
		// example result: usb-TrekStor_USB_Stick_QU_AA04013000007552-0
		name := strings.Split(f.Name(), ":")[0]
		// filter out non-usb drives
		if !strings.HasPrefix(name, "usb-") {
			continue
		}
		// ignore device partitions, e.g: /dev/sdc1
		if name == previousName {
			continue
		}
		names = append(names, f.Name())
		previousName = name
	}
	return names
}

func device(name string) (Devicer, error) {
	path, err := devicePath("/dev/disk/by-id/" + name)
	if err != nil {
		return nil, err
	}
	device, err := details(path)
	if err != nil {
		return nil, err
	}
	return Devicer(&device), nil
}

// devicePath like /dev/sdc for the entry e.
// E.g.: /dev/disk/by-id/usb-TrekStor_USB_Stick_QU_AA04013000007552-0:0
func devicePath(e string) (string, error) {
	path, err := os.Readlink(e)
	if err != nil {
		return "", fmt.Errorf("Empty device path for: %v.\n%v", e, err.Error())
	}
	return fmt.Sprintf("/dev/%s", filepath.Base(path)), nil
}

// details returns the device for the device-path p.
func details(p string) (linuxDevice, error) {
	if len(p) == 0 {
		return linuxDevice{}, ErrPathEmpty
	}
	pre := "/sys/block/" + filepath.Base(p)
	if ro := read(pre + "/ro"); ro == "1" {
		return linuxDevice{}, ErrNotWritable
	}
	if removable := read(pre + "/removable"); removable == "0" {
		return linuxDevice{}, ErrNotRemovable
	}
	bSize, err := strconv.Atoi(read(pre + "/queue/physical_block_size"))
	if err != nil {
		return linuxDevice{}, ErrBlocksizeEmpty
	}
	sectors, err := strconv.Atoi(read(pre + "/size"))
	if err != nil {
		return linuxDevice{}, ErrSectorEmpty
	}
	uuid, err := uuid(p)
	if err != nil {
		return linuxDevice{}, err
	}
	return linuxDevice{
		vendor:    read(pre + "/device/vendor"),
		model:     read(pre + "/device/model"),
		path:      p,
		sectors:   sectors,
		blockSize: bSize,
		uuid:      uuid,
	}, nil
}

// read file with path p and trim space.
func read(p string) string {
	f, err := os.Open(p)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(bytes.TrimSpace(b))
}

// uuid returns the UUID for the device dev. E.g: /dev/sdc
func uuid(dev string) (string, error) {
	dev = filepath.Base(dev)
	byuuid := "/dev/disk/by-uuid"
	files, err := ioutil.ReadDir(byuuid)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		s, err := os.Readlink(byuuid + "/" + f.Name())
		if err != nil {
			return "", err
		}
		if len(s) == 0 {
			continue
		}
		s = filepath.Base(s[:len(s)-1])
		if s == dev {
			return f.Name(), nil
		}
	}
	return "", errors.New("Could not found UUID for: " + dev)
}

func (d *linuxDevice) Vendor() string { return d.vendor }
func (d *linuxDevice) Model() string  { return d.model }
func (d *linuxDevice) Path() string   { return d.path }
func (d *linuxDevice) Sectors() int   { return d.sectors }
func (d *linuxDevice) BlockSize() int { return d.blockSize }

func (d *linuxDevice) IsSameDevice() bool {
	if uuid, err := uuid(d.path); err == nil {
		return d.uuid == uuid
	}
	return false
}

func (d *linuxDevice) Unmount() error {
	arr, err := d.isMounted()
	if err != nil {
		return err
	}
	if len(arr) == 0 {
		return nil
	}
	for _, dev := range arr {
		if err = exec.Command("umount", dev).Run(); err != nil {
			return errors.New("Can not unmount: " + dev)
		}
	}
	return nil
}

// isMounted returns a slice of mounted partitions.
func (d *linuxDevice) isMounted() ([]string, error) {
	s := read("/proc/self/mountinfo")
	if len(s) == 0 {
		return []string{}, errors.New("Mount information not readable")
	}
	var arr []string
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		block := strings.Split(line, " - ")[1]
		for _, segment := range strings.Split(block, " ") {
			if strings.HasPrefix(segment, d.path) {
				arr = append(arr, segment)
			}
		}
	}
	return arr, nil
}

func (d *linuxDevice) String() string {
	return fmt.Sprintf("Vendor: %s\nModel: %s\nPath: %s\nSectors: %d\n"+
		"BlockSize: %d\nUUID: %v",
		d.vendor,
		d.model,
		d.path,
		d.sectors,
		d.blockSize,
		d.uuid)
}
