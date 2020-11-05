package order

type Item struct {

}

type Status string

const (
	StatusPending Status = "pending"
	StatusPaid Status = "paid"
)

type Order struct {
	Items []Item
	Status Status
}

