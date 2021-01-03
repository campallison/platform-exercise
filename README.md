# Platform Exercise
### Discussion:
This auth implementation is written in Go on an AWS SAM `sam init` template with a Postgres database. Because Fender's platform, to my understanding, operates on the serverless application model, I wanted to explore that a bit more. I wrote it in Go because it's what I've been using the past few years and a really excellent language overall.

This service implements the required endpoints and a couple more just for fun, and perhaps a bit of utility. The database contains two tables: a `users` table and an `invalid_tokens` table, along with the migrations table. 
```sql
+----------+------------------+--------+---------+
| Schema   | Name             | Type   | Owner   |
|----------+------------------+--------+---------|
| public   | goose_db_version | table  | root    |
| public   | invalid_tokens   | table  | root    |
| public   | users            | table  | root    |
+----------+------------------+--------+---------+

+------------+--------------------------+--------------------------------------+
| Column     | Type                     | Modifiers                            |
|------------+--------------------------+--------------------------------------|
| created_at | timestamp with time zone |                                      |
| updated_at | timestamp with time zone |                                      |
| deleted_at | timestamp with time zone |                                      |
| id         | uuid                     |  not null default uuid_generate_v4() |
| name       | text                     |  not null                            |
| email      | text                     |  not null                            |
| password   | text                     |  not null                            |
+------------+--------------------------+--------------------------------------+
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)

+------------+--------------------------+-------------+
| Column     | Type                     | Modifiers   |
|------------+--------------------------+-------------|
| created_at | timestamp with time zone |             |
| updated_at | timestamp with time zone |             |
| token      | text                     |  not null   |
+------------+--------------------------+-------------+
Indexes:
    "invalid_tokens_pkey" PRIMARY KEY, btree (token)
```  
Passwords are encrypted with the commonly-used bcrypt package and of course only the hash is stored. The `invalid_tokens` table holds logged-out tokens.  All endpoint that require authorization check that the given token does not exist in the `invalid_tokens` table.


**CreateUser**

`POST /user` endpoint, requires a name, email, and password, all as strings. Does not require authorization. Checks the name for "illegal" characters, though I chose to do that just for the exercise of it. It would be very difficult to prohibit some characters or structures without inadvertently excluding some users. This person shares their opinion on it here: https://www.kalzumeus.com/2010/06/17/falsehoods-programmers-believe-about-names/ . Checks the email for proper formatting and checks that the domain is not on a prohibited list. Email is used as the primary key for easy lookup and as a bonus deal it is then unique. The password is checked for strength using the zxcvbn package https://github.com/dropbox/zxcvbn and the chosen threshold is two on their scale of zero to four. Two is selected because it's pretty strong and in previous user testing it seemed that requiring the or four frustrated users.

Returns user ID, name, and email. ID is useful for `GET` requests

**GetUser**

`GET /user/{id}` endpoint, accepts the user ID in the path and requires an authorization header with a valid token. Returns user ID, name, and email.

**UpdateUser**

`PATCH /user/{id}` endpoint, accepts the user ID in the path and requires an authorization header with a valid token, as well as a JSON request body with any or all of the following: name, email, old password + new password. If a new password is provided, the old password must be present and is checked against the stored hash using bcrypt's built-in tools since its hashing algorithms are nondeterministic. New password is checked for strength.

Returns ID, name, and email for the user, with new values for whichever fields were updated.

**DeleteUser**

`DELETE /user/{id}` endpoint, accepts the user ID in the path and requires an authorization header with a valid token. Of note, the `gorm` ORM performs a soft delete and sets a time in the `deleted_at` column. GORM will normally add `deleted_at IS NOT NULL` to `WHERE` clauses, but it is wise sometimes to explicitly add that into a clause for readability if needed on your team.

Returns the ID of the deleted user, which was a somewhat arbitrary decision. I would generally ask the front end what is most useful to them.

**PasswordStrength**

`POST /password-strength` endpoint, accepts a JSON body with a potential password string and checks its strenght on the aforementioned zxcvbn scale. Kind of just for fun in this instance since I was first figuring out the SAM template. While packages exist that can check password strength on the client side, that is a potential use for this endpoint. Better not to send the password if you don't have to of course. The user value though is to give the user near-immediate feedback as they fill in fields to create an account. Seeing that feedback is preferable to hitting submit and getting an error.

Returns an integer representing the strength on the zxcvbn 0 - 4 scale.

**ValidateEmail**

`POST /validate-email` endpoint, accepts a JSON body with a potential email string and validates it with the same checks used when creating a user. Same user value as a password strength check during account creation, the user can see before submitting if the service will accept their email address.

Returns the given email, a boolean value representing validity, and a descriptive error field.


**Login**

`POST /login` endpoint, accepts an email address and password. Retrieves the user by email, checks the password, and if successful creates a JWT that expires 12 hours after creation. 

Returns the signed token, and the expiration time at the top level.

**Logout**

`POST /logout/{id}` endpoint, takes the user ID in the path as well as the access token in the authorization header. Saves the token to the `invalid_tokens` table, and uses the opportunity to delete any rows in the table created more than 12 hours ago. In a very large service, I would probably opt not to delete stale tokens during this step so the logout request could execute as quickly as possible. Perhaps in the case of a much larger service, a worker could run periodically and clear stale tokens.

Returns a boolean representing success or failure.



## How to run the service

**Prerequisites:**

Go https://golang.org/dl/

Postgres https://www.postgresql.org/download/

**optional:**

PGCLI https://www.pgcli.com/install

Python https://www.python.org/downloads/


## Fire it up

Locally, run `make run-dev` to initialize the database and build the project. As you may be aware, after the project builds you will see output including something like this:

```bash
Mounting GetUserFunction at http://127.0.0.1:1946/user/{id} [GET]
Mounting LoginFunction at http://127.0.0.1:1946/login [POST]
Mounting CheckPasswordStrengthFunction at http://127.0.0.1:1946/password-strength [POST]
Mounting CreateUserFunction at http://127.0.0.1:1946/user [POST]
Mounting DeleteUserFunction at http://127.0.0.1:1946/user/{id} [DELETE]
Mounting LogoutFunction at http://127.0.0.1:1946/logout/{id} [POST]
Mounting ValidateEmailFunction at http://127.0.0.1:1946/validate-email [POST]
Mounting UpdateUserFunction at http://127.0.0.1:1946/user/{id} [PATCH]
```

Those addresses can be used to curl the endpoints locally. Examples are given below.

### Run the tests

In the command line, `go test` will run all the unit tests in the project. Depending on the output your're used to seeing, please note that you will see feedback from the database reporting records not found. That's good, that is output from tests that specifically tested calls to the database that could not find records to assert that a proper error was returned. Further options on the `go test` command, such as `go test -v -run Test_CreateUser`, will run tests with function names matching the last argument. In the case of that example, you would be running `Test_CreateUser` as well as `Test_CreateUserHandler`.

Being familiar with Go, you can also run the tests in your IDE of choice and the `env.json` file should ensure it will run. If the tests will not run locally, it is possible the OS-agnostic host I found, `"postgres://root:postgres@host.docker.internal:5432/postgres?sslmode=disable"`, does not do what I read it would. In that case, you'll need to look up the correct host for your OS and put that here: `"postgres://root:postgres@<some other host>:5432/postgres?sslmode=disable"`. But make sure to let me know, as I didn't have access to multiple operating systems for testing.


If you'd like to run pgcli to look at tables or records in the database, here ye be: `pgcli postgres://root:postgres@localhost:5432/postgres`

## Examples
Note: if you have Python on your system you should be able to format the result by adding `| python -m json.tool` to the examples below.

**Create User**

`curl -X POST 'http://127.0.0.1:1946/user' -d '{"name": "First Last", "email": "firstlast@domain.com", "password": "ArbitraryPassword%^&890"}' -v -H '{"Content-Type": "application/json"}'`

Result looks like:
`{"email": "firstlast@domain.com", "id": <userID>,"name": "First Last"}`


**Log in as this user**

`curl -X POST http://127.0.0.1:1946/login -d '{"email":  "firstlast@domain.com", "password": "ArbitraryPassword%^&890"}'`


**Get User**

`curl GET 'http://127.0.0.1:1946/user/<userID>' -H "Authorization: bearer <token>"`


**Update User**
Name, email, or password if given correct old password and strong enough new password. Or all.

Name and email example
`curl -X PATCH 'http://127.0.0.1:1946/user/<userID>' -d '{"name":"NewFirst NewLast", "email": "newfirstlast@domain.com"}'  -H "Authorization: bearer <token>"}'`

Password example
`curl -X PATCH 'http://127.0.0.1:1946/user/<userID>' -d '{"oldPassword": "ArbitraryPassword%^&890", "newPassword": "NewArbitraryPassword%*23"}'  -H "Authorization: bearer <token>"`


**Delete User**

`curl -X DELETE http://127.0.0.1:1946/user/<userID> -H "Authorization: bearer <token>"`


**Log out as this user**

`curl -X POST http://127.0.0.1:1946/logout/<userID> -H "Authorization: bearer <token>`

**Check password strength**
Should probably be something done on the front end with the zxcvbn package, but can be done like this.

`curl -X POST http://127.0.0.1:1946/password-strength -d '<your choice of password goes here>' -v -H '{"content-type": "text/plain"}'`


**Validate an email address**
Intended for use on the front end to check the validity of the password before submitting.

`curl -X POST http://127.0.0.1:1946/validate-email -d '{"email": <an email address to check>}' -v -H '{"content-type": "application/json"}'`


## Enhancements

The first thing I would want to do is talk to the imaginary front end using this service and ensure that we're returning the data they need in a way they'd like to read it. To that end, I'd probably want to set up a graphQL layer because that makes for a very pleasant front end development experience.

While I'm pretty happy with how it turned out trying the AWS SAM structure, the next immediate enhancement would be to ask for review or collaboration from teammates with more experience and see what "gotchas" I may have naively missed, or what useful optimizations are available or baked into the SAM. Tough spot of course, because you don't know what you're missing until you know what you're missing.

On that note, I noticed that the build process, and separately the time it was taking to run the tests, was greater than I would have liked. I did not choose to spend time optimizing for speed yet.
