package discount

import (
	"time"
)

type Coupon struct {
	ID      string    `json:"id"`
	Code    string    `json:"code"`
	CartId  string    `json:"cart_id"`
	Name    string    `json:"name"`
	Rule    Rule      `json:"rule"`
	Percent int       `json:"percent"`
	Expire  time.Time `json:"expire"`
}

func (c Coupon) GetPercentage() int {
	return c.Percent
}

func (c Coupon) GetName() string {
	return c.Name
}

func (c Coupon) GetExpireTime() time.Time {
	return c.Expire
}

func (c Coupon) IsExpired() bool {
	return c.Expire.Before(time.Now())
}

func (c Coupon) IsUsed() bool {
	return c.CartId != ""
}
