package constant

type OrderStatus int
type OrderStatusString string

const (
	Inital         OrderStatus = 0
	Pending        OrderStatus = 100
	PendingPayment OrderStatus = 200
	PaymentSuccess OrderStatus = 300
	Completed      OrderStatus = 400
	Cancelled      OrderStatus = 500

	InitialString        OrderStatusString = "initial"
	PendingString        OrderStatusString = "pending"
	PendingPaymentString OrderStatusString = "pending payment"
	PaymentSuccessString OrderStatusString = "payment success"
	CompletedString      OrderStatusString = "completed"
	CancelledString      OrderStatusString = "cancelled"
)

var mapOrderStatusIntToString = map[OrderStatus]OrderStatusString{
	Inital:         InitialString,
	Pending:        PendingString,
	PendingPayment: PendingPaymentString,
	PaymentSuccess: PaymentSuccessString,
	Completed:      CompletedString,
	Cancelled:      CancelledString,
}

var mapOrderStatusStringToInt = map[OrderStatusString]OrderStatus{
	InitialString:        Inital,
	PendingString:        Pending,
	PendingPaymentString: PendingPayment,
	PaymentSuccessString: PaymentSuccess,
	CompletedString:      Completed,
	CancelledString:      Cancelled,
}

func (o OrderStatusString) String() string {
	return string(o)
}

func (o OrderStatus) Int() int {
	return int(o)
}

func (o OrderStatus) GetStatusString() OrderStatusString {
	return mapOrderStatusIntToString[o]
}

func (o OrderStatusString) GetStatusInt() OrderStatus {
	return mapOrderStatusStringToInt[o]
}

func (o OrderStatus) String() string {
	return o.GetStatusString().String()
}

func (o OrderStatusString) Int() int {
	return o.GetStatusInt().Int()
}
