package discount

import (
	"time"
)

type Coupon struct {
	ID      string
	Code    string
	Name    string
	Rule    Rule
	Percent int
	Expire  time.Time
}

func (c Coupon) GetPercentage() int {
	return c.Percent
}

func (c Coupon) GetName() string {
	return c.Name
}
