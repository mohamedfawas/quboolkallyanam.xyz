package validation

import "fmt"

type ProfessionType string

const (
	ProfessionTypeFullTime     ProfessionType = "full_time"
	ProfessionTypePartTime     ProfessionType = "part_time"
	ProfessionTypeFreelance    ProfessionType = "freelance"
	ProfessionTypeSelfEmployed ProfessionType = "self_employed"
	ProfessionTypeNotWorking   ProfessionType = "not_working"
	ProfessionTypeAny          ProfessionType = "any"
)

var validProfessionTypes = func() map[ProfessionType]struct{} {
	vals := []ProfessionType{
		ProfessionTypeFullTime,
		ProfessionTypePartTime,
		ProfessionTypeFreelance,
		ProfessionTypeSelfEmployed,
		ProfessionTypeNotWorking,
		ProfessionTypeAny,
	}
	m := make(map[ProfessionType]struct{}, len(vals))
	for _, p := range vals {
		m[p] = struct{}{}
	}
	return m
}()

func (p ProfessionType) IsValid() bool {
	_, ok := validProfessionTypes[p]
	return ok
}

func IsValidProfessionType(professionType string) bool {
	return ProfessionType(professionType).IsValid()
}

func ParsePreferredProfessionTypes(input []string) ([]ProfessionType, error) {
	for _, s := range input {
		if s == string(ProfessionTypeAny) {
			return []ProfessionType{ProfessionTypeAny}, nil
		}
	}

	var out []ProfessionType
	for _, s := range input {
		p := ProfessionType(s)
		if !p.IsValid() {
			return nil, fmt.Errorf("invalid profession type: %q", s)
		}
		out = append(out, p)
	}
	return out, nil
}
