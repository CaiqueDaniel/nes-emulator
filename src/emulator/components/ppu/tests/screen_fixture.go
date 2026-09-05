package tests

import (
	"image"
)

type screenFixture struct {
	image *image.RGBA
}

func NewScreenFixture() *screenFixture {
	return &screenFixture{
		image: image.NewRGBA(image.Rect(0, 0, 256, 240)),
	}
}

func (s *screenFixture) ShowImage(buffer *[][]uint32) {
}
