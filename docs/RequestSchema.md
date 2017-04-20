# Business Endpoints

## Login
Endpoint: POST /api/business/login/ <br />
Request:
```json
{
    "username": string,
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


# Consumer Endpoints
## Login
Endpoint: POST /api/business/login/ <br />
Request:
```json
{
    "username": string,
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