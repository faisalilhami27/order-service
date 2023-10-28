package constant

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSettlement PaymentStatus = "settlement"
	PaymentStatusExpire     PaymentStatus = "expire"
)

func (p PaymentStatus) String() string {
	return string(p)
}
