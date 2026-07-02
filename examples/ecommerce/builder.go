package main

import goda "goda"

func buildByBuilder() *goda.Node {
	root := goda.New("root").
		SetWidth(800).SetHeight(1080).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetPadding(goda.EdgeAll, 16).
		SetGap(goda.GutterAll, 12)

	header := newNodeWithID("header", roleHeader).
		SetFlexDirection(goda.FlexDirectionRow).
		SetAlignItems(goda.AlignCenter).
		SetGap(goda.GutterAll, 16).
		SetHeight(64).
		SetFlexShrink(0).
		SetPadding(goda.EdgeLeft, 20).SetPadding(goda.EdgeRight, 20).
		SetPadding(goda.EdgeTop, 12).SetPadding(goda.EdgeBottom, 12)

	logo := newNodeWithID("logo", roleLogo).
		SetWidth(120).SetHeight(40).
		SetFlexShrink(0)

	search := newNodeWithID("search", roleSearch).
		SetFlexGrow(1).SetHeight(40)

	cart := newNodeWithID("cart", roleCart).
		SetWidth(80).SetHeight(36).
		SetFlexShrink(0)

	header.InsertChildNode(logo, 0)
	header.InsertChildNode(search, 1)
	header.InsertChildNode(cart, 2)

	banner := newNodeWithID("banner", roleBanner).
		SetHeight(150).SetFlexShrink(0).
		SetPadding(goda.EdgeAll, 24)

	content := goda.New("content").
		SetFlexDirection(goda.FlexDirectionRow).
		SetGap(goda.GutterAll, 16).
		SetFlexGrow(1)

	sidebar := newNodeWithID("sidebar", roleSidebar).
		SetWidth(170).SetFlexShrink(0).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetGap(goda.GutterAll, 6).
		SetPadding(goda.EdgeLeft, 14).SetPadding(goda.EdgeRight, 14).
		SetPadding(goda.EdgeTop, 14).SetPadding(goda.EdgeBottom, 14)

	for _, cat := range categories {
		item := goda.New().
			SetFlexGrow(1).
			SetPadding(goda.EdgeTop, 14).SetPadding(goda.EdgeBottom, 14).
			SetPadding(goda.EdgeLeft, 12).SetPadding(goda.EdgeRight, 12)
		item.SetContext(cat.name)
		sidebar.InsertChildNode(item, sidebar.GetChildCount())
	}

	grid := goda.New("grid").
		SetFlexShrink(1).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetGap(goda.GutterAll, 12)

	for i := 0; i < 6; i += 3 {
		row := goda.New().
			SetFlexDirection(goda.FlexDirectionRow).
			SetGap(goda.GutterAll, 12).
			SetFlexShrink(0)
		grid.InsertChildNode(row, grid.GetChildCount())

		for j := 0; j < 3 && i+j < len(products); j++ {
			p := products[i+j]
			card := newNodeWithID("", roleCard).
				SetWidth(180).
				SetFlexDirection(goda.FlexDirectionColumn).
				SetGap(goda.GutterAll, 6).
				SetPadding(goda.EdgeAll, 8)

			img := goda.New().
				SetHeight(120).SetWidthPercent(100).
				SetFlexShrink(0)
			img.SetContext(p.imgColor + "|" + p.imgShape + "|" + p.badge)

			title := goda.New().
				SetHeight(34).SetFlexShrink(0)
			title.SetContext(p.name)

			rating := goda.New().
				SetHeight(16).SetFlexShrink(0)
			rating.SetContext(p.rating)

			price := goda.New().
				SetHeight(22).SetFlexShrink(0)
			price.SetContext(p.price + "|" + p.oldPrice)

			btn := goda.New().
				SetHeight(34).SetWidthPercent(100).
				SetFlexShrink(0).SetMargin(goda.EdgeTop, 4)

			card.InsertChildNode(img, 0)
			card.InsertChildNode(title, 1)
			card.InsertChildNode(rating, 2)
			card.InsertChildNode(price, 3)
			card.InsertChildNode(btn, 4)

			row.InsertChildNode(card, row.GetChildCount())
		}
	}

	content.InsertChildNode(sidebar, 0)
	content.InsertChildNode(grid, 1)

	footer := newNodeWithID("footer", roleFooter).
		SetHeight(40).SetFlexShrink(0)

	root.InsertChildNode(header, 0)
	root.InsertChildNode(banner, 1)
	root.InsertChildNode(content, 2)
	root.InsertChildNode(footer, 3)

	return root
}
