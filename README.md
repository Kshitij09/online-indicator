# Online indicator

## v0 requirements

* Self-contained, in-memory server
* `POST /register` endpoint to add new account. It should accept `name` in the request body 
* `POST /login` endpoint to login and start new session. API request expects credentials (simply name for now), response should return a token for that session
* `POST /ping` endpoint to indicate being online. Request header should include session token acquired by logging in. 
* `GET /status/<name>` should return online status of single account.
* `GET /status/all` should return online status of all the accounts
* `POST /status/batch` should provide list of ids in the request body and API should response with their online status