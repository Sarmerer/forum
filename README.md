# forum

## API response
Content-type of an API response is always application/json.
Response has the following structure:
```json
{
  "status": "success/error", //response status
  "code": 200, //http status
  "message": null, //if response status is "error", this field will contains error message
  "data": {
      "ID": 1,
      "Name":"John Doe",
      "Age": 22
  }
}
```
