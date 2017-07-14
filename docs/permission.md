
# Group 权限

>用户权限API群

## 用户权限 [/permission]

## 角色 [/permission/role]

### 获取用户角色

### 芯片卡功能权限 [GET /permission/is-chipcard-operator]

>直接验证当前用户是否拥有芯片卡余额转移功能的权限

+ Response 01011300 (applycation/json)

  + body

            {
                "status":"01011300",
                "msg":"拉取用户权限列表成功",
                "data":true
            }

  + schema

            {
                "title":"data"
                "type":"boolean"
            }

### 芯片卡功能权限 [PUT /permission/is-chipcard-operator]

>直接验证当前用户是否拥有芯片卡余额转移功能的权限

+ Response 01011300 (applycation/json)

  + body

            {
                "status":"01011300",
                "msg":"拉取用户权限列表成功",
                "data":true
            }

  + schema

            {
                "title":"data"
                "type":"boolean"
            }
