package ascii

import "fmt"

type Canvas [][]*Pixel

type Pixel struct {
	Rune        rune
	Flag        uint16
	InitialFlag uint16
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

func (c *Canvas) drawPixel(x, y int, flag uint16) (err error) {
	pixel := (*c)[y][x]
	var newPixel Pixel
	switch {
	case pixel == nil:
		newPixel.Rune = AsciiBitmask[flag]
		newPixel.Flag = flag
		newPixel.InitialFlag = flag
	case pixel.Flag == 0:
		err = fmt.Errorf("Found colision for pixel `%s` with letter `%s` on x: `%d` y: `%d`\n%s\n",
			resolve(flag), string([]rune{pixel.Rune}), x, y, c)
		return
	default:
		newPixel.Rune = AsciiBitmask[flag|pixel.InitialFlag]
		newPixel.Flag = flag | pixel.InitialFlag
		newPixel.InitialFlag = pixel.InitialFlag
		if AsciiBitmask[flag|pixel.InitialFlag] == 0 {
			err = fmt.Errorf("unexpected symbol intersection new `%s` and old `%s`\n%s\n",
				resolve(flag), resolve(pixel.Flag), c)
			fmt.Println(flag, pixel.Flag, x, y, "---", flag|pixel.Flag)
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

func resolve(flag uint16) string {
	return string([]rune{AsciiBitmask[flag]})
}
