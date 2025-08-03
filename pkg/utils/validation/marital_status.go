package validation

import "fmt"

type MaritalStatus string

const (
	NeverMarried     MaritalStatus = "never_married"
	Divorced         MaritalStatus = "divorced"
	NikkahDivorce    MaritalStatus = "nikkah_divorce"
	Widowed          MaritalStatus = "widowed"
	MaritalStatusAny MaritalStatus = "any"
)

var validMaritalSet = func() map[MaritalStatus]struct{} {
	vals := []MaritalStatus{
		NeverMarried,
		Divorced,
		NikkahDivorce,
		Widowed,
		MaritalStatusAny,
	}
	m := make(map[MaritalStatus]struct{}, len(vals))
	for _, s := range vals {
		m[s] = struct{}{}
	}
	return m
}()

func (m MaritalStatus) IsValid() bool {
	_, ok := validMaritalSet[m]
	return ok
}

func IsValidMaritalStatus(maritalStatus string) bool {
	return MaritalStatus(maritalStatus).IsValid()
}

func ParsePreferredMaritalStatuses(input []string) ([]MaritalStatus, error) {
	for _, s := range input {
		if s == string(MaritalStatusAny) {
			return []MaritalStatus{MaritalStatusAny}, nil
		}
	}

	var out []MaritalStatus
	for _, s := range input {
		ms := MaritalStatus(s)
		if !ms.IsValid() {
			return nil, fmt.Errorf("invalid marital status: %q", s)
		}
		out = append(out, ms)
	}
	return out, nil
}
