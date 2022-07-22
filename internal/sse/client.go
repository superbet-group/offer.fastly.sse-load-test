package sse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	matchStreamUrl       string = "/subscription/events/%d"
	allStreamUrl         string = "/subscription/events/all"
	liveStreamUrl        string = "/subscription/events/live"
	prematchStreamUrl    string = "/subscription/events/prematch"
	sportStreamUrl       string = "/subscription/events/prematch/sport/%d"
	tournamentsStreamUrl string = "/subscription/events/prematch/tournament?ids=%s"
	structureStreamUrl   string = "/subscription/structure"
	liveCountStreamUrl   string = "/subscription/events/live/count"
)

// Client fetches offer or single matches from server
type Client struct {
	offerHost  string
	httpClient httpClient
}

// NewClient returns new offer client
func NewClient(host string, hc httpClient) *Client {
	return &Client{
		offerHost:  host,
		httpClient: hc,
	}
}

func (c *Client) StreamMatch(id int64) error {
	url := c.offerHost + fmt.Sprintf(matchStreamUrl, id)
	//start := time.Now()
	err := c.stream(url)

	//fmt.Printf("stream duration: %f s\n", time.Now().Sub(start).Seconds())
	return err
}

func (c *Client) StreamSport(id int64) error {
	url := c.offerHost + fmt.Sprintf(sportStreamUrl, id)
	return c.stream(url)
}

func (c *Client) StreamTournaments(ids ...int64) error {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.FormatInt(id, 10)
	}

	url := c.offerHost + fmt.Sprintf(tournamentsStreamUrl, strings.Join(idsStr, ","))
	return c.stream(url)
}

func (c *Client) StreamAll() error {
	return c.stream(c.offerHost + allStreamUrl)
}

func (c *Client) StreamLive() error {
	return c.stream(c.offerHost + liveStreamUrl)
}

func (c *Client) StreamPrematch() error {
	return c.stream(c.offerHost + prematchStreamUrl)
}

func (c *Client) StreamLiveCount() error {
	return c.stream(c.offerHost + liveCountStreamUrl)
}

func (c *Client) StreamStructure() error {
	return c.stream(c.offerHost + structureStreamUrl)
}

func (c *Client) stream(url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to generate stream request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get stream "+url)
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read stream"+url)
	}

	resp.Body.Close()

	return nil
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
