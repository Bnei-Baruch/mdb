package batch

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/volatiletech/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type RegexpReplacerSuite struct {
	suite.Suite
	utils.TestDBManager
	app RegexpReplacer
}

func (s *RegexpReplacerSuite) SetupSuite() {
	s.Require().Nil(s.InitTestDB())
}

func (s *RegexpReplacerSuite) TearDownSuite() {
	s.Require().Nil(s.DestroyTestDB())
}

func (s *RegexpReplacerSuite) SetupTest() {
	s.app = RegexpReplacer{
		DB:        s.DB,
		RegStr:    "(http://.{0,5}youtube)",
		NewStr:    "https://www.youtube",
		Limit:     10,
		TableName: "blog_posts",
		ColName:   "content",
	}
}

func TestRegexpReplace(t *testing.T) {
	suite.Run(t, new(RegexpReplacerSuite))
}

func (s *RegexpReplacerSuite) TestHttpToHttps() {
	act := "<p>http://www.youtube.com</p>"
	exp := "<p>https://www.youtube.com</p>"
	post := models.BlogPost{
		BlogID:  2,
		Title:   "test post ",
		Content: act,
	}
	s.NoError(post.Insert(s.DB))
	s.app.Do()
	s.NoError(post.Reload(s.DB))
	s.Equal(exp, post.Content)
}

func (s *RegexpReplacerSuite) TestPersonsPattern() {
	act := "http://www.youtube.com"
	exp := "https://www.youtube.com"
	s.app.TableName = "persons"
	s.app.ColName = "pattern"
	p := models.Person{
		UID:     "12345678",
		Pattern: null.String{String: act, Valid: true},
	}
	s.NoError(p.Insert(s.DB))
	s.app.Do()
	s.NoError(p.Reload(s.DB))
	s.Equal(exp, p.Pattern.String)
}

func (s *RegexpReplacerSuite) TestPostContentHtmlToText() {
	act := `<div attr="sad"><span>Text</span> with <i class="test">html tags</i></div>`
	exp := "Text with html tags"
	s.app.RegStr = "<[^>]*>"
	s.app.NewStr = ""
	post := models.BlogPost{
		BlogID:  2,
		Title:   "test post",
		Content: act,
	}
	s.NoError(post.Insert(s.DB))
	s.app.Do()
	s.NoError(post.Reload(s.DB))
	s.Equal(exp, post.Content)
}
