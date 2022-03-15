package keycloak

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Bnei-Baruch/mdb/common"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/volatiletech/null.v6"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type ExportKC struct {
	mdb  *sql.DB
	resp []*RespItem
}

type KCData struct {
	AccountID string `json:"id"`
	FName     string `json:"firstName"`
	LName     string `json:"lastName"`
}

type RespItem struct {
	email    string
	wasAdded bool
	error    error
}

func (e *ExportKC) Run() {
	e.mdb = e.openDB()
	defer e.mdb.Close()

	users, err := models.Users(e.mdb, qm.Where("account_id IS NULL")).All()
	utils.Must(err)
	e.resp = make([]*RespItem, len(users))
	withErr := make([]*models.User, 0)
	for i, u := range users {
		e.resp[i] = &RespItem{
			email:    u.Email,
			wasAdded: true,
			error:    nil,
		}
		if err := e.updateUser(u); err != nil {
			e.resp[i].wasAdded = false
			e.resp[i].error = err
			withErr = append(withErr, u)
			log.Debugf("Error for email: %s Error: %v ", u.Email, err)
		}
	}
	e.useEmailAsAccountId(withErr)
	e.printToCSV()
}

func (e *ExportKC) updateUser(u *models.User) error {
	url := viper.GetString("keycloak.api-url") + "?email=" + u.Email
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("can't get user data from kc, responce status: %s", resp.Status)
	}
	var data KCData
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return errors.Wrap(err, "error on decode response")
	}
	u.AccountID = null.StringFrom(data.AccountID)
	if !u.Name.Valid {
		u.Name = null.StringFrom(strings.Join([]string{data.FName, data.LName}, " "))
	}
	return u.Update(e.mdb)
}

func (e *ExportKC) useEmailAsAccountId(users []*models.User) {
	e.resp = append(e.resp, &RespItem{
		email:    "Try add email as account id for not kc users",
		wasAdded: true,
		error:    nil,
	})

	for _, u := range users {
		r := &RespItem{
			email:    u.Email,
			wasAdded: true,
			error:    nil,
		}
		u.AccountID = null.StringFrom(u.Email)
		if err := u.Update(e.mdb); err != nil {
			r.wasAdded = false
			r.error = err
			log.Debugf("Error for email: %s Error: %v ", u.Email, err)
		}
		e.resp = append(e.resp, r)
	}
}

func (c *ExportKC) openDB() *sql.DB {
	mdb, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(mdb.Ping())
	boil.SetDB(mdb)
	boil.DebugMode = true
	utils.Must(common.InitTypeRegistries(mdb))
	return mdb
}

func (e *ExportKC) printToCSV() {
	lines := []string{"Email", "Was updated", "Error"}
	for _, d := range e.resp {
		l := fmt.Sprintf("\n%s, %t, %v", d.email, d.wasAdded, d.error)
		lines = append(lines, l)
	}
	b := []byte(strings.Join(lines, ","))
	p := path.Join(viper.GetString("keycloak.os-dir"), "update-results.csv")
	utils.Must(ioutil.WriteFile(p, b, 0644))
}
