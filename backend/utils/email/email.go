package email

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"sync"
	"time"
	"trxd/db"
	"trxd/utils/consts"

	gomail "github.com/wneessen/go-mail"
)

var clientInitMutex sync.Mutex
var client *gomail.Client
var fromAddr *mail.Address

func InitEmailClientFromConfigs(ctx context.Context) error {
	server, err := db.GetConfig(ctx, "email-server")
	if err != nil {
		return fmt.Errorf("failed to get email-server config: %w", err)
	}

	portStr, err := db.GetConfig(ctx, "email-port")
	if err != nil {
		return fmt.Errorf("failed to get email-port config: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid email-port config: %w", err)
	}

	addr, err := db.GetConfig(ctx, "email-addr")
	if err != nil {
		return fmt.Errorf("failed to get email-addr config: %w", err)
	}

	passwd, err := db.GetConfig(ctx, "email-passwd")
	if err != nil {
		return fmt.Errorf("failed to get email-passwd config: %w", err)
	}

	if server == "" || port == 0 || addr == "" {
		return errors.New("email configuration is incomplete")
	}

	err = InitEmailClient(server, port, addr, passwd)
	if err != nil {
		return fmt.Errorf("failed to initialize email client: %w", err)
	}

	return nil
}

func InitEmailClient(server string, port int, addr string, passwd string) error {
	clientInitMutex.Lock()
	defer clientInitMutex.Unlock()

	if client != nil {
		return nil
	}

	var err error

	fromAddr, err = mail.ParseAddress(addr)
	if err != nil {
		return fmt.Errorf("failed to parse address: %w", err)
	}

	client, err = gomail.NewClient(
		server,
		gomail.WithPort(port),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(addr),
		gomail.WithPassword(passwd),
		gomail.WithTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create new Client: %w", err)
	}

	return nil
}

// The client has an internal mutex to make this thread-safe.
func SendEmail(ctx context.Context, to string, subject string, body string) error {
	if client == nil {
		return errors.New(consts.EmailClientNotInitialized)
	}

	toAddr, err := mail.ParseAddress(to)
	if err != nil {
		return fmt.Errorf("failed to parse address: %w", err)
	}

	message := gomail.NewMsg()
	message.FromMailAddress(fromAddr)
	message.ToMailAddress(toAddr)
	message.Subject(subject)
	message.SetBodyString(gomail.TypeTextPlain, body)

	err = client.DialAndSendWithContext(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
