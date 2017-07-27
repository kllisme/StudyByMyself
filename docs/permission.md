
# Group Permission

## Permission [/permission]

### GET [GET /permission/{id}]

+ Parameters
    + id: 404(number,required)

+ Response 200

### POST
+ Response 200

### DELETE [DELETE /permission/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

### PUT [PUT /permission/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

### AssignMenus [PUT /permission/{id}/menus]
+ Parameters
    + id: 404(number,required)

+ Response 200

### AssignElements [PUT /permission/{id}/elements]
+ Parameters
    + id: 404(number,required)

+ Response 200

### AssignActions [PUT /permission/{id}/actions]
+ Parameters
    + id: 404(number,required)

+ Response 200

## Role [/role]

### GET [GET /role/{id}]

+ Parameters
    + id: 404(number,required)

+ Response 200

### POST
+ Response 200

### DELETE [DELETE /role/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

### PUT [PUT /role/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

### AssignPermissions [PUT /role/{id}/permissions]
+ Parameters
    + id: 404(number,required)

+ Response 200


## Menu [/menu]

### GET [GET /menu/{id}]

+ Parameters
    + id: 404(number,required)

+ Response 200

### POST
+ Response 200

### DELETE [DELETE /menu/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

### PUT [PUT /menu/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

## Element [/element]

### GET [GET /element/{id}]

+ Parameters
    + id: 404(number,required)

+ Response 200

### POST
+ Response 200

### DELETE [DELETE /element/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

### PUT [PUT /element/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 200

## Action [/action]

### GET [GET /action/{id}]

+ Parameters
    + id: 404(number,required)

+ Response 27030100

       + Attributes (object)
        + message: ` `(string)
        + data (Action)
        + status: OK(string)

### Query [GET /actions{?method,handler_name}]

+ Parameters
    + method: get(string,optional)
    + handler_name: usercontroller(string,optional) - partial matching query,case insensitive!

+ Response 27030200

    + Attributes (object)
        + message: ` `(string)
        + data (array[Action])
        + status: OK(string)

### POST
+ Request (application/json)

    + Attributes (Action)

+ Response 27030300
 
    + Attributes (object)
        + message: ` `(string)
        + data (Action)
        + status: OK(string)

+ Response 27030301

            {
                "status": "INTERNAL_SERVER_ERROR",
                "data": null,
                "message": ""
            } 

### DELETE [DELETE /action/{id}]
+ Parameters
    + id: 404(number,required)

+ Response 27030400

            {
                "status": "OK",
                "data": null,
                "message": ""
            } 

### PUT [PUT /action/{id}]
+ Parameters
    + id: 404(number,required)

+ Request (application/json)

    + Attributes (Action)    

+ Response 27030500
 
    + Attributes (object)
        + message: ` `(string)
        + data (Action)
        + status: OK(string)

+ Response 27030501

            {
                "status": "INTERNAL_SERVER_ERROR",
                "data": null,
                "message": ""
            } 
