package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

type MaritalStatusEnum string

const (
	MaritalStatusNeverMarried  MaritalStatusEnum = "never_married"
	MaritalStatusDivorced      MaritalStatusEnum = "divorced"
	MaritalStatusNikkahDivorce MaritalStatusEnum = "nikkah_divorce"
	MaritalStatusWidowed       MaritalStatusEnum = "widowed"
	MaritalStatusNotMentioned  MaritalStatusEnum = "not_mentioned"
)

func (m MaritalStatusEnum) Value() (driver.Value, error) {
	return string(m), nil
}

func (m *MaritalStatusEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if s, ok := value.(string); ok {
		*m = MaritalStatusEnum(s)
		return nil
	}
	return fmt.Errorf("cannot scan %T into MaritalStatusEnum", value)
}

func (p *PartnerPreference) SetMaritalStatuses(statuses []MaritalStatusEnum) error {
	stringStatuses := make([]string, len(statuses))
	for i, s := range statuses {
		stringStatuses[i] = string(s)
	}
	data, err := json.Marshal(stringStatuses)
	if err != nil {
		return err
	}
	p.PreferredMaritalStatus = datatypes.JSON(data)
	return nil
}

func (p *PartnerPreference) GetMaritalStatuses() ([]MaritalStatusEnum, error) {
	var stringStatuses []string
	if len(p.PreferredMaritalStatus) == 0 {
		return []MaritalStatusEnum{}, nil
	}

	if err := json.Unmarshal(p.PreferredMaritalStatus, &stringStatuses); err != nil {
		return nil, err
	}

	statuses := make([]MaritalStatusEnum, len(stringStatuses))
	for i, s := range stringStatuses {
		statuses[i] = MaritalStatusEnum(s)
	}
	return statuses, nil
}
