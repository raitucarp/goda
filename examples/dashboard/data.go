package main

import "goda"

type WidgetKind string

const (
	wHeader      WidgetKind = "header"
	wLogo        WidgetKind = "logo"
	wSearch      WidgetKind = "search"
	wBell        WidgetKind = "bell"
	wUser        WidgetKind = "user"
	wSidebar     WidgetKind = "sidebar"
	wMenuItem    WidgetKind = "menu-item"
	wMain        WidgetKind = "main"
	wKPIRow      WidgetKind = "kpi-row"
	wKPICard     WidgetKind = "kpi-card"
	wChartsRow   WidgetKind = "charts-row"
	wChartCard   WidgetKind = "chart-card"
	wBigChart    WidgetKind = "big-chart"
	wOrderTable  WidgetKind = "order-table"
	wStatusBar   WidgetKind = "statusbar"
	wCustTable   WidgetKind = "cust-table"
	wProdCard    WidgetKind = "prod-card"
	wReportCard  WidgetKind = "report-card"
	wSettingRow  WidgetKind = "setting-row"
)

type Widget struct {
	Kind WidgetKind
	Data interface{}
}

type KPIData struct {
	Label    string
	Value    string
	Change   string
	Positive bool
	Color    string
}

type MenuData struct {
	Label  string
	Active bool
	Index  int
}

type ChartSeries struct {
	Type   string
	Title  string
	Data   []float64
	Labels []string
	Color  string
}

type OrderRow struct {
	ID       string
	Customer string
	Date     string
	Amount   string
	Status   string
}

type CustomerRow struct {
	Name   string
	Email  string
	Plan   string
	Joined string
	Spent  string
}

type ProductRow struct {
	Name     string
	Category string
	Price    string
	Stock    int
	Active   bool
	Color    string
}

type ReportCard struct {
	Title       string
	Description string
	Date        string
	Color       string
}

type SettingRow struct {
	Label    string
	Enabled  bool
}

var menuItems = []struct{ Label string }{
	{"Overview"}, {"Analytics"}, {"Customers"}, {"Orders"},
	{"Products"}, {"Reports"}, {"Settings"},
}

func newWidget(id string, kind WidgetKind) *goda.Node {
	n := goda.New(id)
	n.SetContext(Widget{Kind: kind})
	return n
}

var overviewKPIs = []KPIData{
	{Label: "Total Revenue", Value: "$45,231", Change: "+20.1%", Positive: true, Color: "#4F46E5"},
	{Label: "Active Users", Value: "2,847", Change: "+12.5%", Positive: true, Color: "#10B981"},
	{Label: "Conversion", Value: "3.24%", Change: "-0.4%", Positive: false, Color: "#F59E0B"},
	{Label: "Avg. Order", Value: "$59.27", Change: "+8.2%", Positive: true, Color: "#06B6D4"},
}

var overviewOrders = []OrderRow{
	{"#1001", "Alice Chen", "2026-06-28", "$79.99", "Delivered"},
	{"#1002", "Bob Martinez", "2026-06-28", "$199.99", "Processing"},
	{"#1003", "Carol Smith", "2026-06-27", "$49.99", "Shipped"},
	{"#1004", "David Kim", "2026-06-27", "$34.99", "Delivered"},
}

var analyticsKPIs = []KPIData{
	{Label: "Page Views", Value: "1.2M", Change: "+18.3%", Positive: true, Color: "#4F46E5"},
	{Label: "Sessions", Value: "345K", Change: "+12.1%", Positive: true, Color: "#10B981"},
	{Label: "Bounce Rate", Value: "32.1%", Change: "-2.4%", Positive: true, Color: "#F59E0B"},
	{Label: "Avg Duration", Value: "4m 23s", Change: "+8.7%", Positive: true, Color: "#06B6D4"},
}

var customersKPIs = []KPIData{
	{Label: "Total", Value: "8,492", Change: "+5.7%", Positive: true, Color: "#4F46E5"},
	{Label: "New (Month)", Value: "347", Change: "+22.1%", Positive: true, Color: "#10B981"},
	{Label: "Churn Rate", Value: "2.1%", Change: "-0.3%", Positive: true, Color: "#EF4444"},
}

var customersList = []CustomerRow{
	{"Emma Wilson", "emma@example.com", "Enterprise", "2025-03-15", "$12,847"},
	{"James Lee", "james@demo.co", "Pro", "2025-06-01", "$8,421"},
	{"Sofia Garcia", "sofia@mail.com", "Pro", "2025-09-12", "$5,230"},
	{"Lucas Brown", "lucas@demo.co", "Basic", "2026-01-08", "$1,902"},
	{"Mia Taylor", "mia@example.com", "Enterprise", "2025-04-22", "$15,600"},
}

var ordersKPIs = []KPIData{
	{Label: "Total Orders", Value: "1,847", Change: "+15.2%", Positive: true, Color: "#4F46E5"},
	{Label: "Pending", Value: "28", Change: "-5.1%", Positive: true, Color: "#F59E0B"},
	{Label: "Shipped", Value: "1,792", Change: "+17.8%", Positive: true, Color: "#10B981"},
	{Label: "Revenue", Value: "$283K", Change: "+21.3%", Positive: true, Color: "#06B6D4"},
}

var ordersList = []OrderRow{
	{"#2001", "Emma Wilson", "2026-07-01", "$2,450", "Delivered"},
	{"#2002", "Noah Davis", "2026-07-01", "$890", "Processing"},
	{"#2003", "Olivia Martin", "2026-06-30", "$1,200", "Shipped"},
	{"#2004", "Liam Johnson", "2026-06-30", "$340", "Delivered"},
	{"#2005", "Ava Williams", "2026-06-29", "$5,600", "Delivered"},
	{"#2006", "Ethan Jones", "2026-06-29", "$1,750", "Cancelled"},
}

var productsKPIs = []KPIData{
	{Label: "Total SKUs", Value: "156", Change: "+8", Positive: true, Color: "#4F46E5"},
	{Label: "Active", Value: "142", Change: "91%", Positive: true, Color: "#10B981"},
	{Label: "Low Stock", Value: "8", Change: "Alert", Positive: false, Color: "#EF4444"},
}

var productsList = []ProductRow{
	{"Pro Headphones", "Electronics", "$79.99", 142, true, "#4F46E5"},
	{"Smart Watch X", "Electronics", "$199.99", 67, true, "#10B981"},
	{"BT Speaker Mini", "Electronics", "$49.99", 230, true, "#F59E0B"},
	{"USB-C Dock Pro", "Accessories", "$34.99", 89, true, "#06B6D4"},
	{"Mech Keyboard", "Accessories", "$89.99", 5, true, "#EF4444"},
	{"Travel Backpack", "Lifestyle", "$59.99", 0, false, "#8B5CF6"},
}

var reportsList = []ReportCard{
	{"Monthly Sales", "Sales across all channels", "Jun 2026", "#4F46E5"},
	{"Weekly Users", "Acquisition and retention", "Jun 28", "#10B981"},
	{"Revenue Q2", "Quarterly revenue breakdown", "Jun 2026", "#F59E0B"},
	{"Campaign ROI", "Marketing campaign analysis", "Jun 2026", "#06B6D4"},
	{"Inventory Status", "Stock levels and alerts", "Jun 2026", "#EF4444"},
	{"Support Metrics", "Ticket volume and resolution", "Jun 2026", "#8B5CF6"},
}

var settingsGeneral = []SettingRow{
	{"Auto-save reports", true},
	{"Compact view", false},
	{"Show onboarding tips", true},
}

var settingsNotifications = []SettingRow{
	{"Email alerts", true},
	{"Push notifications", true},
	{"Weekly digest", false},
}

var settingsAccount = []SettingRow{
	{"Two-factor auth", true},
	{"Session timeout (30 min)", true},
}
