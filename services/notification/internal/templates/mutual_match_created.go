package templates

import "fmt"

func BuildMutualMatchCommonBody(
	user1FullName string, user1ProfileID int64,
	user2FullName string, user2ProfileID int64,
) string {
	return fmt.Sprintf(
		`Hello,

Great news — it's a mutual match on Qubool Kallyanam!

%s (Profile ID: %d) and %s (Profile ID: %d) have both shown interest in each other's profiles. Congratulations!

You can view their profiles and start a conversation by logging in to your Qubool Kallyanam account.

A few quick tips:
- Only share personal contact details when you feel comfortable.
- Use the platform chat to get to know each other first.
- If you need any help, contact us at support@quboolkallyanam.xyz.

Wishing you the best,
Team Qubool Kallyanam`,
		user1FullName, user1ProfileID,
		user2FullName, user2ProfileID,
	)
}
