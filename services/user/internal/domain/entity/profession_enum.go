package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
)

type ProfessionEnum string

const (
	ProfessionStudent      ProfessionEnum = validation.ConstantProfessionStudent
	ProfessionDoctor       ProfessionEnum = validation.ConstantProfessionDoctor
	ProfessionEngineer     ProfessionEnum = validation.ConstantProfessionEngineer
	ProfessionFarmer       ProfessionEnum = validation.ConstantProfessionFarmer
	ProfessionTeacher      ProfessionEnum = validation.ConstantProfessionTeacher
	ProfessionNotMentioned ProfessionEnum = validation.ConstantProfessionNotMentioned
)

func isValidProfessionEnum(val string) bool {
	switch ProfessionEnum(val) {
	case ProfessionStudent, ProfessionDoctor, ProfessionEngineer,
		ProfessionFarmer, ProfessionTeacher, ProfessionNotMentioned:
		return true
	default:
		return false
	}
}

func (p ProfessionEnum) MarshalJSON() ([]byte, error) {
	if !isValidProfessionEnum(string(p)) {
		return nil, fmt.Errorf("invalid profession: %q", p)
	}
	return json.Marshal(string(p))
}

func (p *ProfessionEnum) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if !isValidProfessionEnum(s) {
		return fmt.Errorf("invalid profession: %q", s)
	}
	*p = ProfessionEnum(s)
	return nil
}

func (p ProfessionEnum) Value() (driver.Value, error) {
	return string(p), nil
}

func (p *ProfessionEnum) Scan(value interface{}) error {
	if value == nil {
		*p = ""
		return nil
	}
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into ProfessionEnum", value)
	}
	if !isValidProfessionEnum(s) {
		return fmt.Errorf("invalid profession: %q", s)
	}
	*p = ProfessionEnum(s)
	return nil
}
