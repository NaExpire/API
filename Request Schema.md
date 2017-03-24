## Reusable schema 
### menuItem: 
```json
{
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
    "foodTypes": [<string>...]
}
```

## Get Menu
GET /api/restaurant/<restaurantId:string>/menu
Response
```json
{
    "menuItems": [<menuItem> ...]
}
```
