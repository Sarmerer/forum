# forum

* [Makefile](##Makefile)
* [API](#API)
  * [Requests](###Requests)
  * [Responses](###Responses)

## Makefile

### `make go`
run backend

### `make vue`
run frontend

### `make push m="commit comment"`
push **ALL** changes to origin master

---

# API
API receives HTTP requests with `JSON` body and no query parameters. But there are some exclusions, [axios](https://github.com/axios/axios), library useed on the frontend, doesn't allow to send request body along with `GET` requests, so when sending a `GET` request, only query parameters are used.

## Requests

### Auth
#### Sign In
```json
{
  url: /api/auth/signin,
  method: "POST",
  request: {
    body: {
    login: String,
    password: String
    }
  }
}
```
#### Sign Up
```json
{
  url: /api/auth/signup,
  method: "POST",
  request: {
    body: {
      login: String,
      email: String,
      password: String
    }
  }
}
```
#### Sign Out
```json
{
  url: /api/auth/signout,
  method: "POST"
}
```
#### Me
Returns data about user account. It is being called after a Sign In attempt, returned data being stored in Vuex store.
```json
{
  url: /api/auth/me,
  method: "GET",
  permissions: ["account owner"]
  response: {
    id: Number,
    login: String,
    password: String, //TODO clean up
    email: String,
    avatar: URL,
    display_name: String,
    created: Stringified Date,
    last_online: String,
    sessiob_id: String,
    Role: Int
  }
}
```

### User
#### Get all users
```json
{
  url: /api/users,
  method: "GET",
  permissions: ["account owner", "moder", "admin"],
  response: [
    {
      id: Number,
      login: String,
      password: String, //TODO clean up
      email: String,
      avatar: URL,
      display_name: String,
      created: Stringified Date,
      last_online: String,
      sessiob_id: String,
      Role: Int
    }
    ...
  ]
}
```
#### Get user by ID
```json
{
  url: /api/user?:id,
  method: "GET",
  permissions: ["everyone"],
  response: {
    id: Number,
    login: String,
    password: String, //TODO clean up
    email: String,
    avatar: URL,
    display_name: String,
    created: Stringified Date,
    last_online: Stringified Date,
    sessiob_id: String,
    Role: Int
  }
}
```
#### Update user
```json
{
  url: /api/user/update?:id,
  method: "PUT",
  permissions: ["account owner", "moder", "admin"],
  request: {
    params: {
      id: Number
      //TODO move id to body
    },
    body: {
      //TODO fill this
    }
  }
}
```
#### Delete user
```json
{
  url: /api/user/delete?:id,
  method: "DELETE",
  permissions: ["account owner", "moder", "admin"],
  request: {
    params: {
      id: Number,
    }
  }
}
```

### Post
```javascript
{
	"url":"/api/posts",
	"methhod":"GET",
	"response":[
		{
			"id":Number,
			"author_id":Number,
			"author_name":String,
			"title":String,
			"content":String,
			"created":Date,
			"Updated":Date,
			"categories":[
				{
					"id":Number,
					"name":String,
					"use_count":Number // indicates how many posts use this category
        },
        ...,
		    "comments":[
				{
					"id":Number,
					"author_id":Number,
					"post_id":Number,
          "author_name":String,
					"content":String,
					"created":Date,
					"edited":Number // 0 or 1 Boolean indicates wheter comment was edited
        },
        ...,
      ],
      "comments_count":Number,
			"rating":Number,
			"your_reaction":Number // -1, 0 or 1, indicates your reaction highlite like dislike or neither
		},
		...
	]
}
```
```json
{
  url: /api/post/find,
  methhod: "GET",
  request: {
    body: {
      by: {
        type: String,
        variants: ["id", "author", "categories"]
      },
      id: Number, // if by == "id"
      author: Number, //if by == "author"
      categories: Array of Strings // if by == "categories"
    }
    response: { // can be Object or Array of Object, depending on search type

    }
  }
}
```
```json
{
  url: /api/post/create,
  methhod: "POST",
}
```
```json
{
  url: /api/post/update,
  methhod: "PUT",
}
```
```json
{
  url: /api/post/delete,
  methhod: "DELETE",
}
```
```json
{
  url: /api/post/rate,
  methhod: "POST",
}
```

### Categories
 ```json
 {
 URI: /api/categories,
 }
 ```

### Comments
 ```json
 {
 URI: /api/comments,
 }
 ```
 ```json
 {
URI:  /api/comment/add,
 }
 ```
 ```json
 {
 URI: /api/comment/update,
 }
 ```
 ```json
 {
 URI: /api/comment/delete,
 }
 ```

 ## Responses

 #### User

 #### Post

 ```json
{
  id: Number,
  author_id: Number,
  author_name: String,
  title: String,
  content: String,
  created: Stringified Date,
  Updated: Stringified Date,
  categories: [
    {
      id: Number,
      name: String,
      use_count: Number // indicates how many posts use this category
    },
    ...
  ],
  comments: [
    {
      id: Number,
      author_id: Number,
      post_id: Number, // post id, to which comment bolongs
      author_name: String,
      content: String,
      created: Stringified Date,
      edited: Number // 0 or 1 Boolean indicates wheter comment was edited
    },
    ...
  ],
  comments_count: Number,
  rating: Number,
  your_reaction: Number // -1, 0 or 1, indicates your reaction to highlite like, dislike or neither
}
 ```