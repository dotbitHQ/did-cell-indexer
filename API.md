* [Query API LIST](#query-api-list)
    * [Did List](#did-list)
    * [Record List](#record-list)
* [Error Code List](#error-code-list)
  

## Did List

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
* did_type: did cell type, 0(default value, search all did cell), 1(search nornal did cell), 2(search expired did cell); Type: Integer; Required: No.

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




## Error Code List

* 10000: request body error
* 10002: system db error
