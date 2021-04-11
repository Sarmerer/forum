package emailverification

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/logger"
	"github.com/sarmerer/forum/api/models"
	"github.com/sarmerer/forum/api/repository"
	"github.com/sarmerer/forum/api/repository/crud"
)

type manager struct {
	Pending map[string]string
	Emails  map[string]string
}

func NewManager() *manager {
	return &manager{Pending: make(map[string]string), Emails: make(map[string]string)}
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
	if code, err = generateCode(6); err != nil {
		return err
	}
	go m.deleteUnverifiedUser(user.ID, code, user.Email, config.VerificationCodeExpiriation)

	m.addPending(code, user.Email)

	message = fmt.Sprintf("Your verification code: %s", code)

	if err = m.sendMail(user.Email, message); err != nil {
		return err
	}
	return nil
}

func (m *manager) ResendVerificationEmail(email string) (int, error) {
	var (
		code    string = m.getCodeByEmail(email)
		message string
		err     error

		httpTokenExpired int = 498
	)
	if code == "" {
		return httpTokenExpired, errors.New("invalid or expired code")
	}
	message = fmt.Sprintf("Your verification code: %s", code)

	if err = m.sendMail(email, message); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (m *manager) Verify(email, code string) (int, error) {
	var httpTokenExpired int = 498
	if m.getCodeByEmail(m.getEmailByCode(code)) != code {
		return httpTokenExpired, errors.New("invalid or expired token")
	}
	return http.StatusOK, nil
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

	if err := smtp.SendMail(smtpServer+":587", auth, from, []string{to}, []byte(message)); err != nil {
		return err
	}
	return nil
}

func (m *manager) getEmailByCode(code string) string {
	return m.Pending[code]
}

func (m *manager) getCodeByEmail(email string) string {
	return m.Emails[email]
}

func (m *manager) addPending(code, email string) {
	m.Pending[code] = email
	m.Emails[email] = code
}

func (m *manager) removePending(code, email string) {
	delete(m.Pending, code)
	delete(m.Emails, email)
}

func (m *manager) deleteUnverifiedUser(userID int64, code, email string, after time.Duration) {
	var (
		errText  string = "deleting unverified user error"
		instance string = "Email verification manager"
	)
	time.Sleep(after)
	var (
		repo repository.UserRepo = crud.NewUserRepoCRUD()
		user *models.User
		err  error
	)

	if user, _, err = repo.FindByID(userID); err != nil {
		logger.CheckErrAndLog(instance, errText, err)
		return
	}

	if !user.Verified {
		m.removePending(code, email)
		if _, err = repo.Delete(user.ID); err != nil {
			logger.CheckErrAndLog(instance, errText, err)
			return
		}
		logger.CheckErrAndLog(instance, "deleted unverified user", nil)
	}
}

func (m *manager) emailPending(email string) bool {
	return m.Emails[email] != ""
}

func generateCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	const otpChars = "1234567890"
	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer), nil
}
