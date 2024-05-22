* [Query API LIST](#query-api-list)
    * [Did List](#did-list)
    * [Record List](#record-list)
* [Operate API LIST](#operate-api-list)
  * [Transfer](#transfer)
  * [Edit Record](#edit-record)
  * [Recycle](#recycle)
  * [Tx Send](#tx-send)
* [Error Code List](#error-code-list)
  

## Did List
Query  did list of a ckb address

**Request Syntax**

```
POST /v1/did/list HTTP/1.1
Content-type: application/json
```
```json
{
  "ckb_address": "string",
  "did_type": 1
}
```
**Request Body**

The request accepts the following data in JSON format.
* ckb_address: ckb address; Type: string; Required: Yes.
* did_type: did cell type, 0(default value, search all did cell), 1(search normal did cell), 2(search expired did cell); Type: Integer; Required: No.

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
    "did_list": [
      {
        "outpoint": "",
        "account_id": "",
        "account": "",
        "args": "",
        "expired_at":  111,
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
* args: did cell args; Type: String
* expired_at: did cell expired_at; Type: Integer
* did_cell_status: did cell status; Type: Integer

## Record List
Query record list of a did

**Request Syntax**

```
POST /v1/record/list HTTP/1.1
Content-type: application/json
```
```json
{
  "account": "aaaaa.bit"
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
        "ttl":  ""
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



## Transfer
Transfer a did cell to other ckb address

**Request Syntax**

```
POST /v1/transfer HTTP/1.1
Content-type: application/json
```
```json
{
  "account": "aaaaa.bit",
  "ckb_addr": "",
  "receive_ckb_addr": ""
}
```
**Request Body**

The request accepts the following data in JSON format.
* account: dotbit account; Type: string; Required: Yes.
* ckb_addr: ckb address of the account`s owner; Type: string; Required: Yes.
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
    ]
  }
}
```
**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.
* sign_key: tx key; Type: String
* sign_list: sign msg list; Type: String



## Edit Record
Edit did record 
**Request Syntax**

```
POST /v1/edit/record HTTP/1.1
Content-type: application/json
```
```json
{
  "account": "aaaaa.bit",
  "ckb_addr": "",
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
* account: dotbit account; Type: string; Required: Yes.
* ckb_addr: ckb address of the account`s owner; Type: string; Required: Yes.
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
    ]
  }
}
```
**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.
* sign_key: tx key; Type: String
* sign_list: sign msg list; Type: String


## Recycle
Recycle a did cell

**Request Syntax**

```
POST /v1/recycle HTTP/1.1
Content-type: application/json
```
```json
{
  "account": "aaaaa.bit",
  "outpoint": ""
}
```
**Request Body**

The request accepts the following data in JSON format.
* account: dotbit account; Type: string; Required: Yes.
* outpoint: outpoint of did cell; Type: string; Required: Yes.

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
    ]
  }
}
```
**Response Elements**

If the action is successful, the service sends back an HTTP 201 response.
The following data is returned in JSON format by the service.
* sign_key: tx key; Type: String
* sign_list: sign msg list; Type: String


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
  ]
}
```
**Request Body**

The request accepts the following data in JSON format.
* sign_key: tx key; Type: String
* sign_list: sign  list; Type: String

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
