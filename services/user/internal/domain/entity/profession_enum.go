package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type ProfessionEnum string

const (
	ProfessionStudent      ProfessionEnum = "student"
	ProfessionDoctor       ProfessionEnum = "doctor"
	ProfessionEngineer     ProfessionEnum = "engineer"
	ProfessionFarmer       ProfessionEnum = "farmer"
	ProfessionTeacher      ProfessionEnum = "teacher"
	ProfessionNotMentioned ProfessionEnum = "not_mentioned"
)

func (p ProfessionEnum) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *ProfessionEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok {
		*p = ProfessionEnum(s)
		return nil
	}
	return fmt.Errorf("cannot scan %T into ProfessionEnum", value)
}

func (p *PartnerPreference) SetProfessions(professions []ProfessionEnum) error {
	stringProfessions := make([]string, len(professions))
	for i, prof := range professions {
		stringProfessions[i] = string(prof)
	}
	data, err := json.Marshal(stringProfessions)
	if err != nil {
		return err
	}
	p.PreferredProfessions = datatypes.JSON(data)
	return nil
}

func (p *PartnerPreference) GetProfessions() ([]ProfessionEnum, error) {
	var stringProfessions []string
	if len(p.PreferredProfessions) == 0 {
		return []ProfessionEnum{}, nil
	}

	if err := json.Unmarshal(p.PreferredProfessions, &stringProfessions); err != nil {
		return nil, err
	}

	professions := make([]ProfessionEnum, len(stringProfessions))
	for i, s := range stringProfessions {
		professions[i] = ProfessionEnum(s)
	}
	return professions, nil
}
