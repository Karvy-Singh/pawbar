package tui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"git.sr.ht/~rockorager/vaxis"
	"github.com/nekorg/katnip"
)

var (
	now       = time.Now()
	currYear  = now.Year()
	currMonth = now.Month()
)

func Panel(k *katnip.Kitty, rw io.ReadWriter) int {
	vx, err := vaxis.New(vaxis.Options{
		WithTTY:         os.Stdout.Name(),
		EnableSGRPixels: true,
	})
	if err != nil {
		return 1
	}
	defer vx.Close()

	draw := func() {
		win := vx.Window()
		win.Clear()
		PrintMonthCal(currYear, currMonth)

		vx.Render()
	}
	draw()

	for ev := range vx.Events() {
		switch ev := ev.(type) {
		case vaxis.Key:
			if ev.EventType == vaxis.EventPress {
				switch ev.Keycode {
				case vaxis.KeyEsc:
					return 0
				case vaxis.KeyLeft:
					if currMonth.String() == "January" {
						currYear -= 1
						currMonth = time.Month(12)
					} else {
						currMonth = currMonth - 1
					}
					draw()
				case vaxis.KeyRight:
					if currMonth.String() == "December" {
						currYear += 1
						currMonth = time.Month(1)
					} else {
						currMonth = currMonth + 1
					}
					draw()
				}
			}
		case vaxis.Resize, vaxis.Redraw:
			draw()
		}
	}
	return 0
}

func PrintMonthCal(year int, month time.Month) {
	const width = 20

	title := fmt.Sprintf("%s %d", month, year)
	fmt.Println(center(title, width))
	fmt.Println("Su Mo Tu We Th Fr Sa")

	loc := time.Now().Location()
	first := time.Date(year, month, 1, 0, 0, 0, 0, loc)
	offset := int(first.Weekday())
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, loc).Day()
	for i := 0; i < offset; i++ {
		fmt.Print("   ") // 2 digits + space
	}
	for day := 1; day <= daysInMonth; day++ {
		fmt.Printf("%2d ", day)
		wd := (offset + (day - 1)) % 7
		if wd == 6 {
			fmt.Println()
		}
	}
	if (offset+daysInMonth-1)%7 != 6 {
		fmt.Println()
	}
}

func center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	left := (width - len(s)) / 2
	return strings.Repeat(" ", left) + s
}
