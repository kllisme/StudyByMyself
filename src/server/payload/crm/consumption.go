package crm

import "time"

type Consumption struct {
	TicketID             string        `json:"ticketId"`
	Agency               string        `json:"agency"`
	Operator             string        `json:"operator"`
	ParentOperator       string        `json:"parentOperator"`
	ParentOperatorMobile string        `json:"parentOperatorMobile"`
	Mobile               string        `json:"mobile"`
	Telephone            string        `json:"telephone"`
	DeviceSerial         string        `json:"deviceSerial"`
	Address              string        `json:"address"`
	CustomerMobile       string        `json:"customerMobile"`
	Password             string        `json:"password"`
	TypeName             string        `json:"typeName"`
	Value                int        `json:"value"`
	Payment              string        `json:"payment"`
	PaymentID            int        `json:"paymentId"`
	CreatedAt            time.Time        `json:"createdAt"`
	Status               int        `json:"status"`
}
