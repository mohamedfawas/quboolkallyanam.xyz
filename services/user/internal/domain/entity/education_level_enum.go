package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
)

type EducationLevelEnum string

const (
	EducationLevelLessThanHighSchool EducationLevelEnum = validation.ConstantEducationLevelLessThanHighSchool
	EducationLevelHighSchool         EducationLevelEnum = validation.ConstantEducationLevelHighSchool
	EducationLevelHigherSecondary    EducationLevelEnum = validation.ConstantEducationLevelHigherSecondary
	EducationLevelUnderGraduation    EducationLevelEnum = validation.ConstantEducationLevelUnderGraduation
	EducationLevelPostGraduation     EducationLevelEnum = validation.ConstantEducationLevelPostGraduation
	EducationLevelNotMentioned       EducationLevelEnum = validation.ConstantEducationLevelNotMentioned
)

func isValidEducationLevelEnum(val string) bool {
	switch EducationLevelEnum(val) {
	case EducationLevelLessThanHighSchool, EducationLevelHighSchool,
		EducationLevelHigherSecondary, EducationLevelUnderGraduation,
		EducationLevelPostGraduation, EducationLevelNotMentioned:
		return true
	}
	return false
}

func (e EducationLevelEnum) MarshalJSON() ([]byte, error) {
	if !isValidEducationLevelEnum(string(e)) {
		return nil, fmt.Errorf("invalid EducationLevelEnum: %q", e)
	}
	return json.Marshal(string(e))
}

func (e *EducationLevelEnum) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !isValidEducationLevelEnum(s) {
		return fmt.Errorf("invalid EducationLevelEnum: %q", s)
	}
	*e = EducationLevelEnum(s)
	return nil
}

func (e EducationLevelEnum) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *EducationLevelEnum) Scan(value interface{}) error {
	if value == nil {
		*e = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into EducationLevelEnum", value)
	}
	if !isValidEducationLevelEnum(s) {
		return fmt.Errorf("invalid EducationLevelEnum: %q", s)
	}
	*e = EducationLevelEnum(s)
	return nil
}
