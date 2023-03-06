package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/beeper/libgmessages/pb"
)

var (
	ErrPairingAlreadyConnected = errors.New("Pair must be called before logging in")
)

const (
	pairingPayload = "\nFMozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0\020\003\032\005Linux0\001"
)

type PairingItem struct {
	// The raw data for the qr code
	Data string

	// An error that occurred. If nil, then a value should be in Data.
	Error error

	// The remaining time that this code is valid
	Timeout time.Duration
}

type pairingAgent struct {
	client *Client
	ctx    context.Context

	output chan PairingItem
}

func (agent *pairingAgent) register() {
	payload := &pb.RegisterPhoneRelayRequest{
		Header: agent.client.newRequest(),
	}

	fmt.Println("payload", payload)
}

func (agent *pairingAgent) start() chan PairingItem {
	agent.output = make(chan PairingItem, 8)

	go agent.register()

	return agent.output
}

func (client *Client) Pair(ctx context.Context) (chan PairingItem, error) {
	if client.IsConnected() {
		return nil, ErrPairingAlreadyConnected
	}

	agent := &pairingAgent{
		client: client,
		ctx:    ctx,
	}

	return agent.start(), nil
}

func (client *Client) Unpair() {

}
