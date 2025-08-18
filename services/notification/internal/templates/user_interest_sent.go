package templates

import "fmt"

func BuildInterestSentBody(receiverEmail string, senderProfileID int64, senderFullName string) string {
	return fmt.Sprintf(
		`Hello %s,

You have received a new interest on Qubool Kallyanam!

%s (Profile ID: %d) has expressed interest in your profile.

To view their profile and respond, please log in to your Qubool Kallyanam account.

If you did not expect this or need help, please contact us at support@quboolkallyanam.xyz.

Warm regards,
Team Qubool Kallyanam`,
		receiverEmail,
		senderFullName,
		senderProfileID,
	)
}
