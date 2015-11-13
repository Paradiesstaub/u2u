package usb

import "fmt"

const (
	KIBIBYTE ByteSize = 1024
	MEBIBYTE ByteSize = 1024 * KIBIBYTE
	GIBIBYTE ByteSize = 1024 * MEBIBYTE
	TEBIBYTE ByteSize = 1024 * GIBIBYTE
)

type ByteSize float64

// ToHuman returns ByteSize in a human readable form (e.g.: 1,1 GiB).
func (b ByteSize) ToHuman() string {
	if b < MEBIBYTE {
		return fmt.Sprintf("%.1f KiB", (b / KIBIBYTE))
	} else if b < GIBIBYTE {
		return fmt.Sprintf("%.1f MiB", (b / MEBIBYTE))
	} else if b < TEBIBYTE {
		return fmt.Sprintf("%.1f GiB", (b / GIBIBYTE))
	}
	return fmt.Sprintf("%.1f TiB", (b / TEBIBYTE))
}
