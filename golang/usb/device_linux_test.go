package usb

import (
	"fmt"
	"testing"
)

// go test -run=BenchmarkList -bench=.
func BenchmarkList(b *testing.B) {
	for n := 0; n < b.N; n++ {
		List()
	}
}

// TODO
func TestList(t *testing.T) {
	l := List()
	fmt.Println(l)

	// for _, d := range FindDevices() {
	// 	fmt.Println(d)
	// }
}
