package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type State int32

const (
	NORMAL State = iota
	SAVING
)

var (
	lineHeight          int32 = 0
	font                rl.Font
	DEFAULT_LEFT_OFFSET int32  = 100
	DEFAULT_TOP_OFFSET  int32  = 10
	insideBracket       bool   = false
	state               State  = NORMAL
	file                string = ""
)

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

func saveFile(path string, editor *MainEditor) {
	file, _ := os.Create(path)
	defer file.Close()

	for _, v := range editor.buffer {
		file.WriteString(v)
		file.WriteString("\n")
	}

}

func main() {
	regexString, _ := regexp.Compile(`"(.*?)"|"(.*?)$`)
	screenWidth := int32(1200)
	screenHeight := int32(800)

	rl.InitWindow(screenWidth, screenHeight, "Simple Text Editor - Golang")
	rl.SetTargetFPS(120)
	font = rl.LoadFontEx("fonts/JetBrainsMono-Regular.ttf", 32, nil)
	DEFAULT_LEFT_OFFSET = 6 * int32(font.Recs.Width)
	fontSize := font.BaseSize
	lineHeight = fontSize
	fontPosition := rl.NewVector2(float32(DEFAULT_LEFT_OFFSET), float32(DEFAULT_TOP_OFFSET))

	// keywords for syntax highlighting
	keywords := make(map[string][]KeywordPos)

	// main editor
	mainEditor := new(MainEditor)
	mainEditor.buffer = append(mainEditor.buffer, "")
	mainEditor.cursor = fontPosition

	// save input when saving file
	saveEditor := new(MainEditor)
	saveEditor.buffer = append(saveEditor.buffer, "")
	saveEditor.cursor = rl.NewVector2(0, float32(DEFAULT_TOP_OFFSET))

	// main blink
	blink := 0.0

	// secondary blink (save input)
	blinkSearch := 0.0
	cursorColor := rl.NewColor(60, 60, 60, 200)

	backspace_timer := 0.0

	for !rl.WindowShouldClose() {
		blink += float64(rl.GetFrameTime())
		blinkSearch += float64(rl.GetFrameTime())
		backspace_timer += float64(rl.GetFrameTime())

		if state == NORMAL {
			mainEditorHandle(mainEditor, font, blink, cursorColor, &backspace_timer)

		} else {
			mainEditorHandle(saveEditor, font, blinkSearch, cursorColor, &backspace_timer)

			if rl.IsKeyPressed(rl.KeyEnter) {
				rootPath, _ := os.Getwd()
				path := rootPath + string(os.PathSeparator) + saveEditor.buffer[0]
				fmt.Println(path)
				saveFile(path, mainEditor)
				fmt.Println("File saved.")
				file = path
				state = NORMAL
			}
		}
		if rl.IsKeyPressed(rl.KeyF5) { // enter saving mode
			if file == "" { // if not file loaded, open input to save new file
				if state == SAVING {
					state = NORMAL
				} else {
					state = SAVING
				}
			} else { // if file is loaded, just save current file
				saveFile(file, mainEditor)
				fmt.Println("File saved.")
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.DrawFPS(screenWidth-100, screenHeight-20)

		if state == SAVING {
			label := "Save File: "
			yPosition := screenHeight - 40
			xStartPosition := int32(font.Recs.Width) * int32(len(label)+2)
			// draw search background
			rl.DrawRectangle(20, yPosition, screenWidth-40, font.Recs.ToInt32().Height, rl.NewColor(200, 200, 45, 255))
			// blink cursor
			blinkCursor(&blinkSearch, &cursorColor)
			// draw cursor block
			rl.DrawRectangle(xStartPosition+int32(saveEditor.cursor.X), yPosition, int32(font.Recs.Width), int32(font.Recs.Height), cursorColor)
			// draw save label
			rl.DrawTextEx(font, label, rl.NewVector2(30, float32(yPosition)), float32(fontSize), 0, rl.Black)
			// draw save buffer
			rl.DrawTextEx(font, saveEditor.buffer[0], rl.NewVector2(float32(xStartPosition+4), float32(yPosition)), float32(fontSize), 0, rl.Black)
		} else {
			blinkCursor(&blink, &cursorColor)
			rl.DrawRectangle(int32(mainEditor.cursor.X), int32(mainEditor.cursor.Y), int32(font.Recs.Width), int32(font.Recs.Height), cursorColor)
		}

		drawLineNumbers(mainEditor)

		for i, v := range mainEditor.buffer {
			// 	// draw text
			linePos := fontPosition

			if i > 0 {
				linePos.Y = float32(DEFAULT_TOP_OFFSET + lineHeight*int32(i))
			}

			checkKeywords(mainEditor, v, keywords)
			checkTypes(mainEditor, v, keywords)

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

				list_pos := regexString.FindAllStringIndex(v, 2)

				for _, pos := range list_pos {
					if len(pos) > 1 && j >= pos[0] && j <= pos[1]-1 {
						text_color = rl.DarkGreen
					}
				}

				rl.DrawTextEx(font, string(ch), linePos, float32(fontSize), 0, text_color)
			}

		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}

func mainEditorHandle(editor *MainEditor, font rl.Font, blink float64, cursorColor rl.Color, backspace_timer *float64) {
	k := rl.GetCharPressed()
	if k > 0 {
		editor.addChar(k)
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

	if rl.IsKeyPressed(rl.KeyEnd) {
		editor.moveCursorBy(len(editor.buffer[editor.line]) - editor.cursorIndex)
	}
	if rl.IsKeyPressed(rl.KeyHome) {
		editor.moveCursorBy(editor.cursorIndex * -1)
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		editor.addNewLine(font.BaseSize)
	}

	if rl.IsKeyDown(rl.KeyBackspace) && *backspace_timer > float64(0.1) {

		if len(editor.buffer[editor.line]) > 0 {
			editor.removeChar()
		} else {
			editor.removeLine()
			fmt.Println(editor.cursor)
		}

		stopBlink(&blink, &cursorColor)
		*backspace_timer = float64(0.0)

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
		//fmt.Println("keywords:", keywords)
		//pos := regexString.FindStringIndex(editor.buffer[editor.line])
		//fmt.Println(pos)
	}
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

func drawLineNumbers(editor *MainEditor) {
	//background
	//rl.DrawRectangle(0, 0, 40, int32(rl.GetScreenHeight()), rl.NewColor(222, 222, 222, 222))
	pos := rl.NewVector2(10, float32(DEFAULT_TOP_OFFSET))
	i := 1

	for i <= len(editor.buffer) {

		rl.DrawTextEx(font, getFormatedLineNumber(i), pos, float32(font.BaseSize), 0, rl.Gray)

		pos.Y = float32(DEFAULT_TOP_OFFSET) + font.Recs.Height*float32(i)

		i++
	}

}

func checkKeywords(editor *MainEditor, text string, keywords map[string][]KeywordPos) {

	//regexString, _ := regexp.Compile("\"\w*\"?")

	//match, _ := regexp.String("p([a-z]+)ch", "peach")

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

func checkTypes(editor *MainEditor, text string, keywords map[string][]KeywordPos) {

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
