package format

import (
	"fmt"
)

// SizeFormat is used to change the way how fmt.Print prints the file size. It's String() method transforms
// the size, which should be in bytes, to Kibibytes, Mebibytes or Gibibytes.
type SizeFormat int64

func (s SizeFormat) String() string {
	switch {
	case s < 1<<10:
		return fmt.Sprintf("%v B", float64(s))
	case s < 1<<20:
		return fmt.Sprintf("%.2f KiB", float64(s)/(1<<10))
	case s < 1<<30:
		return fmt.Sprintf("%.2f MiB", float64(s)/(1<<20))
	default:
		return fmt.Sprintf("%.2f GiB", float64(s)/(1<<30))
	}
}