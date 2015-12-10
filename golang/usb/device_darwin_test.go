package usb

import (
	"fmt"
	"testing"

	"github.com/Paradiesstaub/u2u/golang/root"
)

func TestSudoCmd(t *testing.T) {
	s := root.NewExecutor().Exec("touch /tmp/test/foo2.txt", "")
	fmt.Println(s)
}

func TestMapToDevice(t *testing.T) {
	m := parseDiskutil(TEXT)
	d, err := mapToDevice(m)
	if err != nil {
		t.Error(err)
	}
	if d.vendor != "" {
		t.Error(`d.vendor != ""`)
	}
	if d.model != "Untitled 1" {
		t.Error(`d.model != "Untitled 1"`)
	}
	if d.path != "/dev/disk3s1" {
		t.Error(`d.path != "/dev/disk3s1"`)
	}
	if d.sectors != 31307776 {
		t.Error("d.sectors != 31307776")
	}
	if d.blockSize != 512 {
		t.Error("d.blockSize != 512")
	}
	if !d.writable {
		t.Error("!d.writable")
	}
	if !d.removable {
		t.Error("!d.removable")
	}
	if d.uuid != "36039B3C-2241-4C61-B79C-96212398CD44" {
		t.Error(`d.path != "36039B3C-2241-4C61-B79C-96212398CD44"`)
	}
}

// func TestDeviceName(t *testing.T) {
// 	s := deviceName("/dev/disk3s1")
// 	fmt.Println(s)
// }

// func TestIsMounted(t *testing.T) {
// 	s := deviceInfo("/dev/disk3s1")
// 	m := parseDiskutil(s)
// 	d, err := mapToDevice(m)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	b := d.IsMounted()
// 	fmt.Println(b)
// }
//
// func TestUnmount(t *testing.T) {
// 	//s := deviceInfo(`/Volumes/16\ GB`)
// 	s := deviceInfo("/dev/disk3s1")
// 	m := parseDiskutil(s)
// 	d, err := mapToDevice(m)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if err = d.Unmount(); err != nil {
// 		fmt.Println("#### 1")
// 		t.Error(err)
// 	}
// 	if err = d.Unmount(); err != nil {
// 		fmt.Println("#### 2")
// 		t.Error(err)
// 	}
// }

// func TestList(t *testing.T) {
// 	arr := List()
// 	if len(arr) != 1 {
// 		t.Errorf("List len should be 1, but is %d.\narr: %v", len(arr), arr)
// 	}
// }

const TEXT = `
Device Identifier:        disk3s1
Device Node:              /dev/disk3s1
Whole:                    No
Part of Whole:            disk3
Device / Media Name:      Untitled 1

Volume Name:

Mounted:                  Yes
Mount Point:              /Volumes/Untitled

File System Personality:  NTFS
Type (Bundle):            ntfs
Name (User Visible):      Windows NT File System (NTFS)

Partition Type:           Microsoft Basic Data
OS Can Be Installed:      No
Media Type:               Generic
Protocol:                 USB
SMART Status:             Not Supported
Disk / Partition UUID:    36039B3C-2241-4C61-B79C-96212398CD44

Total Size:               16.0 GB (16029581312 Bytes) (exactly 31307776 512-Byte-Units)
Volume Free Space:        16.0 GB (15961518080 Bytes) (exactly 31174840 512-Byte-Units)
Device Block Size:        512 Bytes
Allocation Block Size:    4096 Bytes

Read-Only Media:          No
Read-Only Volume:         Yes

Device Location:          External
Removable Media:          Yes
Media Removal:            Software-Activated
`
