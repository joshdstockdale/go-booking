# Backend for hotel reservations

# Project environment variables
```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=supersupersecret
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
MONGO_DB_URL_TEST=mongodb://localhost:27017
```

## Prerequisites
### Mongodriver
Documentation
```
```

Installing mongodb client
```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber
Documentation
```
https://gofiber.io
```

Installing gofiber
```
go get github.com/gofiber/fiber/v2
```

## Docker
### Installing mongodb
```
docker run -d -p 27017:27017 --name mongodb mongo:latest
```
