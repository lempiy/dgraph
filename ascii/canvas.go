package ascii

import "fmt"

type Canvas [][]*Pixel

type Pixel struct {
	Rune rune
	Flag uint8
}

func NewCanvas(width int, height int) *Canvas {
	canvas := make(Canvas, height)
	for i := range canvas {
		canvas[i] = make([]*Pixel, width)
	}
	return &canvas
}

func (c *Canvas) drawLetter(x, y int, letter rune) (err error) {
	pixel := (*c)[y][x]
	if pixel != nil {
		err = fmt.Errorf("Found colision for letter `%s` with `%s` on x: `%d` y: `%d`\n%s\n",
			string([]rune{letter}), string([]rune{pixel.Rune}), x, y, c)
		return
	}
	newPixel := Pixel{
		Rune: letter,
	}
	(*c)[y][x] = &newPixel
	return
}

func (c *Canvas) drawPixel(x, y int, flag uint8) (err error) {
	pixel := (*c)[y][x]
	var newPixel Pixel
	switch {
	case pixel == nil:
		newPixel.Rune = AsciiBitmask[flag]
		newPixel.Flag = flag
	case pixel.Flag == 0:
		err = fmt.Errorf("Found colision for pixel `%s` with letter `%s` on x: `%d` y: `%d`\n%s\n",
			resolve(flag), string([]rune{pixel.Rune}), x, y, c)
		return
	default:
		newPixel.Rune = AsciiBitmask[flag|pixel.Flag]
		newPixel.Flag = flag | pixel.Flag
		if AsciiBitmask[flag|pixel.Flag] == 0 {
			err = fmt.Errorf("unexpected symbol intersection `%s` | `%s`\n%s\n",
				resolve(newPixel.Flag), resolve(pixel.Flag), c)
			return
		}
	}
	(*c)[y][x] = &newPixel
	return
}

func (c *Canvas) String() (acc string) {
	for _, row := range *c {
		for _, cell := range row {
			if cell == nil {
				acc += " "
				continue
			}
			acc += string([]rune{cell.Rune})
		}
		acc += "\n"
	}
	return
}

func resolve(flag uint8) string {
	return string([]rune{AsciiBitmask[flag]})
}
