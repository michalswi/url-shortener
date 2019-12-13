```
[in progress]

$ SERVICE_ADDR=:8080 go run main.go

$ curl -H "Content-Type: application/json" -X POST -d '{"longUrl":"https://golang.org/doc/effective_go.html"}' localhost:8080/links | jq
{
  "id": "rlntMCCH",
  "longUrl": "http://google.com/",
  "shortUrl": "http://localhost:8080/rlntMCCH",
  "createdAt": "2019-12-13 12:32:211213"
}

$ curl -i localhost:8080/rlntMCCH
HTTP/1.1 301 Moved Permanently
Content-Type: text/html; charset=utf-8
Location: http://google.com/
Date: Fri, 13 Dec 2019 11:32:31 GMT
Content-Length: 53

<a href="https://golang.org/doc/effective_go.html">Moved Permanently</a>.
```