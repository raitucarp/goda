package main

import (
	"log"

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
	root        *goda.Node
	winW        int
	winH        int
	mouseX      int
	mouseY      int
	activeMenu  int
	needsLayout bool
	prevW       int
	prevH       int
}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	g.mouseX = mx
	g.mouseY = my

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.root != nil && handleClick(g.root, 0, 0, g) {
			g.needsLayout = true
		}
	}
	return nil
}

func handleClick(n *goda.Node, absLeft, absTop float32, g *Game) bool {
	lo := n.LayoutOut()
	x := absLeft + float32(lo.Left)
	y := absTop + float32(lo.Top)
	w := float32(lo.Width)
	h := float32(lo.Height)

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
		if handleClick(n.GetChild(i), absLeft+float32(lo.Left), absTop+float32(lo.Top), g) {
			return true
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.winW <= 0 || g.winH <= 0 {
		return
	}
	if g.needsLayout || g.root == nil {
		g.rebuildLayout()
		g.needsLayout = false
	}
	if g.root == nil {
		return
	}
	screen.Fill(hex("#F0F2F5"))
	renderTree(screen, g.root, 0, 0, g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth <= 0 {
		outsideWidth = 800
	}
	if outsideHeight <= 0 {
		outsideHeight = 600
	}
	if outsideWidth != g.prevW || outsideHeight != g.prevH {
		g.prevW = outsideWidth
		g.prevH = outsideHeight
		g.winW = outsideWidth
		g.winH = outsideHeight
		g.needsLayout = true
	}
	if g.winW == 0 {
		g.winW = outsideWidth
		g.winH = outsideHeight
		g.needsLayout = true
	}
	return outsideWidth, outsideHeight
}

func (g *Game) rebuildLayout() {
	g.root = buildLayout(float32(g.winW), float32(g.winH), g.activeMenu)
	goda.CalculateNodeLayout(g.root, float32(g.winW), float32(g.winH), goda.DirectionLTR)
}

func renderTree(screen *ebiten.Image, n *goda.Node, absLeft, absTop float32, g *Game) {
	lo := n.LayoutOut()
	x := absLeft + float32(lo.Left)
	y := absTop + float32(lo.Top)
	w := float32(lo.Width)
	h := float32(lo.Height)

	mx := float32(g.mouseX)
	my := float32(g.mouseY)
	isHovered := w > 0 && h > 0 && mx >= x && mx < x+w && my >= y && my < y+h

	if w > 0 && h > 0 {
		if widget, ok := n.GetContext().(Widget); ok {
			renderWidget(screen, x, y, w, h, widget, isHovered)
		}
	}

	for i := 0; i < n.GetChildCount(); i++ {
		renderTree(screen, n.GetChild(i), x, y, g)
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
