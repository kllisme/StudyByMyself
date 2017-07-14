FORMAT: 1A
HOST: http://192.168.4.79:8080/api/v1

# 苏打生活平台系统API文档

## 说明

>This is an api doc for Soda Systerm

# Data Structures

## ResetPassword

+ type: 1 (required, number) - 用户账户类型，1为B端用户，2为C用户
+ password: 123456 (required) - 重置后的用户密码，明文
+ account: kefu (required) - 用户登录账号

<!-- include(docs/users.md) -->

<!-- include(docs/permission.md) -->

