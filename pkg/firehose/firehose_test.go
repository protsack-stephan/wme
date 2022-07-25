package firehose_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/protsack-stephan/wme/pkg/firehose"
	"github.com/stretchr/testify/suite"
)

type firehoseClientTestSuite struct {
	suite.Suite
	ctx   context.Context
	cl    *firehose.Client
	since time.Time
	ids   []string
	data  []string
	sts   int
	err   error
	srv   *httptest.Server
}

func (s *firehoseClientTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	hdl := func(c *gin.Context) {
		if s.err != nil {
			c.JSON(s.sts, map[string]interface{}{
				"status":  s.sts,
				"message": s.err.Error(),
			})
			return
		}

		s.Assert().Equal(s.since.UTC().Format(time.RFC3339), c.Query("since"))

		for i := range s.data {
			fmt.Fprintf(c.Writer, "id: %s\n", s.ids[i])
			fmt.Fprintf(c.Writer, "data: %s\n", s.data[i])
			fmt.Fprintf(c.Writer, "\n")
			c.Writer.Flush()
		}

		c.Status(http.StatusOK)
	}

	rtr := gin.New()
	rtr.GET("/page-update", hdl)
	rtr.GET("/page-delete", hdl)
	rtr.GET("/page-visibility", hdl)

	s.ctx = context.Background()
	s.srv = httptest.NewServer(rtr)
	s.cl = firehose.NewClient()
	s.cl.BaseURL = s.srv.URL
}

func (s *firehoseClientTestSuite) assertEvents(evs []*firehose.Event) {
	for i := range s.ids {
		evt := evs[i]
		data, _ := json.Marshal(evt.Data)
		id, _ := json.Marshal(evt.ID)

		s.Assert().Equal(s.ids[i], string(id))
		s.Assert().Equal(s.data[i], string(data))
	}

	s.Assert().Len(evs, len(s.ids))
}

func (s *firehoseClientTestSuite) TearDownSuite() {
	s.srv.Close()
}

func (s *firehoseClientTestSuite) TestPageUpdate() {
	evs := []*firehose.Event{}

	err := s.cl.PageUpdate(s.ctx, s.since, func(evt *firehose.Event) {
		evs = append(evs, evt)
	})

	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
		s.Assert().Contains(err.Error(), fmt.Sprint(s.sts))
	}

	s.assertEvents(evs)
}

func (s *firehoseClientTestSuite) TestPageDelete() {
	evs := []*firehose.Event{}

	err := s.cl.PageDelete(s.ctx, s.since, func(evt *firehose.Event) {
		evs = append(evs, evt)
	})

	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
		s.Assert().Contains(err.Error(), fmt.Sprint(s.sts))
	}

	s.assertEvents(evs)
}

func (s *firehoseClientTestSuite) TestPageVisibility() {
	evs := []*firehose.Event{}

	err := s.cl.PageVisibility(s.ctx, s.since, func(evt *firehose.Event) {
		evs = append(evs, evt)
	})

	if s.err != nil {
		s.Assert().Contains(err.Error(), s.err.Error())
		s.Assert().Contains(err.Error(), fmt.Sprint(s.sts))
	}

	s.assertEvents(evs)
}

func TestFirehoseClient(t *testing.T) {
	for _, testCase := range []*firehoseClientTestSuite{
		{
			ids: []string{
				`[{"topic":"aws.data-service.page-update.3","partition":0,"dt":"2022-07-24T13:03:10.431Z","timestamp":1658667790431,"offset":912025743}]`,
				`[{"topic":"aws.data-service.page-update.3","partition":0,"dt":"2022-07-24T13:03:10.431Z","timestamp":1658667790431,"offset":912025744}]`,
				`[{"topic":"aws.data-service.page-update.3","partition":0,"dt":"2022-07-24T13:03:10.431Z","timestamp":1658667790431,"offset":912025745}]`,
			},
			data: []string{
				`{"name":"Ninja","identifier":1290049,"date_modified":"2022-03-19T16:25:42Z"}`,
				`{"name":"Earth","identifier":1290050,"date_modified":"2022-03-19T16:25:43Z"}`,
				`{"name":"Airport","identifier":1290051,"date_modified":"2022-03-19T16:25:44Z"}`,
			},
			since: time.Now().Add(-1 * time.Hour),
		},
		{
			sts: http.StatusInternalServerError,
			err: errors.New(http.StatusText(http.StatusInternalServerError)),
		},
	} {
		suite.Run(t, testCase)
	}
}
