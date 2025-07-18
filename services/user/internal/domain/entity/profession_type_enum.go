package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type ProfessionTypeEnum string

const (
	ProfessionTypeFullTime     ProfessionTypeEnum = "full_time"
	ProfessionTypePartTime     ProfessionTypeEnum = "part_time"
	ProfessionTypeFreelance    ProfessionTypeEnum = "freelance"
	ProfessionTypeSelfEmployed ProfessionTypeEnum = "self_employed"
	ProfessionTypeNotWorking   ProfessionTypeEnum = "not_working"
	ProfessionTypeNotMentioned ProfessionTypeEnum = "not_mentioned"
)

func (pt ProfessionTypeEnum) Value() (driver.Value, error) {
	return string(pt), nil
}

func (pt *ProfessionTypeEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok {
		*pt = ProfessionTypeEnum(s)
		return nil
	}
	return fmt.Errorf("cannot scan %T into ProfessionTypeEnum", value)
}

func (p *PartnerPreference) SetProfessionTypes(types []ProfessionTypeEnum) error {
	stringTypes := make([]string, len(types))
	for i, t := range types {
		stringTypes[i] = string(t)
	}
	data, err := json.Marshal(stringTypes)
	if err != nil {
		return err
	}
	p.PreferredProfessionTypes = datatypes.JSON(data)
	return nil
}

func (p *PartnerPreference) GetProfessionTypes() ([]ProfessionTypeEnum, error) {
	var stringTypes []string
	if len(p.PreferredProfessionTypes) == 0 {
		return []ProfessionTypeEnum{}, nil
	}

	if err := json.Unmarshal(p.PreferredProfessionTypes, &stringTypes); err != nil {
		return nil, err
	}

	types := make([]ProfessionTypeEnum, len(stringTypes))
	for i, s := range stringTypes {
		types[i] = ProfessionTypeEnum(s)
	}
	return types, nil
}
