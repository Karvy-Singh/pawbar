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

func inc() {
	if currMonth.String() == "December" {
		currYear += 1
		currMonth = time.Month(1)
	} else {
		currMonth = currMonth + 1
	}
}

func dec() {
	if currMonth.String() == "January" {
		currYear -= 1
		currMonth = time.Month(12)
	} else {
		currMonth = currMonth - 1
	}
}

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
				case vaxis.KeyLeft, vaxis.KeyUp:
					dec()
					draw()
				case vaxis.KeyRight, vaxis.KeyDown:
					inc()
					draw()
				}
			}
		case vaxis.Mouse:
			switch ev.Button {

			case vaxis.MouseWheelDown:
				inc()
				draw()

			case vaxis.MouseWheelUp:
				dec()
				draw()

				// case vaxis.EventPress:
				// vx.Notify("press:", fmt.Sprintf("%d%d", ev.Col, ev.Row))
			}
			// 		case vaxis.FocusOut:
			// 			return 0
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
	day := 1
	for week := 0; week < 6; week++ {
		for weekday := 0; weekday < 7; weekday++ {
			cellIndex := week*7 + weekday
			if cellIndex < offset || day > daysInMonth {
				fmt.Print("   ")
			} else {
				if day == now.Day() && month == now.Month() {
					reverseVideo := "\033[7m"
					reset := "\033[0m"
					fmt.Printf("%s%2d%s ", reverseVideo, day, reset)
				} else if day == now.Day() {
					bgOn := "\033[48;5;243m"
					off := "\033[0m"
					fgOn := "\033[97m"
					fmt.Printf("%s%s%2d%s%s ", bgOn, fgOn, day, off, off)
				} else {
					fmt.Printf("%2d ", day)
				}
				day++
			}
		}
	}
}

func center(s string, width int) string {
	if len(s) >= width {
		return s
	}
	left := (width - len(s)) / 2
	return strings.Repeat(" ", left) + s
}
