package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	lineHeight          int32 = 0
	font                rl.Font
	DEFAULT_LEFT_OFFSET int32 = 100
	DEFAULT_TOP_OFFSET  int32 = 10
	insideBracket       bool  = false
)

//var msg strings.Builder

type Editor struct {
	buffer      []string
	line        int
	cursor      rl.Vector2
	cursorIndex int
}

func (e Editor) numberOfLines() int {
	return len(e.buffer)
}

func (e *Editor) addChar(c int32, font rl.Font) {

	// cursor is at the end of the line
	if e.cursorIndex == len(e.buffer[e.line]) {
		e.buffer[e.line] = e.buffer[e.line] + fmt.Sprintf("%c", c)
		e.moveCursorBy(1)
	} else {
		line := e.buffer[e.line]

		beforePart := line[:e.cursorIndex]
		afterPart := line[e.cursorIndex:]

		fmt.Println("beforePart: ", beforePart)
		fmt.Println("afterPart: ", afterPart)

		newText := beforePart + fmt.Sprintf("%c", c) + afterPart

		fmt.Println(newText)

		e.buffer[e.line] = newText
		e.moveCursorBy(1)
	}
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

func (e *Editor) removeLine(font rl.Font) {
	if e.line > 0 {
		e.buffer = append(e.buffer[:e.line], e.buffer[e.line+1:]...)
		e.line -= 1
		e.cursor.Y = float32(DEFAULT_TOP_OFFSET) + float32(font.BaseSize)*float32(e.line)
		e.cursor.X = float32(DEFAULT_LEFT_OFFSET + int32(len(e.buffer[e.line]))*int32(font.Recs.Width))

		newPos := len(e.buffer[e.line]) - 1

		if newPos < 0 {
			newPos = 0
		}

		fmt.Println("newPos:", newPos)
		e.cursorIndex = newPos
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
	screenWidth := int32(1200)
	screenHeight := int32(800)

	keywords := make(map[string][]KeywordPos)

	editor := new(Editor)
	editor.buffer = append(editor.buffer, "")

	rl.InitWindow(screenWidth, screenHeight, "Simple Text Editor - Golang")

	rl.SetTargetFPS(60)

	font = rl.LoadFontEx("fonts/JetBrainsMono-Regular.ttf", 24, nil)

	DEFAULT_LEFT_OFFSET = 6 * int32(font.Recs.Width)

	fontSize := font.BaseSize
	lineHeight = fontSize
	fontPosition := rl.NewVector2(float32(DEFAULT_LEFT_OFFSET), float32(DEFAULT_TOP_OFFSET))

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

			if k == 123 { // {
				fmt.Println("open bracket")
				insideBracket = true
			}

			if k == 125 { // }
				fmt.Println("close bracket")
				insideBracket = false

				// if we close bracket in a new line, remove identation
				if strings.TrimLeft(editor.buffer[editor.line], " ") == "}" {
					editor.buffer[editor.line] = strings.TrimLeft(editor.buffer[editor.line], " ")
					editor.cursor.X = float32(DEFAULT_LEFT_OFFSET) + font.Recs.Width
					editor.cursorIndex = 1
				}
			}
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
			fmt.Println("cursor pos:", editor.cursorIndex)
			fmt.Println("keywords:", keywords)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		blinkCursor(&blink, &cursorColor)
		rl.DrawRectangle(int32(editor.cursor.X), int32(editor.cursor.Y), int32(font.Recs.Width), int32(font.Recs.Height), cursorColor)

		drawLineNumbers(editor)

		for i, v := range editor.buffer {
			// 	// draw text
			linePos := fontPosition

			if i > 0 {
				linePos.Y = float32(DEFAULT_TOP_OFFSET + lineHeight*int32(i))
			}

			checkKeywords(editor, v, keywords)
			checkTypes(editor, v, keywords)

			for j, ch := range v {
				if j > 0 {
					linePos.X = float32(DEFAULT_LEFT_OFFSET) + font.Recs.Width*float32(j)
				}

				text_color := rl.Black

				for _, key := range getKeywords() {
					for _, pos := range keywords[key] {

						if j >= pos.init && j <= pos.end {
							text_color = rl.Purple
						}

					}
				}

				for _, key := range getTypes() {
					for _, pos := range keywords[key] {

						if j >= pos.init && j <= pos.end {
							text_color = rl.Blue
						}

					}
				}
				rl.DrawTextEx(font, string(ch), linePos, float32(fontSize), 0, text_color)
			}

		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func getFormatedLineNumber(i int) string {
	if i < 10 {
		return fmt.Sprintf("   %d", i)
	} else if i < 100 {
		return fmt.Sprintf("  %d", i)
	} else if i < 1000 {
		return fmt.Sprintf(" %d", i)
	} else {
		return fmt.Sprintf("%d", i)
	}
}

func drawLineNumbers(editor *Editor) {
	//background
	//rl.DrawRectangle(0, 0, 40, int32(rl.GetScreenHeight()), rl.NewColor(222, 222, 222, 222))
	pos := rl.NewVector2(10, float32(DEFAULT_TOP_OFFSET))
	i := 1

	for i <= editor.numberOfLines() {

		rl.DrawTextEx(font, getFormatedLineNumber(i), pos, float32(font.BaseSize), 0, rl.Gray)

		pos.Y = float32(DEFAULT_TOP_OFFSET) + font.Recs.Height*float32(i)

		i++
	}

}

func checkKeywords(editor *Editor, text string, keywords map[string][]KeywordPos) {

	for _, key := range getKeywords() {

		if strings.Contains(text, key) {
			index := strings.Index(text, key)

			posAlreadyExists := false

			if keywords[key] == nil {
				keywords[key] = append(keywords[key], KeywordPos{editor.line, index, index + (len(key) - 1)})
			} else {
				i := 0

				for i < len(keywords[key]) {
					if keywords[key][i].init == index {
						posAlreadyExists = true
					}
					i++
				}

				if !posAlreadyExists {
					// add pos
					keywords[key] = append(keywords[key], KeywordPos{editor.line, index, index + (len(key) - 1)})
				}
			}

		} else {
			keywords[key] = nil
		}

	}
}

func checkTypes(editor *Editor, text string, keywords map[string][]KeywordPos) {

	for _, key := range getTypes() {

		if strings.Contains(text, key) {
			index := strings.Index(text, key)

			posAlreadyExists := false

			if keywords[key] == nil {
				keywords[key] = append(keywords[key], KeywordPos{editor.line, index, index + (len(key) - 1)})
			} else {
				i := 0

				for i < len(keywords[key]) {
					if keywords[key][i].init == index {
						posAlreadyExists = true
					}
					i++
				}

				if !posAlreadyExists {
					// add pos
					keywords[key] = append(keywords[key], KeywordPos{editor.line, index, index + (len(key) - 1)})
				}
			}

		} else {
			keywords[key] = nil
		}

	}
}
