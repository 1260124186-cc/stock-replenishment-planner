package domain

import "fmt"

var ErrUnknownSKU = fmt.Errorf("unknown SKU")

type OrderSignal struct {
	SKU          string `json:"sku"`
	OnHand       int    `json:"on_hand"`
	OpenPurchase int    `json:"open_purchase"`
	DailySales   int    `json:"daily_sales"`
	DaysOfCover  int    `json:"days_of_cover"`
}

type ReorderPolicy struct {
	SKU             string `json:"sku"`
	MinimumStock    int    `json:"minimum_stock"`
	ReorderMultiple int    `json:"reorder_multiple"`
	LeadTimeDays    int    `json:"lead_time_days"`
	SafetyStock     *int   `json:"safety_stock,omitempty"`
}

type ReplenishmentPlan struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Reason   string `json:"reason"`
}
