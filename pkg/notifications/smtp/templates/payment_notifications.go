package templates

import (
	"fmt"
	"time"
)

func BuildPaymentSuccessBody(userEmail, subscriptionID, planID, paymentID string, subscriptionEndDate time.Time) string {
	return fmt.Sprintf(
		`Dear Valued User,

Congratulations! Your payment has been successfully processed and your Premium subscription is now active.

Payment Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
• Subscription ID: %s
• Plan ID: %s  
• Payment ID: %s
• Active Until: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Premium Benefits:
✓ Unlimited profile views and contacts
✓ Advanced search filters
✓ Priority customer support
✓ Enhanced privacy settings

You can now enjoy all premium features of Qubool Kallyanam. Thank you for choosing us to be part of your journey in finding the perfect life partner.

If you have any questions, feel free to contact our support team.

Best regards,
Qubool Kallyanam Team
Customer Support: support@qubool-kallyanam.xyz`,
		subscriptionID,
		planID,
		paymentID,
		subscriptionEndDate.Format("January 2, 2006"),
	)
}