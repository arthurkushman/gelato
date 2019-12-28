package gelato

import (
	"gelato/redis"
	"log"
)

type CheckOut interface {
	Scan(uid int64, item *Item)
	Total() uint64
}

var basket map[int64]Items

type Items []*Item

type Item struct {
	SKU            string
	UnitPriceCents uint64
	Counter        uint32
}

type ProductService struct {
	cs          *redis.CacheService
	Items       Items
	offTheTotal []*map[uint64]uint8
}

type ProductPriceData struct {
	Data map[string][]*PricingPerProduct `json:"data"` // map SKU to price rules from Redis persistent db
}

type PricingPerProduct struct {
	Amount uint32 `json:"amount"`
	Price  uint64 `json:"price"`
}

type PricingRules struct {
	OffTheTotal []*map[uint64]uint8 // reverse sorted slice TotalPrice -> PercentOff
}

func NewCheckOut(cs *redis.CacheService, offTheTotal []*map[uint64]uint8) *ProductService {
	return &ProductService{cs: cs, offTheTotal: offTheTotal}
}

func (ps *ProductService) Scan(uid int64, item *Item) {
	same := false
	for _, v := range ps.Items {
		if item.SKU == v.SKU {
			v.Counter++ // in case we need some cache storage for counter (ex.: we have > 1 billion users etc) - Redis with uuid + HyperLogLog algo
			same = true
		}
	}
	if !same {
		ps.Items = append(ps.Items, item)
	}
	basket[uid] = ps.Items
}

func (ps *ProductService) Total() uint64 {
	var sum uint64
	for _, v := range ps.Items {
		if v.Counter > 1 { // getting pricing rules for product individually
			sum += ps.countPricePerProducts(v)
		} else {
			sum += v.UnitPriceCents
		}
	}

	// total off calculation
	for _, tMap := range ps.offTheTotal {
		for total, off := range *tMap { // it will be O(n) because the map is always 1 in slice, just for sorting
			if sum > total {
				return sum - (sum * uint64(off) / 100)
			}
		}
	}
	return sum
}

// compares counters of items and cached amount of goods needed to make a discount
func (ps *ProductService) countPricePerProducts(v *Item) uint64 {
	sum := uint64(0)
	ppd, err := ps.cs.Get(v.SKU)
	if err != nil {
		log.Println(err)
		return 0
	}
	// range over pre-cached prices for particular product
	for _, pricing := range ppd.Data {
		for _, ppp := range pricing {
			if v.Counter >= ppp.Amount {
				v.Counter -= ppp.Amount
				sum += ppp.Price
			}
		}
	}
	return sum
}
