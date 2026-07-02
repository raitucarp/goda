package main

import (
	"fmt"
	"image/color"
	"math"

	gg "github.com/fogleman/gg"
	goda "goda"
)

const (
	roleHeader   = "header"
	roleLogo     = "logo"
	roleSearch   = "search"
	roleCart     = "cart"
	roleBanner   = "banner"
	roleSidebar  = "sidebar"
	roleCategory = "category"
	roleCard     = "card"
	roleFooter   = "footer"
)

func main() {
	const W, H = 800, 1080

	builds := map[string]func() *goda.Node{
		"builder":    buildByBuilder,
		"cssstring":  buildByCSSString,
		"cssmap":     buildByCSSMap,
		"renderfrom": buildByRenderFrom,
	}

	for name, build := range builds {
		fmt.Printf("=== Building with %s ===\n", name)
		root := build()

		goda.CalculateNodeLayout(root, W, H, goda.DirectionLTR)

		dc := gg.NewContext(W, H)
		dc.SetRGB(0.94, 0.94, 0.95)
		dc.Clear()
		renderNode(dc, root, 0, 0)

		filename := fmt.Sprintf("ecommerce_%s.png", name)
		if err := dc.SavePNG(filename); err != nil {
			fmt.Printf("Error saving %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("Saved %s (800x1080)\n", filename)
	}
}

func newNode(role string) *goda.Node {
	n := goda.New()
	n.SetContext(role)
	return n
}

func newNodeWithID(id, role string) *goda.Node {
	n := goda.New(id)
	n.SetContext(role)
	return n
}

var products = []struct {
	name     string
	price    string
	oldPrice string
	rating   string
	badge    string
	imgColor string
	imgShape string
}{
	{"Wireless Headphones Pro", "$79.99", "$129.99", "4.8", "SALE", "#E8F0FE", "headphones"},
	{"Smart Watch Ultra", "$199.99", "$249.99", "4.9", "HOT", "#FCE8E6", "watch"},
	{"Bluetooth Speaker", "$49.99", "$69.99", "4.5", "-29%", "#E6F4EA", "speaker"},
	{"USB-C Hub 7-in-1", "$34.99", "", "4.3", "", "#FFF3E0", "hub"},
	{"Mechanical Keyboard RGB", "$89.99", "$119.99", "4.7", "NEW", "#F3E5F5", "keyboard"},
	{"4K Webcam Pro", "$129.99", "", "4.6", "", "#E8EAF6", "camera"},
}

var categories = []struct{ name, icon string }{
	{"Electronics", ""},
	{"Clothing & Fashion", ""},
	{"Home & Garden", ""},
	{"Sports & Outdoors", ""},
	{"Books & Media", ""},
	{"Toys & Games", ""},
}

func renderNode(dc *gg.Context, n *goda.Node, absLeft, absTop float64) {
	lo := n.LayoutOut()
	x := absLeft + float64(lo.Left)
	y := absTop + float64(lo.Top)
	w := float64(lo.Width)
	h := float64(lo.Height)
	if w <= 0 || h <= 0 {
		return
	}

	role, _ := n.GetContext().(string)

	switch role {
	case roleHeader:
		drawBox(dc, x, y, w, h, hex(0xFFFFFF), hex(0xDDDDDD), 8)

	case roleLogo:
		drawBox(dc, x, y, w, h, hex(0x1A73E8), hex(0x1A73E8), 8)
		dc.SetColor(color.White)
		dc.DrawStringAnchored("ShopStore", x+w/2, y+h/2, 0.5, 0.4)

	case roleSearch:
		drawBox(dc, x, y, w, h, hex(0xF1F3F4), hex(0xDADCE0), h/2)
		dc.SetRGBA255(148, 148, 148, 255)
		dc.DrawStringAnchored("Search products, brands...", x+18, y+h/2, 0, 0.45)

	case roleCart:
		drawBox(dc, x, y, w, h, hex(0xF8F9FA), hex(0x1A73E8), 8)
		dc.SetRGBA255(26, 115, 232, 255)
		dc.DrawStringAnchored("Cart", x+w/2, y+h/2, 0.5, 0.4)

	case roleBanner:
		grad := gg.NewLinearGradient(x, y, x+w, y+h)
		grad.AddColorStop(0, hex(0x1A73E8))
		grad.AddColorStop(0.5, hex(0x4285F4))
		grad.AddColorStop(1, hex(0x174EA6))
		dc.SetFillStyle(grad)
		drawRoundRect(dc, x, y, w, h, 12)
		dc.Fill()
		dc.SetColor(color.White)
		dc.DrawStringAnchored("SUMMER SALE", x+w/2, y+h/2-22, 0.5, 0.5)
		dc.SetRGBA255(189, 208, 255, 255)
		dc.DrawStringAnchored("Get up to 50% off on selected items  Free shipping on orders over $50", x+w/2, y+h/2+14, 0.5, 0.5)

	case roleSidebar:
		drawBox(dc, x, y, w, h, hex(0xFFFFFF), hex(0xE0E0E0), 8)
		dc.SetRGBA255(80, 80, 80, 255)
		dc.DrawStringAnchored("CATEGORIES", x+16, y+22, 0, 0.5)

	case roleCategory:
		hover := (int(y/32))%2 == 0
		bg := hex(0xF5F7FA)
		if hover {
			bg = hex(0xE8F0FE)
		}
		border := hex(0xE8ECF0)
		if hover {
			border = hex(0xD2E3FC)
		}
		drawBox(dc, x, y, w, h, bg, border, 6)
		dc.SetRGBA255(40, 40, 50, 255)
		label, _ := n.GetContext().(string)
		dc.DrawStringAnchored(label, x+12, y+h/2, 0, 0.45)

	case roleCard:
		drawBox(dc, x, y, w, h, hex(0xFFFFFF), hex(0xE8E8E8), 10)

	case roleFooter:
		drawBox(dc, x, y, w, h, hex(0xF1F3F4), hex(0xE0E0E0), 0)
		dc.SetRGBA255(120, 120, 120, 255)
		dc.DrawStringAnchored(" 2026 ShopStore  |  Privacy  |  Terms  |  Contact", x+w/2, y+h/2, 0.5, 0.4)

	default:
		parent := n.GetParent()
		if parent != nil {
			gp, _ := parent.GetContext().(string)
			if gp == roleCard {
				renderCardChild(dc, n, x, y, w, h)
				return
			}
		}
	}

	for i := 0; i < n.GetChildCount(); i++ {
		renderNode(dc, n.GetChild(i), x, y)
	}
}

func renderCardChild(dc *gg.Context, n *goda.Node, x, y, w, h float64) {
	parent := n.GetParent()
	idx := -1
	for i := 0; i < parent.GetChildCount(); i++ {
		if parent.GetChild(i) == n {
			idx = i
			break
		}
	}
	ctx, _ := n.GetContext().(string)
	switch idx {
	case 0:
		clr := hex(0xD0D0D8)
		shape, badge := "generic", ""
		if len(ctx) > 0 {
			i1 := indexOf(ctx, '|')
			if i1 >= 0 {
				if c, ok := parseHexColor(ctx[:i1]); ok {
					clr = c
				}
				rest := ctx[i1+1:]
				i2 := indexOf(rest, '|')
				if i2 >= 0 {
					shape = rest[:i2]
					badge = rest[i2+1:]
				} else {
					shape = rest
				}
			}
		}
		drawRoundRect(dc, x, y, w, h, 8)
		dc.SetColor(clr)
		dc.Fill()
		dc.SetRGBA255(255, 255, 255, 120)
		dc.DrawStringAnchored(shape, x+w/2, y+h-14, 0.5, 0.5)
		drawShape(dc, x+w/2, y+h/2-6, min(w, h)*0.3, shape)
		if badge != "" {
			dc.SetRGBA255(234, 67, 53, 255)
			dc.DrawStringAnchored(badge, x+6, y+4, 0, 0.3)
		}
	case 1:
		dc.SetRGBA255(30, 30, 30, 255)
		dc.DrawStringWrapped(ctx, x+4, y+4, 0, 0, w-8, 1.2, gg.AlignLeft)
	case 2:
		dc.SetRGBA255(245, 124, 0, 255)
		dc.DrawStringAnchored(ctx+" / 5", x+4, y+h/2, 0, 0.45)
	case 3:
		current := ctx
		old := ""
		if i := indexOf(ctx, '|'); i >= 0 {
			current = ctx[:i]
			old = ctx[i+1:]
		}
		dc.SetRGBA255(15, 157, 88, 255)
		dc.DrawStringAnchored(current, x+4, y+h/2+1, 0, 0.45)
		if old != "" {
			tw, _ := dc.MeasureString(current)
			dc.SetRGBA255(148, 148, 148, 255)
			sx := x + 4 + tw + 6
			dc.DrawStringAnchored(old, sx, y+h/2+1, 0, 0.45)
			ow, oh := dc.MeasureString(old)
			dc.SetLineWidth(1)
			dc.DrawLine(sx-2, y+h/2-oh*0.2, sx+ow-2, y+h/2-oh*0.2)
			dc.Stroke()
		}
	case 4:
		drawBox(dc, x, y, w, h, hex(0x1A73E8), hex(0x1557B0), 5)
		dc.SetColor(color.White)
		dc.DrawStringAnchored("Add to Cart", x+w/2, y+h/2, 0.5, 0.4)
	}
}

func drawBox(dc *gg.Context, x, y, w, h float64, fill, stroke color.NRGBA, r float64) {
	drawRoundRect(dc, x, y, w, h, r)
	dc.SetColor(fill)
	dc.FillPreserve()
	dc.SetColor(stroke)
	dc.SetLineWidth(1)
	dc.Stroke()
}

func drawRoundRect(dc *gg.Context, x, y, w, h, r float64) {
	if r <= 0 {
		dc.DrawRectangle(x, y, w, h)
	} else {
		dc.DrawRoundedRectangle(x, y, w, h, r)
	}
}

func hex(c uint32) color.NRGBA {
	return color.NRGBA{uint8(c >> 16), uint8(c >> 8), uint8(c), 255}
}

func parseHexColor(s string) (color.NRGBA, bool) {
	if len(s) != 7 || s[0] != '#' {
		return color.NRGBA{}, false
	}
	var r, g, b uint8
	fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	return color.NRGBA{r, g, b, 255}, true
}

func drawShape(dc *gg.Context, cx, cy, size float64, shape string) {
	dc.Push()
	dc.SetRGBA255(255, 255, 255, 170)
	switch shape {
	case "headphones":
		dc.DrawArc(cx, cy-size*0.3, size*0.6, math.Pi, 0)
		dc.SetLineWidth(4)
		dc.Stroke()
		dc.DrawRoundedRectangle(cx-size*0.55, cy-size*0.1, size*0.45, size*0.5, 6)
		dc.Fill()
		dc.DrawRoundedRectangle(cx+size*0.1, cy-size*0.1, size*0.45, size*0.5, 6)
		dc.Fill()
	case "watch":
		dc.DrawRoundedRectangle(cx-size*0.35, cy-size*0.6, size*0.7, size*1.2, 12)
		dc.Fill()
		dc.SetRGBA255(60, 60, 60, 60)
		dc.DrawCircle(cx, cy, size*0.18)
		dc.Fill()
		dc.SetRGBA255(255, 255, 255, 90)
		dc.DrawCircle(cx, cy, size*0.25)
		dc.SetLineWidth(2)
		dc.Stroke()
	case "speaker":
		dc.DrawRoundedRectangle(cx-size*0.4, cy-size*0.5, size*0.8, size, 8)
		dc.Fill()
		dc.SetRGBA255(60, 60, 60, 50)
		for ly := cy - size*0.2; ly < cy+size*0.35; ly += size * 0.1 {
			dc.DrawLine(cx-size*0.25, ly, cx+size*0.25, ly)
			dc.SetLineWidth(2)
			dc.Stroke()
		}
		dc.SetRGBA255(255, 255, 255, 90)
		dc.DrawCircle(cx, cy+size*0.4, size*0.15)
		dc.Fill()
	case "hub":
		dc.DrawRoundedRectangle(cx-size*0.45, cy-size*0.15, size*0.9, size*0.3, 5)
		dc.Fill()
		for px := cx - size*0.3; px <= cx+size*0.3; px += size * 0.15 {
			dc.SetRGBA255(60, 60, 60, 60)
			dc.DrawRectangle(px, cy+size*0.08, size*0.1, size*0.15)
			dc.Fill()
			dc.SetRGBA255(255, 255, 255, 90)
		}
		dc.DrawRectangle(cx-size*0.05, cy-size*0.5, size*0.1, size*0.4)
		dc.Fill()
	case "keyboard":
		dc.DrawRoundedRectangle(cx-size*0.5, cy-size*0.4, size, size*0.8, 6)
		dc.Fill()
		dc.SetRGBA255(80, 80, 80, 40)
		keyW := size * 0.13
		for row := 0; row < 4; row++ {
			ry := cy - size*0.25 + float64(row)*size*0.17
			for col := 0; col < 6; col++ {
				rx := cx - size*0.4 + float64(col)*(keyW+2)
				dc.DrawRoundedRectangle(rx, ry, keyW, size*0.13, 2)
				dc.Fill()
			}
		}
		dc.SetRGBA255(255, 255, 255, 90)
	case "camera":
		dc.DrawRoundedRectangle(cx-size*0.45, cy-size*0.3, size*0.9, size*0.5, 8)
		dc.Fill()
		dc.SetRGBA255(60, 60, 60, 60)
		dc.DrawCircle(cx, cy-size*0.05, size*0.2)
		dc.Fill()
		dc.SetRGBA255(180, 180, 200, 120)
		dc.DrawCircle(cx, cy-size*0.05, size*0.12)
		dc.Fill()
		dc.SetRGBA255(255, 255, 255, 90)
		dc.DrawRectangle(cx+size*0.1, cy-size*0.5, size*0.15, size*0.2)
		dc.Fill()
	default:
		dc.DrawRoundedRectangle(cx-size*0.45, cy-size*0.45, size*0.9, size*0.9, 6)
		dc.Fill()
		dc.SetRGBA255(60, 60, 60, 40)
		dc.DrawRoundedRectangle(cx-size*0.2, cy-size*0.2, size*0.4, size*0.4, 3)
		dc.Fill()
		dc.SetRGBA255(255, 255, 255, 90)
	}
	dc.Pop()
}

func indexOf(s string, ch byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ch {
			return i
		}
	}
	return -1
}

func init() { _ = math.NaN }
