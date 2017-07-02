package entity

type Result struct {
	//Code   string 	   `json:""`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"message"`
	//Err    interface{} `json:""`
}
