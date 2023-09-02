package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DEFAULT_LEFT_OFFSET = 40
	DEFAULT_TOP_OFFSET  = 40
)

var (
	lineHeight int32 = 0
	font       rl.Font
)

//var msg strings.Builder

type Editor struct {
	buffer      []string
	line        int
	cursor      rl.Vector2
	cursorIndex int
}

func (e *Editor) addChar(c int32, font rl.Font) {
	fmt.Println(e.buffer[e.line])
	e.buffer[e.line] = e.buffer[e.line] + fmt.Sprintf("%c", c)
	e.moveCursorBy(1)
}

func (e *Editor) moveCursorBy(positions int) {
	e.cursor.X += font.Recs.Width * float32(positions)
	e.cursorIndex += positions
}

func (e *Editor) removeChar(font rl.Font) {
	e.buffer[e.line] = e.buffer[e.line][0 : len(e.buffer[e.line])-1]
	e.moveCursorBy(-1)
}

func (e Editor) String() {
	i := 0

	fmt.Println("Line:", e.line)
	fmt.Printf("Cursor (%f, %f)\n", e.cursor.X, e.cursor.Y)

	for i < len(e.buffer) {
		fmt.Print(e.buffer[i])

		i++
	}

	fmt.Println()
}

func (e *Editor) addNewLine(fontSize int32) {
	if e.line < len(e.buffer)-1 {
		e.buffer = append(e.buffer[:e.line+1], e.buffer[e.line:]...)
		e.buffer[e.line+1] = ""
	} else {
		e.buffer = append(e.buffer, "")
	}

	e.line += 1
	e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(fontSize)*float32(e.line)
	e.cursor.X = DEFAULT_LEFT_OFFSET
	e.cursorIndex = 0
}

func (e *Editor) removeLine(font rl.Font) {
	if e.line > 0 {
		e.line -= 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)
		e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))
		e.cursorIndex = len(e.buffer[e.line]) - 1
	}
}

func (e *Editor) moveToLineBelow() {
	if e.line < len(e.buffer)-1 {
		e.line += 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)
		if !(e.cursorIndex <= len(e.buffer[e.line])-1) {
			e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))
			e.cursorIndex = len(e.buffer[e.line]) - 1
		}
	}
}

func (e *Editor) moveToLineAbove() {
	if e.line > 0 {
		e.line -= 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)

		if !(e.cursorIndex <= len(e.buffer[e.line])-1) {
			e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))
			e.cursorIndex = len(e.buffer[e.line]) - 1
		}
	}
}

func blinkCursor(blink *float64, cursorColor *rl.Color) {
	if *blink >= 0.5 {

		if cursorColor.A == 0.0 {
			cursorColor.A = 200
		} else {
			cursorColor.A = 0.0
		}

		*blink = 0.0
	}
}

func stopBlink(blink *float64, color *rl.Color) {
	*blink = 0.0
	*color = rl.NewColor(60, 60, 60, 200)
}

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	editor := new(Editor)
	editor.buffer = append(editor.buffer, "")

	rl.InitWindow(screenWidth, screenHeight, "Simple Text Editor - Golang")

	rl.SetTargetFPS(60)

	font = rl.LoadFontEx("fonts/JetBrainsMono-Regular.ttf", 40, nil)

	fontSize := font.BaseSize
	lineHeight = fontSize
	fontPosition := rl.NewVector2(DEFAULT_LEFT_OFFSET, DEFAULT_TOP_OFFSET)

	editor.cursor = fontPosition

	blink := 0.0
	cursorColor := rl.NewColor(60, 60, 60, 200)

	backspace_timer := 0.0

	for !rl.WindowShouldClose() {
		blink += float64(rl.GetFrameTime())
		backspace_timer += float64(rl.GetFrameTime())

		k := rl.GetCharPressed()
		if k > 0 {
			editor.addChar(k, font)
			stopBlink(&blink, &cursorColor)
			fmt.Println(font.Recs.Width)
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			fmt.Println(editor.cursor)
			editor.addNewLine(fontSize)
		}

		if rl.IsKeyDown(rl.KeyBackspace) && backspace_timer > 0.1 {

			if len(editor.buffer[editor.line]) > 0 {
				editor.removeChar(font)
			} else {
				editor.removeLine(font)
				fmt.Println(editor.cursor)
			}

			stopBlink(&blink, &cursorColor)
			backspace_timer = 0.0

		}

		if rl.IsKeyPressed(rl.KeyLeft) {
			fmt.Println("back", editor.cursorIndex)
			if editor.cursorIndex > 0 {
				editor.moveCursorBy(-1)
			}

		}

		if rl.IsKeyPressed(rl.KeyRight) {
			fmt.Println("forward", editor.cursorIndex)
			if editor.cursorIndex <= len(editor.buffer[editor.line])-1 {
				editor.moveCursorBy(1)
			}
		}

		if rl.IsKeyPressed(rl.KeyUp) {
			editor.moveToLineAbove()
			fmt.Println("up", editor.line)
		}

		if rl.IsKeyPressed(rl.KeyDown) {
			editor.moveToLineBelow()
			fmt.Println("down", editor.line)
		}

		if rl.IsKeyPressed(rl.KeyF1) {
			fmt.Println("buffer:", editor.buffer)
			fmt.Println("line:", editor.line)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		blinkCursor(&blink, &cursorColor)
		rl.DrawRectangle(int32(editor.cursor.X), int32(editor.cursor.Y), int32(font.Recs.Width), int32(font.Recs.Height), cursorColor)

		for i, v := range editor.buffer {
			// 	// draw text
			linePos := fontPosition

			if i > 0 {
				linePos.Y = float32(DEFAULT_TOP_OFFSET + lineHeight*int32(i))
			}

			rl.DrawTextEx(font, v, linePos, float32(fontSize), 0, rl.Black)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
