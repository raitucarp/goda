package main

import goda "goda"

func buildByCSSString() *goda.Node {
	root := newNodeWithID("root", "").
		ApplyStyleString(`
			width: 800; height: 1080;
			flex-direction: column;
			padding: 16; gap: 12;
		`)

	header := newNodeWithID("header", roleHeader).
		ApplyStyleString(`
			flex-direction: row;
			align-items: center;
			gap: 16; height: 64;
			flex-shrink: 0;
		`)
	header.SetPadding(goda.EdgeLeft, 20).SetPadding(goda.EdgeRight, 20).
		SetPadding(goda.EdgeTop, 12).SetPadding(goda.EdgeBottom, 12)

	logo := newNodeWithID("logo", roleLogo).
		ApplyStyleString("width: 120; height: 40; flex-shrink: 0;")

	search := newNodeWithID("search", roleSearch).
		ApplyStyleString("flex-grow: 1; height: 40;")

	cart := newNodeWithID("cart", roleCart).
		ApplyStyleString("width: 80; height: 36; flex-shrink: 0;")

	header.InsertChildNode(logo, 0)
	header.InsertChildNode(search, 1)
	header.InsertChildNode(cart, 2)

	banner := newNodeWithID("banner", roleBanner).
		ApplyStyleString("height: 150; flex-shrink: 0; padding: 24;")

	content := goda.New("content").
		ApplyStyleString("flex-direction: row; gap: 16; flex-grow: 1;")

	sidebar := newNodeWithID("sidebar", roleSidebar).
		ApplyStyleString("width: 170; flex-shrink: 0; flex-direction: column; gap: 6;")
	sidebar.SetPadding(goda.EdgeLeft, 14).SetPadding(goda.EdgeRight, 14).
		SetPadding(goda.EdgeTop, 14).SetPadding(goda.EdgeBottom, 14)

	for _, cat := range categories {
		item := goda.New().
			ApplyStyleString("flex-grow: 1; padding-top: 14; padding-bottom: 14; padding-left: 12; padding-right: 12;")
		item.SetContext(cat.name)
		sidebar.InsertChildNode(item, sidebar.GetChildCount())
	}

	grid := goda.New("grid").
		ApplyStyleString("flex-shrink: 1; flex-direction: column; gap: 12;")

	for i := 0; i < 6; i += 3 {
		row := goda.New().
			ApplyStyleString("flex-direction: row; gap: 12; flex-shrink: 0;")
		grid.InsertChildNode(row, grid.GetChildCount())

		for j := 0; j < 3 && i+j < len(products); j++ {
			p := products[i+j]
			card := newNodeWithID("", roleCard).
				ApplyStyleString("width: 180; flex-direction: column; gap: 6; padding: 8;")

			img := goda.New().
				ApplyStyleString("height: 120; width: 100%; flex-shrink: 0;")
			img.SetContext(p.imgColor + "|" + p.imgShape + "|" + p.badge)

			title := goda.New().
				ApplyStyleString("height: 34; flex-shrink: 0;")
			title.SetContext(p.name)

			rating := goda.New().
				ApplyStyleString("height: 16; flex-shrink: 0;")
			rating.SetContext(p.rating)

			price := goda.New().
				ApplyStyleString("height: 22; flex-shrink: 0;")
			price.SetContext(p.price + "|" + p.oldPrice)

			btn := goda.New().
				ApplyStyleString("height: 34; width: 100%; flex-shrink: 0; margin-top: 4;")

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
		ApplyStyleString("height: 40; flex-shrink: 0;")

	root.InsertChildNode(header, 0)
	root.InsertChildNode(banner, 1)
	root.InsertChildNode(content, 2)
	root.InsertChildNode(footer, 3)

	return root
}
