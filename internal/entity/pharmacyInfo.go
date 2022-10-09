package entity

// PRIMARY KEY(UID)
type PharmacyInfo struct {
	UID       []byte  `spanner:"UID" json:"uid,omitempty" validate:"required,max=16"`
	Day       int64   `spanner:"Day" json:"day,omitempty"`
	OpenHour  float64 `spanner:"OpenHour" json:"open_hour,omitempty"`
	CloseHour float64 `spanner:"CloseHour" json:"close_hour,omitempty"`
}
