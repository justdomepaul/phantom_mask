package spanner

import "phantom_mask/internal/storage"

type Set struct {
	Pharmacy        storage.IPharmacy
	PharmacyInfo    storage.IPharmacyInfo
	Product         storage.IProduct
	User            storage.IUser
	PurchaseHistory storage.IPurchaseHistory
}
