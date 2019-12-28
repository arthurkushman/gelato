## Supermarket checkout service to calculate the total price of a number of items

### Installation 
```go
go get github.com/arthurkushman/gilato
```

### Example of usage
```go
package yourpackage

import "github.com/arthurkushman/gilato"

co := NewCheckOut(cacheService, []*map[uint64]uint8{ 
    0: {1000: 20},
    1: {500: 15},
    2: {300: 10},
    3: {200: 5},
})

co.Scan(obj.uid, item)
co.Scan(obj.uid, item)

total := co.Total()
```

#### To run a local environment with this application
```go
docker-compose --build up
```
Make sure `Dockerfile` and `docker-compose.yml` are in root directory.

#### To stop the service 
```go
docker-compose down
```