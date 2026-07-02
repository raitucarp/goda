package main

import (
	"fmt"

	goda "goda"
)

func buildByRenderFrom() *goda.Node {
	source := `
		#root {
			width: 800;
			height: 1080;
			flex-direction: column;
			padding: 16;
			gap: 12;

			#header {
				flex-direction: row;
				align-items: center;
				gap: 16;
				height: 64;
				flex-shrink: 0;

				#logo {
					width: 120; height: 40; flex-shrink: 0;
				}
				#search {
					flex-grow: 1; height: 40;
				}
				#cart {
					width: 80; height: 36; flex-shrink: 0;
				}
			}

			#banner {
				height: 150;
				flex-shrink: 0;
				padding: 24;
			}

			#content {
				flex-direction: row;
				gap: 16;
				flex-grow: 1;

				#sidebar {
					width: 170;
					flex-shrink: 0;
					flex-direction: column;
					gap: 6;
				}

				#grid {
					flex-shrink: 1;
					flex-direction: column;
					gap: 12;
				}
			}

			#footer {
				height: 40;
				flex-shrink: 0;
			}
		}
	`

	roots, err := goda.RenderFrom(source)
	if err != nil {
		panic(fmt.Sprintf("RenderFrom error: %v", err))
	}
	root := roots[0]

	// Set contexts by walking tree
	contextByID := map[string]string{
		"header":  roleHeader,
		"logo":    roleLogo,
		"search":  roleSearch,
		"cart":    roleCart,
		"banner":  roleBanner,
		"sidebar": roleSidebar,
		"footer":  roleFooter,
	}
	walkAndSetContext(root, contextByID)

	// ── Header padding ──
	header := findByID(root, "header")
	if header != nil {
		header.SetPadding(goda.EdgeLeft, 20).SetPadding(goda.EdgeRight, 20).
			SetPadding(goda.EdgeTop, 12).SetPadding(goda.EdgeBottom, 12)
	}

	// ── Sidebar padding + categories ──
	sidebar := findByID(root, "sidebar")
	if sidebar != nil {
		sidebar.SetPadding(goda.EdgeLeft, 14).SetPadding(goda.EdgeRight, 14).
			SetPadding(goda.EdgeTop, 14).SetPadding(goda.EdgeBottom, 14)

		for _, cat := range categories {
			item := goda.New().
				ApplyStyleString("flex-grow: 1; padding-top: 14; padding-bottom: 14; padding-left: 12; padding-right: 12;")
			item.SetContext(cat.name)
			sidebar.InsertChildNode(item, sidebar.GetChildCount())
		}
	}

	// ── Product grid ──
	grid := findByID(root, "grid")
	if grid != nil {
		for i := 0; i < 6; i += 3 {
			row := goda.New().
				ApplyStyleString("flex-direction: row; gap: 12; flex-shrink: 0;")
			grid.InsertChildNode(row, grid.GetChildCount())

			for j := 0; j < 3 && i+j < len(products); j++ {
				p := products[i+j]
				card := goda.New().
					ApplyStyleString("width: 180; flex-direction: column; gap: 6; padding: 8;")
				card.SetContext(roleCard)

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
	}

	return root
}

func findByID(n *goda.Node, id string) *goda.Node {
	if n.GetID() == id {
		return n
	}
	for i := 0; i < n.GetChildCount(); i++ {
		if found := findByID(n.GetChild(i), id); found != nil {
			return found
		}
	}
	return nil
}

func walkAndSetContext(n *goda.Node, ctxMap map[string]string) {
	if role, ok := ctxMap[n.GetID()]; ok && role != "" {
		n.SetContext(role)
	}
	for i := 0; i < n.GetChildCount(); i++ {
		walkAndSetContext(n.GetChild(i), ctxMap)
	}
}
