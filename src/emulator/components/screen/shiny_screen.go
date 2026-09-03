package screen

import (
	"image/color"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/paint"
)

type shinyScreen struct {
	window *screen.Window
	buffer *screen.Buffer
}

func NewShinyScreen(window *screen.Window, buffer *screen.Buffer) *shinyScreen {
	return &shinyScreen{
		window: window,
		buffer: buffer,
	}
}

func (s *shinyScreen) ShowImage(buffer *[][]uint32) {
	img := (*s.buffer).RGBA()

	for y, scanline := range *buffer {
		for x, pixel := range scanline {
			c := color.RGBA{

				R: uint8(pixel >> 16),
				G: uint8(pixel >> 8),
				B: uint8(pixel),
				A: 255,
			}
			img.SetRGBA(x, y, c)
		}
	}

	(*s.window).Send(paint.Event{})
}
