package ondemand_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/wme/pkg/ondemand"
	"github.com/protsack-stephan/wme/schema/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (r *Response) SetStatus(status int) {
	if status == 0 {
		r.Status = http.StatusOK
	} else {
		r.Status = status
	}
}

func (r *Response) SetMessage(message string) {
	r.Message = message
}

func NewHandler(assert *assert.Assertions, status int, params interface{}, value interface{}, err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err != nil {
			res := new(Response)
			res.SetStatus(status)
			res.SetMessage(err.Error())
			c.JSON(res.Status, res)
			return
		}

		if params != nil {
			req := reflect.New(reflect.TypeOf(params).Elem())
			assert.NoError(json.NewDecoder(c.Request.Body).Decode(req.Interface()))
			assert.Equal(params, req.Interface())
		}

		if value != nil {
			c.JSON(status, value)
			return
		}

		c.Status(status)
	}
}

type odClientArticleTestSuite struct {
	suite.Suite
	srv  *httptest.Server
	odc  *ondemand.Client
	ctx  context.Context
	req  *ondemand.ArticleRequest
	page *schema.Page
	sts  int
	err  error
}

func (s *odClientArticleTestSuite) TestArticle() {
	gin.SetMode(gin.TestMode)

	rtr := gin.New()
	rtr.GET(fmt.Sprintf("/pages/meta/%s/%s", s.req.Project, s.req.Name), NewHandler(s.Assertions, s.sts, nil, s.page, s.err))

	s.ctx = context.Background()
	s.srv = httptest.NewServer(rtr)
	s.odc = &ondemand.Client{
		HTTPClient: &http.Client{},
		BaseURL:    s.srv.URL,
	}

	res, err := s.odc.Article(s.ctx, s.req)
	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
		s.Assert().Contains(err.Error(), http.StatusText(s.sts))
	} else {
		s.Assert().NoError(err)
	}

	s.Assert().Equal(s.page, res)
	s.srv.Close()
}

func TestOndemandArticle(t *testing.T) {
	for _, testCase := range []*odClientArticleTestSuite{
		{
			req: &ondemand.ArticleRequest{
				Project: "enwiki",
				Name:    "Steamship",
			},
			page: &schema.Page{
				Name:       "Steamship",
				Identifier: 525252,
				URL:        "https://en.wikipedia.org/wiki/Steamship",
				ArticleBody: &schema.ArticleBody{
					HTML:     "<html> some html </html>",
					Wikitext: "some text here...",
				},
			},
			sts: http.StatusOK,
		},
		{
			req: &ondemand.ArticleRequest{
				Project: "enwiki",
				Name:    "Steamship",
			},
			page: nil,
			err:  errors.New("Not valid request!"),
			sts:  http.StatusUnprocessableEntity,
		},
	} {
		suite.Run(t, testCase)
	}
}

type odClientProjectTestSuite struct {
	suite.Suite
	srv      *httptest.Server
	odc      *ondemand.Client
	ctx      context.Context
	projects []*schema.Project
	sts      int
	err      error
}

func (s *odClientProjectTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	rtr := gin.New()
	rtr.GET("/projects", NewHandler(s.Assertions, s.sts, nil, s.projects, s.err))

	s.ctx = context.Background()
	s.srv = httptest.NewServer(rtr)
	s.odc = &ondemand.Client{
		HTTPClient: &http.Client{},
		BaseURL:    s.srv.URL,
	}
}

func (s *odClientProjectTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *odClientProjectTestSuite) TestProjects() {
	res, err := s.odc.Projects(s.ctx)

	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
		s.Assert().Contains(err.Error(), http.StatusText(s.sts))
	} else {
		s.Assert().NoError(err)
	}

	s.Assert().Equal(s.projects, res)
}

func TestOndemandProjects(t *testing.T) {
	for _, testCase := range []*odClientProjectTestSuite{
		{
			projects: []*schema.Project{
				{
					Name:       "Wikivir",
					Identifier: "slwikisource",
					URL:        "https://sl.wikisource.org",
				},
				{
					Name:       "Wikiverza",
					Identifier: "slwikiversity",
					URL:        "https://sl.wikiversity.org",
				},
			},
			sts: http.StatusOK,
		},
		{
			projects: nil,
			err:      errors.New("Not valid request!"),
			sts:      http.StatusUnprocessableEntity,
		},
	} {
		suite.Run(t, testCase)
	}
}
