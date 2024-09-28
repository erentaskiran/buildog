package utils

import (
	"api/internal/models"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func SendEmail(recipient, password string, cfg *aws.Config) error {

	sendEmailCredentials := &models.SendEmailCredentials{
		Sender:    os.Getenv("SES_MESSAGE_SENDER_EMAIL"),
		Recipient: recipient,
		Subject:   "dummy subject",
		HtmlBody: "<h1 style='color:blue'>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>" +
			"And your login password is " + password,
		TextBody: "This email was sent with Amazon SES using the AWS SDK for Go.",
		Charset:  "UTF-8",
	}

	client := ses.NewFromConfig(*cfg)

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{
				sendEmailCredentials.Recipient,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String(sendEmailCredentials.Charset),
					Data:    aws.String(sendEmailCredentials.HtmlBody),
				},
				Text: &types.Content{
					Charset: aws.String(sendEmailCredentials.Charset),
					Data:    aws.String(sendEmailCredentials.TextBody),
				},
			},
			Subject: &types.Content{
				Charset: aws.String(sendEmailCredentials.Charset),
				Data:    aws.String(sendEmailCredentials.Subject),
			},
		},
		Source: aws.String(sendEmailCredentials.Sender),
	}

	_, err := client.SendEmail(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("e-mail could not be sent: %v", err)
	}

	return nil
}
