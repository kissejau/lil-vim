package main

import (
	"strconv"
	"time"

	term "github.com/nsf/termbox-go"
)

var (
	SIDE_INDENT int = 2
	NUM_INDENT  int = 3
	CUR_ROW     int = 0
	CUR_COL     int = SIDE_INDENT + NUM_INDENT + 1 // first row index (1)

	BUFFER map[int][]byte
)

func localizeData() {
	w, h := term.Size()

	for row := 0; row < h; row++ {
		BUFFER[row] = make([]byte, w)
		numberingBufferRows(row)
		for col := 0; col < w; col++ {
		}
	}
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
			term.SetCell(col, row, rune(BUFFER[row][col]), term.AttrBlink, term.AttrBlink)
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

	default:
		if e.Ch != 0 {
			printSymbol(e.Ch)
		}
	}

	return false
}

func printSymbol(ch rune) {
	BUFFER[CUR_ROW][CUR_COL] = byte(ch)
	recalcCursorPos(1, 0)
}

func recalcCursorPos(col, row int) {
	w, _ := term.Size()
	if col != 0 {
		if CUR_COL+col < NUM_INDENT+SIDE_INDENT+len(strconv.FormatInt(int64(CUR_ROW), 10)) && CUR_ROW > 0 {
			CUR_ROW--
			CUR_COL = w - 1
		} else if CUR_COL+col >= w {
			CUR_ROW++
			CUR_COL = NUM_INDENT + SIDE_INDENT + len(strconv.FormatInt(int64(CUR_ROW), 10))
		} else {
			CUR_COL += col
		}
	} else {
		if CUR_ROW+row >= 0 {
			CUR_ROW += row
		}
	}
}

func numberingBufferRows(rowIndex int) {
	colIndex := 0
	for _, e := range strconv.FormatInt(int64(rowIndex), 10) {
		BUFFER[rowIndex][colIndex+SIDE_INDENT] = byte(e)
		colIndex++
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
