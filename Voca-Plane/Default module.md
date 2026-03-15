---
title: Default module
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# Default module

Base URLs:

* <a href="http://localhost:8000/api/v1">Local: http://localhost:8000/api/v1</a>

# Authentication

- HTTP Authentication, scheme: bearer

# Default

## GET Test

GET /

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## POST Seed

POST /system/seed

> Body Parameters

```json
{
  "app_password": "ambatutu"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|password-app|header|string| no |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "database reset and seeded successfully"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# User

## POST Login Super Admin

POST /auth/login

> Body Parameters

```json
{
  "email": "superadmin@flightbooking.com",
  "password": "admin123"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "login successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGZsaWdodGJvb2tpbmcuY29tIiwiZXhwIjoxNzcyODY2MzMxLCJpYXQiOjE3NzI3Nzk5MzEsImlkIjoxLCJyb2xlIjoic3VwZXJfYWRtaW4iLCJ0eXBlIjoiYWNjZXNzIn0.x2E_es_zwx_Go0bTgeRZBWpjxVw-M33UKN4pp2VLNY4",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzMzODQ3MzEsImlkIjoxLCJqdGkiOiIwZDRlNGEyOC1kMzU4LTRlZGItYjQ2NS1hNjFkNTY0N2JkNzMiLCJyb2xlIjoic3VwZXJfYWRtaW4iLCJ0eXBlIjoicmVmcmVzaCJ9.pQTEDsiKYVTeNmD8GxJEQdcp2vdY4xrcHiUdDv0yhX0",
    "expires_in": 86400,
    "token_type": "Bearer"
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## POST Register

POST /auth/register

> Body Parameters

```json
{
  "name": "Pasha",
  "email": "pashaprabasakti@gmail.com",
  "password": "360589"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": false,
  "message": "Key: 'RegisterRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## POST Refresh Token

POST /auth/refresh

> Body Parameters

```json
{
  "refresh_token": "{{rt}}"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": false,
  "message": "EOF"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Profile

## GET Get Profile

GET /user/profile

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "profile retrieved",
  "data": {
    "id": 1,
    "name": "s",
    "email": "admin@flightbooking.com",
    "role": "super_admin",
    "created_at": "2026-03-06T19:08:18.443122+07:00",
    "updated_at": "2026-03-06T19:47:06.35676+07:00"
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PATCH Update Profile

PATCH /user/profile

> Body Parameters

```json
{
  "name": "",
  "email": "",
  "password": ""
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "profile updated"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get My Device

GET /user/device-info

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "device retrieved successfully",
  "data": {
    "ip": "::1",
    "device_name": "DESKTOP-THMLE2N",
    "os": "windows",
    "ram_usage": "8.36 GB Total (terpakai: 89.0%)",
    "cpu_model": "11th Gen Intel(R) Core(TM) i5-1155G7 @ 2.50GHz",
    "internet_status": "Request Timeout (RTO)",
    "suggestion": "Peringatan: RAM hampir penuh, performa aplikasi mungkin melambat."
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Airlines

## POST Create Airlines

POST /admin/airlines

> Body Parameters

```yaml
name: Garuda Tester
code: GAS
logo: ""

```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|
|» name|body|string| no |none|
|» code|body|string| no |none|
|» logo|body|string(binary)| no |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PUT Edit Airlines

PUT /admin/airlines/{id}

> Body Parameters

```yaml
name: ""
code: ""
logo: ""

```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|
|» name|body|string| no |none|
|» code|body|string| no |none|
|» logo|body|string(binary)| no |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## DELETE Delete Airlines

DELETE /admin/airlines/{id}

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get AIrlines

GET /airlines

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |none|
|limit|query|string| no |none|
|sort_by|query|string| no |none|
|order|query|string| no |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "airlines retrieved",
  "data": [
    {
      "id": 1,
      "name": "Air Excursion, LLC",
      "code": "X4",
      "logo_url": "https://pics.avs.io/200/200/X4@2x.png"
    },
    {
      "id": 2,
      "name": "Vanilla Air",
      "code": "JW",
      "logo_url": "https://pics.avs.io/200/200/JW@2x.png"
    },
    {
      "id": 3,
      "name": "Olympic Air",
      "code": "OA",
      "logo_url": "https://pics.avs.io/200/200/OA@2x.png"
    },
    {
      "id": 4,
      "name": "Toki Air",
      "code": "BV",
      "logo_url": "https://pics.avs.io/200/200/BV@2x.png"
    },
    {
      "id": 5,
      "name": "Brussels Airlines",
      "code": "SN",
      "logo_url": "https://pics.avs.io/200/200/SN@2x.png"
    },
    {
      "id": 6,
      "name": "Iberia",
      "code": "IB",
      "logo_url": "https://pics.avs.io/200/200/IB@2x.png"
    },
    {
      "id": 7,
      "name": "Thai Airways",
      "code": "TG",
      "logo_url": "https://pics.avs.io/200/200/TG@2x.png"
    },
    {
      "id": 8,
      "name": "Small Planet Airlines",
      "code": "P7",
      "logo_url": "https://pics.avs.io/200/200/P7@2x.png"
    },
    {
      "id": 9,
      "name": "Germanwings",
      "code": "4U",
      "logo_url": "https://pics.avs.io/200/200/4U@2x.png"
    },
    {
      "id": 10,
      "name": "Jet Airways",
      "code": "9W",
      "logo_url": "https://pics.avs.io/200/200/9W@2x.png"
    }
  ],
  "meta": {
    "limit": 10,
    "page": 1,
    "total": 1149
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Airports

## POST Create Airports

POST /admin/airports

> Body Parameters

```json
{
  "code": "",
  "name": "",
  "city": ""
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get Airports

GET /admin/airports

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "airports retrieved",
  "data": [
    {
      "id": 1,
      "code": "CGK",
      "name": "Soekarno-Hatta",
      "city": "Jakarta",
      "created_at": "2026-03-06T19:08:18.320542+07:00",
      "updated_at": "2026-03-06T19:08:18.320542+07:00"
    },
    {
      "id": 2,
      "code": "DPS",
      "name": "Ngurah Rai",
      "city": "Bali",
      "created_at": "2026-03-06T19:08:18.326029+07:00",
      "updated_at": "2026-03-06T19:08:18.326029+07:00"
    },
    {
      "id": 3,
      "code": "SUB",
      "name": "Juanda",
      "city": "Surabaya",
      "created_at": "2026-03-06T19:08:18.3271+07:00",
      "updated_at": "2026-03-06T19:08:18.3271+07:00"
    },
    {
      "id": 4,
      "code": "KNO",
      "name": "Kualanamu",
      "city": "Medan",
      "created_at": "2026-03-06T19:08:18.327708+07:00",
      "updated_at": "2026-03-06T19:08:18.327708+07:00"
    },
    {
      "id": 5,
      "code": "UPG",
      "name": "Sultan Hasanuddin",
      "city": "Makassar",
      "created_at": "2026-03-06T19:08:18.328775+07:00",
      "updated_at": "2026-03-06T19:08:18.328775+07:00"
    },
    {
      "id": 6,
      "code": "",
      "name": "",
      "city": "",
      "created_at": "2026-03-09T10:46:00.618214+07:00",
      "updated_at": "2026-03-09T10:46:30.729698+07:00"
    }
  ],
  "meta": {
    "limit": 10,
    "page": 1,
    "total": 6
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PUT Edit Airports

PUT /admin/airports/{id}

> Body Parameters

```json
{
  "code": "",
  "name": "",
  "city": ""
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## DELETE Delete Airports

DELETE /admin/airports/{id}

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get Airports

GET /airports

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |none|
|limit|query|string| no |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "airports retrieved",
  "data": [
    {
      "id": 1,
      "code": "AKY",
      "name": "Sittwe Airport",
      "city": "Sittwe"
    },
    {
      "id": 2,
      "code": "KMJ",
      "name": "Kumamoto Airport",
      "city": "Kumamoto"
    },
    {
      "id": 3,
      "code": "NYK",
      "name": "Nanyuki Airport",
      "city": "Nanyuki"
    },
    {
      "id": 4,
      "code": "CBV",
      "name": "Coban Airport",
      "city": "Coban"
    },
    {
      "id": 5,
      "code": "GVR",
      "name": "Governador Valadares Airport",
      "city": "Governador Valadares"
    },
    {
      "id": 6,
      "code": "JQE",
      "name": "Jaque Airport",
      "city": "Jaque"
    },
    {
      "id": 7,
      "code": "BCB",
      "name": "Virginia Tech Montgomery Executive Airport",
      "city": "Blacksburg"
    },
    {
      "id": 8,
      "code": "SME",
      "name": "Lake Cumberland Regional Airport",
      "city": "Somerset"
    },
    {
      "id": 9,
      "code": "YXD",
      "name": "Edmonton City Centre (Blatchford Field) Airport",
      "city": "Edmonton"
    },
    {
      "id": 10,
      "code": "SAC",
      "name": "Sacramento Executive Airport",
      "city": "Sacramento"
    },
    {
      "id": 11,
      "code": "ESE",
      "name": "Ensenada Airport",
      "city": "Default"
    },
    {
      "id": 12,
      "code": "VLP",
      "name": "Vila Rica Airport",
      "city": "Vila Rica"
    },
    {
      "id": 13,
      "code": "BAX",
      "name": "Barnaul Airport",
      "city": "Barnaul"
    },
    {
      "id": 14,
      "code": "RRK",
      "name": "Rourkela Airport",
      "city": "Default"
    },
    {
      "id": 15,
      "code": "SYM",
      "name": "Simao Airport",
      "city": "Simao"
    },
    {
      "id": 16,
      "code": "GIY",
      "name": "Giyani Airport",
      "city": "Giyani"
    },
    {
      "id": 17,
      "code": "KKS",
      "name": "Kashan Airport",
      "city": "Default"
    },
    {
      "id": 18,
      "code": "NFR",
      "name": "Nafurah 1 Airport",
      "city": "Nafurah 1"
    },
    {
      "id": 19,
      "code": "ESS",
      "name": "Essen Mulheim Airport",
      "city": "Default"
    },
    {
      "id": 20,
      "code": "BBB",
      "name": "Benson Municipal Airport",
      "city": "Benson"
    },
    {
      "id": 21,
      "code": "XRY",
      "name": "Jerez Airport",
      "city": "Jerez de la Forntera"
    },
    {
      "id": 22,
      "code": "WTN",
      "name": "RAF Waddington",
      "city": "Waddington"
    },
    {
      "id": 23,
      "code": "SDN",
      "name": "Sandane Airport Anda",
      "city": "Sandane"
    },
    {
      "id": 24,
      "code": "NCH",
      "name": "Nachingwea Airport",
      "city": "Nachingwea"
    },
    {
      "id": 25,
      "code": "AGU",
      "name": "Jesus Teran International Airport",
      "city": "Aguascalientes"
    },
    {
      "id": 26,
      "code": "PTF",
      "name": "Malolo Lailai Island Airport",
      "city": "Malolo Lailai Island"
    },
    {
      "id": 27,
      "code": "MQK",
      "name": "San Matias Airport",
      "city": "San Matias"
    },
    {
      "id": 28,
      "code": "YLW",
      "name": "Kelowna Airport",
      "city": "Kelowna"
    },
    {
      "id": 29,
      "code": "SEE",
      "name": "Gillespie Field",
      "city": "San Diego/El Cajon"
    },
    {
      "id": 30,
      "code": "AEU",
      "name": "Abumusa Island Airport",
      "city": "Default"
    },
    {
      "id": 31,
      "code": "OKC",
      "name": "Will Rogers World Airport",
      "city": "Oklahoma City"
    },
    {
      "id": 32,
      "code": "NQN",
      "name": "Presidente Peron Airport",
      "city": "Neuquen"
    },
    {
      "id": 33,
      "code": "PVS",
      "name": "Provideniya Bay Airport",
      "city": "Chukotka"
    },
    {
      "id": 34,
      "code": "MGS",
      "name": "Mangaia Island Airport",
      "city": "Mangaia Island"
    },
    {
      "id": 35,
      "code": "VBC",
      "name": "Chanmyathazi Airport",
      "city": "Mandalay"
    },
    {
      "id": 36,
      "code": "HOA",
      "name": "Hola Airport",
      "city": "Hola"
    },
    {
      "id": 37,
      "code": "ENO",
      "name": "Encarnacion Airport",
      "city": "Encarnacion"
    },
    {
      "id": 38,
      "code": "NOD",
      "name": "Norden-Norddeich Airport",
      "city": "Norddeich"
    },
    {
      "id": 39,
      "code": "GOV",
      "name": "Gove Airport",
      "city": "Nhulunbuy"
    },
    {
      "id": 40,
      "code": "RBS",
      "name": "Orbost Airport",
      "city": "Default"
    },
    {
      "id": 41,
      "code": "SYK",
      "name": "Stykkisholmur Airport",
      "city": "Stykkisholmur"
    },
    {
      "id": 42,
      "code": "CVT",
      "name": "Coventry Airport",
      "city": "Coventry"
    },
    {
      "id": 43,
      "code": "INK",
      "name": "Winkler County Airport",
      "city": "Wink"
    },
    {
      "id": 44,
      "code": "IKU",
      "name": "Issyk-Kul International Airport",
      "city": "Tamchy"
    },
    {
      "id": 45,
      "code": "ZMT",
      "name": "Masset Airport",
      "city": "Masset"
    },
    {
      "id": 46,
      "code": "SUA",
      "name": "Witham Field",
      "city": "Stuart"
    },
    {
      "id": 47,
      "code": "HCZ",
      "name": "Chenzhou Beihu Airport",
      "city": "Chenzhou"
    },
    {
      "id": 48,
      "code": "OGD",
      "name": "Ogden Hinckley Airport",
      "city": "Ogden"
    },
    {
      "id": 49,
      "code": "PSM",
      "name": "Portsmouth International at Pease Airport",
      "city": "Portsmouth"
    },
    {
      "id": 50,
      "code": "TOG",
      "name": "Togiak Airport",
      "city": "Togiak Village"
    },
    {
      "id": 51,
      "code": "IKK",
      "name": "Greater Kankakee Airport",
      "city": "Kankakee"
    },
    {
      "id": 52,
      "code": "KOV",
      "name": "Kokshetau Airport",
      "city": "Kokshetau"
    },
    {
      "id": 53,
      "code": "LOS",
      "name": "Murtala Muhammed International Airport",
      "city": "Lagos"
    },
    {
      "id": 54,
      "code": "QUS",
      "name": "Gusau Airport",
      "city": "Gusau"
    },
    {
      "id": 55,
      "code": "CQS",
      "name": "Costa Marques Airport",
      "city": "Costa Marques"
    },
    {
      "id": 56,
      "code": "BIS",
      "name": "Bismarck Municipal Airport",
      "city": "Bismarck"
    },
    {
      "id": 57,
      "code": "SQR",
      "name": "Soroako Airport",
      "city": "Soroako-Celebes Island"
    },
    {
      "id": 58,
      "code": "OFJ",
      "name": "Olafsfjordur Airport",
      "city": "Olafsfjordur"
    },
    {
      "id": 59,
      "code": "GXX",
      "name": "Yagoua Airport",
      "city": "Yagoua"
    },
    {
      "id": 60,
      "code": "KGL",
      "name": "Kigali International Airport",
      "city": "Kigali"
    },
    {
      "id": 61,
      "code": "YQH",
      "name": "Watson Lake Airport",
      "city": "Watson Lake"
    },
    {
      "id": 62,
      "code": "OTL",
      "name": "Boutilimit Airport",
      "city": "Boutilimit"
    },
    {
      "id": 63,
      "code": "BNO",
      "name": "Burns Municipal Airport",
      "city": "Burns"
    },
    {
      "id": 64,
      "code": "ESN",
      "name": "Easton Newnam Field",
      "city": "Easton"
    },
    {
      "id": 65,
      "code": "PAH",
      "name": "Barkley Regional Airport",
      "city": "Paducah"
    },
    {
      "id": 66,
      "code": "KCO",
      "name": "Cengiz Topel Airport",
      "city": "Default"
    },
    {
      "id": 67,
      "code": "MHA",
      "name": "Mahdia Airport",
      "city": "Mahdia"
    },
    {
      "id": 68,
      "code": "CZM",
      "name": "Cozumel International Airport",
      "city": "Cozumel"
    },
    {
      "id": 69,
      "code": "DNK",
      "name": "Dnipropetrovsk International Airport",
      "city": "Dnipropetrovsk"
    },
    {
      "id": 70,
      "code": "TKC",
      "name": "Tiko Airport",
      "city": "Tiko"
    },
    {
      "id": 71,
      "code": "MLS",
      "name": "Frank Wiley Field",
      "city": "Miles City"
    },
    {
      "id": 72,
      "code": "ULM",
      "name": "New Ulm Municipal Airport",
      "city": "New Ulm"
    },
    {
      "id": 73,
      "code": "VIJ",
      "name": "Virgin Gorda Airport",
      "city": "Spanish Town"
    },
    {
      "id": 74,
      "code": "QOB",
      "name": "Ansbach-Petersdorf Airport",
      "city": "Ansbach"
    },
    {
      "id": 75,
      "code": "AOH",
      "name": "Lima Allen County Airport",
      "city": "Lima"
    },
    {
      "id": 76,
      "code": "GRK",
      "name": "Robert Gray  Army Air Field Airport",
      "city": "Fort Cavazos/Killeen"
    },
    {
      "id": 77,
      "code": "BMK",
      "name": "Borkum Airport",
      "city": "Borkum"
    },
    {
      "id": 78,
      "code": "NLC",
      "name": "Lemoore Naval Air Station (Reeves Field) Airport",
      "city": "Lemoore"
    },
    {
      "id": 79,
      "code": "TBI",
      "name": "New Bight Airport",
      "city": "Cat Island"
    },
    {
      "id": 80,
      "code": "KUV",
      "name": "Kunsan Air Base",
      "city": "Kunsan"
    },
    {
      "id": 81,
      "code": "BWQ",
      "name": "Brewarrina Airport",
      "city": "Default"
    },
    {
      "id": 82,
      "code": "KDX",
      "name": "Kadugli Airport",
      "city": "Kadugli"
    },
    {
      "id": 83,
      "code": "AUQ",
      "name": "Hiva Oa-Atuona Airport",
      "city": "Default"
    },
    {
      "id": 84,
      "code": "DBV",
      "name": "Dubrovnik Airport",
      "city": "Dubrovnik"
    },
    {
      "id": 85,
      "code": "KTY",
      "name": "Katukurunda Air Force Base",
      "city": "Kalutara"
    },
    {
      "id": 86,
      "code": "SCM",
      "name": "Scammon Bay Airport",
      "city": "Scammon Bay"
    },
    {
      "id": 87,
      "code": "VAL",
      "name": "Valenca Airport",
      "city": "Valenca"
    },
    {
      "id": 88,
      "code": "NUR",
      "name": "Nullabor Motel Airport",
      "city": "Default"
    },
    {
      "id": 89,
      "code": "CFS",
      "name": "Coffs Harbour Airport",
      "city": "Coffs Harbour"
    },
    {
      "id": 90,
      "code": "JYR",
      "name": "Jiroft Airport",
      "city": "Default"
    },
    {
      "id": 91,
      "code": "AGN",
      "name": "Angoon Seaplane Base",
      "city": "Angoon"
    },
    {
      "id": 92,
      "code": "EYP",
      "name": "El Yopal Airport",
      "city": "El Yopal"
    },
    {
      "id": 93,
      "code": "YMS",
      "name": "Moises Benzaquen Rengifo Airport",
      "city": "Yurimaguas"
    },
    {
      "id": 94,
      "code": "IGS",
      "name": "Ingolstadt Manching Airport",
      "city": "Manching"
    },
    {
      "id": 95,
      "code": "YAY",
      "name": "St. Anthony Airport",
      "city": "St. Anthony"
    },
    {
      "id": 96,
      "code": "YWY",
      "name": "Wrigley Airport",
      "city": "Wrigley"
    },
    {
      "id": 97,
      "code": "ERI",
      "name": "Erie International Tom Ridge Field",
      "city": "Erie"
    },
    {
      "id": 98,
      "code": "PHK",
      "name": "Palm Beach Co Glades Airport",
      "city": "Pahokee"
    },
    {
      "id": 99,
      "code": "YEY",
      "name": "Amos Magny Airport",
      "city": "Amos"
    },
    {
      "id": 100,
      "code": "DCA",
      "name": "Ronald Reagan Washington National Airport",
      "city": "Washington"
    }
  ],
  "meta": {
    "limit": 100,
    "page": 1,
    "total": 100
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Flights

## GET Search Flights

GET /flights/search

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|origin|query|string| no |none|
|destination|query|string| no |none|
|date|query|array[string]| no |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "string",
  "data": [
    {
      "id": 0,
      "airline_id": 0,
      "airline": {
        "id": 0,
        "name": "string",
        "code": "string",
        "logo_url": "string",
        "created_at": "string",
        "updated_at": "string"
      },
      "origin_id": 0,
      "origin": {
        "id": 0,
        "code": "string",
        "name": "string",
        "city": "string",
        "created_at": "string",
        "updated_at": "string"
      },
      "destination_id": 0,
      "destination": {
        "id": 0,
        "code": "string",
        "name": "string",
        "city": "string",
        "created_at": "string",
        "updated_at": "string"
      },
      "departure_time": "string",
      "arrival_time": "string",
      "flight_number": "string",
      "classes": [
        {
          "id": 0,
          "flight_id": 0,
          "class_type": "string",
          "price": 0,
          "seats": [
            {}
          ],
          "created_at": "string",
          "updated_at": "string"
        }
      ],
      "created_at": "string",
      "updated_at": "string"
    }
  ],
  "meta": {
    "limit": 0,
    "page": 0,
    "total": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|
|» message|string|true|none||none|
|» data|[object]|true|none||none|
|»» id|integer|false|none||none|
|»» airline_id|integer|false|none||none|
|»» airline|object|false|none||none|
|»»» id|integer|true|none||none|
|»»» name|string|true|none||none|
|»»» code|string|true|none||none|
|»»» logo_url|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|»» origin_id|integer|false|none||none|
|»» origin|object|false|none||none|
|»»» id|integer|true|none||none|
|»»» code|string|true|none||none|
|»»» name|string|true|none||none|
|»»» city|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|»» destination_id|integer|false|none||none|
|»» destination|object|false|none||none|
|»»» id|integer|true|none||none|
|»»» code|string|true|none||none|
|»»» name|string|true|none||none|
|»»» city|string|true|none||none|
|»»» created_at|string|true|none||none|
|»»» updated_at|string|true|none||none|
|»» departure_time|string|false|none||none|
|»» arrival_time|string|false|none||none|
|»» flight_number|string|false|none||none|
|»» classes|[object]|false|none||none|
|»»» id|integer|false|none||none|
|»»» flight_id|integer|false|none||none|
|»»» class_type|string|false|none||none|
|»»» price|integer|false|none||none|
|»»» seats|[object]|false|none||none|
|»»»» id|integer|true|none||none|
|»»»» flight_class_id|integer|true|none||none|
|»»»» seat_number|string|true|none||none|
|»»»» is_available|boolean|true|none||none|
|»»»» created_at|string|true|none||none|
|»»»» updated_at|string|true|none||none|
|»»» created_at|string|false|none||none|
|»»» updated_at|string|false|none||none|
|»» created_at|string|false|none||none|
|»» updated_at|string|false|none||none|
|» meta|object|true|none||none|
|»» limit|integer|true|none||none|
|»» page|integer|true|none||none|
|»» total|integer|true|none||none|

## POST Create Flights

POST /admin/flights

> Body Parameters

```json
{
  "airline_id": 1,
  "origin_id": 1,
  "destination_id": 2,
  "departure_time": "2026-03-10T08:00:00Z",
  "arrival_time": "2026-03-10T10:00:00Z",
  "flight_number": "GA001",
  "total_seats": 25,
  "total_rows": 5,
  "total_columns": 5,
  "class_count": 3,
  "class_prices": [
    {
      "class_type": "First",
      "price": 5000000
    },
    {
      "class_type": "Business",
      "price": 3000000
    },
    {
      "class_type": "Economy",
      "price": 1000000
    }
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PUT Update Flights

PUT /admin/flights/{id}

> Body Parameters

```json
{
  "airline_id": 1,
  "origin_id": 1,
  "destination_id": 2,
  "departure_time": "2026-03-10T08:00:00Z",
  "arrival_time": "2026-03-10T10:00:00Z",
  "flight_number": "GA001",
  "total_seats": 25,
  "total_rows": 5,
  "total_columns": 5,
  "class_count": 3,
  "class_prices": [
    {
      "class_type": "First",
      "price": 5000000
    },
    {
      "class_type": "Business",
      "price": 3000000
    },
    {
      "class_type": "Economy",
      "price": 1000000
    }
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## DELETE Delete Flights

DELETE /admin/flights/{id}

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "flight deleted"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get All Flights

GET /flights

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |none|
|limit|query|string| no |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "flights retrieved",
  "data": [
    {
      "id": 2,
      "airline_id": 1,
      "airline": {
        "id": 1,
        "name": "Garuda Indonesia",
        "code": "GA",
        "logo_url": "",
        "created_at": "2026-03-09T11:55:30.505461+07:00",
        "updated_at": "2026-03-09T11:55:30.505461+07:00"
      },
      "origin_id": 1,
      "origin": {
        "id": 1,
        "code": "CGK",
        "name": "Soekarno-Hatta",
        "city": "Jakarta",
        "created_at": "2026-03-09T11:55:30.515251+07:00",
        "updated_at": "2026-03-09T11:55:30.515251+07:00"
      },
      "destination_id": 2,
      "destination": {
        "id": 2,
        "code": "DPS",
        "name": "Ngurah Rai",
        "city": "Bali",
        "created_at": "2026-03-09T11:55:30.51653+07:00",
        "updated_at": "2026-03-09T11:55:30.51653+07:00"
      },
      "departure_time": "2026-03-10T11:55:30.539724+07:00",
      "arrival_time": "2026-03-10T13:55:30.539724+07:00",
      "flight_number": "GA400",
      "classes": [
        {
          "id": 1,
          "flight_id": 2,
          "class_type": "Economy",
          "price": 1500000,
          "seats": null,
          "created_at": "2026-03-09T11:55:30.547611+07:00",
          "updated_at": "2026-03-09T11:55:30.550244+07:00"
        }
      ],
      "created_at": "2026-03-09T11:55:30.530285+07:00",
      "updated_at": "2026-03-09T11:55:30.530285+07:00"
    }
  ],
  "meta": {
    "limit": 10,
    "page": 1,
    "total": 1
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get Flight By ID

GET /flights/{id}

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "flight found",
  "data": {
    "id": 2,
    "airline_id": 1,
    "airline": {
      "id": 1,
      "name": "Garuda Indonesia",
      "code": "GA",
      "logo_url": "",
      "created_at": "2026-03-09T11:55:30.505461+07:00",
      "updated_at": "2026-03-09T11:55:30.505461+07:00"
    },
    "origin_id": 1,
    "origin": {
      "id": 1,
      "code": "CGK",
      "name": "Soekarno-Hatta",
      "city": "Jakarta",
      "created_at": "2026-03-09T11:55:30.515251+07:00",
      "updated_at": "2026-03-09T11:55:30.515251+07:00"
    },
    "destination_id": 2,
    "destination": {
      "id": 2,
      "code": "DPS",
      "name": "Ngurah Rai",
      "city": "Bali",
      "created_at": "2026-03-09T11:55:30.51653+07:00",
      "updated_at": "2026-03-09T11:55:30.51653+07:00"
    },
    "departure_time": "2026-03-10T11:55:30.539724+07:00",
    "arrival_time": "2026-03-10T13:55:30.539724+07:00",
    "flight_number": "GA400",
    "classes": [
      {
        "id": 1,
        "flight_id": 2,
        "class_type": "Economy",
        "price": 1500000,
        "seats": [
          {
            "id": 1,
            "flight_class_id": 1,
            "seat_number": "A0",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.558465+07:00",
            "updated_at": "2026-03-09T11:55:30.558465+07:00"
          },
          {
            "id": 3,
            "flight_class_id": 1,
            "seat_number": "A2",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.563301+07:00",
            "updated_at": "2026-03-09T11:55:30.563301+07:00"
          },
          {
            "id": 4,
            "flight_class_id": 1,
            "seat_number": "A3",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.56436+07:00",
            "updated_at": "2026-03-09T11:55:30.56436+07:00"
          },
          {
            "id": 5,
            "flight_class_id": 1,
            "seat_number": "A4",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.565394+07:00",
            "updated_at": "2026-03-09T11:55:30.565394+07:00"
          },
          {
            "id": 6,
            "flight_class_id": 1,
            "seat_number": "A5",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.566439+07:00",
            "updated_at": "2026-03-09T11:55:30.566439+07:00"
          },
          {
            "id": 7,
            "flight_class_id": 1,
            "seat_number": "A6",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.567521+07:00",
            "updated_at": "2026-03-09T11:55:30.567521+07:00"
          },
          {
            "id": 8,
            "flight_class_id": 1,
            "seat_number": "A7",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.568601+07:00",
            "updated_at": "2026-03-09T11:55:30.568601+07:00"
          },
          {
            "id": 9,
            "flight_class_id": 1,
            "seat_number": "A8",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.569719+07:00",
            "updated_at": "2026-03-09T11:55:30.569719+07:00"
          },
          {
            "id": 10,
            "flight_class_id": 1,
            "seat_number": "A9",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.570789+07:00",
            "updated_at": "2026-03-09T11:55:30.570789+07:00"
          },
          {
            "id": 11,
            "flight_class_id": 1,
            "seat_number": "B0",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.571318+07:00",
            "updated_at": "2026-03-09T11:55:30.571318+07:00"
          },
          {
            "id": 12,
            "flight_class_id": 1,
            "seat_number": "B1",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.572377+07:00",
            "updated_at": "2026-03-09T11:55:30.572377+07:00"
          },
          {
            "id": 13,
            "flight_class_id": 1,
            "seat_number": "B2",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.572921+07:00",
            "updated_at": "2026-03-09T11:55:30.572921+07:00"
          },
          {
            "id": 14,
            "flight_class_id": 1,
            "seat_number": "B3",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.573446+07:00",
            "updated_at": "2026-03-09T11:55:30.573446+07:00"
          },
          {
            "id": 15,
            "flight_class_id": 1,
            "seat_number": "B4",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.57397+07:00",
            "updated_at": "2026-03-09T11:55:30.57397+07:00"
          },
          {
            "id": 16,
            "flight_class_id": 1,
            "seat_number": "B5",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.576049+07:00",
            "updated_at": "2026-03-09T11:55:30.576049+07:00"
          },
          {
            "id": 17,
            "flight_class_id": 1,
            "seat_number": "B6",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.576575+07:00",
            "updated_at": "2026-03-09T11:55:30.576575+07:00"
          },
          {
            "id": 18,
            "flight_class_id": 1,
            "seat_number": "B7",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.577639+07:00",
            "updated_at": "2026-03-09T11:55:30.577639+07:00"
          },
          {
            "id": 19,
            "flight_class_id": 1,
            "seat_number": "B8",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.578691+07:00",
            "updated_at": "2026-03-09T11:55:30.578691+07:00"
          },
          {
            "id": 20,
            "flight_class_id": 1,
            "seat_number": "B9",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.579221+07:00",
            "updated_at": "2026-03-09T11:55:30.579221+07:00"
          },
          {
            "id": 21,
            "flight_class_id": 1,
            "seat_number": "C0",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.58028+07:00",
            "updated_at": "2026-03-09T11:55:30.58028+07:00"
          },
          {
            "id": 22,
            "flight_class_id": 1,
            "seat_number": "C1",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.580792+07:00",
            "updated_at": "2026-03-09T11:55:30.580792+07:00"
          },
          {
            "id": 23,
            "flight_class_id": 1,
            "seat_number": "C2",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.58184+07:00",
            "updated_at": "2026-03-09T11:55:30.58184+07:00"
          },
          {
            "id": 24,
            "flight_class_id": 1,
            "seat_number": "C3",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.582373+07:00",
            "updated_at": "2026-03-09T11:55:30.582373+07:00"
          },
          {
            "id": 25,
            "flight_class_id": 1,
            "seat_number": "C4",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.583425+07:00",
            "updated_at": "2026-03-09T11:55:30.583425+07:00"
          },
          {
            "id": 26,
            "flight_class_id": 1,
            "seat_number": "C5",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.583948+07:00",
            "updated_at": "2026-03-09T11:55:30.583948+07:00"
          },
          {
            "id": 27,
            "flight_class_id": 1,
            "seat_number": "C6",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.585043+07:00",
            "updated_at": "2026-03-09T11:55:30.585043+07:00"
          },
          {
            "id": 28,
            "flight_class_id": 1,
            "seat_number": "C7",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.58689+07:00",
            "updated_at": "2026-03-09T11:55:30.58689+07:00"
          },
          {
            "id": 29,
            "flight_class_id": 1,
            "seat_number": "C8",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.588545+07:00",
            "updated_at": "2026-03-09T11:55:30.588545+07:00"
          },
          {
            "id": 30,
            "flight_class_id": 1,
            "seat_number": "C9",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.589596+07:00",
            "updated_at": "2026-03-09T11:55:30.589596+07:00"
          },
          {
            "id": 31,
            "flight_class_id": 1,
            "seat_number": "D0",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.590138+07:00",
            "updated_at": "2026-03-09T11:55:30.590138+07:00"
          },
          {
            "id": 32,
            "flight_class_id": 1,
            "seat_number": "D1",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.591358+07:00",
            "updated_at": "2026-03-09T11:55:30.591358+07:00"
          },
          {
            "id": 34,
            "flight_class_id": 1,
            "seat_number": "D3",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.593504+07:00",
            "updated_at": "2026-03-09T11:55:30.593504+07:00"
          },
          {
            "id": 35,
            "flight_class_id": 1,
            "seat_number": "D4",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.59514+07:00",
            "updated_at": "2026-03-09T11:55:30.59514+07:00"
          },
          {
            "id": 36,
            "flight_class_id": 1,
            "seat_number": "D5",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.596308+07:00",
            "updated_at": "2026-03-09T11:55:30.596308+07:00"
          },
          {
            "id": 37,
            "flight_class_id": 1,
            "seat_number": "D6",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.597355+07:00",
            "updated_at": "2026-03-09T11:55:30.597355+07:00"
          },
          {
            "id": 38,
            "flight_class_id": 1,
            "seat_number": "D7",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.598983+07:00",
            "updated_at": "2026-03-09T11:55:30.598983+07:00"
          },
          {
            "id": 39,
            "flight_class_id": 1,
            "seat_number": "D8",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.602183+07:00",
            "updated_at": "2026-03-09T11:55:30.602183+07:00"
          },
          {
            "id": 40,
            "flight_class_id": 1,
            "seat_number": "D9",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.604924+07:00",
            "updated_at": "2026-03-09T11:55:30.604924+07:00"
          },
          {
            "id": 41,
            "flight_class_id": 1,
            "seat_number": "E0",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.606824+07:00",
            "updated_at": "2026-03-09T11:55:30.606824+07:00"
          },
          {
            "id": 42,
            "flight_class_id": 1,
            "seat_number": "E1",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.609208+07:00",
            "updated_at": "2026-03-09T11:55:30.609208+07:00"
          },
          {
            "id": 43,
            "flight_class_id": 1,
            "seat_number": "E2",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.610301+07:00",
            "updated_at": "2026-03-09T11:55:30.610301+07:00"
          },
          {
            "id": 44,
            "flight_class_id": 1,
            "seat_number": "E3",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.611898+07:00",
            "updated_at": "2026-03-09T11:55:30.611898+07:00"
          },
          {
            "id": 45,
            "flight_class_id": 1,
            "seat_number": "E4",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.613056+07:00",
            "updated_at": "2026-03-09T11:55:30.613056+07:00"
          },
          {
            "id": 46,
            "flight_class_id": 1,
            "seat_number": "E5",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.6141+07:00",
            "updated_at": "2026-03-09T11:55:30.6141+07:00"
          },
          {
            "id": 47,
            "flight_class_id": 1,
            "seat_number": "E6",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.614631+07:00",
            "updated_at": "2026-03-09T11:55:30.614631+07:00"
          },
          {
            "id": 48,
            "flight_class_id": 1,
            "seat_number": "E7",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.615157+07:00",
            "updated_at": "2026-03-09T11:55:30.615157+07:00"
          },
          {
            "id": 49,
            "flight_class_id": 1,
            "seat_number": "E8",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.616253+07:00",
            "updated_at": "2026-03-09T11:55:30.616253+07:00"
          },
          {
            "id": 50,
            "flight_class_id": 1,
            "seat_number": "E9",
            "is_available": true,
            "created_at": "2026-03-09T11:55:30.616784+07:00",
            "updated_at": "2026-03-09T11:55:30.616784+07:00"
          },
          {
            "id": 2,
            "flight_class_id": 1,
            "seat_number": "A1",
            "is_available": false,
            "created_at": "2026-03-09T11:55:30.5622+07:00",
            "updated_at": "2026-03-09T11:56:02.511988+07:00"
          },
          {
            "id": 33,
            "flight_class_id": 1,
            "seat_number": "D2",
            "is_available": false,
            "created_at": "2026-03-09T11:55:30.592431+07:00",
            "updated_at": "2026-03-09T11:56:02.513607+07:00"
          }
        ],
        "created_at": "2026-03-09T11:55:30.547611+07:00",
        "updated_at": "2026-03-09T11:55:30.550244+07:00"
      }
    ],
    "created_at": "2026-03-09T11:55:30.530285+07:00",
    "updated_at": "2026-03-09T11:55:30.530285+07:00"
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Transaction

## POST Create Transaction

POST /transactions

> Body Parameters

```json
{
  "flight_id": 1,
  "class_id": 1,
  "promo_code": "",
  "passengers": [
    {
      "full_name": "Budi Santoso",
      "nationality": "Indonesia",
      "passport_no": "A1234567",
      "seat_number": "1A"
    },
    {
      "full_name": "Siti Aminah",
      "nationality": "Indonesia",
      "passport_no": "B7654321",
      "seat_number": "1C"
    }
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "transaction created",
  "data": {
    "id": 1,
    "code": "f1ec0ac0-c492-4dc7-a7b3-3fb8ea996ffd",
    "user_id": 1,
    "flight_id": 2,
    "flight": {
      "id": 0,
      "airline_id": 0,
      "airline": {
        "id": 0,
        "name": "",
        "code": "",
        "logo_url": "",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
      },
      "origin_id": 0,
      "origin": {
        "id": 0,
        "code": "",
        "name": "",
        "city": "",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
      },
      "destination_id": 0,
      "destination": {
        "id": 0,
        "code": "",
        "name": "",
        "city": "",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
      },
      "departure_time": "0001-01-01T00:00:00Z",
      "arrival_time": "0001-01-01T00:00:00Z",
      "flight_number": "",
      "classes": null,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z"
    },
    "total_price": 3000000,
    "payment_status": "PENDING",
    "promo_code": "LIBURAN2026",
    "discount": 0,
    "expires_at": "2026-03-06T20:02:47.8740775+07:00",
    "created_at": "2026-03-06T19:52:47.8740775+07:00",
    "updated_at": "2026-03-06T19:52:47.8740775+07:00",
    "passengers": null
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get User Transactions

GET /transactions

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "transaction retrieved",
  "data": [
    {
      "code": "ceb0ff14-517e-4dd6-9908-458193af089e",
      "flight_number": "FL014",
      "class_name": "",
      "total_price": 7100000,
      "payment_status": "PENDING",
      "payment_url": "https://app.sandbox.midtrans.com/snap/v4/redirection/77a2518e-5fe6-4c73-a7d5-d0eb9bd48a24",
      "promo_code": "",
      "discount": 0,
      "expires_at": "2026-03-11T11:44:20.955013+07:00",
      "created_at": "2026-03-11T11:34:20.955539+07:00",
      "passengers": [
        {
          "full_name": "werqr",
          "nationality": "rerwe",
          "passport_no": "2313123",
          "seat_number": "1A"
        }
      ]
    },
    {
      "code": "75a84e34-9d41-45b6-a569-7b3bd3346aec",
      "flight_number": "FL003",
      "class_name": "",
      "total_price": 16350000,
      "payment_status": "PENDING",
      "payment_url": "https://app.sandbox.midtrans.com/snap/v4/redirection/87a6c66d-2cd9-4fe8-a9d1-6251fcbca1f0",
      "promo_code": "",
      "discount": 0,
      "expires_at": "2026-03-11T11:33:19.704241+07:00",
      "created_at": "2026-03-11T11:23:19.704241+07:00",
      "passengers": [
        {
          "full_name": "adasdsad",
          "nationality": "Indonesia",
          "passport_no": "sdsadads",
          "seat_number": "1B"
        },
        {
          "full_name": "dssa",
          "nationality": "Indonesia",
          "passport_no": "dasdasd",
          "seat_number": "1C"
        },
        {
          "full_name": "fafasdfa",
          "nationality": "Indonesia",
          "passport_no": "fadfasf",
          "seat_number": "1D"
        }
      ]
    },
    {
      "code": "f159a568-13e8-49af-94e5-faeece2243d4",
      "flight_number": "FL003",
      "class_name": "",
      "total_price": 5450000,
      "payment_status": "PAID",
      "payment_url": "https://app.sandbox.midtrans.com/snap/v4/redirection/83bd75fa-37dd-4608-8296-dc8a98352df8",
      "promo_code": "",
      "discount": 0,
      "expires_at": "2026-03-11T11:29:27.865768+07:00",
      "created_at": "2026-03-11T11:19:27.865768+07:00",
      "passengers": [
        {
          "full_name": "Udin",
          "nationality": "Indonesia",
          "passport_no": "3213123123123",
          "seat_number": "1A"
        }
      ]
    }
  ],
  "meta": {
    "limit": 10,
    "page": 1,
    "total": 3
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get Transaction By Code

GET /transactions/{code}

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|code|path|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "transaction found",
  "data": {
    "id": 12,
    "code": "7ad0665c-305d-4cbd-aa68-b54754b319a1",
    "user_id": 1,
    "flight_id": 2,
    "flight": {
      "id": 2,
      "airline_id": 1,
      "airline": {
        "id": 1,
        "name": "Garuda Indonesia",
        "code": "GA",
        "logo_url": "",
        "created_at": "2026-03-06T19:08:18.312323+07:00",
        "updated_at": "2026-03-06T19:08:18.312323+07:00"
      },
      "origin_id": 1,
      "origin": {
        "id": 1,
        "code": "CGK",
        "name": "Soekarno-Hatta",
        "city": "Jakarta",
        "created_at": "2026-03-06T19:08:18.320542+07:00",
        "updated_at": "2026-03-06T19:08:18.320542+07:00"
      },
      "destination_id": 2,
      "destination": {
        "id": 2,
        "code": "DPS",
        "name": "Ngurah Rai",
        "city": "Bali",
        "created_at": "2026-03-06T19:08:18.326029+07:00",
        "updated_at": "2026-03-06T19:08:18.326029+07:00"
      },
      "departure_time": "2026-03-07T19:08:18.339621+07:00",
      "arrival_time": "2026-03-07T21:08:18.339621+07:00",
      "flight_number": "GA400",
      "classes": null,
      "created_at": "2026-03-06T19:08:18.332602+07:00",
      "updated_at": "2026-03-06T19:08:18.332602+07:00"
    },
    "total_price": 3000000,
    "payment_status": "PENDING",
    "promo_code": "LIBURAN2026",
    "discount": 0,
    "expires_at": "2026-03-06T20:34:53.435013+07:00",
    "created_at": "2026-03-06T20:24:53.435013+07:00",
    "updated_at": "2026-03-06T20:24:53.435013+07:00",
    "passengers": [
      {
        "id": 23,
        "transaction_id": 12,
        "full_name": "Budi Santoso",
        "nationality": "Indonesia",
        "passport_no": "A1234567",
        "seat_number": "A1",
        "flight_class_id": 1,
        "created_at": "2026-03-06T20:24:53.435637+07:00",
        "updated_at": "2026-03-06T20:24:53.435637+07:00"
      },
      {
        "id": 24,
        "transaction_id": 12,
        "full_name": "Siti Aminah",
        "nationality": "Indonesia",
        "passport_no": "B7654321",
        "seat_number": "D2",
        "flight_class_id": 1,
        "created_at": "2026-03-06T20:24:53.435637+07:00",
        "updated_at": "2026-03-06T20:24:53.435637+07:00"
      }
    ]
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PATCH Pay Transaction

PATCH /transactions/{uid}/pay

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|uid|path|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": false,
  "message": "already paid"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Admin

## PATCH Update Role

PATCH /admin/users/{id}/role

> Body Parameters

```json
{
  "role": "user"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "users retrieved",
  "data": [
    {
      "id": 1,
      "name": "Super Admin",
      "email": "superadmin@flightbooking.com",
      "role": "super_admin",
      "created_at": "2026-03-11T14:06:11.608406+07:00"
    },
    {
      "id": 2,
      "name": "Admin User",
      "email": "admin@flightbooking.com",
      "role": "admin",
      "created_at": "2026-03-11T14:06:11.674298+07:00"
    },
    {
      "id": 3,
      "name": "John Doe",
      "email": "user@flightbooking.com",
      "role": "user",
      "created_at": "2026-03-11T14:06:11.737357+07:00"
    }
  ],
  "meta": {
    "limit": 10,
    "page": 1,
    "total": 3
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Get All User

GET /admin/users

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |none|
|limit|query|string| no |none|
|sort_by|query|string| no |none|
|order|query|string| no |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "users retrieved",
  "data": [
    {
      "id": 1,
      "name": "Super Admin",
      "email": "superadmin@flightbooking.com",
      "role": "super_admin",
      "created_at": "2026-03-11T14:06:11.608406+07:00"
    },
    {
      "id": 2,
      "name": "Admin User",
      "email": "admin@flightbooking.com",
      "role": "admin",
      "created_at": "2026-03-11T14:06:11.674298+07:00"
    },
    {
      "id": 3,
      "name": "John Doe",
      "email": "user@flightbooking.com",
      "role": "user",
      "created_at": "2026-03-11T14:06:11.737357+07:00"
    }
  ],
  "meta": {
    "limit": 10,
    "page": 1,
    "total": 3
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## GET Dashboard

GET /admin/dashboard

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "dashboard stats retrieved",
  "data": {
    "total_users": 3,
    "total_flights": 100,
    "total_transactions": 3,
    "total_revenue": 9300000,
    "pending_payments": 1,
    "completed_bookings": 1
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PATCH Ban User

PATCH /admin/users/{id}/ban

> Body Parameters

```json
{
  "reason": "Ambatutu"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|
|body|body|object| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "user banned"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PATCH Unban User

PATCH /admin/users/{id}/unban

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "success": true,
  "message": "user unbanned"
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## DELETE Delete User

DELETE /admin/users/{id}

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

## PATCH Restore User

PATCH /admin/users/{id}/restore

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|id|path|string| yes |none|

> Response Examples

> 200 Response

```json
{}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

# Data Schema

