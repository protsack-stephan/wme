package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/protsack-stephan/wme/pkg/api"
	"github.com/stretchr/testify/suite"
)

func createHandler(sts int, dta string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(sts)
		_, _ = w.Write([]byte(dta))
	}
}

type apiTestSuite struct {
	suite.Suite
	clt api.API
	ctx context.Context
	sra *httptest.Server
	srr *httptest.Server
	req *api.Request
	cid string
	cds string
	cde string
	lid string
	lgs string
	lgg string
	pid string
	pgs string
	pjt string
	nid int
	nss string
	nsp string
	sts int
	act string
	err error
}

func (s *apiTestSuite) createAPIServer() http.Handler {
	rtr := http.NewServeMux()

	rtr.HandleFunc("/v2/codes", createHandler(s.sts, s.cds))
	rtr.HandleFunc(fmt.Sprintf("/v2/codes/%s", s.cid), createHandler(s.sts, s.cde))

	rtr.HandleFunc("/v2/languages", createHandler(s.sts, s.lgs))
	rtr.HandleFunc(fmt.Sprintf("/v2/languages/%s", s.lid), createHandler(s.sts, s.lgg))

	rtr.HandleFunc("/v2/projects", createHandler(s.sts, s.pgs))
	rtr.HandleFunc(fmt.Sprintf("/v2/projects/%s", s.pid), createHandler(s.sts, s.pjt))

	rtr.HandleFunc("/v2/namespaces", createHandler(s.sts, s.nss))
	rtr.HandleFunc(fmt.Sprintf("/v2/namespaces/%d", s.nid), createHandler(s.sts, s.nsp))

	return rtr
}

func (s *apiTestSuite) createRealtimeServer() http.Handler {
	rtr := http.NewServeMux()
	return rtr
}

func (s *apiTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.sra = httptest.NewServer(s.createAPIServer())
	s.srr = httptest.NewServer(s.createRealtimeServer())
	s.clt = api.NewClient(func(clt *api.Client) {
		clt.BaseUrl = s.sra.URL
		clt.BaseUrl = fmt.Sprintf("%s/", s.sra.URL)
		clt.RealtimeURL = fmt.Sprintf("%s/", s.srr.URL)
	})
}

func (s *apiTestSuite) TearDownSuite() {
	s.sra.Close()
	s.srr.Close()
}

func (s *apiTestSuite) TestSetAccessToken() {
	s.clt.SetAccessToken(s.act)

	s.Assert().Equal(s.act, s.clt.(*api.Client).AccessToken)
}

func (s *apiTestSuite) TestGetCodes() {
	cds, err := s.clt.GetCodes(s.ctx, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(cds)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(cds)
	}
}

func (s *apiTestSuite) TestGetCode() {
	cde, err := s.clt.GetCode(s.ctx, s.cid, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(cde)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(cde)
	}
}

func (s *apiTestSuite) TestGetLanguages() {
	lgs, err := s.clt.GetLanguages(s.ctx, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(lgs)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(lgs)
	}
}

func (s *apiTestSuite) TestGetLanguage() {
	lgn, err := s.clt.GetLanguage(s.ctx, s.lid, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(lgn)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(lgn)
	}
}

func (s *apiTestSuite) TestGetProjects() {
	pgs, err := s.clt.GetProjects(s.ctx, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(pgs)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(pgs)
	}
}

func (s *apiTestSuite) TestGetProject() {
	pjt, err := s.clt.GetProject(s.ctx, s.pid, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(pjt)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(pjt)
	}
}

func (s *apiTestSuite) TestGetNamespaces() {
	nss, err := s.clt.GetNamespaces(s.ctx, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(nss)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(nss)
	}
}

func (s *apiTestSuite) TestGetNamespace() {
	nps, err := s.clt.GetNamespace(s.ctx, s.nid, s.req)

	if s.err != nil {
		s.Assert().Error(err)
		s.Assert().Empty(nps)
	} else {
		s.Assert().NoError(err)
		s.Assert().NotEmpty(nps)
	}
}

func (s *apiTestSuite) TestGetBatches() {}

func (s *apiTestSuite) TestGetBatch() {}

func (s *apiTestSuite) TestHeadBatch() {}

func (s *apiTestSuite) TestReadBatch() {}

func (s *apiTestSuite) TestDownloadBatch() {}

func (s *apiTestSuite) TestGetSnapshots() {}

func (s *apiTestSuite) TestGetSnapshot() {}

func (s *apiTestSuite) TestHeadSnapshot() {}

func (s *apiTestSuite) TestReadSnapshot() {}

func (s *apiTestSuite) TestDownloadSnapshot() {}

func (s *apiTestSuite) TestGetArticles() {}

func (s *apiTestSuite) TestStreamArticles() {}

func (s *apiTestSuite) TestReadAll() {}

func TestAPI(t *testing.T) {
	for _, testCase := range []*apiTestSuite{
		{
			act: "test",
			cid: "wiki",
			cds: `[
				{
					"identifier": "wiki",
					"name": "Wikipedia",
					"description": "..."
				},
				{
					"identifier": "wikibooks",
					"name": "Wikibooks",
					"description": "..."
				}
			]`,
			cde: `{
				"identifier": "wiki",
				"name": "Wikipedia",
				"description": "..."
			}`,
			lid: "en",
			lgs: `[
				{
					"identifier": "en",
					"name": "...",
					"alternate_name": "...",
					"direction": "ltr"
				},
				{
					"identifier": "ff",
					"name": "...",
					"alternate_name": "...",
					"direction": "ltr"
				}
			]`,
			lgg: `{
				"identifier": "en",
				"name": "...",
				"alternate_name": "...",
				"direction": "ltr"
			}`,
			pid: "enwiki",
			pgs: `[
				{
					"name": "...",
					"identifier": "fiwikivoyage",
					"url": "https://fi.wikivoyage.org",
					"code": "wikivoyage",
					"in_language": {
						"identifier": "fi"
					}
				},
				{
					"name": "...",
					"identifier": "enwiki",
					"url": "https://en.wikipedia.org",
					"code": "wiki",
					"in_language": {
						"identifier": "en"
					}
				}
			]`,
			pjt: `{
				"name": "...",
				"identifier": "enwiki",
				"url": "https://en.wikipedia.org",
				"code": "wiki",
				"in_language": {
					"identifier": "en"
				}
			}`,
			nid: 6,
			nss: `[
				{
					"name": "...",
					"identifier": 6,
					"description": "..."
				},
				{
					"name": "...",
					"identifier": 0,
					"description": "..."
				}
			]`,
			nsp: `{
				"name": "...",
				"identifier": 6,
				"description": "..."
			}`,
			sts: http.StatusOK,
		},
	} {
		suite.Run(t, testCase)
	}
}
