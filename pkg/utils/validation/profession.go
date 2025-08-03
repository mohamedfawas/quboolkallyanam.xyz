package validation

import "fmt"

type Profession string

const (
	Student       Profession = "student"
	Doctor        Profession = "doctor"
	Engineer      Profession = "engineer"
	Farmer        Profession = "farmer"
	Teacher       Profession = "teacher"
	ProfessionAny Profession = "any"
)

var validProfessions = func() map[Profession]struct{} {
	vals := []Profession{
		Student,
		Doctor,
		Engineer,
		Farmer,
		Teacher,
		ProfessionAny,
	}
	m := make(map[Profession]struct{}, len(vals))
	for _, p := range vals {
		m[p] = struct{}{}
	}
	return m
}()

func (p Profession) IsValid() bool {
	_, ok := validProfessions[p]
	return ok
}

func IsValidProfession(profession string) bool {
	return Profession(profession).IsValid()
}

func ParsePreferredProfessions(input []string) ([]Profession, error) {
	for _, s := range input {
		if s == string(ProfessionAny) {
			return []Profession{ProfessionAny}, nil
		}
	}

	var out []Profession
	for _, s := range input {
		p := Profession(s)
		if !p.IsValid() {
			return nil, fmt.Errorf("invalid profession: %q", s)
		}
		out = append(out, p)
	}
	return out, nil
}
