# MicroUser
Small microservice to manage Users, written in Go and using mongoDb as database.
The all application as been dockerized.

## Installation

The installation required to have `docker` and `docker-compose` installed.  
You can find the instruction here:  
   - `docker` : https://docs.docker.com/get-docker/
   - `docker-compose` : https://docs.docker.com/compose/install/

Once they are installed you can start the app with:
    `docker-compose up -d`

The API is listening on port `8080` on `localhost` by default.

### Check
You can run a check the server is running with:  
    `curl http://localhost:8080/ping` (or just reached http://localhost:8080/ping in a browser)  
The response should be a `pong`

You can monitor the logs with:  
    `docker-compose logs -f`  
Or for only the API:  
    `docker-compose logs -f go-api`

A webui directly connected to the db can be access at:
    `http://localhost:8081`

### Run the tests
You can run the tests with:  
`docker exec -it go-api go test -v ./...`

### Tests coverage and remarques

- Tests use another db in the same mongoDb container as a mock database. It could be better if donne with a fully mocked db saved locally.  
- It could be possible to write tests for usersResource, they will be really close to duplicate of those for usersStore.  
- A personal time constraints made me focus on unit test at the usersStore level.  
- Another point is that there could be more integration tests.

## API

### Add a new User

- Add a new User by sending the corresponding json by a POST request to `http://localhost:8080/users`.  
`first_name`, `last_name`, `password`, `nickname`, `email`, and `country` are required.  


- All other fields will be ignored.  
`email` should be a valid email address.


- The response will be the complete user schema, with `id`, `created_at` and `updated_at`.

#### Example

```
curl -X POST http://localhost:8080/users \
-H 'Content-Type: application/json' \
-d '{"first_name":"Mike","last_name":"Tyson","password":"dad154", "nickname":"Myki mike","email":"miky@ggmail.com","country":"US"}'
```

_response:_
```
{"id":"61e41ed578752c5997718aff","first_name":"Mike","last_name":"Tyson","nickname":"Myki%20mike","password":"dad154","email":"miky@ggmail.com","country":"US","created_at":"2022-01-16T13:34:13.684Z","updated_at":"2022-01-16T13:34:13.684Z"}
```


### Update an existing User

- Update an existing User by sending the new user as json by a PUT request to `http://localhost:8080/users/{userId}`.  
`first_name`, `last_name`, `password`, `nickname`, `email`, and `country` are required.  


- All other fields will be ignored.  
`email` should be a valid email address.


- The response will be the updated user schema, with `id`, `created_at` and `updated_at`.

_**Warning**_: All fields must be present, it's currently closer to a replaceOne than a updateOne. But it keeps created_at and the id unmodified.

#### Example

```
curl -X PUT http://localhost:8080/users/61e41ed578752c5997718aff \
-H 'Content-Type: application/json' \
-d '{"first_name":"Mike","last_name":"Longbow","password":"dad154", "nickname":"Myki mike","email":"miky@ggmail.com","country":"US"}'
```

_response:_
```
{"id":"61e41ed578752c5997718aff","first_name":"Mike","last_name":"Longbow","nickname":"Myki%20mike","password":"dad154","email":"miky@ggmail.com","country":"US","created_at":"2022-01-16T13:34:13.684Z","updated_at":"2022-01-16T13:34:13.684Z"}
```

### Remove a User

Delete an existing User by sending a DELETE request to `http://localhost:8080/users/{userId}`.
The response will be a user schema, with `id`. (And `created_at` and `updated_at`) 

#### Example
```
curl -X DELETE http://localhost:8080/users/61e41ed578752c5997718aff
```

_response:_
```
{"id":"61e41ed578752c5997718aff", "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}
```

### Search Users

Return paginated list of Users, with possibly some filtering by certain criteria, with a GET request at `http://localhost:8080/users`.

- All filters are pass by query parameters.  


- `first_name`, `last_name`, `password`, `nickname`, `email`, and `country` can be filtered by regex and can be only a substring of the value.  
_example_: `first_name=To` returns every User with `To` in their `first_name`.


- An extra parameter `text` is a regex value that can be found in any of those previously mentioned fields.  
_example_: `text=To` returns every User with `To` in their `first_name` or `last_name` or ... .


- `created_at` and `updated_at` are filtered with a start and end dates.  `startdcreated`, `enddcreated`, `startdupdated`, `enddupdated`  
_example_: `startdcreated=2022-01-15T12:30:00.00Z` returns every User with `created_at` greater or equal than `startdcreated`.  


- `id` are exact matches.


- The pagination is done with the parameters `page` for the page number (starting at 0) and `page_size` for the number of Users per page, order by `created_at`.  
Without those page parameters all filtered Users are returned.  
_example_:`page=1&page_size=5` return the second page and up to five results.  


- The result are an array of User `users`, and a `count` value (the number of result). If a no User are found the result, or if the pagination request is too far for example the result will be a `null` value as the `users`.

_**Rmq**_: In the database every string has been escaped before being saved. You maybe need to unescape the result on the front.  
Plus if you want to research a special character you need to escape it. (example: %20 for space, or %40 for @)

#### Examples

####1]
```
curl -X GET http://localhost:8080/users?page=1&page_size=2
```

_response:_
```
{
"users":[
{"id":"61e41ed578752c5997718aff","first_name":"Mike","last_name":"Longbow","nickname":"Myki%20mike","password":"dad154","email":"miky@ggmail.com","country":"US","created_at":"2022-01-16T13:34:13.684Z","updated_at":"2022-01-16T13:34:13.684Z"},
{"id":"61e6788f78987008888888ff","first_name":"Tike","last_name":"Tongbow","nickname":"Tyki%20mike","password":"Tad154","email":"tiky@ggmail.com","country":"UK","created_at":"2022-01-16T13:35:13.684Z","updated_at":"2022-01-16T13:35:13.684Z"}
],
"count":2
}
```


####2]
```
curl -X GET http://localhost:8080/users?page=0&page_size=1&startdcreated=2022-01-16T13:35:00.000Z
```

_response:_
```
{
"users":[
{"id":"61e6788f78987008888888ff","first_name":"Tike","last_name":"Tongbow","nickname":"Tyki%20mike","password":"Tad154","email":"tiky@ggmail.com","country":"UK","created_at":"2022-01-16T13:35:13.684Z","updated_at":"2022-01-16T13:35:13.684Z"}
],
"count":1
}
```


####3]
```
curl -X GET http://localhost:8080/users?page=1&page_size=3&first_name=ike&country=UK
```

_response:_
```
{
"users":[
{"id":"61e41ed578752c5997718aff","first_name":"Mike","last_name":"Longbow","nickname":"Myki%20mike","password":"dad154","email":"miky@ggmail.com","country":"UK","created_at":"2022-01-16T13:34:13.684Z","updated_at":"2022-01-16T13:34:13.684Z"},
{"id":"6abc988ee0908dd297718aff","first_name":"Rike","last_name":"Rongbow","nickname":"Ryki%20mike","password":"dad154","email":"riky@ggmail.com","country":"UK","created_at":"2022-01-16T13:36:13.684Z","updated_at":"2022-01-16T13:36:13.684Z"},
],
"count":2
}
```

## Choices and structure explanations

- The User is the central object of the api and so a struct User as been design to maintain consistency.
- Files, folders and packages are various for clarity and to enable flexibility and modifications in potential future.


- The database is designed so no duplicate nickname or email are authorized. 
- It was assumed that this service was used by admins, so the id, password and other data considered sensitives are not encrypted and are present in the search requests.  
- It was assumed the service is the principal manager of the users and so manage the id, create_at and updated_at. Those field can't be initialized or modified manually through the api.
- For security, `first_name`, `last_name`, `password`, `nickname`, `email`, and `country` are escaped before being saved in the database.

### The Design Pattern
```
.
├── api                                 -- Routing for API logic
│   ├── api.go                              -- Root API view
│   └── server.go                           -- Root server view
├── errors                              -- Routing for Authentication logic
│   └── errors.go                           -- Errors logics
├── user                                -- All user controllers
│   ├── userModel.go                        -- Defines the User schema as a struc
│   ├── userModel_test.go                   -- userModel Unit tests
│   ├── usersResource.go                    -- Defines User management handler
│   ├── usersStore.go                       -- Defines DB operations on Users
│   └── usersStore_test.go                  -- usersStore Unit tests
├── utils                               -- Define utils functions and struct usable through all the app
│   ├── utils.go                           -- Define utils functions and struct
│   └── utils_test.go                      -- Utils Unit tests
├── main.go                             -- Connect to the DB, init the Server then start it
└── main_test.go                             -- Integration tests
```

## For more

Notification of User modification to other services can be easily send by creating a function like SendNotification. 
An example is written in `usersResource.go`.

Some limitations on the characters used or the length of some fields can be added.