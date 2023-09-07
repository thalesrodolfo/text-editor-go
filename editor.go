package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Inputable interface {
	addChar(int32)
	moveCursorBy(int)
	removeChar()
}

type MainEditor struct {
	buffer      []string
	line        int
	cursor      rl.Vector2
	cursorIndex int
	state       State
	currentFile string
}

func (e *MainEditor) addChar(c int32) {
	if e.cursorIndex == len(e.buffer[e.line]) {
		e.buffer[e.line] = e.buffer[e.line] + fmt.Sprintf("%c", c)
		e.moveCursorBy(1)
	} else {
		line := e.buffer[e.line]

		beforePart := line[:e.cursorIndex]
		afterPart := line[e.cursorIndex:]

		newText := beforePart + fmt.Sprintf("%c", c) + afterPart

		e.buffer[e.line] = newText
		e.moveCursorBy(1)
	}
}

func (e *MainEditor) moveCursorBy(positions int) {
	e.cursor.X += font.Recs.Width * float32(positions)
	e.cursorIndex += positions
}

func (e *MainEditor) removeChar() {
	fmt.Println("e.cursorIndex: ", e.cursorIndex)
	fmt.Println("len line: ", len(e.buffer[e.line]))

	if e.cursorIndex == len(e.buffer[e.line]) {
		e.buffer[e.line] = e.buffer[e.line][0 : len(e.buffer[e.line])-1]

		e.moveCursorBy(-1)
	} else if e.cursorIndex > 0 {
		line := e.buffer[e.line]

		beforePart := line[:e.cursorIndex-1]
		afterPart := line[e.cursorIndex:]

		e.buffer[e.line] = beforePart + afterPart
		e.moveCursorBy(-1)
	}
}

func (e *MainEditor) addNewLine(fontSize int32) {
	if e.line < len(e.buffer)-1 { // if its not the last line add new line between two lines
		e.buffer = append(e.buffer[:e.line+1], e.buffer[e.line:]...)
		e.buffer[e.line+1] = ""
	} else { // add line at the end
		e.buffer = append(e.buffer, "")
	}

	e.line += 1
	e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(fontSize)*float32(e.line)

	if insideBracket {
		fmt.Println("Inside bracket")
		e.buffer[e.line] = e.buffer[e.line] + "    "
		e.cursorIndex = 4
		e.cursor.X = float32(DEFAULT_LEFT_OFFSET) + float32(e.cursorIndex)*font.Recs.Width
	} else {
		e.cursor.X = float32(DEFAULT_LEFT_OFFSET)
		e.cursorIndex = 0
	}

}

func (e *MainEditor) removeLine() {
	if e.line > 0 {
		e.buffer = append(e.buffer[:e.line], e.buffer[e.line+1:]...)
		e.line -= 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)
		e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))

		newPos := len(e.buffer[e.line])

		if newPos < 0 {
			newPos = 0
		}

		fmt.Println("newPos:", newPos)
		e.cursorIndex = newPos
	}
}

func (e *MainEditor) moveToLineBelow() {
	if e.line < len(e.buffer)-1 {
		e.line += 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)
		if !(e.cursorIndex <= len(e.buffer[e.line])-1) {
			e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))
			e.cursorIndex = len(e.buffer[e.line]) - 1
		}
	}
}

func (e *MainEditor) moveToLineAbove() {
	if e.line > 0 {
		e.line -= 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)

		if !(e.cursorIndex <= len(e.buffer[e.line])-1) {
			e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))
			e.cursorIndex = len(e.buffer[e.line]) - 1
		}
	}
}

type SaveInput struct {
	buffer      []string
	line        int
	cursor      rl.Vector2
	cursorIndex int
}

func (e *SaveInput) addChar(c int32) {
	if e.cursorIndex == len(e.buffer[0]) {
		e.buffer[0] = e.buffer[0] + fmt.Sprintf("%c", c)
		e.moveCursorBy(1)
	} else {
		line := e.buffer[0]

		beforePart := line[:e.cursorIndex]
		afterPart := line[e.cursorIndex:]

		newText := beforePart + fmt.Sprintf("%c", c) + afterPart

		e.buffer[0] = newText
		e.moveCursorBy(1)
	}
}

func (e *SaveInput) moveCursorBy(positions int) {
	e.cursor.X += font.Recs.Width * float32(positions)
	e.cursorIndex += positions
}

func (e *SaveInput) removeChar() {
	fmt.Println("e.cursorIndex: ", e.cursorIndex)
	fmt.Println("len line: ", len(e.buffer[0]))

	if e.cursorIndex == len(e.buffer[0]) {
		e.buffer[0] = e.buffer[0][0 : len(e.buffer[0])-1]

		e.moveCursorBy(-1)
	} else if e.cursorIndex > 0 {
		line := e.buffer[0]

		beforePart := line[:e.cursorIndex-1]
		afterPart := line[e.cursorIndex:]

		e.buffer[0] = beforePart + afterPart
		e.moveCursorBy(-1)
	}
}
