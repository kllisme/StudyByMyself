package entity

type PaginationData struct {
	Pagination Pagination  `json:"pagination"`
	Objects    interface{} `json:"objects"`
}

type Pagination struct {
	Total int        `json:"total"`
	From  int        `json:"from"`
	To    int        `json:"to"`
}
