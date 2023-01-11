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
)

func display() {
	var fl bool = false
mainloop:
	for {
		switch e := term.PollEvent(); e.Type {
		case term.EventKey:
			fl = keysHandler(e.Key)
			if fl == true {
				break mainloop
			}
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
		numbering(row)
		for col := 0; col < w; col++ {

		}
	}
	term.Flush()
}

func keysHandler(key term.Key) bool {
	if key == term.KeyEsc {
		return true
	}

	if key == term.KeyArrowUp && CUR_ROW > 0 {
		recalcCursorPos(0, -1)
	} else if key == term.KeyArrowDown {
		recalcCursorPos(0, 1)
	} else if key == term.KeyArrowLeft {
		recalcCursorPos(-1, 0)
	} else if key == term.KeyArrowRight {
		recalcCursorPos(1, 0)
	}

	term.SetCursor(CUR_COL, CUR_ROW)
	return false
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

func numbering(rowIndex int) {
	colIndex := 0
	for _, e := range strconv.FormatInt(int64(rowIndex), 10) {
		term.SetCell(colIndex+SIDE_INDENT, rowIndex, e, term.AttrBlink, term.AttrBlink)
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

	draw()
	display()
}
