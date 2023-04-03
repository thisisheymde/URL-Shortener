## URL-Shortener

A simple URL Shortener written in Go by using only net/http.

Tools: Go, MongoDB, Docker

To Run: `docker-compose up` and open `localhost:8080` on the browser

API can be accessed at `localhost:8081`

Using API:

- To shorten

POST to localhost:8081/api/shorten
```
{
  "url": "http://example.com"
}
```
RESPONSE
```
{
  "id": "shortened_code"
}
```

- To unshorten

GET 

```
localhost:8081/s/shortened_code
```
