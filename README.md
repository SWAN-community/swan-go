# ![Secured Web Addressability Network](https://raw.githubusercontent.com/SWAN-community/swan/main/images/swan.128.pxls.100.dpi.png)

# Secured Web Addressability Network (SWAN) - Go Access Layer

Secure Web Addressability Network (SWAN) - an open source secure and privacy
supporting cross domain identity network implemented in Go.

## Introduction

This project contains all the server side components that are needed to
implement SWAN in Go for a publisher, Consent Management Platform (CMP) or
OpenRTB party.

The project wraps calls to SWAN API endpoints in a simple to use data model. A
connection is used to configure the mandatory parameters such as the host domain
for the SWAN Operator access node, access keys and default optional parameters
such as the message to display to the user when performing multiple node storage
operations. 

Once the connection is created functions prefixed New are used to initiate 
storage operations to fetch and update SWAN data. Methods to decrypt SWAN data 
are provided directly by the swan.Connection structure.

## Prerequisites

The reader should be familiar with the 
[SWAN concepts](https://github.com/SWAN-community/swan).

## Connection

An instance of a swan.Connection is created via a call to swan.NewConnection.
The default values to use for operations are provided to this method. These 
values can be overridden for each fetch and update storage operation. 

The following code shows how this would be achieved where an instance of the 
structure swan.Operation provides the default values to use.

```go
connection := swan.NewConnection(swan.Operation{
    Client: swan.Client{
        SWAN: swan.SWAN{
            AccessKey: "AccessKey",
            Operator:  "swan-access-node.org",
            Scheme:    "https"}},
    BackgroundColor:       "white",
    Message:               "Hello. SWAN operation in progress",
    MessageColor:          "Green",
    ProgressColor:         "Blue",
    NodeCount:             10,
    DisplayUserInterface:  true,
    PostMessageOnComplete: false,
    JavaScript:            false,
    UseHomeNode:           true})
```

See the 
[Go source code](https://github.com/SWAN-community/swan-go/blob/main/connection.go)
for the meaning of the different parameters.

## Operations

Once the connection is created with the defaults to be used for all requests
to SWAN the following operations are supported.

In describing the operations available the following common variables and used
and have the following meaning type.

| Name | Type | Description |
|-|-|-|
| encrypted | string | base 64 encoded encrypted data returned from a SWAN storage operation |
| request | *http.Request | http request associated with a user's web browser - used to set the home node |
| returnUrl | *url.URL | the URL to return to after the browser has been directed to the URL provided by the GetURL function |

### Fetch

Provides a URL that the browser should be immediately directed to. The return
URL will have the encrypted SWAN data appended ready to be used with the decrypt
functions.

```go
url := connection.NewFetch(request, returnUrl).GetURL()
```

### Update

Provides a URL that the browser should be immediately directed to. The return
URL will have the encrypted SWAN data appended ready to be used with the decrypt
functions. 

Note: The return URL and optional Access Node could have been passed by the
publisher to the User Interface Provider in the state array of the original 
Fetch operation. See the SWAN demo for passing data between different parties
via the Operation.State array parameter.

The members Pref, Email, Salt and RID should be set to the values provided by 
the user before the GetURL function is called. If they are left blank the 
existing values are removed from SWAN.

```go

// Get the OWID signer which is needed to sign the raw SWAN data. The User
// Interface Provider will be responsible for obtaining their OWID signer.
signer, err := YourMethodToGetOWIDSigner()
if err != nil { return err }

// Create a new Update operation with the request from the web browser and the
// return URL.
u := connection.NewUpdate(request, returnUrl)

// Set the raw SWAN data from the form associated with the request. Pass the 
// OWID signer to each of the methods that generates the OWID signed data 
// structure. 
// Check the err indicator incase there was a problem generating the OWID from 
// the input data or if the input data did not pass validation.
if err != nil { return err }
u.Pref = NewPreferences(signer, r.Form.Get("pref") == "on")
if err != nil { return err }
u.Email, err = NewEmail(signer, r.Form.Get("email"))
if err != nil { return err }
u.Salt, err = NewSaltFromString(signer, r.Form.Get("salt"))
if err != nil { return err }
u.RID, err = IdentifierFromBase64(r.Form.Get("rid"))
if err != nil { return err }

// Get the storage operation URL to redirect the web browser to.
url := u.GetURL()
```

### Stop

Provides a URL that the browser should be immediately directed to. The return
URL will have the encrypted SWAN data appended ready to be used with the decrypt
functions.

The host parameter is the host domain associated with the advert that should be
stopped.

```go
host := r.Form.Get("host")
url := connection.NewSWANStop(r, returnUrl, host).GetURL()
```

### Decrypt

Returns the decrypted SWAN data from the base 64 encoded encrypted data 
provided.

```go
swanPairs := connection.Decrypt(encrypted)
```

Example result in JSON format prior to conversion to SWAN pairs.

```json
[
    {
        "Key": "pref",
        "Created": "2021-05-10T00:00:00Z",
        "Expires": "2021-08-05T00:00:00Z",
        "Value": "AmNtcC...m23avB"
    },
    {
        "Key": "sid",
        "Created": "2021-05-07T00:00:00Z",
        "Expires": "2021-08-05T00:00:00Z",
        "Value": "AjUxZG...2an7jM"
    },
    {
        "Key": "stop",
        "Created": "2021-05-10T00:00:00Z",
        "Expires": "2086-08-03T00:00:00Z",
        "Value": "cool-creams.uk cool-bikes.uk"
    },
    {
        "Key": "rid",
        "Created": "2021-05-10T00:00:00Z",
        "Expires": "2021-08-05T00:00:00Z",
        "Value": "AjUxZGI...xjtRBQ"
    },
    {
        "Key": "val",
        "Created": "2021-05-10T09:15:42.7843197Z",
        "Expires": "2086-08-03T00:00:00Z",
        "Value": "2021-05-10T09:31:42Z"
    }
]
```

The returned keys map to the SWAN data. 

The keys `rid`, `sid` and `pref` have related OWIDs. 

The key `val` contains the time when the caller should revalidate the SWAN data
with SWAN via a call to Fetch. It is possible another tab in the same web
browser has been used to update the SWAN data and the current domain will not be
aware of these changes until it validates the data is still current.

### DecryptRaw

Returns the decrypted raw SWAN data as a map of string keys to values from the 
base 64 encoded encrypted data provided. Must only be used by User Interface 
Providers to update SWAN data.

```go
raw := connection.DecryptRaw(encrypted)
```

Example result.

```json
{
    "backgroundColor": "#f5f5f5",
    "email": "test@test.com",
    "message": "Hang tight. We're getting things ready.",
    "messageColor": "darkslategray",
    "pref": "off",
    "progressColor": "darkgreen",
    "salt": "qqo",
    "state": [
        "Example state"
    ],
    "rid": "AjUxZGIudWsAgdMK...TaK/AWD4tDXxjtRBQ",
    "title": "SWAN Demo"
}
```

### CreateRID

Returns a new RID with an OWID from from the SWAN Operator. Only SWAN operators
can create RIDs.

```go
rid := connection.CreateRID()
```

### HomeNode

Returns the domain of the home node for the web browser associated with the 
request.

```go
homeNode := connection.HomeNode(request)
```