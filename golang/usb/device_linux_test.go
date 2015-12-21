package usb

import (
	"fmt"
	"testing"
)

func newStub() *linuxDevice {
	return &linuxDevice{
		vendor:    "TrekStor",
		model:     "USB Stick QU",
		path:      "/dev/sdc",
		sectors:   31326208,
		blockSize: 512,
		uuid:      "6416385D2C49736A",
	}
}

func BenchmarkList(b *testing.B) {
	for n := 0; n < b.N; n++ {
		List()
	}
}

func TestList(t *testing.T) {
	arr := List()
	if len(arr) == 0 {
		t.Fail()
	}
	for _, d := range arr {
		fmt.Println(d)
	}
}

func TestUUID(t *testing.T) {
	uuid, err := uuid("/dev/sdc")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(uuid)
}

func TestNames(t *testing.T) {
	names := names()
	if len(names) == 0 {
		t.Fail()
	}
	for _, n := range names {
		fmt.Println(n)
	}
}

func TestIsMounted(t *testing.T) {
	d := newStub()
	arr, err := d.isMounted()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(arr)
}

func TestUnmount(t *testing.T) {
	d := newStub()
	if err := d.Unmount(); err != nil {
		t.Error(err)
	}
}
