package validation

const (
	ConstantEducationLevelLessThanHighSchool = "less_than_high_school"
	ConstantEducationLevelHighSchool         = "high_school"
	ConstantEducationLevelHigherSecondary    = "higher_secondary"
	ConstantEducationLevelUnderGraduation    = "under_graduation"
	ConstantEducationLevelPostGraduation     = "post_graduation"
	ConstantEducationLevelNotMentioned       = "not_mentioned"
)

func IsValidEducationLevel(educationLevel string) bool {
	switch educationLevel {
	case ConstantEducationLevelLessThanHighSchool, ConstantEducationLevelHighSchool,
		ConstantEducationLevelHigherSecondary, ConstantEducationLevelUnderGraduation,
		ConstantEducationLevelPostGraduation, ConstantEducationLevelNotMentioned:
		return true
	default:
		return false
	}
}
