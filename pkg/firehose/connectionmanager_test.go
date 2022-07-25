package firehose_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/protsack-stephan/wme/pkg/firehose"
	"github.com/stretchr/testify/suite"
)

type connectionManagerTestSuite struct {
	suite.Suite
	since time.Time
	err   error
	evt   *firehose.Event
}

func (s *connectionManagerTestSuite) TestConnectionManagerNoErrors() {
	ctx, cancel := context.WithCancel(context.Background())
	cmr := firehose.NewConnectionManger()

	cmr.Add(&firehose.Connection{
		Since: s.since,
		Stream: func(ctx context.Context, since time.Time, cb func(evt *firehose.Event)) error {
			s.Assert().Equal(s.since, since)
			cb(s.evt)
			cancel()
			return s.err
		},
		Handler: func(evt *firehose.Event) {
			if s.evt != nil {
				s.Assert().Equal(s.evt, evt)
			}
		},
	})

	cmr.Connect(ctx, nil)
}

func (s *connectionManagerTestSuite) TestConnectionManagerWithErrors() {
	ctx, cancel := context.WithCancel(context.Background())
	cmr := firehose.NewConnectionManger()

	cmr.Add(&firehose.Connection{
		Since: s.since,
		Stream: func(ctx context.Context, since time.Time, cb func(evt *firehose.Event)) error {
			s.Assert().Equal(s.since, since)
			cb(s.evt)
			cancel()
			return s.err
		},
		Handler: func(evt *firehose.Event) {
			if s.evt != nil {
				s.Assert().Equal(s.evt, evt)
			}
		},
	})

	errs := make(chan error)
	go cmr.Connect(ctx, errs)

	for err := range errs {
		if s.err != nil {
			s.Assert().Equal(s.err, err)
		}
	}
}

func TestConnectionManager(t *testing.T) {
	for _, testCase := range []*connectionManagerTestSuite{
		{
			since: time.Now().Add(-1),
			err:   nil,
			evt:   &firehose.Event{},
		},
		{
			since: time.Now().Add(1),
			err:   errors.New(http.StatusText(http.StatusInternalServerError)),
			evt:   &firehose.Event{},
		},
	} {
		suite.Run(t, testCase)
	}
}
