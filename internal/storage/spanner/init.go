package spanner

import (
	"phantom_mask/internal/storage"
)

var (
	DBCreatedTime = "CreatedTime"
)

func withTimeOrder(orderEnum storage.OrderListEnum) string {
	return map[storage.OrderListEnum]string{
		storage.CreatedTimeASC:  " ORDER BY CreatedTime ASC",
		storage.CreatedTimeDESC: " ORDER BY CreatedTime DESC",
		storage.UpdatedTimeASC:  " ORDER BY UpdatedTime ASC",
		storage.UpdatedTimeDESC: " ORDER BY UpdatedTime DESC",
		storage.PharmacyNameASC: " ORDER BY Name ASC",
		storage.ProductNameASC:  " ORDER BY Name ASC",
		storage.ProductPriceASC: " ORDER BY Price ASC",
		storage.PharmacyProduct: " ORDER BY PharmacyName ASC, ProductName ASC",
	}[orderEnum]
}
