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
Endpoint: POST /api/login
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
Endpoint: POST /api/register
Request
```json
{
    "name": <string>,
    "address": <string>,
    "phoneNumber": <string>,
    "username": <string>,
    "description": <string>,
    "email": <string>,
    "personalPhoneNumber": <string>,
    "password": <string>,
    "foodTypes": [<string>...],
    "restaurantId": <string>,
    optional "menuItems": [<menuItem> ...]
}
```

Response
```json
{
    "ok": <boolean>
}
```

## Get Restaurant Details
GET /api/restaurant/<restaurantId:string>
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
GET /api/restaurant///menu
Response
```json
{
    "menuItems": [<menuItem> ...]
}
```

## Update Menu Item
GET /api/restaurant/<restaurantId:string>/menu