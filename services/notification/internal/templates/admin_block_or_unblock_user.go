package templates

import "fmt"


func BuildAdminBlockOrUnblockUserStatusBody(email string, shouldBlock bool) string {
	action := "blocked"
	message := `Your account on Qubool Kallyanam has been temporarily blocked by our administrators due to a policy violation or unusual activity.

During this time, you will not be able to access your account or use our services.

If you believe this is a mistake or have any questions, please contact our support team at support@quboolkallyanam.xyz.`

	if !shouldBlock {
		action = "unblocked"
		message = `Your account on Qubool Kallyanam has been successfully unblocked by our administrators.

You may now log in and continue using our services as usual.

If you face any issues, please do not hesitate to contact our support team at support@quboolkallyanam.xyz.`
	}

	return fmt.Sprintf(
		`Hello %s,

We'd like to inform you that your account has been %s.

%s

Thank you for being a part of our community.

Warm regards,  
Team Qubool Kallyanam`,
		email,
		action,
		message,
	)
}
