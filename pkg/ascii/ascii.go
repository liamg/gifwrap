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
	fill   bool
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

func (r *Renderer) SetFill(fill bool) {
	r.fill = fill
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
			if ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'q':
					return ErrQuit
				}
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

	termWidth := r.width
	termHeight := r.height

	if !r.fill {

		imgRatio := float64(width) / float64(height)
		termRatio := float64(r.width) / float64(r.height)

		if termRatio > imgRatio {
			termWidth = int(float64(termHeight) * imgRatio)
		} else {
			termHeight = int(float64(termWidth) / imgRatio)
		}

	}

	pixPerCellX := width / termWidth
	pixPerCellY := height / termHeight

	var red uint64
	var green uint64
	var blue uint64

	count := uint64(pixPerCellX * pixPerCellY)
	if count == 0 {
		return nil
	}

	r.screen.Clear()

	for x := 0; x < termWidth; x++ {
		for y := 0; y < termHeight; y++ {
			red, green, blue = 0, 0, 0
			for pX := x * pixPerCellX; pX < (x*pixPerCellX)+pixPerCellX; pX++ {
				for pY := y * pixPerCellY; pY < (y*pixPerCellY)+pixPerCellY; pY++ {
					colour := img.At(pX, pY)
					ir, ig, ib, ia := colour.RGBA()
					ia = ia / 0xff
					r := (uint64(ia) * uint64(ir)) / 0xffff
					g := (uint64(ia) * uint64(ig)) / 0xffff
					b := (uint64(ia) * uint64(ib)) / 0xffff
					if r > 0xff {
						r = 0xff
					}
					if g > 0xff {
						g = 0xff
					}
					if b > 0xff {
						b = 0xff
					}
					red += r
					green += g
					blue += b
				}
			}

			cr, cg, cb := int32(red/count), int32(green/count), int32(blue/count)
			r.screen.SetCell(x, y, tcell.StyleDefault.Background(tcell.NewRGBColor(cr, cg, cb)), ' ')
		}
	}

	r.screen.Show()

	return nil
}
