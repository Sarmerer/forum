package emailverification

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
	uuid "github.com/satori/go.uuid"
)

type manager struct {
	Pending map[string]string
	Emails  map[string]bool
}

func NewManager() *manager {
	return &manager{Pending: make(map[string]string), Emails: make(map[string]bool)}
}

var Manager = NewManager()

func (m *manager) SendVerificationEmail(user *models.User) error {
	var (
		code    string
		message string
		err     error
	)
	if m.emailPending(user.Email) {
		return errors.New("email already pending verification")
	}
	code = uuid.NewV4().String()
	go m.deleteUnverifiedUser(user.ID, code, user.Email, config.VerificationCodeExpiriation)

	m.addPending(code, user.Email)

	message = fmt.Sprintf("Follow the link to verily your email: http://localhost:8081/auth/verify?code=%s", code)

	if err = m.sendMail(user.Email, message); err != nil {
		return err
	}
	return nil
}

func (m *manager) sendMail(to, body string) error {
	var (
		username   string = os.Getenv("AWS_SES_USERNAME")
		password   string = os.Getenv("AWS_SES_PASSWORD")
		smtpServer string = os.Getenv("AWS_SES_SMTP_SERVER")

		from    string = "verify@sarmerer.ml"
		subject string = "Verify Email"
		message string

		headers     map[string]string = make(map[string]string)
		mimeVersion string            = "1.0"
		contentType string            = "text/plain; charset=\"utf-8\""
		encoding    string            = "base64"
	)

	if username == "" || password == "" || smtpServer == "" {
		return errors.New("could not find credinentials environment vars")
	}

	auth := smtp.PlainAuth("", username, password, smtpServer)

	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = mimeVersion
	headers["Content-Type"] = contentType
	headers["Content-Transfer-Encoding"] = encoding

	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	if err := smtp.SendMail(smtpServer+":25", auth, from, []string{to}, []byte(message)); err != nil {
		return err
	}
	return nil
}

func (m *manager) GetEmail(code string) string {
	return m.Pending[code]
}

func (m *manager) addPending(code, email string) {
	m.Pending[code] = email
	m.Emails[email] = true
}

func (m *manager) removePending(code, email string) {
	delete(m.Pending, code)
	delete(m.Emails, email)
}

func (m *manager) deleteUnverifiedUser(userID int64, code, email string, after time.Duration) {
	for {
		<-time.After(after)
		var (
			repo repository.UserRepo = crud.NewUserRepoCRUD()
			user *models.User
			err  error
		)
		if user, _, err = repo.FindByID(userID); err != nil {
			logger.CheckErrAndLog("Garbage collector", "deleting unverified user error", err)
			return
		}
		if !user.Verified {
			m.removePending(code, email)
			if _, err = repo.Delete(user.ID); err != nil {
				logger.CheckErrAndLog("Garbage collector", "deleting unverified user error", err)
				return
			}
			logger.CheckErrAndLog("Garbage collector", "deleted unverified user", nil)
		}
	}
}

func (m *manager) emailPending(email string) bool {
	return m.Emails[email] == true
}
