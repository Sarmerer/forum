# forum

## How to run the app?

In order to run the application you need to create a ``.env`` file, in the root directory of the app. The file should have these fields:
```env
API_KEY=<your api key>
```

### API response
Content-type of an API response is always application/json.
For the request:
```
https://localhost:4433/users/1
```
response will have the following structure:
```json
{
  "status": "success", //response status. Can also be "error"
  "code": 200,
  "message": null, //contains error, message if response status is "error"
  "data": {
      "ID": 1,
      "Name":"John Doe",
      "Email":"johndoe@gmail.com",
      "Password":"CAF39A5B635B838EEF70BE1E280875FC",
      "SessionID":"a37fd706-df07-4cd2-94cc-e8516fa476df"
  }
}
```
