package offer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/superbet-group/offer.fastly.sse-load-test/internal/models"
)

const (
	getOffer        = "/offer/getOfferByDate?offerState=prematch&"
	getMatchID      = "/matches/byId?matchIds="
	dateQueryFormat = "startDate=%s"
	dateFormat      = "2006-01-02+15:04:05"
)

// Client fetches offer or single matches from server
type Client struct {
	offerHost  string
	httpClient httpClient
	location   *time.Location
}

// NewClient returns new offer client
func NewClient(host string, hc httpClient, loc *time.Location) *Client {
	return &Client{
		offerHost:  host,
		httpClient: hc,
		location:   loc,
	}
}

// GetOffer gets the offer from server
func (cl *Client) GetOffer() (*models.Offer, error) {
	url := cl.offerHost + getOffer + cl.generateDateQueryParams()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate offer request")
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get offer")
	}

	var offer models.Offer

	err = json.NewDecoder(resp.Body).Decode(&offer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall offer response")
	}

	return &offer, nil
}

// GetMatch fetches a single match from server
func (cl *Client) GetMatch(ID int64) (*models.Match, error) {
	url := cl.offerHost + getMatchID + strconv.FormatInt(ID, 10)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate offer request")
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get match")
	}

	var offer models.Offer

	err = json.NewDecoder(resp.Body).Decode(&offer)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshall offer response")
	}

	if len(offer.Data) == 0 {
		return nil, errors.Errorf("empty response for matchID:%d", ID)
	}

	return &offer.Data[0], nil
}

func (cl *Client) generateDateQueryParams() string {
	return fmt.Sprintf(dateQueryFormat, time.Now().UTC().Format(dateFormat))
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
