package server

import (
	"fmt"
	"gelato"
	"gelato/redis"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", Run)
	err := http.ListenAndServe(":8080", nil)
	if err != nil { // in real-world it would be logrus or zap
		log.Println(err)
	}
}

func Run(w http.ResponseWriter, r *http.Request) {
	rConn := redis.NewConn(&redis.RedisConf{ // in real-world app it will be ENV vars
		Host: "127.0.0.1:3306",
		Pwd:  "",
		Db:   0,
	})
	cacheService := redis.NewCacheService(rConn)
	co := gelato.NewCheckOut(cacheService, []*map[uint64]uint8{ // assume we got those from cache/mongodb
		0: {1000: 20},
		1: {500: 15},
		2: {300: 10},
		3: {200: 5},
	})

	uids, ok := r.URL.Query()["uid"] // it would be JWT decrypted message
	if !ok {
		log.Println("no param uid in uri")
	}
	uid, err := strconv.ParseInt(uids[0], 10, 64)
	if err != nil {
		log.Println(err)
	}
	co.Scan(uid, &gelato.Item{
		SKU:            "A",
		UnitPriceCents: 1000,
	})
	co.Scan(uid, &gelato.Item{
		SKU:            "B",
		UnitPriceCents: 3000,
	})
	t := co.Total()

	_, err = w.Write([]byte(fmt.Sprintf("Total price: %d", t)))
	if err != nil {
		log.Println(err)
	}
}
