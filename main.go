package main

import (
	"strconv"
	"time"

	term "github.com/nsf/termbox-go"
)

var (
	CUR_ROW int = 0
	CUR_COL int = 0 // first row index (1)

	BUFFER map[int][]byte
)

func localizeData() {
	w, h := term.Size()

	for row := 0; row < h; row++ {
		BUFFER[row] = make([]byte, w)
		for col := 0; col < w; col++ {
		}
	}
	additionalInfo(BUFFER[h-1])
}

func display() {
	var fl bool = false
mainloop:
	for {
		switch e := term.PollEvent(); e.Type {
		case term.EventKey:
			fl = keysHandler(e)
			if fl == true {
				break mainloop
			}
			term.SetCursor(CUR_COL, CUR_ROW)
			additionalInfo(BUFFER[len(BUFFER)-1])
			draw()
		case term.EventResize:
			draw()
		default:
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

func draw() {
	w, h := term.Size()

	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			term.SetCell(col, row, rune(BUFFER[row][col]), term.ColorBlack, term.ColorCyan)
		}
	}
	term.Flush()
}

func keysHandler(e term.Event) bool {
	switch e.Key {
	case term.KeyEsc:
		return true

	case term.KeyArrowUp:
		recalcCursorPos(0, -1)

	case term.KeyArrowDown:
		recalcCursorPos(0, 1)

	case term.KeyArrowLeft:
		recalcCursorPos(-1, 0)

	case term.KeyArrowRight:
		recalcCursorPos(1, 0)

	case term.KeyTab:
		printTab()

	case term.KeySpace:
		printSymbol(' ')
	//fix switch to last symbol of prev line when deleting first item of current line
	case term.KeyBackspace:
		leftShift()
		recalcCursorPos(-1, 0)

	default:
		if e.Ch != 0 {
			printSymbol(e.Ch)
		}
	}

	return false
}

func leftShift() {
	var temp byte = BUFFER[CUR_ROW][len(BUFFER[CUR_ROW])-1]
	var ll int = CUR_COL
	if ll != 0 {
		ll--
	}
	if CUR_COL == len(BUFFER[CUR_ROW])-1 {
		temp = 0
	}
	for i := len(BUFFER[CUR_ROW]) - 1; i >= ll; i-- {
		BUFFER[CUR_ROW][i], temp = temp, BUFFER[CUR_ROW][i]
	}
}

func rightShift() {
	var temp byte = BUFFER[CUR_ROW][CUR_COL]
	for i := CUR_COL + 1; i < len(BUFFER[CUR_ROW]); i++ {
		BUFFER[CUR_ROW][i], temp = temp, BUFFER[CUR_ROW][i]
	}
}

func printTab() {
	for i := 0; i < 4; i++ {
		printSymbol(' ')
	}
}

func printSymbol(ch rune) {
	rightShift()
	BUFFER[CUR_ROW][CUR_COL] = byte(ch)
	recalcCursorPos(1, 0)
}

func recalcCursorPos(col, row int) {
	w, _ := term.Size()
	if col != 0 {
		if CUR_ROW > 0 || CUR_COL+col >= 0 {
			if CUR_COL+col < 0 {
				CUR_COL = w - 1
				CUR_ROW--
			} else if CUR_COL+col >= w {
				CUR_COL = 0
				CUR_ROW++
			} else {
				CUR_COL += col
			}
		}
	} else {
		if CUR_ROW+row >= 0 {
			CUR_ROW += row
		}
	}
}

func additionalInfo(arr []byte) {
	for i, _ := range arr {
		arr[i] = 0
	}
	out := strconv.FormatInt(int64(CUR_ROW), 10) + "," + strconv.FormatInt(int64(CUR_COL), 10)
	for i, _ := range out {
		arr[len(arr)/2+i] = out[i]
	}
}

func main() {
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	term.SetInputMode(term.InputCurrent | term.InputEsc)
	term.SetCursor(CUR_COL, CUR_ROW)

	BUFFER = make(map[int][]byte)
	localizeData()
	draw()
	display()
}
