package constant

type OrderStatus string

const (
	Draft 		OrderStatus = "draft"
	Pending 	OrderStatus = "pending"
	Paid    	OrderStatus = "paid"
	Cancelled 	OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
    switch s {
    case Draft, Pending, Paid, Cancelled:
        return true
    }
    return false
}