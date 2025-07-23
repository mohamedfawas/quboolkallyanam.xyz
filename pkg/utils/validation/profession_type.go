package validation

const (
	ConstantProfessionTypeFullTime     = "full_time"
	ConstantProfessionTypePartTime     = "part_time"
	ConstantProfessionTypeFreelance    = "freelance"
	ConstantProfessionTypeSelfEmployed = "self_employed"
	ConstantProfessionTypeNotWorking   = "not_working"
	ConstantProfessionTypeNotMentioned = "not_mentioned"
)

func IsValidProfessionType(professionType string) bool {
	switch professionType {
	case ConstantProfessionTypeFullTime, ConstantProfessionTypePartTime,
		ConstantProfessionTypeFreelance, ConstantProfessionTypeSelfEmployed,
		ConstantProfessionTypeNotWorking, ConstantProfessionTypeNotMentioned:
		return true
	default:
		return false
	}
}
