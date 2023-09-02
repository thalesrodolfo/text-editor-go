package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var msg strings.Builder

type Editor struct {
	buffer []string
	line int
	cursor rl.Vector2
}



func (e *Editor) addChar(c int32, font rl.Font) {
	e.buffer = append(e.buffer, fmt.Sprintf("%c", c))
	e.cursor.X += font.Recs.Width
}

func (e *Editor) removeChar(font rl.Font) {
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

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)

	editor := new(Editor)

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")

	rl.SetTargetFPS(60)

	font := rl.LoadFontEx("fonts/JetBrainsMono-Regular.ttf", 96, nil)

	fontSize := font.BaseSize
	fontPosition := rl.NewVector2(40, float32(10))

	editor.cursor = fontPosition

	fmt.Println(font.Recs.Width)

	blink := 0.0
	cursorColor := rl.NewColor(60, 60, 60, 200)

	for !rl.WindowShouldClose() {
		blink += float64(rl.GetFrameTime())


		k := rl.GetCharPressed()
		if k > 0 {
			msg.WriteString(fmt.Sprintf("%c", k))
			editor.addChar(k, font)
			editor.String()
		}

		if rl.IsKeyPressed(rl.KeyEnter) {
			msg.WriteString("\n")
		}

		if rl.IsKeyPressed(rl.KeyBackspace) {

			if msg.Len() > 0 {
				a := msg.String()[0 : msg.Len()-1]
				msg.Reset()
				msg.WriteString(a)
				fmt.Println(a)
				editor.removeChar(font)
			}
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		fmt.Println(blink)
		if blink >= 0.5 {

			if cursorColor.A == 0.0 {
				cursorColor.A = 200
			} else {
				cursorColor.A = 0.0
			}

			blink = 0.0
		}

		rl.DrawRectangle(int32(editor.cursor.X), int32(editor.cursor.Y), fontSize/2, fontSize, cursorColor)

		rl.DrawTextEx(font, msg.String(), fontPosition, float32(fontSize), 0, rl.Black)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
