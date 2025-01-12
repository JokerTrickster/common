package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type emailType string

const (
	emailTypePassword         emailType = "password"
	emailTypeAuth             emailType = "authCode"
	emailTypeReport           emailType = "report"
	emailTypeSignup           emailType = "signup"
	emailTypeFoodNameReport   emailType = "foodNameReport"
	emailTypeFoodUploadReport emailType = "foodUploadReport"
)

type EmailTemplate struct {
	Name       string
	Recipients []string
	Data       map[string]interface{}
	Type       emailType
}

type EmailService struct {
	client  *sesv2.Client
	from    string
	mailReq chan EmailTemplate
}

var emailService *EmailService

// Initialize the EmailService
func InitAwsSes(client *sesv2.Client, from string, bufferSize int) {
	emailService = &EmailService{
		client:  client,
		from:    from,
		mailReq: make(chan EmailTemplate, bufferSize),
	}
	go emailService.processEmails()
}

// Process emails from the channel
func (s *EmailService) processEmails() {
	for email := range s.mailReq {
		err := s.send(email)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}
}

// Send email using SES
func (s *EmailService) send(email EmailTemplate) error {
	// Serialize template data to JSON
	templateData, err := json.Marshal(email.Data)
	if err != nil {
		return fmt.Errorf("failed to serialize template data: %w", err)
	}

	// Send email
	_, err = s.client.SendEmail(context.TODO(), &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Template: &types.Template{
				TemplateName: aws.String(email.Name),
				TemplateData: aws.String(string(templateData)),
			},
		},
		Destination: &types.Destination{
			ToAddresses: email.Recipients,
		},
		EmailTags: []types.MessageTag{{
			Name:  aws.String("type"),
			Value: aws.String(string(email.Type)),
		}},
		FromEmailAddress: aws.String(s.from),
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// Send email (push to channel)
func (s *EmailService) QueueEmail(email EmailTemplate) {
	select {
	case s.mailReq <- email:
	default:
		log.Println("Email queue is full. Dropping email.")
	}
}

// Helper functions to send specific types of emails
func EmailSendFoodInfoEmptyReport(imageEmpty, nutrientEmpty []string) {
	emailService.QueueEmail(EmailTemplate{
		Name: "foodInfoEmptyReport",
		Recipients: []string{
			"pkjhj485@gmail.com", "dtw7225@naver.com", "ohhyejin1213@naver.com",
		},
		Data: map[string]interface{}{
			"imageMissingFoods":    imageEmpty,
			"nutrientMissingFoods": nutrientEmpty,
		},
		Type: emailTypeFoodNameReport,
	})
}

func EmailSendAuthCode(email, validateCode string) {
	emailService.QueueEmail(EmailTemplate{
		Name:       "foodAuth",
		Recipients: []string{email},
		Data: map[string]interface{}{
			"code": validateCode,
		},
		Type: emailTypeAuth,
	})
}

type ReqReportSES struct {
	UserID string
	Reason string
}

func EmailSendReport(email []string, req *ReqReportSES) {
	emailService.QueueEmail(EmailTemplate{
		Name:       "foodReport",
		Recipients: email,
		Data: map[string]interface{}{
			"userID": req.UserID,
			"reason": req.Reason,
		},
		Type: emailTypeReport,
	})
}
