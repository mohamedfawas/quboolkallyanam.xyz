package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type EducationLevelEnum string

const (
	EducationLevelLessThanHighSchool EducationLevelEnum = "less_than_high_school"
	EducationLevelHighSchool         EducationLevelEnum = "high_school"
	EducationLevelHigherSecondary    EducationLevelEnum = "higher_secondary"
	EducationLevelUnderGraduation    EducationLevelEnum = "under_graduation"
	EducationLevelPostGraduation     EducationLevelEnum = "post_graduation"
	EducationLevelNotMentioned       EducationLevelEnum = "not_mentioned"
)

func (e EducationLevelEnum) Value() (driver.Value, error) {
	return string(e), nil
}

func (e *EducationLevelEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok {
		*e = EducationLevelEnum(s)
		return nil
	}
	return fmt.Errorf("cannot scan %T into EducationLevelEnum", value)
}

func (p *PartnerPreference) SetEducationLevels(levels []EducationLevelEnum) error {
	stringLevels := make([]string, len(levels))
	for i, l := range levels {
		stringLevels[i] = string(l)
	}
	data, err := json.Marshal(stringLevels)
	if err != nil {
		return err
	}
	p.PreferredEducationLevels = datatypes.JSON(data)
	return nil
}

func (p *PartnerPreference) GetEducationLevels() ([]EducationLevelEnum, error) {
	var stringLevels []string
	if len(p.PreferredEducationLevels) == 0 {
		return []EducationLevelEnum{}, nil
	}

	if err := json.Unmarshal(p.PreferredEducationLevels, &stringLevels); err != nil {
		return nil, err
	}

	levels := make([]EducationLevelEnum, len(stringLevels))
	for i, s := range stringLevels {
		levels[i] = EducationLevelEnum(s)
	}
	return levels, nil
}
