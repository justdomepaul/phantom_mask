package entity

type PurchaseHistoryJSON struct {
	PharmacyName      string  `json:"pharmacyName,omitempty"`
	MaskName          string  `json:"maskName,omitempty"`
	TransactionAmount float64 `json:"transactionAmount,omitempty"`
	TransactionDate   string  `json:"transactionDate,omitempty"`
}

type UserJSON struct {
	Name              string                `json:"name,omitempty"`
	CashBalance       float64               `json:"cashBalance,omitempty"`
	PurchaseHistories []PurchaseHistoryJSON `json:"purchaseHistories,omitempty"`
}

type MaskJSON struct {
	Name  string  `json:"name,omitempty"`
	Price float64 `json:"price,omitempty"`
}

type PharmacyJSON struct {
	Name         string     `json:"name,omitempty"`
	CashBalance  float64    `json:"cashBalance,omitempty"`
	OpeningHours string     `json:"openingHours,omitempty"`
	Masks        []MaskJSON `json:"masks,omitempty"`
}
