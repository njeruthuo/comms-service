package emails

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type DBHandler struct {
	DB *sql.DB
}

func SendEmailService(to, subject, token string) error {
	resetURL := os.Getenv("PASSWORD_RESET_URL")
	if resetURL == "" {
		resetURL = "#"
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<body style="margin:0; padding:0; background-color:#f4f4f7; font-family:Arial, Helvetica, sans-serif;">
	<table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f7; padding:24px 0;">
		<tr>
			<td align="center">
				<table role="presentation" width="480" cellpadding="0" cellspacing="0" style="background-color:#ffffff; border-radius:8px; padding:32px; box-shadow:0 1px 3px rgba(0,0,0,0.1);">
					<tr>
						<td>
							<h2 style="color:#1a1a1a; margin-top:0;">Reset your password</h2>
							<p style="color:#333333; font-size:15px; line-height:1.5;">
								We received a request to reset the password for your account. Click the button below to continue.
							</p>
							<p style="text-align:center; margin:32px 0;">
								<a href="%s" style="background-color:#2d6cdf; color:#ffffff; text-decoration:none; padding:12px 24px; border-radius:6px; font-size:15px; display:inline-block;">
									Reset Password
								</a>
							</p>
							<p style="color:#333333; font-size:14px; line-height:1.5;">
								If the button doesn't work, copy your reset token below and paste it into the reset password page manually:
							</p>
							<p style="background-color:#f4f4f7; border:1px solid #dddddd; border-radius:6px; padding:12px; font-family:monospace; font-size:14px; color:#1a1a1a; word-break:break-all;">
								%s
							</p>
							<p style="color:#888888; font-size:13px; line-height:1.5; margin-top:32px;">
								If you did not request a password reset, you can safely ignore this email — no changes will be made to your account.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>
`, resetURL, token)

	message := gomail.NewMessage()
	message.SetHeader("From", "juliusn411@gmail.com")
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SOURCE_EMAIL"), os.Getenv("APP_PASSWORD"))
	return d.DialAndSend(message)
}
