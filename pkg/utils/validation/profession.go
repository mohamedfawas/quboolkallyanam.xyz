package validation

const (
	ConstantProfessionStudent      = "student"
	ConstantProfessionDoctor       = "doctor"
	ConstantProfessionEngineer     = "engineer"
	ConstantProfessionFarmer       = "farmer"
	ConstantProfessionTeacher      = "teacher"
	ConstantProfessionNotMentioned = "not_mentioned"
)

func IsValidProfession(profession string) bool {
	switch profession {
	case ConstantProfessionStudent, ConstantProfessionDoctor, ConstantProfessionEngineer,
		ConstantProfessionFarmer, ConstantProfessionTeacher, ConstantProfessionNotMentioned:
		return true
	default:
		return false
	}
}
