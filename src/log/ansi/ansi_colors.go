// ANSI escape code sequences for pretty logging
package ansi

import (
	"fmt"
)

type Color string

const (
	ColorReset         Color = "\033[0m"
	ColorBlack         Color = "30"
	ColorRed           Color = "31"
	ColorGreen         Color = "32"
	ColorYellow        Color = "33"
	ColorBlue          Color = "34"
	ColorMagenta       Color = "35"
	ColorCyan          Color = "36"
	ColorLightGrey     Color = "37"
	ColorGrey          Color = "90"
	ColorBrightRed     Color = "91"
	ColorBrightGreen   Color = "92"
	ColorBrightYellow  Color = "93"
	ColorBrightBlue    Color = "94"
	ColorBrightMagenta Color = "95"
	ColorBrightCyan    Color = "96"
	ColorWhite         Color = "97"
)

func Colorize(colorCode Color, message string) string {
	return fmt.Sprintf("\033[%sm%s%s", colorCode, message, ColorReset)
}
