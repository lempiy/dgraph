package ascii

import "fmt"

type Canvas [][]uint8

func NewCanvas(width uint, height uint) *Canvas {
	canvas := make(Canvas, height)
	for i := range canvas {
		canvas[i] = make([]uint8, width)
	}
	return &canvas
}

func (c *Canvas) drawPixel(x, y uint, p uint8) (err error) {
	pixel := (*c)[y][x]
	newPixel := p | pixel
	if AsciiBitmask[newPixel] == 0 {
		err = fmt.Errorf("unexpected symbol intersection `%s` | `%s`",
			resolve(p), resolve(pixel))
		return
	}
	(*c)[y][x] = newPixel
	return
}

func resolve(flag uint8) string {
	return string([]rune{AsciiBitmask[flag]})
}
