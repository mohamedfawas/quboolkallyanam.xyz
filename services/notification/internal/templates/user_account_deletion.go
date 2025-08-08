package templates

import "fmt"


func BuildAccountDeletionBody(email string) string {
    return fmt.Sprintf(
        `Hello %s,

We are writing to confirm that your account on Qubool Kallyanam has been successfully deleted.

Thank you for having been part of our community. We are sorry to see you go, and we hope our paths cross again in the future. If you have any feedback or questions, feel free to reach out to us at support@quboolkallyanam.xyz.

Wishing you all the best in your journey ahead!

Warm regards,  
Team Qubool Kallyanam`,
        email,
    )
}
