package main

import goda "goda"

func buildOverviewView(main *goda.Node) {
	kpiRow := buildKPIRow(overviewKPIs)
	main.InsertChildNode(kpiRow, 0)

	chartsRow := newWidget("charts-row", wChartsRow).
		SetFlexDirection(goda.FlexDirectionRow).
		SetGap(goda.GutterAll, 14).
		SetFlexGrow(1).SetFlexShrink(1).SetMinHeight(240)

	line := buildChartCard("Revenue Trend", "line",
		[]float64{12, 19, 15, 27, 22, 35, 30, 42, 38, 55, 48, 62},
		[]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		"#4F46E5")

	bar := buildChartCard("Sales by Category", "bar",
		[]float64{45, 68, 52, 39},
		[]string{"Electronics", "Clothing", "Home", "Sports"},
		"#10B981")

	chartsRow.InsertChildNode(line, 0)
	chartsRow.InsertChildNode(bar, 1)
	main.InsertChildNode(chartsRow, 1)

	table := buildOrderTable("Recent Orders", overviewOrders)
	main.InsertChildNode(table, 2)
}

func buildAnalyticsView(main *goda.Node) {
	kpiRow := buildKPIRow(analyticsKPIs)
	main.InsertChildNode(kpiRow, 0)

	bigChart := newWidget("big-chart", wBigChart).
		SetFlexGrow(1).SetFlexShrink(1).SetMinHeight(280).
		SetPadding(goda.EdgeAll, 16)
	bigChart.SetContext(Widget{Kind: wBigChart, Data: ChartSeries{
		Type: "line", Title: "User Growth",
		Data:   []float64{8, 16, 24, 35, 48, 62, 78, 95, 115, 138, 162, 190},
		Labels: []string{"Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb", "Mar", "Apr", "May", "Jun"},
		Color:  "#4F46E5",
	}})
	main.InsertChildNode(bigChart, 1)

	sources := buildChartCard("Traffic Sources", "bar",
		[]float64{42, 28, 18, 12},
		[]string{"Organic", "Direct", "Social", "Referral"},
		"#10B981")
	main.InsertChildNode(sources, 2)
}

func buildCustomersView(main *goda.Node) {
	kpiRow := buildKPIRow(customersKPIs)
	main.InsertChildNode(kpiRow, 0)

	table := newWidget("cust-table", wCustTable).
		SetFlexGrow(1).SetFlexShrink(1).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetPadding(goda.EdgeAll, 16)
	table.SetContext(Widget{Kind: wCustTable, Data: customersList})
	main.InsertChildNode(table, 1)
}

func buildOrdersView(main *goda.Node) {
	kpiRow := buildKPIRow(ordersKPIs)
	main.InsertChildNode(kpiRow, 0)

	table := buildOrderTable("All Orders", ordersList)
	main.InsertChildNode(table, 1)
}

func buildProductsView(main *goda.Node) {
	kpiRow := buildKPIRow(productsKPIs)
	main.InsertChildNode(kpiRow, 0)

	grid := goda.New("product-grid").
		SetFlexDirection(goda.FlexDirectionRow).
		SetFlexWrap(goda.WrapWrap).
		SetGap(goda.GutterAll, 12).
		SetFlexGrow(1).SetFlexShrink(1)

	for _, p := range productsList {
		card := newWidget("", wProdCard).
			SetWidth(210).SetMinWidth(180).SetHeight(170).SetFlexShrink(0).
			SetFlexDirection(goda.FlexDirectionColumn).
			SetPadding(goda.EdgeAll, 14).SetGap(goda.GutterAll, 6)
		card.SetContext(Widget{Kind: wProdCard, Data: p})
		grid.InsertChildNode(card, grid.GetChildCount())
	}
	main.InsertChildNode(grid, 1)
}

func buildReportsView(main *goda.Node) {
	grid := goda.New("report-grid").
		SetFlexDirection(goda.FlexDirectionRow).
		SetFlexWrap(goda.WrapWrap).
		SetGap(goda.GutterAll, 14).
		SetFlexGrow(1).SetFlexShrink(1)

	for _, r := range reportsList {
		card := newWidget("", wReportCard).
			SetWidth(300).SetHeight(130).SetFlexShrink(0).
			SetFlexDirection(goda.FlexDirectionColumn).
			SetPadding(goda.EdgeAll, 16).SetGap(goda.GutterAll, 6)
		card.SetContext(Widget{Kind: wReportCard, Data: r})
		grid.InsertChildNode(card, grid.GetChildCount())
	}
	main.InsertChildNode(grid, 0)
}

func buildSettingsView(main *goda.Node) {
	sections := []struct {
		title string
		items []SettingRow
	}{
		{"General", settingsGeneral},
		{"Notifications", settingsNotifications},
		{"Account", settingsAccount},
	}

	for si, sec := range sections {
		section := goda.New("setting-group").
			SetFlexDirection(goda.FlexDirectionColumn).
			SetGap(goda.GutterAll, 4).
			SetFlexShrink(0).
			SetPadding(goda.EdgeAll, 16)

		label := goda.New("").
			SetHeight(24).SetFlexShrink(0)
		label.SetContext(Widget{Kind: WidgetKind(""), Data: nil})
		section.InsertChildNode(label, 0)

		for _, row := range sec.items {
			item := newWidget("", wSettingRow).
				SetHeight(40).SetFlexShrink(0).
				SetFlexDirection(goda.FlexDirectionRow).
				SetAlignItems(goda.AlignCenter).
				SetPadding(goda.EdgeLeft, 8).SetPadding(goda.EdgeRight, 14)
			item.SetContext(Widget{Kind: wSettingRow, Data: row})
			section.InsertChildNode(item, section.GetChildCount())
		}
		main.InsertChildNode(section, si)
	}
}

func buildKPIRow(cards []KPIData) *goda.Node {
	row := newWidget("kpi-row", wKPIRow).
		SetFlexDirection(goda.FlexDirectionRow).
		SetGap(goda.GutterAll, 14).
		SetFlexShrink(0).SetHeight(110)

	n := len(cards)
	minCardW := float32(140)
	if n > 0 {
		available := float32(designWidth) - 220 - 40 - float32((n-1)*14)
		minCardW = available / float32(n)
		if minCardW < 120 {
			minCardW = 120
		}
	}

	for i, k := range cards {
		card := newWidget("", wKPICard).
			SetFlexGrow(1).SetMinWidth(minCardW).SetMinHeight(90).
			SetFlexDirection(goda.FlexDirectionColumn).
			SetPadding(goda.EdgeAll, 14).SetGap(goda.GutterAll, 4)
		card.SetContext(Widget{Kind: wKPICard, Data: k})
		row.InsertChildNode(card, i)
	}
	return row
}

func buildChartCard(title, chartType string, data []float64, labels []string, color string) *goda.Node {
	card := newWidget("chart-"+title, wChartCard).
		SetFlexGrow(1).SetFlexShrink(1).SetMinWidth(250).
		SetPadding(goda.EdgeAll, 16)
	card.SetContext(Widget{Kind: wChartCard, Data: ChartSeries{
		Type:   chartType,
		Title:  title,
		Data:   data,
		Labels: labels,
		Color:  color,
	}})
	return card
}

func buildOrderTable(title string, orders []OrderRow) *goda.Node {
	table := newWidget("order-table", wOrderTable).
		SetFlexShrink(0).SetMinHeight(140).
		SetFlexDirection(goda.FlexDirectionColumn).
		SetPadding(goda.EdgeAll, 16)
	table.SetContext(Widget{Kind: wOrderTable, Data: orders})
	return table
}
