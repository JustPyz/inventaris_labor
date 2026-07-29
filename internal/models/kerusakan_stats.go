package models

type KerusakanStats struct {
	ItemInstanceID uint  `json:"id_item_instance"`
	TotalPerbaikan int64 `json:"total_perbaikan"`
	TotalBiaya     int64 `json:"total_biaya"`
}
