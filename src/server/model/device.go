package model

type Device struct {
	Preset      string        `json:"preset"`
	Name        string        `json:"name"`
	Value       float64       `json:"value"`
	Unit        string        `json:"unit"`
	Description string        `json:"description"`
}

const (
	LAUNDRY_DRY int = iota + 1
	LAUNDRY_EXPRESS
	LAUNDRY_STANDARD
	LAUNDRY_HEAVY
	LAUNDRY_EXPRESS_WITH_DETERGENT
	LAUNDRY_STANDARD_WITH_DETERGENT
	LAUNDRY_HEAVY_WITH_DETERGENT
)
