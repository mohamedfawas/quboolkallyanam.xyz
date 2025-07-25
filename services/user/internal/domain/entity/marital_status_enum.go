package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
)

type MaritalStatusEnum string

const (
	MaritalStatusNeverMarried  MaritalStatusEnum = validation.ConstantMaritalStatusNeverMarried
	MaritalStatusDivorced      MaritalStatusEnum = validation.ConstantMaritalStatusDivorced
	MaritalStatusNikkahDivorce MaritalStatusEnum = validation.ConstantMaritalStatusNikkahDivorce
	MaritalStatusWidowed       MaritalStatusEnum = validation.ConstantMaritalStatusWidowed
	MaritalStatusNotMentioned  MaritalStatusEnum = validation.ConstantMaritalStatusNotMentioned
)

func isValidMaritalStatusEnum(s string) bool {
	switch MaritalStatusEnum(s) {
	case MaritalStatusNeverMarried, MaritalStatusDivorced,
		MaritalStatusNikkahDivorce, MaritalStatusWidowed, MaritalStatusNotMentioned:
		return true
	default:
		return false
	}
}

func (m MaritalStatusEnum) MarshalJSON() ([]byte, error) {
	if !isValidMaritalStatusEnum(string(m)) {
		return nil, fmt.Errorf("invalid marital status: %q", m)
	}
	return json.Marshal(string(m))
}

func (m *MaritalStatusEnum) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !isValidMaritalStatusEnum(s) {
		return fmt.Errorf("invalid marital status: %q", s)
	}
	*m = MaritalStatusEnum(s)
	return nil
}

func (m MaritalStatusEnum) Value() (driver.Value, error) {
	return string(m), nil
}

func (m *MaritalStatusEnum) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into MaritalStatusEnum", value)
	}
	if !isValidMaritalStatusEnum(s) {
		return fmt.Errorf("invalid value for MaritalStatusEnum: %s", s)
	}
	*m = MaritalStatusEnum(s)
	return nil
}
