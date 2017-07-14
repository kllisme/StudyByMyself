
# Group 用户

>B端用户API群

## 无身份操作 [/]

### 登录	[POST /login]

### 登出 [GET /logout]

### 注册 [POST /sign_up]

### 验证码 [GET /captcha]

## 用户信息 [/user]

>该接口负责用户密码重置，会自动从session中提取当前登录用户的ID

### 拉取基本信息 [GET /user/{id}]
 
### 更新信息 [PUT /user/{id}]

### 重置 C & B 端用户密码 [POST /user/reset-password]

+ Request (application/json)
  + body

            {
                "type": 1,
                "password": "123456",
                "account":"kefu"
            }
  + Schema

            {
                "type":"object",
                "properties":{
                    "type": {
                        "description":"用户账户类型"
                        "type":"enum",
                        "enum":[{"B端用户":1},{"C端用户":2}]
                    },
                    "password": {
                        "type":"string"
                    },
                    "account":{
                        "type":"string"
                    }
                },
                "required":["type","password","account"]
            }

+ Response 01011700 (application/json)
  + Body

            {
                "status": 01011700,
                "msg":"重置密码成功",
                "data":{}
            }

### 重置 C & B 端用户密码 [PUT /user/reset-password]

参数 | 类型 | 描述
--:| ---- | -----------
type | required, number  | 用户账户类型，1为B端用户，2为C用户
password | required,string  | 重置后的用户密码，明文
account | required,string  | 用户登录账号

+ Request (application/json)
  + body

            {
                "type": 1,
                "password": "123456",
                "account":"kefu"
            }
  + Schema

            {
                "type":"object",
                "properties":{
                    "type": {
                        "description":"用户账户类型"
                        "type":"enum",
                        "enum":[{"B端用户":1},{"C端用户":2}]
                    },
                    "password": {
                        "type":"string"
                    },
                    "account":{
                        "type":"string"
                    }
                },
                "required":["type","password","account"]
            }

+ Response 01011700 (application/json)
  + Body

            {
                "status": 01011700,
                "msg":"重置密码成功",
                "data":{}
            }
