package realtime_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/wme/pkg/realtime"
	"github.com/protsack-stephan/wme/schema/v2"
	"github.com/stretchr/testify/suite"
)

type realtimeTestSuite struct {
	suite.Suite
	ctx context.Context
	srv *httptest.Server
	cli *realtime.Client
	dat []string
	err error
}

func (s *realtimeTestSuite) createServer() http.Handler {
	gin.SetMode(gin.TestMode)
	rtr := gin.New()

	rtr.POST("/articles", func(gcx *gin.Context) {
		for _, art := range s.dat {
			fmt.Fprintf(gcx.Writer, "%s\n", art)
			gcx.Writer.Flush()
		}
	})

	return rtr
}

func (s *realtimeTestSuite) SetupSuite() {
	s.srv = httptest.NewServer(s.createServer())
	s.cli = realtime.NewClient()
	s.cli.BaseURL = s.srv.URL
	s.ctx = context.Background()
}

func (s *realtimeTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *realtimeTestSuite) TestArticles() {
	dat := []string{}

	err := s.cli.Articles(s.ctx, nil, func(art *schema.Article) error {
		dar, err := json.Marshal(art)

		if err != nil {
			return err
		}

		dat = append(dat, string(dar))
		return s.err
	})

	s.Assert().Equal(s.err, err)
	s.Assert().Equal(s.dat, dat)
}

func TestClient(t *testing.T) {
	for _, testCase := range []*realtimeTestSuite{
		{
			dat: []string{
				`{"name":"Earth"}`,
			},
			err: errors.New("test error"),
		},
		{
			dat: []string{
				`{"name":"Earth"}`,
				`{"identifier":100}`,
			},
		},
	} {
		suite.Run(t, testCase)
	}
}
