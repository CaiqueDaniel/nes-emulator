package tests

import (
	"image"
)

type screenFixture struct {
	image           *image.RGBA
	ShowImageCalls int
	LastBuffer      *[][]uint32
}

func NewScreenFixture() *screenFixture {
	return &screenFixture{
		image: image.NewRGBA(image.Rect(0, 0, 256, 240)),
	}
}

func (s *screenFixture) ShowImage(buffer *[][]uint32) {
	s.ShowImageCalls++
	s.LastBuffer = buffer
}
