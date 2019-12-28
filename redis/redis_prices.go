package redis

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"gelato"
	"github.com/go-redis/redis/v7"
)

type CacheService struct {
	rClnt *redis.Client
}

func NewCacheService(r *redis.Client) *CacheService {
	return &CacheService{rClnt: r}
}

func (p *CacheService) Get(SKU string) (*gelato.ProductPriceData, error) {
	val, err := p.rClnt.Get(generateHash(SKU)).Result()
	if err != nil {
		return nil, err
	}

	ppd := &gelato.ProductPriceData{}
	err = json.Unmarshal([]byte(val), ppd)
	if err != nil {
		return nil, err
	}

	return ppd, nil
}

func generateHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
