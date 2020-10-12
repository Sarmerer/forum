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
```yaml
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
```yaml
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
```yaml
{
  url: /api/auth/signout,
  method: "POST"
}
```
#### Me
Returns data about user account. It is being called after a Sign In attempt, returned data being stored in Vuex store.
```yaml
{
  url: /api/auth/me,
  method: "GET",
  permissions: ["account owner"]
  response: {
    id: int,
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
```yaml
{
  url: /api/users,
  method: "GET",
  permissions: ["account owner", "moder", "admin"],
  response: [
    {
      id: int,
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
```yaml
{
  url: /api/user?:id,
  method: "GET",
  permissions: ["everyone"],
  response: {
    id: int,
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
```yaml
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
```yaml
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
```yaml
{
  url: /api/posts,
  methhod: "GET",
  response: [
    {
     id: int,
     author_id: int,
     author_name: String,
     title: String,
     content: String,
     created: Stringified Date,
     Updated: Stringified Date,
     categories: [
       {
         id: int,
         name: String,
         use_count: int // indicates how many posts use this category
       }
       ...
     ]
     comments: [
       {
         id: int,
         author_id: int,
         post_id: int, // post id, to which comment bolongs
         author_name: String,
         content: String,
         created: Stringified Date,
         edited: int // 0 or 1 Boolean indicates wheter comment was edited
       }
       ...
     ]
     comments_count: int,
     rating: int,
     your_reaction: int // -1, 0 or 1, indicates your reaction to highlite like, dislike or neither
    },
    ...
  ]
}
```
```yaml
{
  url: /api/post/find,
  methhod: "GET",
  request: {
    body: {
      by: {
        type: String,
        variants: ["id", "author", "categories"]
      },
      id: int, // if by == "id"
      author: int, //if by == "author"
      categories: Array of Strings // if by == "categories"
    }
    response: { // can be Object or Array of Object, depending on search type

    }
  }
}
```
```yaml
{
  url: /api/post/create,
  methhod: "POST",
}
```
```yaml
{
  url: /api/post/update,
  methhod: "PUT",
}
```
```yaml
{
  url: /api/post/delete,
  methhod: "DELETE",
}
```
```yaml
{
  url: /api/post/rate,
  methhod: "POST",
}
```

### Categories
 ```yaml
 {
 URI: /api/categories,
 }
 ```

### Comments
 ```yaml
 {
 URI: /api/comments,
 }
 ```
 ```yaml
 {
URI:  /api/comment/add,
 }
 ```
 ```yaml
 {
 URI: /api/comment/update,
 }
 ```
 ```yaml
 {
 URI: /api/comment/delete,
 }
 ```

 ## Responses

 #### User

 #### Post

 ```yaml
{
  id: int,
  author_id: int,
  author_name: String,
  title: String,
  content: String,
  created: Stringified Date,
  Updated: Stringified Date,
  categories: [
    {
      id: int,
      name: String,
      use_count: int // indicates how many posts use this category
    },
    ...
  ],
  comments: [
    {
      id: int,
      author_id: int,
      post_id: int, // post id, to which comment bolongs
      author_name: String,
      content: String,
      created: Stringified Date,
      edited: int // 0 or 1 Boolean indicates wheter comment was edited
    },
    ...
  ],
  comments_count: int,
  rating: int,
  your_reaction: int // -1, 0 or 1, indicates your reaction to highlite like, dislike or neither
}
 ```