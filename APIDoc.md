FORMAT: 1A
HOST: http://api.erp.sodalife.xyz/v1

# Sodalife ERP API Doc

## Resume

>This is an api doc for Soda Systerm

* 系统架构
  * 本API Server使用GO 1.8.3 编写， 基于[iris v5](https://github.com/kataras/iris/tree/5.0.0)框架开发，使用 JWT 验证，除 Common 类外，所有 API 调用都需提供 Authorization 请求头，附带服务器颁发的 token。登录成功后将返回 token。

  * 所有请求依赖于 Cookie ，注意不要禁用，以及设置 XMLHttpRequest 的 WithCredentials 属性为 true。

# Data Structures

## Result (object)

+ message: `权限不足，拒绝访问！`(string) - Brief description of the Response status.
+ data - Can be anything for your wish.
+ status: FORBBIDEN(string) - Response status.

## Pagination (object)

+ total: 15536(number)
+ from: 20(number)
+ to: 40(number)

## Model (object)

+ id: 404(number)
+ createdAt: `2017-04-19T12:24:56+08:00`(string)
+ deletedAt: `0001-01-01T00:00:00Z`(string)
+ updatedAt: `2017-05-09T12:04:07+08:00`(string)

## User (object)

+ Include Model
+ name: martin(string)
+ concact: martini(string)
+ address: 科兴科学园B3(string)
+ mobile: 13260644577(string)
+ account: martin(string)
+ password: e10adc3949ba59abbe56e057f20f883e(string)
+ telephone: +86 18575534464(string)
+ email: martin@hyx.com(string)
+ parentId: 5(number)
+ gender: 0(number)
+ age: 24(number)
+ status: 0(number)

## Menu (object)

+ Include Model
+ name: 运营商管理(string)
+ icon: laptop(string)
+ url: `/user`(string)
+ parentId: 0(number)
+ level: 1(number)
+ status: 0(number)

## Element (object)

+ Include Model
+ name: 删除用户按钮(string)
+ reference: `/user/#deleteBtn`(string)

## Action (object)

+ Include Model
+ handlerName: `maizuo.com/soda/erp/api/src/server/controller/api(*UserController)Login-fm`(string)
+ Method (enum)
    + get(string)
    + post(string)
    + put(string)
    + delete(string)

<!-- include(docs/users.md) -->

<!-- include(docs/permission.md) -->

