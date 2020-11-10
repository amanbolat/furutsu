package cart

import (
	"encoding/json"
	"time"

	"github.com/amanbolat/furutsu/internal/product"
)

type Item struct {
	Id      string          `json:"id"`
	Product product.Product `json:"product"`
	Amount  int             `json:"amount"`
}

// ItemSet is set of items in cart which could have discounts
// applied or not. Used for calculating totals and convenient
// representation on client side
type ItemsSet struct {
	// Set represent set of multiple items and their amounts
	// which share one discount
	Set             map[string]int `json:"set"`
	DiscountPercent int            `json:"discount_percent"`
	DiscountName    string         `json:"discount_name"`
}

type Coupon interface {
	GetPercentage() int
	GetName() string
	GetCode() string
	GetExpireTime() time.Time
	IsApplicableToItems(items map[string]Item) bool
	IsExpired() bool
	IsUsed() bool
	IsAppliedToCart(cartId string) bool
}

type Alias Cart

type jsonCart struct {
	TotalSaving     int           `json:"total_saving"`
	Total           int           `json:"total"`
	TotalForPayment int           `json:"total_for_payment"`
	DiscountSets    []interface{} `json:"discount_sets"`
	NonDiscountSet  []jsonItem    `json:"non_discount_set"`
	*Alias
}

type jsonItem struct {
	Id              string          `json:"id"`
	Product         product.Product `json:"product"`
	Amount          int             `json:"amount"`
	DiscountPercent int             `json:"discount_percent"`
	DiscountName    string          `json:"discount_name"`
}

func (c Cart) MarshalJSON() ([]byte, error) {
	var discountSets []interface{}

	for _, set := range c.DiscountSets {
		d := set.DiscountPercent
		var setItems []jsonItem
		for pId, amount := range set.Set {
			item := jsonItem{
				Id:              c.Items[pId].Id,
				Product:         c.Items[pId].Product,
				Amount:          amount,
				DiscountPercent: d,
				DiscountName:    set.DiscountName,
			}
			setItems = append(setItems, item)
		}
		discountSets = append(discountSets, setItems)
	}

	var nonDiscountSet []jsonItem

	for pId, amount := range c.NonDiscountSet.Set {
		if amount < 1 {
			continue
		}
		item := jsonItem{
			Id:      c.Items[pId].Id,
			Product: c.Items[pId].Product,
			Amount:  amount,
		}
		nonDiscountSet = append(nonDiscountSet, item)
	}

	ps := &jsonCart{
		DiscountSets:    discountSets,
		NonDiscountSet:  nonDiscountSet,
		TotalForPayment: c.TotalForPayment(),
		Total:           c.Total(),
		TotalSaving:     c.TotalSavings(),
		Alias:           (*Alias)(&c),
	}

	return json.Marshal(ps)
}

type Cart struct {
	Id     string `json:"id"`
	UserId string `json:"user_id"`
	// Items is a map of items
	// <product_id:CartItem>
	Items          map[string]Item `json:"items"`
	DiscountSets   []ItemsSet      `json:"discount_sets"`
	NonDiscountSet ItemsSet        `json:"non_discount_set"`
	Coupons        []Coupon        `json:"coupons"`
}

// TotalSavings is a sum of money which could be saved
// if discounts are applied
func (c Cart) TotalSavings() int {
	var total int
	for _, set := range c.DiscountSets {
		for productId, amount := range set.Set {
			price := c.Items[productId].Product.Price
			toPay := price * amount
			saved := toPay * set.DiscountPercent / 100
			total += saved
		}
	}

	return total
}

// Total is a sum of money that has to be payed for
// items in the cart WITHOUT discounts
func (c Cart) Total() int {
	var total int

	for _, item := range c.Items {
		toPay := item.Product.Price * item.Amount
		total += toPay
	}

	return total
}

// TotalForPayment is a sum of money that
// has to be payed. Discounts are applied
func (c Cart) TotalForPayment() int {
	return c.Total() - c.TotalSavings()
}

func (c *Cart) SetProductAmount(p product.Product, amount int) {
	if c.Items == nil {
		c.Items = make(map[string]Item)
	}

	if amount < 1 {
		delete(c.Items, p.ID)
	}

	il := Item{
		Product: p,
		Amount:  amount,
	}

	if oldItem, ok := c.Items[p.ID]; ok {
		il.Amount += oldItem.Amount
		c.Items[p.ID] = il
	}

	c.Items[p.ID] = il
}

func (c *Cart) NonDiscountSetItems() map[string]Item {
	m := make(map[string]Item)

	for k, v := range c.NonDiscountSet.Set {
		p := c.Items[k].Product
		m[k] = Item{
			Product: p,
			Amount:  v,
		}
	}

	return m
}
