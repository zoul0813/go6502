package Debug

import (
	"fmt"
)

const (
	EmptyColor = "\033[1;90m%s\033[0m"
	HiColor    = "\033[1;31m%s\033[0m"
	LoColor    = "\033[1;91m%s\033[0m"
	AddrColor  = "\033[1;32m%s\033[0m"
	DebugColor = "\033[0;37m%s\033[0m"
)

func Colorize(color string, format string, v any) string {
	f := fmt.Sprintf(format, v)
	return fmt.Sprintf(color, f)
}
