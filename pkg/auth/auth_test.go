package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/wme/pkg/auth"
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

type authClientLoginTestSuite struct {
	suite.Suite
	srv *httptest.Server
	acl *auth.Client
	ctx context.Context
	req *auth.LoginRequest
	res *auth.LoginResponse
	sts int
	err error
}

func (s *authClientLoginTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	rtr := gin.New()
	rtr.POST("/login", NewHandler(s.Assertions, s.sts, s.req, s.res, s.err))

	s.ctx = context.Background()
	s.srv = httptest.NewServer(rtr)
	s.acl = &auth.Client{
		HTTPClient: &http.Client{},
		BaseURL:    s.srv.URL,
	}
}

func (s *authClientLoginTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *authClientLoginTestSuite) TestLogin() {
	res, err := s.acl.Login(s.ctx, s.req)

	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
	} else {
		s.Assert().NoError(err)
	}

	s.Assert().Equal(s.res, res)
}

func TestAuthClientLogin(t *testing.T) {
	for _, testCase := range []*authClientLoginTestSuite{
		{
			req: &auth.LoginRequest{
				Username: "test",
				Password: "password",
			},
			res: &auth.LoginResponse{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				IDToken:      "id_token",
				ExpiresIn:    10,
			},
			sts: http.StatusOK,
		},
		{
			req: &auth.LoginRequest{
				Username: "test",
				Password: "password",
			},
			res: nil,
			err: errors.New("Not valid request!"),
			sts: http.StatusUnprocessableEntity,
		},
	} {
		suite.Run(t, testCase)
	}
}

type authClientRefreshTokenTestSuite struct {
	suite.Suite
	srv *httptest.Server
	acl *auth.Client
	ctx context.Context
	req *auth.RefreshTokenRequest
	res *auth.RefreshTokenResponse
	sts int
	err error
}

func (s *authClientRefreshTokenTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	rtr := gin.New()
	rtr.POST("/token-refresh", NewHandler(s.Assertions, s.sts, s.req, s.res, s.err))

	s.ctx = context.Background()
	s.srv = httptest.NewServer(rtr)
	s.acl = &auth.Client{
		HTTPClient: &http.Client{},
		BaseURL:    s.srv.URL,
	}
}

func (s *authClientRefreshTokenTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *authClientRefreshTokenTestSuite) TestRefreshToken() {
	res, err := s.acl.RefreshToken(s.ctx, s.req)

	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
	} else {
		s.Assert().NoError(err)
	}

	s.Assert().Equal(s.res, res)
}

func TestAuthClientRefreshToken(t *testing.T) {
	for _, testCase := range []*authClientRefreshTokenTestSuite{
		{
			req: &auth.RefreshTokenRequest{
				Username:     "test",
				RefreshToken: "refresh_token",
			},
			res: &auth.RefreshTokenResponse{
				AccessToken: "access_token",
				IDToken:     "id_token",
				ExpiresIn:   10,
			},
			sts: http.StatusOK,
		},
		{
			req: &auth.RefreshTokenRequest{
				Username:     "test",
				RefreshToken: "password",
			},
			res: nil,
			err: errors.New("Not valid request!"),
			sts: http.StatusUnprocessableEntity,
		},
	} {
		suite.Run(t, testCase)
	}
}
