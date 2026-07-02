package main

import goda "goda"

func buildLayout(fullWidth, fullHeight float32, activeMenu int) *goda.Node {
	root := goda.New("root").
		SetWidth(fullWidth).SetHeight(fullHeight).
		SetFlexDirection(goda.FlexDirectionColumn)

	header := buildHeader()
	body := buildBody(activeMenu)

	root.InsertChildNode(header, 0)
	root.InsertChildNode(body, 1)

	statusBar := newWidget("statusbar", wStatusBar).
		SetFlexDirection(goda.FlexDirectionRow).
		SetAlignItems(goda.AlignCenter).
		SetHeight(24).SetFlexShrink(0).
		SetPadding(goda.EdgeLeft, 16).SetPadding(goda.EdgeRight, 16)
	root.InsertChildNode(statusBar, 2)

	return root
}

func buildHeader() *goda.Node {
	h := newWidget("header", wHeader).
		SetFlexDirection(goda.FlexDirectionRow).
		SetAlignItems(goda.AlignCenter).
		SetHeight(56).SetFlexShrink(0).
		SetPadding(goda.EdgeLeft, 20).SetPadding(goda.EdgeRight, 20).
		SetGap(goda.GutterAll, 16)

	logo := newWidget("logo", wLogo).
		SetWidth(140).SetHeight(32).SetFlexShrink(0)
	spacer := goda.New("header-spacer").SetFlexGrow(1)
	search := newWidget("search", wSearch).
		SetWidth(200).SetHeight(34).SetFlexShrink(0)
	bell := newWidget("bell", wBell).
		SetWidth(36).SetHeight(36).SetFlexShrink(0)
	user := newWidget("user", wUser).
		SetWidth(130).SetHeight(36).SetFlexShrink(0)

	h.InsertChildNode(logo, 0)
	h.InsertChildNode(spacer, 1)
	h.InsertChildNode(search, 2)
	h.InsertChildNode(bell, 3)
	h.InsertChildNode(user, 4)
	return h
}

func buildBody(activeMenu int) *goda.Node {
	body := goda.New("body").
		SetFlexDirection(goda.FlexDirectionRow).
		SetFlexGrow(1)

	sidebar := buildSidebar(activeMenu)
	mainArea := newWidget("main", wMain).
		SetFlexGrow(1).SetFlexShrink(1).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetPadding(goda.EdgeAll, 20).
		SetGap(goda.GutterAll, 16)

	buildView(mainArea, activeMenu)

	body.InsertChildNode(sidebar, 0)
	body.InsertChildNode(mainArea, 1)
	return body
}

func buildSidebar(activeMenu int) *goda.Node {
	sidebar := newWidget("sidebar", wSidebar).
		SetWidth(220).SetFlexShrink(0).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetGap(goda.GutterAll, 2).
		SetPadding(goda.EdgeTop, 12).SetPadding(goda.EdgeBottom, 12).
		SetPadding(goda.EdgeLeft, 8).SetPadding(goda.EdgeRight, 8)

	for i, m := range menuItems {
		item := newWidget("", wMenuItem).
			SetHeight(40).SetFlexShrink(0).
			SetFlexDirection(goda.FlexDirectionRow).
			SetAlignItems(goda.AlignCenter).
			SetPadding(goda.EdgeLeft, 14).SetPadding(goda.EdgeRight, 14)
		item.SetContext(Widget{Kind: wMenuItem, Data: MenuData{
			Label:  m.Label,
			Active: i == activeMenu,
			Index:  i,
		}})
		sidebar.InsertChildNode(item, i)
	}
	return sidebar
}

func buildView(main *goda.Node, activeMenu int) {
	switch activeMenu {
	case 0:
		buildOverviewView(main)
	case 1:
		buildAnalyticsView(main)
	case 2:
		buildCustomersView(main)
	case 3:
		buildOrdersView(main)
	case 4:
		buildProductsView(main)
	case 5:
		buildReportsView(main)
	case 6:
		buildSettingsView(main)
	}
}
