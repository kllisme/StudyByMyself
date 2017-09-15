package model

type ReferenceDevice struct {
	Model
	Name string	`json:"name"`
}

func (ReferenceDevice) TableName() string {
	return "reference_device"
}
