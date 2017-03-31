# Reusable schema 
## menuItem: 
```json
{
    "itemId": <string>,
    "name": <string>,
    "price": <string>,
    "description": <string>
}
```

## Login
Endpoint: POST /api/business/login/
Request:
```json
{
    "username": <string>,
    "password": <string>
}
```

Response:
```json
{
    "name": <string>,
    "username": <string>,
    "menuItems": [<menuItem> ...],
    "restaurantId": <string>
}
```

## Register
Endpoint: POST /api/business/register/
Request
```json
{
    "firstName": <string>,
    "lastName": <string>,
    "email": <string>,
    "password": <string>,
    "personalPhoneNumber": <string>,
    "restaurantName": <string>,
    "addressLine1": <string>,
    "addressLine2": <string>,
    "city": <string>,
    "state": <string>,
    "zip": <string>,
    "businessPhoneNumber": <string>,
    "description": <string>
}
```

Response
```json
{
    "ok": <boolean>
}
```

## Get Restaurant Details
GET /api/business/restaurant/<restaurantId:string>/
Response
```json
{
    "name": <string>,
    "address": <string>,
    "phoneNumber": <string>,
    "description": <string>,
    "restaurantId": <string>,
    "foodTypes": [<string>...]
}
```

## Get Menu
GET /api/business/restaurant/<restaurantId:string>/menu/<menuItemId:string>/
Response
```json
{
    "menuItems": [<menuItem> ...]
}
```

## Update Menu Item
POST /api/business/restaurant/<restaurantId:string>/menu/<menuItemId:string>/update/
Request
```json
{
    "name": <string>,
    "price": <string>,
    "description": <string>
}
```
Response
```json
{
    "ok": <boolean>
}
```