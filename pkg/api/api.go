package api

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/pgzip"
	"github.com/protsack-stephan/wme/schema/v2"
)

const dateFormat = "2006-01-02"

type ReadCallback func(art *schema.Article) error

type Request struct {
	Fields  []string         `json:"fields,omitempty"`
	Filters []*schema.Filter `json:"filters,omitempty"`
	Limit   int              `json:"limit,omitempty"`
}

type CodesGetter interface {
	GetCodes(ctx context.Context, req *Request) ([]*schema.Code, error)
}

type CodeGetter interface {
	GetCode(ctx context.Context, idr string, req *Request) (*schema.Code, error)
}

type LanguagesGetter interface {
	GetLanguages(ctx context.Context, req *Request) ([]*schema.Language, error)
}

type LanguageGetter interface {
	GetLanguage(ctx context.Context, idr string, req *Request) (*schema.Language, error)
}

type ProjectsGetter interface {
	GetProjects(ctx context.Context, req *Request) ([]*schema.Project, error)
}

type ProjectGetter interface {
	GetProject(ctx context.Context, idr string, req *Request) (*schema.Project, error)
}

type NamespacesGetter interface {
	GetNamespaces(ctx context.Context, req *Request) ([]*schema.Namespace, error)
}

type NamespaceGetter interface {
	GetNamespace(ctx context.Context, idr int, req *Request) (*schema.Namespace, error)
}

type BatchesGetter interface {
	GetBatches(ctx context.Context, dte *time.Time, req *Request) ([]*schema.Batch, error)
}

type BatchGetter interface {
	GetBatch(ctx context.Context, dte *time.Time, idr string, req *Request) (*schema.Batch, error)
}

type BatchHeader interface {
	HeadBatch(ctx context.Context, dte *time.Time, idr string) (*schema.Headers, error)
}

type BatchReader interface {
	ReadBatch(ctx context.Context, dte *time.Time, idr string, cbk ReadCallback) error
}

type BatchDownloader interface{}

type SnapshotsGetter interface {
	GetSnapshots(ctx context.Context, req *Request) ([]*schema.Snapshot, error)
}

type SnapshotGetter interface {
	GetSnapshot(ctx context.Context, idr string, req *Request) (*schema.Snapshot, error)
}

type SnapshotHeader interface {
	HeadSnapshot(ctx context.Context, idr string) (*schema.Headers, error)
}

type SnapshotDownloader interface{}

type SnapshotReader interface {
	ReadSnapshot(ctx context.Context, idr string, cbk ReadCallback) error
}

type AccessTokenSetter interface {
	SetAccessToken(tkn string)
}

type API interface {
	CodesGetter
	CodeGetter
	LanguagesGetter
	LanguageGetter
	ProjectsGetter
	ProjectGetter
	NamespacesGetter
	NamespaceGetter
	BatchesGetter
	BatchGetter
	BatchHeader
	BatchReader
	BatchDownloader
	SnapshotsGetter
	SnapshotGetter
	SnapshotHeader
	SnapshotDownloader
	SnapshotReader
	AccessTokenSetter
}

func NewClient(ops ...func(clt *Client)) API {
	clt := &Client{
		HTTPClient: &http.Client{},
		UserAgent:  "",
		BaseUrl:    "https://api-beta.enterprise.wikimedia.com/",
	}

	for _, opt := range ops {
		opt(clt)
	}

	return clt
}

type Client struct {
	HTTPClient  *http.Client
	UserAgent   string
	BaseUrl     string
	AccessToken string
}

func (c *Client) newRequest(ctx context.Context, mtd string, pth string, req *Request) (*http.Request, error) {
	dta := []byte{}

	if req != nil {
		bdy, err := json.Marshal(req)

		if err != nil {
			return nil, err
		}

		dta = bdy
	}

	hrq, err := http.NewRequestWithContext(ctx, mtd, fmt.Sprintf("%sv2/%s", c.BaseUrl, pth), bytes.NewReader(dta))

	if err != nil {
		return nil, err
	}

	hrq.Header.Set("User-Agent", c.UserAgent)
	hrq.Header.Set("Content-Type", "application/json")
	hrq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	return hrq, nil
}

func (c *Client) do(hrq *http.Request) (*http.Response, error) {
	res, err := c.HTTPClient.Do(hrq)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		dta, err := io.ReadAll(res.Body)
		defer res.Body.Close()

		if err != nil {
			return nil, err
		}

		if len(string(dta)) == 0 {
			return nil, errors.New(res.Status)
		}

		return nil, errors.New(string(dta))
	}

	return res, nil
}

func (c *Client) getEntity(ctx context.Context, req *Request, pth string, val interface{}) error {
	hrq, err := c.newRequest(ctx, http.MethodPost, pth, req)

	if err != nil {
		return err
	}

	res, err := c.do(hrq)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(val)
}

func (c *Client) readAll(ctx context.Context, rdr io.Reader, cbk ReadCallback) error {
	gzr, err := pgzip.NewReader(rdr)

	if err != nil {
		return err
	}

	trr := tar.NewReader(gzr)

	for {
		_, err := trr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		scn := bufio.NewScanner(trr)
		scn.Buffer([]byte{}, 20971520)

		for scn.Scan() {
			art := new(schema.Article)

			if err := json.Unmarshal(scn.Bytes(), art); err != nil {
				return err
			}

			if err := cbk(art); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) readEntity(ctx context.Context, pth string, cbk ReadCallback) error {
	hrq, err := c.newRequest(ctx, http.MethodGet, pth, nil)

	if err != nil {
		return err
	}

	res, err := c.do(hrq)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return c.readAll(ctx, res.Body, cbk)
}

func (c *Client) headEntity(ctx context.Context, pth string) (*schema.Headers, error) {
	hrq, err := c.newRequest(ctx, http.MethodHead, pth, nil)

	if err != nil {
		return nil, err
	}

	res, err := c.do(hrq)

	if err != nil {
		return nil, err
	}

	hdr := &schema.Headers{
		ETag:         strings.Trim(res.Header.Get("ETag"), "\""),
		ContentType:  res.Header.Get("Content-Type"),
		AcceptRanges: res.Header.Get("Accept-Ranges"),
	}

	log.Println(res.Header)

	if lmf := res.Header.Get("Last-Modified"); len(lmf) > 0 {
		lmd, err := time.Parse(time.RFC1123, lmf)

		if err != nil {
			return nil, err
		}

		hdr.LastModified = &lmd
	}

	if ctl := res.Header.Get("Content-Length"); len(ctl) > 0 {
		cti, err := strconv.Atoi(ctl)

		if err != nil {
			return nil, err
		}

		hdr.ContentLength = cti
	}

	return hdr, nil
}

func (c *Client) SetAccessToken(tkn string) {
	c.AccessToken = tkn
}

func (c *Client) GetCodes(ctx context.Context, req *Request) ([]*schema.Code, error) {
	cds := []*schema.Code{}
	return cds, c.getEntity(ctx, req, "codes", &cds)
}

func (c *Client) GetCode(ctx context.Context, idr string, req *Request) (*schema.Code, error) {
	cde := new(schema.Code)
	return cde, c.getEntity(ctx, req, fmt.Sprintf("codes/%s", idr), cde)
}

func (c *Client) GetLanguages(ctx context.Context, req *Request) ([]*schema.Language, error) {
	lgs := []*schema.Language{}
	return lgs, c.getEntity(ctx, req, "languages", &lgs)
}

func (c *Client) GetLanguage(ctx context.Context, idr string, req *Request) (*schema.Language, error) {
	lng := new(schema.Language)
	return lng, c.getEntity(ctx, req, fmt.Sprintf("languages/%s", idr), lng)
}

func (c *Client) GetProjects(ctx context.Context, req *Request) ([]*schema.Project, error) {
	prs := []*schema.Project{}
	return prs, c.getEntity(ctx, req, "projects", &prs)
}

func (c *Client) GetProject(ctx context.Context, idr string, req *Request) (*schema.Project, error) {
	prj := new(schema.Project)
	return prj, c.getEntity(ctx, req, fmt.Sprintf("projects/%s", idr), prj)
}

func (c *Client) GetNamespaces(ctx context.Context, req *Request) ([]*schema.Namespace, error) {
	nss := []*schema.Namespace{}
	return nss, c.getEntity(ctx, req, "namespaces", &nss)
}

func (c *Client) GetNamespace(ctx context.Context, idr int, req *Request) (*schema.Namespace, error) {
	nsp := new(schema.Namespace)
	return nsp, c.getEntity(ctx, req, fmt.Sprintf("namespaces/%d", idr), nsp)
}

func (c *Client) GetBatches(ctx context.Context, dte *time.Time, req *Request) ([]*schema.Batch, error) {
	bts := []*schema.Batch{}
	return bts, c.getEntity(ctx, req, fmt.Sprintf("batches/%s", dte.Format(dateFormat)), &bts)
}

func (c *Client) GetBatch(ctx context.Context, dte *time.Time, idr string, req *Request) (*schema.Batch, error) {
	bth := new(schema.Batch)
	return bth, c.getEntity(ctx, req, fmt.Sprintf("batches/%s/%s", dte.Format(dateFormat), idr), bth)
}

func (c *Client) HeadBatch(ctx context.Context, dte *time.Time, idr string) (*schema.Headers, error) {
	return c.headEntity(ctx, fmt.Sprintf("batches/%s/%s/download", dte.Format(dateFormat), idr))
}

func (c *Client) ReadBatch(ctx context.Context, dte *time.Time, idr string, cbk ReadCallback) error {
	return c.readEntity(ctx, fmt.Sprintf("batches/%s/%s/download", dte.Format(dateFormat), idr), cbk)
}

func (c *Client) GetSnapshots(ctx context.Context, req *Request) ([]*schema.Snapshot, error) {
	sps := []*schema.Snapshot{}
	return sps, c.getEntity(ctx, req, "snapshots", &sps)
}

func (c *Client) GetSnapshot(ctx context.Context, idr string, req *Request) (*schema.Snapshot, error) {
	snp := new(schema.Snapshot)
	return snp, c.getEntity(ctx, req, fmt.Sprintf("snapshots/%s", idr), snp)
}

func (c *Client) HeadSnapshot(ctx context.Context, idr string) (*schema.Headers, error) {
	return c.headEntity(ctx, fmt.Sprintf("snapshots/%s/download", idr))
}

func (c *Client) ReadSnapshot(ctx context.Context, idr string, cbk ReadCallback) error {
	return c.readEntity(ctx, fmt.Sprintf("snapshots/%s/download", idr), cbk)
}
