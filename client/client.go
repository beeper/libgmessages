package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/http2"

	"github.com/beeper/libgmessages/pb"
)

type Client struct {
	http *http.Client

	authToken     []byte
	authTokenLock sync.Mutex
}

func New() *Client {
	return &Client{
		http: &http.Client{
			Timeout:   30 * time.Second,
			Transport: &http2.Transport{},
		},
	}
}

func (client *Client) IsConnected() bool {
	return false
}

func (client *Client) Authenticated() bool {
	client.authTokenLock.Lock()
	defer client.authTokenLock.Unlock()

	return (len(client.authToken) > 0)
}

func (client *Client) newRequest() *pb.RequestHeader {
	client.authTokenLock.Lock()
	defer client.authTokenLock.Unlock()

	return &pb.RequestHeader{
		Id:               client.generateRequestID(),
		App:              client.appName(),
		AuthTokenPayload: client.authToken,
	}
}

// generateRequestID returns a string pointer to a new uuid4
func (client *Client) generateRequestID() *string {
	id := uuid.New().String()

	return &id
}

// appName returns a string pointer to the constant string for the app name
func (client *Client) appName() *string {
	name := appName

	return &name
}

func (client *Client) Connect() error {
	// Build our protobuf request message
	pbReq := &pb.ReceiveMessagesRequest{
		Header: client.newRequest(),
	}

	// Get an http request for our protobuf message
	httpReq, err := client.messagingRequest(pbReq, "ReceiveMessages")
	if err != nil {
		return err
	}

	fmt.Println("req", httpReq)

	// Make the http request.
	resp, err := client.http.Do(httpReq)
	if err != nil {
		return nil
	}

	// Process the responses
	go func(resp *http.Response) {
		fmt.Println("here!")
		fmt.Println("status code", resp.StatusCode)
		fmt.Printf("%#v\n", resp.Header)
		d, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("body:", string(d))
		resp.Body.Close()
	}(resp)

	return nil
}
