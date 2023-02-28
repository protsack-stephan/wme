// Package api holds an API client for Wikimedia Enterprise API(s) version two.
// Here you can find helper function to get started.
// Not that this client is only in beta at them moment, and was not production tested.
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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/klauspost/pgzip"
	"github.com/protsack-stephan/wme/schema/v2"
)

const dateFormat = "2006-01-02"

// ReadCallback is a function that will be called with each Article object that is read from a batch or snapshot.
// You can return a custom error to stop the reading.
type ReadCallback func(art *schema.Article) error

// Request contains properties that are used to apply filters to the API.
type Request struct {
	// Fields represents a list of fields to retrieve from the API.
	// This is an optional argument.
	Fields []string `json:"fields,omitempty"`

	// Filters represents a list of filters to apply to the response.
	// This is an optional argument.
	Filters []*schema.Filter `json:"filters,omitempty"`

	// Limits the amount of results from the API (for now works only with Articles API).
	// This is an optional argument.
	Limit int `json:"limit,omitempty"`
}

// CodesGetter is an interface that retrieves codes from the API.
type CodesGetter interface {
	GetCodes(ctx context.Context, req *Request) ([]*schema.Code, error)
}

// CodeGetter is an interface that retrieves a code by ID from the API.
type CodeGetter interface {
	GetCode(ctx context.Context, idr string, req *Request) (*schema.Code, error)
}

// LanguagesGetter is an interface that retrieves languages from the API.
type LanguagesGetter interface {
	GetLanguages(ctx context.Context, req *Request) ([]*schema.Language, error)
}

// LanguageGetter is an interface that retrieves a language by ID from the API.
type LanguageGetter interface {
	GetLanguage(ctx context.Context, idr string, req *Request) (*schema.Language, error)
}

// ProjectsGetter is an interface that retrieves projects from the API.
type ProjectsGetter interface {
	GetProjects(ctx context.Context, req *Request) ([]*schema.Project, error)
}

// ProjectGetter is an interface that retrieves a project by ID from the API.
type ProjectGetter interface {
	GetProject(ctx context.Context, idr string, req *Request) (*schema.Project, error)
}

// NamespacesGetter is an interface that retrieves namespaces from the API.
type NamespacesGetter interface {
	GetNamespaces(ctx context.Context, req *Request) ([]*schema.Namespace, error)
}

// NamespaceGetter is an interface that retrieves a namespace by ID from the API.
type NamespaceGetter interface {
	GetNamespace(ctx context.Context, idr int, req *Request) (*schema.Namespace, error)
}

// BatchesGetter is an interface that retrieves batches from the API.
type BatchesGetter interface {
	GetBatches(ctx context.Context, dte *time.Time, req *Request) ([]*schema.Batch, error)
}

// BatchGetter is an interface that retrieves a realtime batch by ID from the API.
type BatchGetter interface {
	GetBatch(ctx context.Context, dte *time.Time, idr string, req *Request) (*schema.Batch, error)
}

// BatchHeader is an interface that retrieves the header of a realtime batch by ID from the API.
type BatchHeader interface {
	HeadBatch(ctx context.Context, dte *time.Time, idr string) (*schema.Headers, error)
}

// BatchReader is an interface that reads a realtime batch data by ID from the API.
type BatchReader interface {
	ReadBatch(ctx context.Context, dte *time.Time, idr string, cbk ReadCallback) error
}

// BatchDownloader is an interface that downloads a realtime batch `tar.gz` by ID file from the API.
type BatchDownloader interface {
	DownloadBatch(ctx context.Context, dte *time.Time, idr string, wsk io.WriteSeeker) error
}

// SnapshotsGetter is an interface for getting multiple snapshots.
type SnapshotsGetter interface {
	GetSnapshots(ctx context.Context, req *Request) ([]*schema.Snapshot, error)
}

// SnapshotGetter is an interface for getting a single snapshot by ID.
type SnapshotGetter interface {
	GetSnapshot(ctx context.Context, idr string, req *Request) (*schema.Snapshot, error)
}

// SnapshotHeader is an interface for getting the headers of a single snapshot by ID.
type SnapshotHeader interface {
	HeadSnapshot(ctx context.Context, idr string) (*schema.Headers, error)
}

// SnapshotDownloader is an interface for downloading a single snapshot by ID to a writer.
type SnapshotDownloader interface {
	DownloadSnapshot(ctx context.Context, idr string, wsk io.WriteSeeker) error
}

// SnapshotReader is an interface for reading the contents of a single snapshot by ID with a callback function.
type SnapshotReader interface {
	ReadSnapshot(ctx context.Context, idr string, cbk ReadCallback) error
}

// AllReader is an interface for reading all the contents of a reader with a callback function.
type AllReader interface {
	ReadAll(ctx context.Context, rdr io.Reader, cbk ReadCallback) error
}

// AccessTokenSetter is an interface for setting an access token.
type AccessTokenSetter interface {
	SetAccessToken(tkn string)
}

// API interface tha encapsulates the whole functionality of the client.
// Can be used with composition in unit testing.
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
	AllReader
	AccessTokenSetter
}

// NewClient returns a new instance of the Client that implements the API interface.
// The function takes in optional functional options that allow the caller to configure
// the client with custom settings.
func NewClient(ops ...func(clt *Client)) API {
	clt := &Client{
		HTTPClient:           &http.Client{},
		DownloadMinChunkSize: 5242880,
		DownloadChunkSize:    5242880 * 5,
		DownloadConcurrency:  10,
		UserAgent:            "",
		BaseUrl:              "https://api-beta.enterprise.wikimedia.com/",
	}

	for _, opt := range ops {
		opt(clt)
	}

	return clt
}

// Client is a struct that represents an HTTP client used to interact with the API.
type Client struct {
	HTTPClient           *http.Client // HTTP client used to send requests.
	UserAgent            string       // User-agent header value sent with each request.
	BaseUrl              string       // Base URL for all API requests.
	AccessToken          string       // Access token used to authenticate requests.
	DownloadMinChunkSize int          // Minimum chunk size used for downloading resources.
	DownloadChunkSize    int          // Chunk size used for downloading resources.
	DownloadConcurrency  int          // Number of simultaneous downloads allowed.
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

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
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

type chunk struct {
	start int
	end   int
	data  []byte
}

func (c *Client) downloadEntity(ctx context.Context, pth string, wrr io.WriteSeeker) error {
	hds, err := c.headEntity(ctx, pth)

	if err != nil {
		return err
	}

	csz := c.DownloadChunkSize

	if hds.ContentLength < c.DownloadMinChunkSize {
		csz = c.DownloadMinChunkSize
	}

	cks := []*chunk{}

	for i := 0; true; i++ {
		cnk := &chunk{
			start: i * csz,
			end:   (i * csz) + csz,
		}

		if cnk.end > hds.ContentLength {
			cnk.end = hds.ContentLength
		}

		cks = append(cks, cnk)

		if cnk.end == hds.ContentLength {
			break
		}
	}

	ers := make(chan error, len(cks)*2)
	cds := make(chan *chunk, len(cks))

	go func() {
		for cnk := range cds {
			if _, err := wrr.Seek(int64(cnk.start), 0); err != nil {
				ers <- err
				return
			}

			if _, err := io.CopyN(wrr, bytes.NewReader(cnk.data), int64(cnk.end-cnk.start)); err != nil {
				ers <- err
				return
			}

			ers <- nil
		}
	}()

	dcs := c.DownloadConcurrency
	smr := make(chan struct{}, dcs)

	for _, cnk := range cks {
		go func(cnk *chunk) {
			smr <- struct{}{}
			defer func() {
				ers <- nil
				<-smr
			}()

			hrq, err := c.newRequest(ctx, http.MethodGet, pth, nil)
			hrq.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", cnk.start, cnk.end))

			if err != nil {
				ers <- err
				return
			}

			res, err := c.do(hrq)

			if err != nil {
				ers <- err
				return
			}

			defer res.Body.Close()
			cnk.data, err = io.ReadAll(res.Body)

			if err != nil {
				ers <- err
				return
			}

			cds <- cnk
		}(cnk)
	}

	for i := 0; i < cap(ers); i++ {
		if err := <-ers; err != nil {
			return err
		}
	}

	close(cds)

	return nil
}

// SetAccessToken sets the access token for the client.
func (c *Client) SetAccessToken(tkn string) {
	c.AccessToken = tkn
}

// GetCodes retrieves a list of codes, and returns an error if any.
func (c *Client) GetCodes(ctx context.Context, req *Request) ([]*schema.Code, error) {
	cds := []*schema.Code{}
	return cds, c.getEntity(ctx, req, "codes", &cds)
}

// GetCode retrieves a code by ID, and returns an error if any.
func (c *Client) GetCode(ctx context.Context, idr string, req *Request) (*schema.Code, error) {
	cde := new(schema.Code)
	return cde, c.getEntity(ctx, req, fmt.Sprintf("codes/%s", idr), cde)
}

// GetLanguages retrieves a list of languages, and returns an error if any.
func (c *Client) GetLanguages(ctx context.Context, req *Request) ([]*schema.Language, error) {
	lgs := []*schema.Language{}
	return lgs, c.getEntity(ctx, req, "languages", &lgs)
}

// GetLanguage retrieves a language by ID, and returns an error if any.
func (c *Client) GetLanguage(ctx context.Context, idr string, req *Request) (*schema.Language, error) {
	lng := new(schema.Language)
	return lng, c.getEntity(ctx, req, fmt.Sprintf("languages/%s", idr), lng)
}

// GetProjects retrieves a list of projects, and returns an error if any.
func (c *Client) GetProjects(ctx context.Context, req *Request) ([]*schema.Project, error) {
	prs := []*schema.Project{}
	return prs, c.getEntity(ctx, req, "projects", &prs)
}

// GetProject retrieves a project by ID, and returns an error if any.
func (c *Client) GetProject(ctx context.Context, idr string, req *Request) (*schema.Project, error) {
	prj := new(schema.Project)
	return prj, c.getEntity(ctx, req, fmt.Sprintf("projects/%s", idr), prj)
}

// GetNamespaces retrieves a list of namespaces, and returns an error if any.
func (c *Client) GetNamespaces(ctx context.Context, req *Request) ([]*schema.Namespace, error) {
	nss := []*schema.Namespace{}
	return nss, c.getEntity(ctx, req, "namespaces", &nss)
}

// GetNamespaces retrieves a namespaces by ID, and returns an error if any.
func (c *Client) GetNamespace(ctx context.Context, idr int, req *Request) (*schema.Namespace, error) {
	nsp := new(schema.Namespace)
	return nsp, c.getEntity(ctx, req, fmt.Sprintf("namespaces/%d", idr), nsp)
}

// GetBatches retrieves a list of batches for a specific date and request, and returns an error if any.
func (c *Client) GetBatches(ctx context.Context, dte *time.Time, req *Request) ([]*schema.Batch, error) {
	bts := []*schema.Batch{}
	return bts, c.getEntity(ctx, req, fmt.Sprintf("batches/%s", dte.Format(dateFormat)), &bts)
}

// GetBatch retrieves a single batch for a specific date and ID, and returns an error if any.
func (c *Client) GetBatch(ctx context.Context, dte *time.Time, idr string, req *Request) (*schema.Batch, error) {
	bth := new(schema.Batch)
	return bth, c.getEntity(ctx, req, fmt.Sprintf("batches/%s/%s", dte.Format(dateFormat), idr), bth)
}

// HeadBatch retrieves only the headers of a single batch for a specific date and ID, and returns an error if any.
func (c *Client) HeadBatch(ctx context.Context, dte *time.Time, idr string) (*schema.Headers, error) {
	return c.headEntity(ctx, fmt.Sprintf("batches/%s/%s/download", dte.Format(dateFormat), idr))
}

// ReadBatch reads the contents of a single batch for a specific date and ID, and invokes the specified callback function for each chunk read.
func (c *Client) ReadBatch(ctx context.Context, dte *time.Time, idr string, cbk ReadCallback) error {
	return c.readEntity(ctx, fmt.Sprintf("batches/%s/%s/download", dte.Format(dateFormat), idr), cbk)
}

// DownloadBatch downloads the contents of a single batch for a specific date and ID, and writes the data to the specified WriteSeeker.
func (c *Client) DownloadBatch(ctx context.Context, dte *time.Time, idr string, wsk io.WriteSeeker) error {
	return c.downloadEntity(ctx, fmt.Sprintf("batches/%s/%s/download", dte.Format(dateFormat), idr), wsk)
}

// GetSnapshots retrieves a list of all snapshots and returns an error if any.
func (c *Client) GetSnapshots(ctx context.Context, req *Request) ([]*schema.Snapshot, error) {
	sps := []*schema.Snapshot{}
	return sps, c.getEntity(ctx, req, "snapshots", &sps)
}

// GetSnapshot retrieves a single snapshot for a specific ID and returns an error if any.
func (c *Client) GetSnapshot(ctx context.Context, idr string, req *Request) (*schema.Snapshot, error) {
	snp := new(schema.Snapshot)
	return snp, c.getEntity(ctx, req, fmt.Sprintf("snapshots/%s", idr), snp)
}

// HeadSnapshot retrieves only the headers of a single snapshot for a specific ID, and returns an error if any.
func (c *Client) HeadSnapshot(ctx context.Context, idr string) (*schema.Headers, error) {
	return c.headEntity(ctx, fmt.Sprintf("snapshots/%s/download", idr))
}

// ReadSnapshot reads the contents of a single snapshots for a specific ID, and invokes the specified callback function for each chunk read.
func (c *Client) ReadSnapshot(ctx context.Context, idr string, cbk ReadCallback) error {
	return c.readEntity(ctx, fmt.Sprintf("snapshots/%s/download", idr), cbk)
}

// DownloadSnapshot downloads the contents of a single snapshot for a specific ID, and writes the data to the specified WriteSeeker.
func (c *Client) DownloadSnapshot(ctx context.Context, idr string, wsk io.WriteSeeker) error {
	return c.downloadEntity(ctx, fmt.Sprintf("snapshots/%s/download", idr), wsk)
}

// ReadAll reads the contents of the given io.Reader and calls the given ReadCallback function
// with each chunk of data read.
func (c *Client) ReadAll(ctx context.Context, rdr io.Reader, cbk ReadCallback) error {
	return c.readAll(ctx, rdr, cbk)
}
