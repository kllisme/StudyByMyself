package model

import "time"

type User struct {
	ID        int        `json:"id"`
	Name      string      `json:"name"`      //名称
	Contact   string      `json:"contact"`   //联系人
	Address   string      `json:"address"`   //联系地址
	Mobile    string      `json:"mobile"`    //手机
	Account   string      `json:"account"`   //账号
	Password  string      `json:"password"`  //密码
	Telephone string      `json:"telephone"` //服务电话
	Email     string      `json:"email"`     //邮箱
	ParentId  int         `json:"parentId"`  //父id
	Gender    int         `json:"gender"`    //性别
	Age       int         `json:"age"`       //年龄
	Status    int         `json:"status"`    //状态
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (User) TableName() string {
	return "user"
}
