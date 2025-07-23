package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
)

type ProfessionTypeEnum string

const (
	ProfessionTypeFullTime     ProfessionTypeEnum = validation.ConstantProfessionTypeFullTime
	ProfessionTypePartTime     ProfessionTypeEnum = validation.ConstantProfessionTypePartTime
	ProfessionTypeFreelance    ProfessionTypeEnum = validation.ConstantProfessionTypeFreelance
	ProfessionTypeSelfEmployed ProfessionTypeEnum = validation.ConstantProfessionTypeSelfEmployed
	ProfessionTypeNotWorking   ProfessionTypeEnum = validation.ConstantProfessionTypeNotWorking
	ProfessionTypeNotMentioned ProfessionTypeEnum = validation.ConstantProfessionTypeNotMentioned
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
