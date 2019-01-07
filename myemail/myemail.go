package myemail

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

// Identiti - идентити по умолчанию
var Identiti string

// Username - имя пользователя по умолч
var Username string

// Password - пароль по умолч
var Password string

// Host - хост по умолч
var Host string

// Port - порт по умолч
var Port string

// From - от кого по умолч
var From mail.Address

// Auth - пробуем сохранить полученный ауф
var Auth smtp.Auth

// TLSConfig - сохраненный конфиг
var TLSConfig tls.Config

// SetParams - сохраняет параметры по умолчанию для дальнейшего спользования
func SetParams(identiti, username, password, host, port string, from mail.Address) {
	Identiti = identiti
	Username = username
	Password = password
	Host = host
	Port = port
	From = from

	// TLS config
	TLSConfig = tls.Config{
		InsecureSkipVerify: true, // try disable in prod
		ServerName:         Host,
	}

}

// SendMailWithDefaultParams - отправка письма с параметрами по умолчанию
func SendMailWithDefaultParams(to mail.Address, subj string, body string) error {
	serverAddress := fmt.Sprintf("%s:%s", Host, Port)

	// TODO: есть ли готовые модули для установки хеадеров и боди?
	headers := make(map[string]string)
	headers["From"] = From.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	if Auth == nil {
		Auth = smtp.PlainAuth("", Username, Password, Host)
	}

	conn, err := tls.Dial("tcp", serverAddress, &TLSConfig)

	c, err := smtp.NewClient(conn, Host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(Auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(From.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return err
}
