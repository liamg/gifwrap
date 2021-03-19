package ascii

import (
	"fmt"
	"image"
	"image/gif"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Renderer struct {
	screen tcell.Screen
	image  *gif.GIF
	width  int
	height int
}

var ErrQuit = fmt.Errorf("user quit")

func (r *Renderer) init() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		return err
	}
	r.screen = screen
	r.width, r.height = r.screen.Size()
	return nil
}

func (r *Renderer) close() {
	r.screen.Fini()
}

func (r *Renderer) Play() error {
	if err := r.init(); err != nil {
		return err
	}
	defer r.close()

	for {
		if err := r.cycleFrames(); err != nil {
			return err
		}
	}
}

func (r *Renderer) PlayOnce() error {
	if err := r.init(); err != nil {
		return err
	}
	defer r.close()

	return r.cycleFrames()
}

func (r *Renderer) cycleFrames() error {
	for i, frame := range r.image.Image {
		if err := r.drawFrame(frame); err != nil {
			return err
		}

		_ = r.screen.PostEvent(nil)
		switch ev := r.screen.PollEvent().(type) {
		case *tcell.EventResize:
			r.width, r.height = ev.Size()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return ErrQuit
			}
		}

		delay := time.Duration(r.image.Delay[i]) * time.Millisecond * 10
		time.Sleep(delay)
	}
	return nil
}

func (r *Renderer) drawFrame(img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	pixPerCellX := width / r.width
	pixPerCellY := height / r.height

	var red uint64
	var green uint64
	var blue uint64

	count := uint64(pixPerCellX * pixPerCellY)

	r.screen.Clear()

	for x := 0; x < r.width; x++ {
		for y := 0; y < r.height; y++ {
			red, green, blue = 0, 0, 0
			for pX := x * pixPerCellX; pX < (x*pixPerCellX)+pixPerCellX; pX++ {
				for pY := y * pixPerCellY; pY < (y*pixPerCellY)+pixPerCellY; pY++ {
					colour := img.At(pX, pY)
					r, g, b, a := colour.RGBA()
					red += uint64((r * 255) / a)
					green += uint64((g * 255) / a)
					blue += uint64((b * 255) / a)
				}
			}

			r.screen.SetCell(x, y, tcell.StyleDefault.Background(tcell.NewRGBColor(int32(red/count), int32(green/count), int32(blue/count))), ' ')
		}
	}

	r.screen.Show()

	return nil
}
