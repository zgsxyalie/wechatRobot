package model

type Contact struct {
	WxId     string `json:"wx_id" binding:"required"`
	Code     string `json:"code"`
	Remark   string `json:"remark"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Gender   int32  `json:"gender"`
}

func (Contact) TableName() string {
	return "contact"
}
