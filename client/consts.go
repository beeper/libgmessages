package client

const (
	baseURL = "https://instantmessaging-pa.googleapis.com/$rpc/"

	messagingURL = baseURL + "google.internal.communications.instantmessaging.v1.Messaging/"
	pairingURL   = baseURL + "google.internal.communications.instantmessaging.v1.Pairing/"

	googleAPIKey = "AIzaSyCA4RsOZUFrm9whhtGosPlJLmVPnfSHKz8"

	appName = "Bugle"

	headerAuthority     = "instantmessaging-pa.googleapis.com"
	headerOrigin        = "https://messages.google.com"
	headerReferer       = "https://messages.google.com/"
	headerSecFetchDest  = "empty"
	headerSecFetchMode  = "cors"
	headersSecFetchSite = "cross-site"
	headerTE            = "trailers"
	headerUserAgent     = "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0"
	headerXUserAgent    = "grpc-web-javascript/0.1"
)
