
# Group User

## Common [/]

### Login [GET /login{?account,password,captcha}]

+ Parameters
    + account: martin (string,required)
    + password: e10adc3949ba59abbe56e057f20f883e (string,required)
    + captcha: 1234 (string,required)

+ Response 27010100

    + Body

            {
                "status": "OK",
                "data": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDEwMzU3NjAsImlzcyI6ImFwaS5lcnAuc29kYWxpZmUueHl6Iiwic2Vzc2lvbklkIjoiVEVodGNXTlhWbVJHVW1wd2JIaEJSa2RvVG01WGVFdHdZbU5VUVV4SmFFND0ifQ.BeCOfvumGNL6ubsW0c6uakN7CRPvpgxAeShpAzRpoJ0",
                "message": ""
            }           

+ Response 27010106

    + Body

            {
                "status": "CAPTCHA_REQUIRED",
                "data": null,
                "message": "验证码超时，请重新输入"
            } 

+ Response 27010107

    + Body

            {
                "status": "CAPTCHA_REQUIRED",
                "data": null,
                "message": "验证码错误"
            } 

+ Response 27010108

    + Body

            {
                "status": "NOT_FOUND_ENTITY",
                "data": null,
                "message": "找不到该账户"
            } 

+ Response 27010109

    + Body

            {
                "status": "UNPROCESSABLE_ENTITY",
                "data": null,
                "message": "登录密码错误，请检查"
            } 

### SignUp [POST /sign-up]

+ Response 200

### Captcha [GET /captcha.png]

+ Response 200 (image/png)

### Logout [GET /logout]

+ Request
    
    + Headers

            Cookie:sess=TEhtcWNXVmRGUmpwbHhBRkdoTm5XeEtwYmNUQUxJaE4=; path=/; domain=.api.erp.sodalife.xyz; HttpOnly; Expires=Wed Jul 26 2017 02:22:29 GMT+0800 (CST);
            Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDEwMzU3NjAsImlzcyI6ImFwaS5lcnAuc29kYWxpZmUueHl6Iiwic2Vzc2lvbklkIjoiVEVodGNXTlhWbVJHVW1wd2JIaEJSa2RvVG01WGVFdHdZbU5VUVV4SmFFND0ifQ.BeCOfvumGNL6ubsW0c6uakN7CRPvpgxAeShpAzRpoJ0


+ Response 27010200

    + Body

            {
                "status": "OK",
                "data": null,
                "message": ""
            } 

## Profile [/profile]

### Session [GET /profile/session]
> Get current user profile, permission lists from server session.

+ Request
    
    + Headers

            Cookie:sess=TEhtcWNXVmRGUmpwbHhBRkdoTm5XeEtwYmNUQUxJaE4=; path=/; domain=.api.erp.sodalife.xyz; HttpOnly; Expires=Wed Jul 26 2017 02:22:29 GMT+0800 (CST);
            Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDEwMzU3NjAsImlzcyI6ImFwaS5lcnAuc29kYWxpZmUueHl6Iiwic2Vzc2lvbklkIjoiVEVodGNXTlhWbVJHVW1wd2JIaEJSa2RvVG01WGVFdHdZbU5VUVV4SmFFND0ifQ.BeCOfvumGNL6ubsW0c6uakN7CRPvpgxAeShpAzRpoJ0


+ Response 200

    + Attributes
        + user (User)
        + menuList (array[Menu])
        + elementList (array[Element])
        + actionList: [](array[Action]) - unavaliable for front end user

### Password [PUT /profile/password]
>reset password

+ Request
    
    + Headers

            Cookie:sess=TEhtcWNXVmRGUmpwbHhBRkdoTm5XeEtwYmNUQUxJaE4=; path=/; domain=.api.erp.sodalife.xyz; HttpOnly; Expires=Wed Jul 26 2017 02:22:29 GMT+0800 (CST);
            Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDEwMzU3NjAsImlzcyI6ImFwaS5lcnAuc29kYWxpZmUueHl6Iiwic2Vzc2lvbklkIjoiVEVodGNXTlhWbVJHVW1wd2JIaEJSa2RvVG01WGVFdHdZbU5VUVV4SmFFND0ifQ.BeCOfvumGNL6ubsW0c6uakN7CRPvpgxAeShpAzRpoJ0

    + Attributes
        + currentPassword: sodaErpApi(string,required)
        + newPassword: SodaErpApi(string,required)
+ Response 200

## User [/user]

+ Model 
    
            {
              "id": 404,
              "createdAt": "2017-04-19T12:24:56+08:00",
              "deletedAt": "0001-01-01T00:00:00Z",
              "updatedAt": "2017-05-09T12:04:07+08:00",
              "name": "martin",
              "concact": "martini",
              "address": "科兴科学园B3",
              "mobile": "13260644577",
              "account": "martin",
              "password": "e10adc3949ba59abbe56e057f20f883e",
              "telephone": "+86 18575534464",
              "email": "martin@hyx.com",
              "parentId": 5,
              "gender": 0,
              "age": 24,
              "status": 0
            }

## Create [POST /user] 

+ Request

    + Attributes (User)

+ Response 27010200

    + Attributes (object)
        + message: ` `(string)
        + data (User)
        + status: OK(string)
    
### Get [GET /user/{id}]

+ Parameters
    + id: 404(number,required)

+ Response 200

    + Attributes (object)
        + message: ` `(string)
        + data (User)
        + status: OK(string)
 
### Update [PUT /user/{id}]

+ Parameters
    + id: 404(number,required)

+ Request(application/json)

    + Attributes (User)

+ Response 200

    + Attributes (object)
        + message: ` `(string)
        + data (User)
        + status: OK(string)

### Password [PUT /user/{id}/password{?password}]

+ Request

    + Parameters
        + id: 404(number,required)
        + password: e10adc3949ba59abbe56e057f20f883e(string,required)

+ Response 01011700

    + Body

            {
                "status": 01011700,
                "msg":"重置密码成功",
                "data":{}
            }

## Users [/users]