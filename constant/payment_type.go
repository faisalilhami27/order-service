package constant

type PaymentType string
type PaymentTypeTitle string

const (
	PTDownPayment PaymentType = "down_payment"
	PTHalfPayment PaymentType = "half_payment"
	PTFullPayment PaymentType = "full_payment"

	PTDownPaymentTitle PaymentTypeTitle = "Down Payment"
	PTHalfPaymentTitle PaymentTypeTitle = "50% Payment"
	PTFullPaymentTitle PaymentTypeTitle = "100% Payment"
)

var mapPaymentTypeToTitle = map[PaymentType]PaymentTypeTitle{
	PTDownPayment: PTDownPaymentTitle,
	PTHalfPayment: PTHalfPaymentTitle,
	PTFullPayment: PTFullPaymentTitle,
}

func (pt PaymentType) String() string {
	return string(pt)
}

func (pt PaymentTypeTitle) String() string {
	return string(pt)
}

func (pt PaymentType) Title() PaymentTypeTitle {
	return mapPaymentTypeToTitle[pt]
}
