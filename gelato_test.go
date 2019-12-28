package gelato

import (
	"gelato/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tests = []struct {
	uid    int64
	items  []*Item
	result uint64
}{
	{1, []*Item{{SKU: "A", UnitPriceCents: 123}, {SKU: "B", UnitPriceCents: 12}}, 200},
	{2, []*Item{{SKU: "A", UnitPriceCents: 123, Counter: 0}}, 200},
}

func TestCheckOut(t *testing.T) {
	rConn := redis.NewConn(&redis.RedisConf{ // in real-world app it will be ENV vars
		Host: "127.0.0.1:3306",
		Pwd:  "",
		Db:   0,
	})
	cacheService := redis.NewCacheService(rConn)
	co := NewCheckOut(cacheService, []*map[uint64]uint8{ // assume we got those from cache/mongodb
		0: {1000: 20},
		1: {500: 15},
		2: {300: 10},
		3: {200: 5},
	})

	for _, obj := range tests {
		for _, item := range obj.items {
			co.Scan(obj.uid, item)
		}
		total := co.Total()
		assert.Equal(t, obj.result, total)
	}
}
