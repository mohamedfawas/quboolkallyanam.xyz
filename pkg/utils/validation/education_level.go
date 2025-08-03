package validation

import "fmt"

type EducationLevel string

const (
	LessThanHighSchool EducationLevel = "less_than_high_school"
	HighSchool         EducationLevel = "high_school"
	HigherSecondary    EducationLevel = "higher_secondary"
	UnderGraduation    EducationLevel = "under_graduation"
	PostGraduation     EducationLevel = "post_graduation"
	EducationLevelAny  EducationLevel = "any"
)

var validEducationSet = func() map[EducationLevel]struct{} {
	vals := []EducationLevel{
		LessThanHighSchool,
		HighSchool,
		HigherSecondary,
		UnderGraduation,
		PostGraduation,
		EducationLevelAny,
	}
	m := make(map[EducationLevel]struct{}, len(vals))
	for _, lvl := range vals {
		m[lvl] = struct{}{}
	}
	return m
}()

func (e EducationLevel) IsValid() bool {
	_, ok := validEducationSet[e]
	return ok
}

func IsValidEducationLevel(educationLevel string) bool {
	return EducationLevel(educationLevel).IsValid()
}

func ParsePreferredEducationLevels(input []string) ([]EducationLevel, error) {
	// Check for presence of "any"
	for _, s := range input {
		if s == string(EducationLevelAny) {
			return []EducationLevel{EducationLevelAny}, nil
		}
	}

	// Validate and collect all other levels
	var out []EducationLevel
	for _, s := range input {
		level := EducationLevel(s)
		if !level.IsValid() {
			return nil, fmt.Errorf("invalid education level: %q", s)
		}
		out = append(out, level)
	}

	return out, nil
}
