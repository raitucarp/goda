package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"

	goda "goda"
)

var (
	fontFace   font.Face
	fontFaceSm font.Face
)

func initFont() {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatalf("font parse: %v", err)
	}
	fontFace, _ = opentype.NewFace(tt, &opentype.FaceOptions{Size: 13, DPI: 72})
	fontFaceSm, _ = opentype.NewFace(tt, &opentype.FaceOptions{Size: 11, DPI: 72})
}

type Game struct {
	root       *goda.Node
	winW       int
	winH       int
	designW    float32
	designH    float32
	mouseX     int
	mouseY     int
	activeMenu int
	needsLayout bool
}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	g.mouseX = mx
	g.mouseY = my

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		scaleX := float32(g.winW) / g.designW
		scaleY := float32(g.winH) / g.designH
		if handleClick(g.root, 0, 0, scaleX, scaleY, g) {
			g.needsLayout = true
		}
	}
	return nil
}

func handleClick(n *goda.Node, absLeft, absTop, scaleX, scaleY float32, g *Game) bool {
	lo := n.LayoutOut()
	x := (absLeft + float32(lo.Left)) * scaleX
	y := (absTop + float32(lo.Top)) * scaleY
	w := float32(lo.Width) * scaleX
	h := float32(lo.Height) * scaleY

	mx := float32(g.mouseX)
	my := float32(g.mouseY)

	if w > 0 && h > 0 && mx >= x && mx < x+w && my >= y && my < y+h {
		if widget, ok := n.GetContext().(Widget); ok {
			if widget.Kind == wMenuItem {
				if md, ok := widget.Data.(MenuData); ok && !md.Active {
					g.activeMenu = md.Index
					return true
				}
			}
		}
	}

	for i := 0; i < n.GetChildCount(); i++ {
		if handleClick(n.GetChild(i), absLeft+float32(lo.Left), absTop+float32(lo.Top), scaleX, scaleY, g) {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.winW <= 0 || g.winH <= 0 {
		return
	}
	if g.needsLayout {
		g.rebuildLayout()
		g.needsLayout = false
	}

	scaleX := float32(g.winW) / g.designW
	scaleY := float32(g.winH) / g.designH
	screen.Fill(hex("#F0F2F5"))
	renderTree(screen, g.root, 0, 0, scaleX, scaleY, g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.winW = outsideWidth
	g.winH = outsideHeight
	if outsideWidth <= 0 {
		outsideWidth = 1280
	}
	if outsideHeight <= 0 {
		outsideHeight = 800
	}
	g.designW = float32(outsideWidth)
	g.designH = float32(outsideHeight)
	g.rebuildLayout()
	return outsideWidth, outsideHeight
}

func (g *Game) rebuildLayout() {
	g.root = buildLayout(g.designW, g.designH, g.activeMenu)
	goda.CalculateNodeLayout(g.root, g.designW, g.designH, goda.DirectionLTR)
}

func renderTree(screen *ebiten.Image, n *goda.Node, absLeft, absTop, scaleX, scaleY float32, g *Game) {
	lo := n.LayoutOut()
	x := absLeft + float32(lo.Left)
	y := absTop + float32(lo.Top)
	w := float32(lo.Width)
	h := float32(lo.Height)

	rx := x * scaleX
	ry := y * scaleY
	rw := w * scaleX
	rh := h * scaleY

	mx := float32(g.mouseX)
	my := float32(g.mouseY)
	isHovered := rw > 0 && rh > 0 && mx >= rx && mx < rx+rw && my >= ry && my < ry+rh

	if rw > 1 && rh > 1 {
		if widget, ok := n.GetContext().(Widget); ok {
			renderWidget(screen, rx, ry, rw, rh, widget, isHovered)
		}
	}

	for i := 0; i < n.GetChildCount(); i++ {
		renderTree(screen, n.GetChild(i), x, y, scaleX, scaleY, g)
	}
}

func main() {
	initFont()
	ebiten.SetWindowSize(1280, 800)
	ebiten.SetWindowTitle("Dashboard — goda + Ebitengine")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func init() { _ = math.NaN }
