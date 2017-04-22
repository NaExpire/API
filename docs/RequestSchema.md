# Business Endpoints

## Login
Endpoint: POST /api/business/login/ <br />
Request:
```json
{
    "email": string,
    "password": string
}
```

Response:
```json
{
    "ok": boolean,
    "sessionID": string
}
```

## Logout
Endpoint: POST /api/business/logout/ <br />
Headers:
```
session: sessionID
```

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Register
Endpoint: POST /api/business/register/ <br />
Request
```json
{
    "firstName": string,
    "lastName": string,
    "email": string,
    "password": string,
    "personalPhoneNumber": string,
    "restaurantName": string,
    "addressLine1": string,
    "addressLine2": string,
    "city": string,
    "state": string,
    "zip": string,
    "businessPhoneNumber": string,
    "description": string
}
```

Response
```json
{
    "ok": boolean
}
```

## Get Restaurant details
Endpoint: GET /api/business/restaurant/<restaurantID:int>/ <br />

Request: no request schema necessary <br />

Response
```json
{
    "name": string,
    "description": string,
    "address": string,
    "city": string,
    "state": string,
    "pickupTime": string,
    "phoneNumber": string,
    "items": string,
    "meals": [{
        "name": string,
        "description": string,
        "restaurantID": int,
        "price": double,
        "type": string
    }],
    "deals": [{
        "mealID": int,
        "dealPrice": double,
        "quantity": int,
        "restaurantID": int
    }]
}
```

## Update Restaurant details
Endpoint: POST /api/business/restaurant/<restaurantID:int>/update/ <br />

Request
```json
{
    "name": string,
    "description": string,
    "address": string,
    "city": string,
    "state": string,
    "pickupTime": string,
    "phoneNumber": string,
    "items": string
}
```

Response 
```json
{
    "ok": boolean
}

## Confirm Registration
Endpoint: POST /api/business/register/confirm/ <br />
Request
```json
{
    "emailAddress": string,
    "confirmationCode": string
}
```

Response
```json
{
    "ok": boolean
}
```

## Create meal
Endpoint: POST /api/business/meal/create/

Request:
```json
{
    "name": string,
    "description": string,
    "restaurantID": int,
    "price": double,
    "type": string
}
```
type is either "menu-item" or "grab-bag"

Response:
```json
{
    "ok": boolean,
    "id": int
}
```

## Get meal info
Endpoint: GET /api/business/meal/<mealID:int>/ <br />

Request: no request schema necessary <br />

Response: 
```json
{
    "name": string,
    "description": string,
    "restaurantID": int,
    "price": double,
    "type": string
}
```
type is either "menu-item" or "grab-bag"

## Update meal info
Endpoint: PUT /api/business/meal/<mealID:int>/update/ <br />

Request:
```json
{
    "name": string,
    "description": string,
    "restaurantID": int,
    "price": double
}
```

Response:
```json
{
    "ok": boolean
}
```

## Delete meal
Endpoint DELETE /api/business/meal/<mealID:int>/delete/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Create deal
Endpoint: POST /api/business/deal/create/

Request:
```json
{
    "mealID": int,
    "deal-price": double,
    "quantity": int
}
```

Response:
```json
{
    "ok": boolean,
    "id": int
}
```

## Get deal info
Endpoint: GET /api/business/deal/<dealID:int>/ <br />

Request: no request schema necessary <br />

Response: 
```json
{
    "mealID": int,
    "dealPrice": double,
    "quantity": int,
    "restaurantID": int
}
```

## Update deal info
Endpoint: PUT /api/business/deal/<dealID:int>/update/ <br />

Request:
```json
{
    "mealID": int,
    "dealPrice": double,
    "quantity": int
}
```

Response:
```json
{
    "ok": boolean
}
```

## Delete deal
Endpoint DELETE /api/business/deal/<dealID:int>/delete/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Accept Transaction
Endpoint PUT /api/business/transaction/<transactionID:int>/accept/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Reject Transaction 
Endpoint PUT /api/business/transaction/<transactionID:int>/reject/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Get Review
Endpoint GET /review/<reviewID:int>/ <br />

Request: no request schema necessary <br />

Response:
```json
	"restaurantID": int
	"score": int
	"reviewBody": string
```

# Consumer Endpoints 
## Login
Endpoint: POST /api/consumer/login/ <br />
Request:
```json
{
    "email": string,
    "password": string
}
```

Response:
```json
{
    "ok": boolean,
    "sessionID": string
}
```

## Logout
Endpoint: POST /api/consumer/logout/ <br />
Headers:
```
session: sessionID
```

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Register
Endpoint: POST /api/consumer/register/ <br />
Request
```json
{
    "firstName": string,
    "lastName": string,
    "email": string,
    "password": string,
    "personalPhoneNumber": string,
}
```

Response
```json
{
    "ok": boolean
}
```

## Confirm Registration
Endpoint: POST /api/consumer/register/confirm/ <br />
Request
```json
{
    "emailAddress": string,
    "confirmationCode": string
}
```

Response
```json
{
    "ok": boolean
}
```

## Get Restaurant details
Endpoint: GET /api/consumer/restaurant/<restaurantID:int>/ <br />

Response
```json
{
    "name": string,
    "description": string,
    "address": string,
    "city": string,
    "state": string
}
```

## Get meal info
Endpoint: GET /api/consumer/meal/<mealID:int>/ <br />

Request: no request schema necessary <br />

Response: 
```json
{
    "name": string,
    "description": string,
    "restaurantID": int,
    "price": double,
    "type": string
}
```
type is either "menu-item" or "grab-bag"

## Get deal info
Endpoint: GET /api/consumer/deal/<dealID:int>/ <br />

Request: no request schema necessary <br />

Response: 
```json
{
    "mealID": int,
    "dealPrice": double,
    "quantity": int,
    "restaurantID": int
}
```

## Issue Transaction
blocked

## Cancel Transaction
Endpoint PUT /api/consumer/transaction/<transactionID:int>/cancel/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Fulfil Transaction
Endpoint PUT /api/consumer/transaction/<transactionID:int>/fulfill/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Get Cart
Endpoint GET /api/consumer/cart/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "menuitems": [
        ...
        {
            "name": string,
            "description": string,
            "restaurantID": int,
            "price": double,
            "type": string
        },
        ...
    ],
    "deals": [
        ...
        {
            "mealID": int,
            "dealPrice": double,
            "quantity": int,
        },
        ...
    ]
}
```

## Add Meal to Cart
Endpoint POST /api/consumer/cart/add/meal/ <br />

Request:
```json
{
    "mealID": int,
    "quantity": int
}
```

Response:
```json
{
    "ok": boolean
}
```

## Add Deal to Cart
Endpoint POST /api/consumer/cart/add/deal/ <br />

Request:
```json
{
    "dealID": int
}
```

Response:
```json
{
    "ok": boolean
}
```

## Delete Cart Contents
Endpoint DELETE /api/consumer/cart/empty/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Update Quantity Meal for Cart
Endpoint PUT /api/consumer/cart/meal/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Delete Meal for Cart
Endpoint DELETE /api/consumer/cart/meal/empty/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
```

## Delete Meal for Cart
Endpoint DELETE /api/consumer/cart/delete/deal/ <br />

Request: 
{
    "mealID": int
}

Response:
```json
{
    "ok": boolean
}
```

## Delete Deal for Cart
Endpoint DELETE /api/consumer/cart/delete/deal/ <br />

Request: 
{
    "dealID": int
}

Response:
```json
{
    "ok": boolean
}
```

## Get Review
Endpoint GET /review/<reviewID:int>/ <br />

Request: no request schema necessary <br />

Response:
```json
{
	"restaurantID": int,
	"score": int,
	"reviewBody": string
}
```

## Create Review
Endpoint POST /review/create/ <br />

Request:
```json
{
	"restaurantID": int,
	"score": int,
	"reviewBody": string
}
```

Response:
```json
{
    "ok": boolean
}
```

## Update Review
Endpoint PUT /review/<reviewID:int>/update/ <br />

Request: 
```json
{
	"restaurantID": int,
	"score": int,
	"reviewBody": string
}
```

Response:
```json
{
    "ok": boolean
}
```

## Delete Review
Endpoint DELETE /review/<reviewID:int>/delete/ <br />

Request: no request schema necessary <br />

Response:
```json
{
    "ok": boolean
}
