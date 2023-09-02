package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//var msg strings.Builder

type Editor struct {
	buffer []string
	line   int
	cursor rl.Vector2
}

func (e *Editor) addChar(c int32, font rl.Font) {
	fmt.Println(e.buffer[e.line])
	e.buffer[e.line] = e.buffer[e.line] + fmt.Sprintf("%c", c)
	e.cursor.X += font.Recs.Width
}

func (e *Editor) removeChar(font rl.Font) {
	e.buffer[e.line] = e.buffer[e.line][0 : len(e.buffer[e.line])-1]
	e.cursor.X -= font.Recs.Width
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

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)

	font := rl.LoadFontEx("fonts/JetBrainsMono-Regular.ttf", 96, nil)

	fontSize := font.BaseSize
	fontPosition := rl.NewVector2(40, float32(10))

	editor.cursor = fontPosition

	blink := 0.0
	cursorColor := rl.NewColor(60, 60, 60, 200)

	for !rl.WindowShouldClose() {
		blink += float64(rl.GetFrameTime())

		k := rl.GetCharPressed()
		if k > 0 {
			//msg.WriteString(fmt.Sprintf("%c", k))
			editor.addChar(k, font)
			stopBlink(&blink, &cursorColor)
			//editor.String()
			fmt.Println(editor.buffer[0])
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			//msg.WriteString("\n")
		}

		if rl.IsKeyPressed(rl.KeyBackspace) {

			if len(editor.buffer[editor.line]) > 0 {
				editor.removeChar(font)
				stopBlink(&blink, &cursorColor)
			}

		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		blinkCursor(&blink, &cursorColor)
		rl.DrawRectangle(int32(editor.cursor.X), int32(editor.cursor.Y), fontSize/2, fontSize, cursorColor)

		for _, v := range editor.buffer {
			// 	// draw text
			rl.DrawTextEx(font, v, fontPosition, float32(fontSize), 0, rl.Black)
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
