package email

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func SendVerificationEmail(email string, verificationToken string) (string, error) {
	url := "http://localhost:3000/verify-email?token=" + verificationToken
	RESEND_API_KEY := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(RESEND_API_KEY)

	params := &resend.SendEmailRequest{
		From: "Acme <noreply@mikaelsoninitiative.org>",
		To:   []string{email},
		Html: `<div style="max-width: 500px; margin: 0 auto; font-family: Arial, sans-serif; background-color: #ffffff; padding: 30px; border-radius: 8px; border: 1px solid #e5e7eb;">
  
  <h2 style="color: #111827; text-align: center; margin-bottom: 10px;">
    Confirm Your Signup
  </h2>

  <p style="color: #374151; font-size: 15px; text-align: center;">
    Hey there 👋
  </p>

  <p style="color: #374151; font-size: 15px; text-align: center; line-height: 1.5;">
    Thanks for joining <b>Reddit</b>! Please confirm your email address to activate your account.
  </p>

  <div style="text-align: center; margin: 30px 0;">
    <a href="` + url + `"
       style="
         background-color: #2563eb;
         color: #ffffff;
         padding: 14px 30px;
         text-decoration: none;
         border-radius: 6px;
         font-weight: bold;
         display: inline-block;
         font-size: 16px;
       ">
      Confirm Email
    </a>
  </div>

  <p style="color: #6b7280; font-size: 14px; text-align: center; line-height: 1.4;">
    If you didn’t sign up, you can safely ignore this email.
  </p>

</div>`,
		Subject: "Hello from Golang",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(sent.Id)
	return sent.Id, nil
}
