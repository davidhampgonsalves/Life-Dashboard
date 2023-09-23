package draw

import (
	"davidhampgonsalves/lifedashboard/pkg/event"
	"fmt"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/tdewolff/canvas"
)

const Width = 600.0
const Height = 800.0
const EventGap = 20
const Margin = 40.0
const FontSize = 90

func Init() (*canvas.Canvas, *canvas.Context, *canvas.FontFamily) {
	c := canvas.New(Width, Height)
	ctx := canvas.NewContext(c)

	font := canvas.NewFontFamily("AH")
	if err := font.LoadFontFile("static/Atkinson-Hyperlegible.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	return c, ctx, font
}

func Background(ctx *canvas.Context) {
	ctx.SetFillColor(canvas.White)
	ctx.MoveTo(0, 0)
	ctx.LineTo(Width, 0)
	ctx.LineTo(Width, Height)
	ctx.LineTo(0, Height)
	ctx.LineTo(0, 0)
	ctx.Close()
	ctx.FillStroke()

	bMargin := 5.0
	ctx.SetStrokeColor(canvas.Black)
	ctx.MoveTo(bMargin, bMargin)
	ctx.LineTo(Width-bMargin, float64(bMargin))
	ctx.LineTo(Width-bMargin, Height-bMargin)
	ctx.LineTo(bMargin, Height-bMargin)
	ctx.LineTo(bMargin, bMargin)
	ctx.Close()
	ctx.FillStroke()
}

func Date(ctx *canvas.Context, font *canvas.FontFamily, y float64) {
	face := font.Face(120, canvas.Black, canvas.FontRegular)
	numberFace := font.Face(200, canvas.FontRegular, canvas.FontNormal)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	files, err := ioutil.ReadDir("static/pokemon-sprites")
	if err != nil {
		fmt.Printf("Could not open sprite directory\n")
	}
	name := files[r1.Intn(len(files))].Name()
	icon, err := os.Open(fmt.Sprintf("static/pokemon-sprites/%s", name))
	if err != nil {
		fmt.Printf("sprite not found - %s\n", name)
	}
	img, err := png.Decode(icon)
	if err != nil {
		fmt.Printf("Couldn't decode sprite: %s\n", name)
	}
	imgDPMM := 0.75
	imgWidth := float64(img.Bounds().Max.X) / imgDPMM
	imgHeight := float64(img.Bounds().Max.Y) / imgDPMM

	now := time.Now()
	rt := canvas.NewRichText(numberFace)
	rt.WriteString(fmt.Sprintf("%d\n", now.Day()))
	rt.SetFace(face)
	rt.WriteString(now.Month().String()[:3])
	text := rt.ToText(0.0, 0.0, canvas.Center, canvas.Bottom, 0.0, 0.0)

	ctx.DrawText(Width-Margin-(text.Bounds().W/2), Margin+((imgHeight-text.Bounds().H)/2), text)
	ctx.DrawImage(Width-text.Bounds().W-imgWidth-Margin, Margin, img, canvas.DPMM(imgDPMM))
}

func Events(ctx *canvas.Context, font *canvas.FontFamily, events []event.Event) float64 {
	face := font.Face(FontSize, canvas.Black, canvas.FontRegular)
	y := float64(Height - Margin)
	for idx, e := range events {
		fmt.Printf("drawing event %d @%f\n", idx+1, y)
		y -= Event(face, ctx, e, y) + EventGap
	}
	return y
}

func Event(face *canvas.FontFace, ctx *canvas.Context, e event.Event, y float64) float64 {
	rt := canvas.NewRichText(face)
	for _, r := range e.Text {
		if r < 128 {
			rt.WriteString(string(r))
		} else {
			fmt.Printf("searching for png for rune %x\n", r)
			icon, err := os.Open(fmt.Sprintf("static/noto-emoji/emoji_u%x.png", r))
			if err != nil {
				fmt.Printf("rune not found - %x\n", r)
				continue
			}
			pngIcon, err := png.Decode(icon)
			if err != nil {
				fmt.Printf("icon couldn't be opened as png - %x\n", r)
				continue
			}
			// DPMM controlls the icon size (lower means bigger)
			rt.WriteImage(pngIcon, canvas.DPMM(0.8), canvas.FontMiddle)
		}
	}

	// should this be henight or some fixed max cut off?
	text := rt.ToText(Width-(Margin*2), Height, canvas.Left, canvas.Right, 0.0, 0.0)
	ctx.DrawText(Margin, y, text)

	return text.Bounds().H
}
