

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

`curl PATCH 'http://127.0.0.1:1946/user' -d '{"name": "NewFirst NewLast", "email": "newfirstlast@newdomain.com"}' -v -H '{"Content-Type": "application/json", "Authorization": "bearer <token>"}'`


**Delete User**

`curl  DELETE http://127.0.0.1:1946/user -d '{"id":"<userID>"}' -v -H '{"content-type":"application/json", "authorization": "bearer <token>"}'`


**Log out as this user**

`curl -X POST http://127.0.0.1:1946/login -d '{"email":  "firstlast@domain.com", "password": "ArbitraryPassword%^&890"}'`

**Check password strength**
Should probably be something done on the front end with the zxcvbn package, but can be done like this.

`curl -X POST http://127.0.0.1:1946/password-strength -d '<your choice of password goes here>' -v -H '{"content-type": "text/plain"}'`


**Validate an email address**
Intended for use on the front end to check the validity of the password before submitting.

`curl -X POST http://127.0.0.1:1946/validate-email -d '{"email": <an email address to check>}' -v -H '{"content-type": "application/json"}'`


**README**

Please include:
- a readme file that explains your thinking
- how to setup and run the project
- if you chose to use a database, include instructions on how to set that up
- if you have tests, include instructions on how to run them
- a description of what enhancements you might make if you had more time.
