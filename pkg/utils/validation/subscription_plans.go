package validation

const (
	SubscriptionPlanIDLength          = 100
	SubscriptionPlanDurationDaysMin   = 1
	SubscriptionPlanDurationDaysMax   = 365
	SubscriptionPlanAmountMin         = 0.01
	SubscriptionPlanAmountMax         = 1000000
	SubscriptionPlanCurrencyLength    = 3
	SubscriptionPlanDescriptionLength = 255
	CurrencyINR                       = "INR"
)

func IsValidSubscriptionPlanID(id string) bool {
	return len(id) <= SubscriptionPlanIDLength
}

func IsValidSubscriptionPlanDurationDays(durationDays int) bool {
	return durationDays >= SubscriptionPlanDurationDaysMin && durationDays <= SubscriptionPlanDurationDaysMax
}

func IsValidSubscriptionPlanAmount(amount float64) bool {
	return amount >= SubscriptionPlanAmountMin && amount <= SubscriptionPlanAmountMax
}

func IsValidSubscriptionPlanCurrency(currency string) bool {
	if len(currency) != SubscriptionPlanCurrencyLength {
		return false
	}

	if currency != CurrencyINR {
		return false
	}

	return true
}

func IsValidSubscriptionPlanDescription(description string) bool {
	return len(description) <= SubscriptionPlanDescriptionLength
}
