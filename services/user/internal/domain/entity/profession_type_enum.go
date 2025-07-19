package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

func isValidProfessionType(val string) bool {
	switch ProfessionTypeEnum(val) {
	case ProfessionTypeFullTime, ProfessionTypePartTime,
		ProfessionTypeFreelance, ProfessionTypeSelfEmployed,
		ProfessionTypeNotWorking, ProfessionTypeNotMentioned:
		return true
	default:
		return false
	}
}

func (pt ProfessionTypeEnum) MarshalJSON() ([]byte, error) {
	if !isValidProfessionType(string(pt)) {
		return nil, fmt.Errorf("invalid profession type: %q", pt)
	}
	return json.Marshal(string(pt))
}

func (pt *ProfessionTypeEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if !isValidProfessionType(s) {
		return fmt.Errorf("invalid profession type: %q", s)
	}
	*pt = ProfessionTypeEnum(s)
	return nil
}

func (pt ProfessionTypeEnum) Value() (driver.Value, error) {
	return string(pt), nil
}

func (pt *ProfessionTypeEnum) Scan(value interface{}) error {
	if value == nil {
		*pt = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into ProfessionTypeEnum", value)
	}
	if !isValidProfessionType(s) {
		return fmt.Errorf("invalid profession type: %q", s)
	}
	*pt = ProfessionTypeEnum(s)
	return nil
}
