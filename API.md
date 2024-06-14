* [Query API LIST](#query-api-list)
    * [Account List](#account-list)
    * [Record List](#record-list)
* [Operate API LIST](#operate-api-list)
    * [Edit Owner](#edit-owner)
    * [Edit Record](#edit-record)
    * [Recycle](#recycle)
    * [Tx Send](#tx-send)
* [Error Code List](#error-code-list)

## Account List

Query did cell account list by a ckb address

**Request Syntax**

```
POST /v1/account/list HTTP/1.1
Content-type: application/json
```

```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "309",
    "key": "ckbxxx..."
  },
  "page": 1,
  "size": 20,
  "keyword": "xxxxx.bit",
  "did_type": 1
}
```

**Request Body**

The request accepts the following data in JSON format.

* key_info: owner of account; Type: string; Required: Yes.
* did_type: did cell type, 0(default value, search all did cell), 1(search normal did cell), 2(search expired did cell);
  Type: Integer; Required: No.

**Response Syntax**

```
HTTP/1.1 201
Content-type: application/json
```

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "total": 0,
    "did_list": [
      {
        "outpoint": "",
        "account_id": "",
        "account": "",
        "expired_at": 111,
        "did_cell_status": 1
      }
    ]
  }
}
```

**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.

* outpoint: did cell outpoint; Type: String
* account_id: did cell account_id; Type: String
* account: did cell account; Type: String
* expired_at: did cell expired_at; Type: Integer
* did_cell_status: did cell status; Type: Integer

## Record List

Query record list of a did cell account

**Request Syntax**

```
POST /v1/record/list HTTP/1.1
Content-type: application/json
```

```json
{
  "account": "xxxxx.bit"
}
```

**Request Body**

The request accepts the following data in JSON format.

* account: dotbit account; Type: string; Required: Yes.

**Response Syntax**

```
HTTP/1.1 201
Content-type: application/json
```

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "records": [
      {
        "key": "",
        "type": "",
        "label": "",
        "value": "",
        "ttl": ""
      }
    ]
  }
}
```

**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.

* key: record key; Type: String
* type: record type; Type: String
* label: record label; Type: String
* value: record value; Type: String
* ttl: record ttl; Type: String

## Edit Owner

Transfer a did cell account to other ckb address

**Request Syntax**

```
POST /v1/transfer HTTP/1.1
Content-type: application/json
```

```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "309",
    "key": "ckbxxx..."
  },
  "account": "aaaaa.bit",
  "receive_ckb_addr": ""
}
```

**Request Body**

The request accepts the following data in JSON format.

* key_info: owner of account; Type: string; Required: Yes.
* account: dotbit account; Type: string; Required: Yes.
* receive_ckb_addr: ckb address of the receiver; Type: string; Required: Yes.

**Response Syntax**

```
HTTP/1.1 201
Content-type: application/json
```

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "sign_key": "",
    "sign_list": [
      {
        "sign_type": 5,
        "sign_msg": ""
      }
    ],
    "ckb_tx": ""
  }
}
```

**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.

* sign_key: tx key; Type: String
* sign_list: sign msg list; Type: String
* ckb_tx: tx of transfer account; Type: String

## Edit Record

Edit did record
**Request Syntax**

```
POST /v1/edit/record HTTP/1.1
Content-type: application/json
```

```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "309",
    "key": "ckbxxx..."
  },
  "account": "aaaaa.bit",
  "raw_param": {
    "records": [
      {
        "type": "profile",
        "key": "twitter",
        "label": "",
        "value": "111",
        "ttl": "300",
        "action": "add"
      }
    ]
  }
}
```

**Request Body**

The request accepts the following data in JSON format.

* key_info: owner of account; Type: string; Required: Yes.
* account: dotbit account; Type: string; Required: Yes.
* raw_param: record list; Type: string; Required: Yes.

**Response Syntax**

```
HTTP/1.1 201
Content-type: application/json
```

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "sign_key": "",
    "sign_list": [
      {
        "sign_type": 5,
        "sign_msg": ""
      }
    ],
    "ckb_tx": ""
  }
}
```

**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.

* sign_key: tx key; Type: String
* sign_list: sign msg list; Type: String
* ckb_tx: tx of edit record; Type: String

## Recycle

Recycle a did cell

**Request Syntax**

```
POST /v1/recycle HTTP/1.1
Content-type: application/json
```

```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "309",
    "key": "ckbxxx..."
  },
  "account": "aaaaa.bit",
}
```

**Request Body**

The request accepts the following data in JSON format.

* key_info: owner of account; Type: string; Required: Yes.
* account: dotbit account; Type: string; Required: Yes.

**Response Syntax**

```
HTTP/1.1 201
Content-type: application/json
```

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "sign_key": "",
    "sign_list": [
      {
        "sign_type": 99,
        "sign_msg": ""
      }
    ],
    "ckb_tx": ""
  }
}
```

**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.

* sign_key: tx key; Type: String
* sign_list: sign msg list; Type: String
* ckb_tx: tx of recycle account; Type: String

## Tx Send

Send a Transaction

**Request Syntax**

```
POST /v1/tx/send HTTP/1.1
Content-type: application/json
```

```json
{
  "sign_key": "",
  "sign_list": [
    {
      "sign_type": 99,
      "sign_msg": ""
    }
  ],
  "ckb_tx": ""
}
```

**Request Body**

The request accepts the following data in JSON format.

* sign_key: tx key; Type: String
* sign_list: sign list; Type: String
* ckb_tx: signed transactions; Type: String

**Response Syntax**

```
HTTP/1.1 201
Content-type: application/json
```

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "hash": ""
  }
}
```

**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.

* hash: tx hash; Type: String

## Error Code List

* 10000: request body error
* 10002: system db error
