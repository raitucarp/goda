package goda

// Rect holds the computed position and size of a laid-out node.
// These values are only meaningful after CalculateNodeLayout is called.
type Rect struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
	Width  float32
	Height float32
}

// Edges holds the computed margin, border, or padding values for all four sides.
type Edges struct {
	Top    float32
	Right  float32
	Bottom float32
	Left   float32
}

// LayoutOut is the public layout output for a node after CalculateNodeLayout.
// It bundles position, size, box-model edges, and layout metadata into one
// struct for easy consumption by GUI libraries.
//
// Example:
//
//	goda.CalculateNodeLayout(root, 800, 600, goda.DirectionLTR)
//	lo := root.LayoutOut()
//	renderer.DrawBox(lo.Left, lo.Top, lo.Width, lo.Height, lo.Margin, lo.Padding)
type LayoutOut struct {
	Rect
	Margin     Edges
	Border     Edges
	Padding    Edges
	Direction  Direction
	HadOverflow bool
}

// LayoutOut returns the computed layout as a single convenience struct.
// Call this after CalculateNodeLayout to get position, dimensions, and
// box-model edges in one shot.
func (n *Node) LayoutOut() LayoutOut {
	return LayoutOut{
		Rect: Rect{
			Left:   n.GetLeft(),
			Top:    n.GetTop(),
			Right:  n.GetRight(),
			Bottom: n.GetBottom(),
			Width:  n.GetWidth(),
			Height: n.GetHeight(),
		},
		Margin: Edges{
			Top:    n.GetLayoutMargin(EdgeTop),
			Right:  n.GetLayoutMargin(EdgeRight),
			Bottom: n.GetLayoutMargin(EdgeBottom),
			Left:   n.GetLayoutMargin(EdgeLeft),
		},
		Border: Edges{
			Top:    n.GetLayoutBorder(EdgeTop),
			Right:  n.GetLayoutBorder(EdgeRight),
			Bottom: n.GetLayoutBorder(EdgeBottom),
			Left:   n.GetLayoutBorder(EdgeLeft),
		},
		Padding: Edges{
			Top:    n.GetLayoutPadding(EdgeTop),
			Right:  n.GetLayoutPadding(EdgeRight),
			Bottom: n.GetLayoutPadding(EdgeBottom),
			Left:   n.GetLayoutPadding(EdgeLeft),
		},
		Direction:   n.GetLayoutDirection(),
		HadOverflow: n.GetHadOverflow(),
	}
}
