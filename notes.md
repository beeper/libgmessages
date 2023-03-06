# Basics

Each request contains the `RequestHeader` message. This message has a number of
fields that are explained below:

* **id**: A uuid like identifier that doesn't specify a version or a variant.
  it's unknown if using a uuid4 is okay or not long term.
* **app**: set to `Bugle` the assumed internal name of google messages.
* **auth_token_payload**: This is determined during pairing.

The `RequestHeader` message has a nested `ClientInfo` structure. The current
values from the web app are as follows:

* **major**: 2023
* **minor**: 2
* **point**: 13
* **api_version**: 4
* **platform_type**: "Desktop"

# Authentication

Authentication is done via the following base url.

```
https://instantmessaging-pa.googleapis.com/$rpc/google.internal.communications.instantmessaging.v1.Pairing/
```

The first step is to call `baseURL + "RegisterPhoneRelay"` using the
`RegisterPhoneRelayRequest` protobuf. This adds two new fields on top of the
nested `RequestHeader` message.

* **pairing_payload**: This appears to be basically just the user-agent, but
  there's some extra bytes around it as well. Not sure if they're necessary or
  not yet.
* **public_key**: A nested message.

The `PublicKey` message has the following fields:

* **type**: Currently set to 2, this appears to specify ECDH.
* **key**: This is 91 bytes long in my testing that suggests it is an ECDH
  public key in DER format.

The server response with a `RegisterPhoneRelayResponse` with as a
`temporary_access_token` field that is used to connect to the messages stream
for events and in all requests until the user has logged in and it is replaced
by the permanent access token.
