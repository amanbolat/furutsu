package discount

import "time"

type Coupon struct {
	ID string
	Code string
	Discount Discount
	Expire time.Time
}