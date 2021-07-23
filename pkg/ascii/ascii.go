package ascii

import (
	"fmt"
	"image"
	"image/color"
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
	canvas map[uint64]struct{ r, g, b int32 }
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
	r.canvas = map[uint64]struct{ r, g, b int32 }{}
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

		if err := r.drawFrame(frame, i); err != nil {
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

func (r *Renderer) drawFrame(img image.Image, i int) error {

	bounds := img.Bounds()
	_ = bounds
	width := r.image.Config.Width
	height := r.image.Config.Height

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
	var alpha uint64

	count := uint64(pixPerCellX * pixPerCellY)
	if count == 0 {
		return nil
	}

	//r.screen.Clear()

	var disposal byte
	if i == 0 {
		disposal = r.image.Disposal[len(r.image.Image)-1]
	} else {
		disposal = r.image.Disposal[i-1]
	}

	background := r.image.Config.ColorModel.(color.Palette)[r.image.BackgroundIndex]

	for x := 0; x < termWidth; x++ {
		for y := 0; y < termHeight; y++ {
			red, green, blue, alpha = 0, 0, 0, 0
			for pX := x * pixPerCellX; pX < (x*pixPerCellX)+pixPerCellX; pX++ {
				for pY := y * pixPerCellY; pY < (y*pixPerCellY)+pixPerCellY; pY++ {
					if pX < bounds.Min.X || pY < bounds.Min.Y {
						continue
					}
					colour := img.At(pX-bounds.Min.X, pY-bounds.Min.Y) //-bounds.Min.X, pY-bounds.Min.Y)
					ir, ig, ib, ia := colour.RGBA()
					a := ia / 0xff
					r := ir / 0xff
					g := ig / 0xff
					b := ib / 0xff
					if r > 0xff {
						r = 0xff
					}
					if g > 0xff {
						g = 0xff
					}
					if b > 0xff {
						b = 0xff
					}
					if a > 0xff {
						a = 0xff
					}
					red += uint64(r)
					green += uint64(g)
					blue += uint64(b)
					alpha += uint64(a)
				}
			}

			cr, cg, cb, ca := int32(red/count), int32(green/count), int32(blue/count), int32(alpha/count)
			var force bool
			if disposal == gif.DisposalBackground {
				cr, cg, cb, _ := background.RGBA()
				r.screen.SetCell(x, y, tcell.StyleDefault.Background(tcell.NewRGBColor(int32(cr), int32(cg), int32(cb))), '?')
			} else if disposal == gif.DisposalPrevious {
				// lol?
				force = true
				//r.screen.SetCell(x, y, tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 0, 0)), '7')
			} else {
				//r.screen.SetCell(x, y, tcell.StyleDefault.Foreground(tcell.NewRGBColor(0, 255, 0)), '8')
			}
			// if prev, found := r.canvas[(uint64(x)<<16)+uint64(y)]; found && ca < 0x16 {
			// 	//				cr = (int32((uint64(cr) * uint64(ca)) / 0xff)) + (int32((uint64(prev.r) * (0xff - uint64(ca))) / 0xff))
			// 	//		cg = (int32((uint64(cg) * uint64(ca)) / 0xff)) + (int32((uint64(prev.g) * (0xff - uint64(ca))) / 0xff))
			// 	//cb = (int32((uint64(cb) * uint64(ca)) / 0xff)) + (int32((uint64(prev.b) * (0xff - uint64(ca))) / 0xff))
			// 	cr = prev.r
			// 	cg = prev.g
			// 	cb = prev.b
			// 	_ = ca
			// }

			r.canvas[(uint64(x)<<16)+uint64(y)] = struct{ r, g, b int32 }{cr, cg, cb}
			if ca > 0x80 || force {
				r.screen.SetCell(x, y, tcell.StyleDefault.Background(tcell.NewRGBColor(cr, cg, cb)), ' ')
			}
		}
	}

	r.screen.Show()

	return nil
}
