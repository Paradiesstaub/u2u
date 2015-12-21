package usb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Device contains USB device information.
type macDevice struct {
	vendor    string
	model     string
	path      string // device path
	sectors   int    // numbers of sectors
	blockSize int    // physical block size in bytes
	writable  bool
	removable bool
	uuid      string // currently only set on Mac
}

const (
	Vendor    = "Volume Name"
	Model     = "Device / Media Name"
	Path      = "Device Node"
	Sector    = "Total Size"
	BlockSize = "Device Block Size"
	ReadOnly  = "Read-Only Media"
	Removable = "Removable Media"
	UUID      = "Disk / Partition UUID"
	True      = "Yes"
	False     = "No"
	IsUsb     = "Protocol"
	Mounted   = "Mounted"
)

// http://www.ubuntu.com/download/desktop/create-a-usb-stick-on-mac-osx

// TODO conv code to use interface Devicer

// List returns a slice of writable USB devices.
func List() []Devicer {
	files, err := ioutil.ReadDir("/Volumes")
	if err != nil {
		return make([]Devicer, 0)
	}
	c := make(chan result, 5)
	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go func(info os.FileInfo) {
			c <- device(info)
			wg.Done()
		}(f)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	var devices []Devicer
	for r := range c {
		if r.ok {
			devices = append(devices, Devicer(&r.d))
		}
	}
	return devices
}

type result struct {
	d  macDevice
	ok bool
}

func device(f os.FileInfo) result {
	if !f.IsDir() {
		return result{}
	}
	p := fmt.Sprintf("/Volumes/%s", f.Name())
	s := deviceInfo(p)
	if len(s) == 0 {
		return result{}
	}
	m := parseDiskutil(s)
	if m[IsUsb] != "USB" {
		return result{}
	}
	d, err := mapToDevice(m)
	if err != nil {
		return result{}
	}
	if err := d.check(); err != nil {
		return result{}
	}
	if model, err := deviceName(d.path); err == nil {
		d.model = model
	}
	return result{d: d, ok: true}
}

func deviceInfo(p string) string {
	cmd := fmt.Sprintf("%s %s %s", "diskutil", "info", p)
	b, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

func parseDiskutil(s string) map[string]string {
	m := make(map[string]string)
	arr := strings.Split(s, "\n")
	for _, row := range arr {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		var key, value string
		for i, cell := range strings.Split(row, ":") {
			cell = strings.TrimSpace(cell)
			switch i {
			case 0:
				key = cell
			case 1:
				value = cell
			default:
				break
			}
		}
		m[key] = value
	}
	return m
}

func mapToDevice(m map[string]string) (macDevice, error) {
	s := strings.Split(m[Sector], " ")[5]
	sectors, err := strconv.Atoi(s)
	if err != nil {
		return macDevice{}, fmt.Errorf("E: Could not convert %q to int.", s)
	}
	s = strings.Split(m[BlockSize], " ")[0]
	blocksize, err := strconv.Atoi(s)
	if err != nil {
		return macDevice{}, fmt.Errorf("E: Could not convert %q to int.", s)
	}
	return macDevice{
		vendor:    m[Vendor],
		model:     m[Model],
		path:      m[Path],
		sectors:   sectors,
		blockSize: blocksize,
		writable:  m[ReadOnly] == False,
		removable: m[Removable] == True,
		uuid:      m[UUID],
	}, nil
}

// check if device is readable & writable.
func (d *macDevice) check() error {
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
	case len(d.uuid) == 0:
		return ErrNoUUID
	}
	return nil
}

// deviceName returns the name of the device of the partition p.
func deviceName(p string) (string, error) {
	p, err := devicePath(p)
	if err != nil {
		return "", err
	}
	s := deviceInfo(p)
	m := parseDiskutil(s)
	return m[Model], nil
}

// DevicePath returns the path of the device for the partition p.
// E.g.: /dev/disk3s1 -> /dev/disk3
func devicePath(p string) (string, error) {
	if len(p) < 3 {
		return "", errors.New("Device path too short.")
	}
	return p[:len(p)-2], nil
}

// isMounted returns true if the device is mounted.
func (d *macDevice) isMounted() bool {
	s := deviceInfo(d.path)
	m := parseDiskutil(s)
	return m[Mounted] == True
}

// getUUID returns the ID of the device with the given path d.Path.
func (d *macDevice) getUUID() string {
	s := deviceInfo(d.path)
	m := parseDiskutil(s)
	return m[UUID]
}

func (d *macDevice) Vendor() string { return d.vendor }
func (d *macDevice) Model() string  { return d.model }
func (d *macDevice) Path() string   { return d.path }
func (d *macDevice) Sectors() int   { return d.sectors }
func (d *macDevice) BlockSize() int { return d.blockSize }

func (d *macDevice) IsSameDevice() bool {
	return d.uuid == d.getUUID()
}

func (d *macDevice) Unmount() error {
	if !d.isMounted() {
		return nil
	}
	cmd := fmt.Sprintf("diskutil umount %s", d.path)
	return exec.Command("sh", "-c", cmd).Run()
}

// system_profiler SPUSBDataType -xml
// diskutil info /Volumes/Untitled
