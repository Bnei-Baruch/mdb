package alerting

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/Bnei-Baruch/mdb/monitor/interfaces"
	"github.com/Bnei-Baruch/mdb/monitor/internal/viewModels"
	"github.com/Bnei-Baruch/mdb/monitor/plugins/outputs"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/pkg/errors"
)

// ContentUnitDecorator is an object representing the view model object.
type ContentUnitDecorator struct {
	ID         int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	UID        string `boil:"uid" json:"uid" toml:"uid" yaml:"uid"`
	TypeID     int64  `boil:"type_id" json:"type_id" toml:"type_id" yaml:"type_id"`
	TypeName   string
	CreatedAt  time.Time              `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	Properties map[string]interface{} `boil:"properties" json:"properties,omitempty" toml:"properties" yaml:"properties,omitempty"`
	Secure     int16                  `boil:"secure" json:"secure" toml:"secure" yaml:"secure"`
	Published  bool                   `boil:"published" json:"published" toml:"published" yaml:"published"`
}

type Alerting struct {
	emailTemplateFile string
	outputEmailFile   string

	writers []io.Writer
	closers []io.Closer
}

func (d *Alerting) TryParseConfigurations(outputConfigs map[string]interface{}) error {
	log.Printf("Parsing alerting output configurations")
	d.emailTemplateFile = outputConfigs["emailTemplateFile"].(string)
	d.outputEmailFile = outputConfigs["outputEmailFile"].(string)
	log.Printf("Using emailTemplateFile: %s", d.emailTemplateFile)
	log.Printf("Using outputEmailFile: %s", d.outputEmailFile)
	return nil
}
func (d *Alerting) Connect() error {
	var of *os.File
	var err error
	if _, err := os.Stat(d.outputEmailFile); os.IsNotExist(err) {
		of, err = os.Create(d.outputEmailFile)
	} else {
		of, err = os.OpenFile(d.outputEmailFile, os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	}

	if err != nil {
		return err
	}
	d.writers = append(d.writers, of)
	d.closers = append(d.closers, of)
	return nil
}
func (d *Alerting) Close() error {
	var errS string
	for _, c := range d.closers {
		if err := c.Close(); err != nil {
			errS += err.Error() + "\n"
		}
	}
	if errS != "" {
		return fmt.Errorf(errS)
	}
	return nil
}
func (d *Alerting) SampleConfig() string { return "" }
func (d *Alerting) Description() string  { return "Send alerts about abnormal metrics" }
func (d *Alerting) Write(metrics []interfaces.Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	for _, metric := range metrics {
		if metric.Name() == "not_published_content_units" {
			var fields = metric.Fields()
			if fields != nil {
				var totalMeals = fields["totalMeals"].(int64)
				var meals []*viewModels.ContentUnitViewModel
				var mealsJSON = []byte(fields["meals"].(string))
				err := json.Unmarshal(mealsJSON, &meals)
				utils.Must(err)
				if totalMeals > 0 {
					var mealsDecorators []*ContentUnitDecorator
					for _, meal := range meals {
						mealDecorator := new(ContentUnitDecorator)
						mealDecorator.ID = meal.ID
						mealDecorator.UID = meal.UID
						mealDecorator.CreatedAt = meal.CreatedAt
						mealDecorator.Published = meal.Published
						mealDecorator.Secure = meal.Secure
						mealDecorator.TypeID = meal.TypeID
						mealDecorator.TypeName = meal.TypeName
						err := json.Unmarshal(meal.Properties.JSON, &mealDecorator.Properties)
						if err != nil {
							return errors.Wrapf(err, "json.Unmarshal properties %s", meal.UID)
						}
						mealsDecorators = append(mealsDecorators, mealDecorator)
					}
					data := map[string]interface{}{
						"NotPublishedContentUnitsTitle": "Unpublished meals",
						"ContentUnits":                  mealsDecorators,
					}
					tmpl := template.Must(template.ParseFiles(d.emailTemplateFile))
					for _, writer := range d.writers {
						err := tmpl.Execute(writer, data)
						utils.Must(err)
					}

					/*
						content, err := ioutil.ReadFile(d.outputEmailFile)
						utils.Must(err)
						body := string(content)
						var subject = "Not published content units alert"
						sendEmail(subject, body)
					*/

					log.Printf("Alerting about not_published_content_units metric with total meals %v completed", totalMeals)
				}
			}
		}

	}
	return nil
}

func init() {
	outputs.Add("alerting", func() interfaces.Output { return &Alerting{} })
}

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

func sendEmail(subject string, body string) {
	mail := Mail{}
	mail.senderId = "testaccount@gmail.com"
	mail.toIds = []string{"dmitry.khaymov@gmail.com", "edoshor@gmail.com "}
	mail.subject = subject
	mail.body = body

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{host: "smtp.gmail.com", port: "465"}

	log.Println(smtpServer.host)
	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, "password", smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	utils.Must(err)

	client, err := smtp.NewClient(conn, smtpServer.host)
	utils.Must(err)

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	utils.Must(err)

	_, err = w.Write([]byte(messageBody))
	utils.Must(err)

	err = w.Close()
	utils.Must(err)

	client.Quit()

	log.Println("Mail sent successfully")
}
