# CRUD Application for managing book's list

### Stack
- gorilla/mux
- postgres 
- swagger ("http://localhost:8080/swagger/index.html")

### Running
```go build -o app cmd/main.go && ./app```
```source .env && go build -o app cmd/main.go && ./app```

For postgres we can use Docker

```docker run -d --name ninja-db -e POSTGRES_PASSWORD=12345 -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5432:5432 postgres```