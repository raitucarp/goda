package main

import (
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

func renderWidget(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	white := color.NRGBA{255, 255, 255, 255}
	switch widget.Kind {
	case wHeader:
		drawRect(screen, x, y, w, h, hex("#FFFFFF"), hex("#E2E8F0"))
	case wLogo:
		drawRect(screen, x, y, w, h, hex("#4F46E5"), hex("#4F46E5"))
		centeredText(screen, "Dashboard", fontFace, x+w/2, y+h/2, white)
	case wSearch:
		bg := hex("#F1F5F9")
		if hovered {
			bg = hex("#E2E8F0")
		}
		drawRect(screen, x, y, w, h, bg, hex("#CBD5E1"))
		txt(screen, "Search...", fontFace, x+12, y+h/2, hexText("#94A3B8"))
	case wBell:
		bg := hex("#F1F5F9")
		if hovered {
			bg = hex("#E2E8F0")
		}
		drawRoundedRect(screen, x, y, w, h, 8, bg, hex("#CBD5E1"))
		if hovered {
			vector.DrawFilledCircle(screen, x+w-10, y+10, 4, hex("#EF4444"), true)
		}
	case wUser:
		r := h * 0.36
		clr := hex("#4F46E5")
		if hovered {
			clr = hex("#6366F1")
		}
		vector.DrawFilledCircle(screen, x+20, y+h/2, r, clr, true)
		centeredText(screen, "Admin", fontFace, x+40+r*2, y+h/2, hexText("#334155"))
	case wSidebar:
		drawRect(screen, x, y, w, h, hex("#1E293B"), hex("#1E293B"))
	case wMenuItem:
		renderMenuItem(screen, x, y, w, h, widget, hovered)
	case wKPICard:
		renderKPICard(screen, x, y, w, h, widget, hovered)
	case wChartCard:
		renderChartCard(screen, x, y, w, h, widget, hovered)
	case wBigChart:
		renderChartCard(screen, x, y, w, h, widget, hovered)
	case wOrderTable:
		renderOrderTable(screen, x, y, w, h, widget, hovered)
	case wCustTable:
		renderCustomerTable(screen, x, y, w, h, widget, hovered)
	case wProdCard:
		renderProductCard(screen, x, y, w, h, widget, hovered)
	case wReportCard:
		renderReportCard(screen, x, y, w, h, widget, hovered)
	case wSettingRow:
		renderSettingRow(screen, x, y, w, h, widget, hovered)
	case wStatusBar:
		drawRect(screen, x, y, w, h, hex("#E2E8F0"), hex("#CBD5E1"))
		txt(screen, "Dashboard v1.0 | Click sidebar to navigate", fontFaceSm, x+12, y+h/2, hexText("#64748B"))
	}
}

func renderMenuItem(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	md, _ := widget.Data.(MenuData)
	bg := hex("#1E293B")
	textClr := hexText("#94A3B8")
	if md.Active {
		bg = hex("#334155")
		textClr = hexText("#F8FAFC")
		drawRect(screen, x, y, 4, h, hex("#4F46E5"), hex("#4F46E5"))
	} else if hovered {
		bg = hex("#283548")
		textClr = hexText("#E2E8F0")
	}
	drawRect(screen, x, y, w, h, bg, bg)
	txt(screen, md.Label, fontFace, x+18, y+h/2, textClr)
}

func renderKPICard(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	kd, _ := widget.Data.(KPIData)
	accColor := hex(kd.Color)
	border := hex("#E5E7EB")
	if hovered {
		border = accColor
	}
	drawRoundedRect(screen, x, y, w, h, 8, hex("#FFFFFF"), border)

	drawRect(screen, x+12, y+12, 24, 4, accColor, accColor)
	txt(screen, kd.Label, fontFaceSm, x+12, y+28, hexText("#6B7280"))
	txt(screen, kd.Value, fontFace, x+12, y+52, hexText("#1F2937"))

	changeClr := hexText("#10B981")
	if !kd.Positive {
		changeClr = hexText("#EF4444")
	}
	txt(screen, kd.Change, fontFaceSm, x+12, y+74, changeClr)

	vector.DrawFilledCircle(screen, x+w-18, y+18, 5, accColor, true)
}

func renderChartCard(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	cs, _ := widget.Data.(ChartSeries)
	border := hex("#E5E7EB")
	if hovered {
		border = hex("#94A3B8")
	}
	drawRoundedRect(screen, x, y, w, h, 8, hex("#FFFFFF"), border)

	titleH := float32(22)
	txt(screen, cs.Title, fontFace, x+12, y+12, hexText("#1F2937"))

	chartX := x + 40
	chartY := y + titleH + 18
	chartW := w - 80
	chartH := h - titleH - 36

	if chartW < 40 || chartH < 40 {
		return
	}

	accColor := hex(cs.Color)

	switch cs.Type {
	case "line":
		renderLineChart(screen, chartX, chartY, chartW, chartH, cs, accColor)
	case "bar":
		renderBarChart(screen, chartX, chartY, chartW, chartH, cs, accColor)
	}

	for i, lbl := range cs.Labels {
		denom := float32(len(cs.Labels) - 1)
		if denom <= 0 {
			denom = 1
		}
		lx := chartX + float32(i)*(chartW/denom)
		txt(screen, lbl, fontFaceSm, lx, chartY+chartH+14, hexText("#6B7280"))
	}
}

func renderOrderTable(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	orders, _ := widget.Data.([]OrderRow)
	border := hex("#E5E7EB")
	if hovered {
		border = hex("#94A3B8")
	}
	drawRoundedRect(screen, x, y, w, h, 8, hex("#FFFFFF"), border)

	rowH := float32(28)
	headerY := y + 6

	mx := float32(0)
	my := float32(0)
	drawHeader(screen, "ID", x+10, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Customer", x+90, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Date", x+250, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Amount", x+370, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Status", x+450, headerY, rowH, hexText("#475569"))
	_ = mx
	_ = my
	vector.StrokeLine(screen, x+10, headerY+rowH+2, x+w-10, headerY+rowH+2, 1, hex("#E5E7EB"), true)

	for i, o := range orders {
		ry := headerY + rowH + 8 + float32(i)*rowH
		rowBg := hex("#FFFFFF")
		if i%2 == 0 {
			rowBg = hex("#F8FAFC")
		}
		drawRect(screen, x+4, ry-2, w-8, rowH, rowBg, rowBg)
		txt(screen, o.ID, fontFaceSm, x+10, ry+rowH/2, hexText("#1F2937"))
		txt(screen, o.Customer, fontFaceSm, x+90, ry+rowH/2, hexText("#1F2937"))
		txt(screen, o.Date, fontFaceSm, x+250, ry+rowH/2, hexText("#6B7280"))
		txt(screen, o.Amount, fontFaceSm, x+370, ry+rowH/2, hexText("#1F2937"))
		txt(screen, o.Status, fontFaceSm, x+450, ry+rowH/2, statusColor(o.Status))
	}
}

func renderCustomerTable(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	customers, _ := widget.Data.([]CustomerRow)
	border := hex("#E5E7EB")
	if hovered {
		border = hex("#94A3B8")
	}
	drawRoundedRect(screen, x, y, w, h, 8, hex("#FFFFFF"), border)

	rowH := float32(28)
	headerY := y + 6
	drawHeader(screen, "Name", x+10, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Email", x+150, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Plan", x+330, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Joined", x+430, headerY, rowH, hexText("#475569"))
	drawHeader(screen, "Spent", x+530, headerY, rowH, hexText("#475569"))
	vector.StrokeLine(screen, x+10, headerY+rowH+2, x+w-10, headerY+rowH+2, 1, hex("#E5E7EB"), true)

	for i, c := range customers {
		ry := headerY + rowH + 8 + float32(i)*rowH
		if i%2 == 0 {
			drawRect(screen, x+4, ry-2, w-8, rowH, hex("#F8FAFC"), hex("#F8FAFC"))
		}
		txt(screen, c.Name, fontFaceSm, x+10, ry+rowH/2, hexText("#1F2937"))
		txt(screen, c.Email, fontFaceSm, x+150, ry+rowH/2, hexText("#6B7280"))
		txt(screen, c.Plan, fontFaceSm, x+330, ry+rowH/2, hexText("#4F46E5"))
		txt(screen, c.Joined, fontFaceSm, x+430, ry+rowH/2, hexText("#6B7280"))
		txt(screen, c.Spent, fontFaceSm, x+530, ry+rowH/2, hexText("#1F2937"))
	}
}

func renderProductCard(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	p, _ := widget.Data.(ProductRow)
	border := hex("#E5E7EB")
	if hovered {
		border = hex("#4F46E5")
	}
	drawRoundedRect(screen, x, y, w, h, 8, hex("#FFFFFF"), border)

	bannerH := h * 0.3
	clr := hex(p.Color)
	if !p.Active {
		clr = hex("#D1D5DB")
	}
	drawRect(screen, x, y, w, bannerH, clr, clr)
	centeredText(screen, p.Category, fontFaceSm, x+w/2, y+bannerH/2, hex("#FFFFFF"))

	txt(screen, p.Name, fontFace, x+10, y+bannerH+16, hexText("#1F2937"))

	stockClr := hexText("#10B981")
	stockLabel := "In Stock (" + strconv.Itoa(p.Stock) + ")"
	if p.Stock < 10 {
		stockClr = hexText("#EF4444")
		stockLabel = "Low: " + strconv.Itoa(p.Stock)
	}
	if p.Stock == 0 {
		stockClr = hexText("#9CA3AF")
		stockLabel = "Out of Stock"
	}
	txt(screen, stockLabel, fontFaceSm, x+10, y+bannerH+36, stockClr)
	txt(screen, p.Price, fontFace, x+10, y+h-12, hexText("#1F2937"))
}

func renderReportCard(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	r, _ := widget.Data.(ReportCard)
	border := hex("#E5E7EB")
	bg := hex("#FFFFFF")
	if hovered {
		border = hex(r.Color)
		bg = hex("#F8FAFC")
	}
	drawRoundedRect(screen, x, y, w, h, 8, bg, border)

	leftStripe := float32(4)
	drawRect(screen, x+1, y+1, leftStripe, h-2, hex(r.Color), hex(r.Color))

	txt(screen, r.Title, fontFace, x+16, y+20, hexText("#1F2937"))
	txt(screen, r.Description, fontFaceSm, x+16, y+46, hexText("#6B7280"))
	txt(screen, r.Date, fontFaceSm, x+16, y+68, hexText("#94A3B8"))

	btnX := x + w - 100
	btnY := y + h - 34
	btnW := float32(84)
	btnH := float32(24)
	btnClr := hex(r.Color)
	if hovered {
		btnClr = darken(btnClr, 0.15)
	}
	drawRoundedRect(screen, btnX, btnY, btnW, btnH, 4, btnClr, btnClr)
	centeredText(screen, "Download", fontFaceSm, btnX+btnW/2, btnY+btnH/2, hex("#FFFFFF"))
}

func renderSettingRow(screen *ebiten.Image, x, y, w, h float32, widget Widget, hovered bool) {
	sr, _ := widget.Data.(SettingRow)
	bg := hex("#FFFFFF")
	if hovered {
		bg = hex("#F8FAFC")
	}
	drawRoundedRect(screen, x, y, w, h, 6, bg, hex("#E5E7EB"))

	txt(screen, sr.Label, fontFace, x+14, y+h/2, hexText("#1F2937"))

	toggleW := float32(40)
	toggleH := float32(22)
	tx := x + w - toggleW - 14
	ty := y + h/2 - toggleH/2

	toggleBg := hex("#D1D5DB")
	toggleKnobX := tx + 2
	if sr.Enabled {
		toggleBg = hex("#4F46E5")
		toggleKnobX = tx + toggleW - toggleH + 2
	}
	drawRoundedRect(screen, tx, ty, toggleW, toggleH, toggleH/2, toggleBg, toggleBg)
	vector.DrawFilledCircle(screen, toggleKnobX+toggleH/2-2, ty+toggleH/2, toggleH/2-3, hex("#FFFFFF"), true)
}

func renderLineChart(screen *ebiten.Image, x, y, w, h float32, cs ChartSeries, clr color.NRGBA) {
	n := len(cs.Data)
	if n == 0 || w < 20 || h < 20 {
		return
	}
	maxV, minV := cs.Data[0], cs.Data[0]
	for _, v := range cs.Data {
		if v > maxV {
			maxV = v
		}
		if v < minV {
			minV = v
		}
	}
	if maxV == minV {
		maxV = minV + 1
	}
	rangeV := maxV - minV

	gridClr := hex("#F1F5F9")
	for i := 0; i <= 4; i++ {
		gy := y + h - (float32(i) * h / 4)
		vector.StrokeLine(screen, x, gy, x+w, gy, 1, gridClr, true)
		val := minV + rangeV*float64(i)/4
		txt(screen, fmtFloat(val), fontFaceSm, x-36, gy, hexText("#94A3B8"))
	}

	for i := 0; i < n-1; i++ {
		dn := float64(n - 1)
		sx := x + float32(float64(i)*(float64(w)/dn))
		sy := y + h - float32((cs.Data[i]-minV)/rangeV*float64(h))
		ex := x + float32(float64(i+1)*(float64(w)/dn))
		ey := y + h - float32((cs.Data[i+1]-minV)/rangeV*float64(h))
		vector.StrokeLine(screen, sx, sy, ex, ey, 2.5, clr, true)
	}

	for i := 0; i < n; i++ {
		dn := float64(n - 1)
		if dn <= 0 {
			dn = 1
		}
		px := x + float32(float64(i)*(float64(w)/dn))
		py := y + h - float32((cs.Data[i]-minV)/rangeV*float64(h))
		vector.DrawFilledCircle(screen, px, py, 4, clr, true)
		vector.DrawFilledCircle(screen, px, py, 2.5, hex("#FFFFFF"), true)
	}
}

func renderBarChart(screen *ebiten.Image, x, y, w, h float32, cs ChartSeries, clr color.NRGBA) {
	n := len(cs.Data)
	if n == 0 || w < 20 || h < 20 {
		return
	}
	var maxV float64
	for _, v := range cs.Data {
		if v > maxV {
			maxV = v
		}
	}
	if maxV == 0 {
		maxV = 1
	}

	gridClr := hex("#F1F5F9")
	for i := 0; i <= 4; i++ {
		gy := y + h - (float32(i) * h / 4)
		vector.StrokeLine(screen, x, gy, x+w, gy, 1, gridClr, true)
		val := maxV * float64(i) / 4
		txt(screen, fmtFloat(val), fontFaceSm, x-36, gy, hexText("#94A3B8"))
	}

	fn := float32(n)
	barSpacing := w * 0.12
	totalSpacing := barSpacing * (fn - 1)
	barW := (w - totalSpacing) / fn
	if barW < 4 {
		barW = 4
	}
	for i, val := range cs.Data {
		barH := float32(val/maxV) * h
		bx := x + float32(i)*(barW+barSpacing)
		by := y + h - barH

		drawRect(screen, bx, by, barW, barH, clr, clr)
		valStr := fmtFloat(val)
		tw := textWidth(valStr, fontFaceSm)
		txt(screen, valStr, fontFaceSm, bx+barW/2-tw/2, by-6, hexText("#475569"))
	}
}

func drawHeader(screen *ebiten.Image, label string, x, y, h float32, clr color.NRGBA) {
	txt(screen, label, fontFace, x, y+h/2, clr)
}

func statusColor(status string) color.NRGBA {
	switch status {
	case "Delivered":
		return hexText("#10B981")
	case "Processing":
		return hexText("#F59E0B")
	case "Shipped":
		return hexText("#3B82F6")
	case "Cancelled":
		return hexText("#EF4444")
	default:
		return hexText("#6B7280")
	}
}

func drawRect(screen *ebiten.Image, x, y, w, h float32, fill, stroke color.NRGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	vector.DrawFilledRect(screen, x, y, w, h, fill, true)
	if stroke != fill {
		vector.StrokeRect(screen, x, y, w, h, 1, stroke, true)
	}
}

func drawRoundedRect(screen *ebiten.Image, x, y, w, h, r float32, fill, stroke color.NRGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	if r <= 0 || r*2 > w || r*2 > h {
		drawRect(screen, x, y, w, h, fill, stroke)
		return
	}
	vector.DrawFilledRect(screen, x+r, y, w-2*r, h, fill, true)
	vector.DrawFilledRect(screen, x, y+r, w, h-2*r, fill, true)
	vector.DrawFilledCircle(screen, x+r, y+r, r, fill, true)
	vector.DrawFilledCircle(screen, x+w-r, y+r, r, fill, true)
	vector.DrawFilledCircle(screen, x+r, y+h-r, r, fill, true)
	vector.DrawFilledCircle(screen, x+w-r, y+h-r, r, fill, true)

	if stroke != fill {
		vector.StrokeLine(screen, x+r, y, x+w-r, y, 1, stroke, true)
		vector.StrokeLine(screen, x+r, y+h, x+w-r, y+h, 1, stroke, true)
		vector.StrokeLine(screen, x, y+r, x, y+h-r, 1, stroke, true)
		vector.StrokeLine(screen, x+w, y+r, x+w, y+h-r, 1, stroke, true)
	}
}

func txt(screen *ebiten.Image, str string, face font.Face, x, y float32, clr color.NRGBA) {
	m := face.Metrics()
	ascent := float32(m.Ascent) / 64
	descent := float32(m.Descent) / 64
	baseline := y + (ascent-descent)/2
	text.Draw(screen, str, face, int(x), int(baseline), clr)
}

func centeredText(screen *ebiten.Image, str string, face font.Face, cx, cy float32, clr color.NRGBA) {
	tw := textWidth(str, face)
	m := face.Metrics()
	ascent := float32(m.Ascent) / 64
	descent := float32(m.Descent) / 64
	x := cx - tw/2
	baseline := cy + (ascent-descent)/2
	text.Draw(screen, str, face, int(x), int(baseline), clr)
}

func textWidth(str string, face font.Face) float32 {
	var totalW float32
	for _, r := range str {
		a, ok := face.GlyphAdvance(r)
		if ok {
			totalW += float32(a) / 64
		}
	}
	return totalW
}

func hex(s string) color.NRGBA {
	if len(s) == 0 || s[0] != '#' {
		return color.NRGBA{128, 128, 128, 255}
	}
	s = strings.TrimPrefix(s, "#")
	var r, g, b uint8
	switch len(s) {
	case 6:
		v, _ := strconv.ParseUint(s, 16, 32)
		r = uint8(v >> 16)
		g = uint8(v >> 8)
		b = uint8(v)
	case 3:
		v, _ := strconv.ParseUint(s, 16, 16)
		r = uint8(v>>8) * 17
		g = uint8((v>>4)&0xF) * 17
		b = uint8(v&0xF) * 17
	default:
		return color.NRGBA{128, 128, 128, 255}
	}
	return color.NRGBA{r, g, b, 255}
}

func hexText(s string) color.NRGBA { return hex(s) }

func darken(c color.NRGBA, amount float64) color.NRGBA {
	return color.NRGBA{
		uint8(math.Max(0, float64(c.R)*(1-amount))),
		uint8(math.Max(0, float64(c.G)*(1-amount))),
		uint8(math.Max(0, float64(c.B)*(1-amount))),
		255,
	}
}

func fmtFloat(v float64) string {
	if v == math.Trunc(v) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', 1, 64)
}
