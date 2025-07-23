package validation

const (
	ConstantMaritalStatusNeverMarried  = "never_married"
	ConstantMaritalStatusDivorced      = "divorced"
	ConstantMaritalStatusNikkahDivorce = "nikkah_divorce"
	ConstantMaritalStatusWidowed       = "widowed"
	ConstantMaritalStatusNotMentioned  = "not_mentioned"
)

func IsValidMaritalStatus(maritalStatus string) bool {
	switch maritalStatus {
	case ConstantMaritalStatusNeverMarried, ConstantMaritalStatusDivorced,
		ConstantMaritalStatusNikkahDivorce, ConstantMaritalStatusWidowed,
		ConstantMaritalStatusNotMentioned:
		return true
	default:
		return false
	}
}
