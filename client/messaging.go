package client

import (
	"bytes"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/beeper/libgmessages/pblite"
)

func (client *Client) messagingHeaders() http.Header {
	h := http.Header{}

	h.Set("Authority", headerAuthority)
	h.Set("Origin", headerOrigin)
	h.Set("Referer", headerReferer)
	h.Set("Content-Type", "application/json+protobuf")
	h.Set("Sec-Fetch-Dest", headerSecFetchDest)
	h.Set("Sec-Fetch-Mode", headerSecFetchMode)
	h.Set("Sec-Fetch-Site", headersSecFetchSite)
	h.Set("TE", headerTE)
	h.Set("User-Agent", headerUserAgent)
	h.Set("X-User-Agent", headerXUserAgent)
	h.Set("x-goog-api-key", googleAPIKey)

	return h
}

func (client *Client) messagingRequest(message proto.Message, endpoint string) (*http.Request, error) {
	rawBody, err := pblite.Marshal(message)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(rawBody)
	req, err := http.NewRequest(http.MethodPost, messagingURL+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header = client.messagingHeaders()

	return req, nil
}
