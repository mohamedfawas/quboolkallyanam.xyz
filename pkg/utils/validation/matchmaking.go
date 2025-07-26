package validation

const (
	// MatchMaking options liked, passed, mutual matches
	MatchMakingOptionLiked string = "liked"
	MatchMakingOptionPassed string = "passed"
	MatchMakingOptionMutual string = "mutual"
)

func IsValidMatchMakingOption(option string) bool {
	return option == MatchMakingOptionLiked || option == MatchMakingOptionPassed || option == MatchMakingOptionMutual
}