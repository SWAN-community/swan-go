# swan-go

Secure Web Addressability Network (SWAN) - an open source secure and privacy
supporting cross domain identity network implemented in Go

## Introduction

This project contains all the server side components that are needed to
implement SWAN in Go for a publisher, Consent Management Platform (CMP) or
OpenRTB party.

## Connection

An instance of a swan.Connection is created via a call to swan.NewConnection. 
The following code shows how this would be achieved where an instance of a 
struct d contains members that relate to swan.Operation members.

```
connection := swan.NewConnection(swan.Operation{
    Client: swan.Client{
        SWAN: swan.SWAN{
            AccessKey: d.SwanAccessKey,
            Operator:  d.SwanAccessNode,
            Scheme:    d.SwanScheme}},
    BackgroundColor:       d.SwanBackgroundColor,
    Message:               d.SwanMessage,
    MessageColor:          d.SwanMessageColor,
    ProgressColor:         d.SwanProgressColor,
    NodeCount:             d.SwanNodeCount,
    DisplayUserInterface:  d.SwanDisplayUserInterface,
    PostMessageOnComplete: d.SwanPostMessage,
    JavaScript:            d.SwanJavaScript,
    UseHomeNode:           d.SwanUseHomeNode})
```

## Operations

Once the connection is created with the defaults to be used for all requests
to SWAN the following operations are supported.

The following variables have the following meaning and type.

| Name | Type | Description |
|-|-|-|
| encrypted | string | base 64 encoded encrypted data returned from a SWAN storage operation |
| request | *http.Request | http request associated with a user's web browser - used to set the home node |
| returnUrl | *url.URL | the URL to return to after the browser has been directed to the URL provided by the GetURL function |

### Fetch

Provides a URL that the browser should be immediately directed to. The return
URL will have the encrypted SWAN data appended ready to be used with the decrypt
functions.

```
url := connection.NewFetch(request, returnUrl).GetURL()
```

### Update

Provides a URL that the browser should be immediately directed to. The return
URL will have the encrypted SWAN data appended ready to be used with the decrypt
functions.

The members Pref, Email, Salt and SWID should be set to the values provided by 
the user before the GetURL function is called. If they are left blank the 
existing values are removed from SWAN.

```
u := connection.NewUpdate(request, returnUrl)

// Set the raw SWAN data from the form associated with the request.
u.Pref = r.Form.Get("pref") == "on"
u.Email = r.Form.Get("email")
u.Salt = []byte(r.Form.Get("salt"))
u.SWID = r.Form.Get("swid")

url := u.GetURL()
```

### Stop

Provides a URL that the browser should be immediately directed to. The return
URL will have the encrypted SWAN data appended ready to be used with the decrypt
functions.

The host parameter is the host domain associated with the advert that should be
stopped.

```
host := r.Form.Get("host")
url := connection.NewSWANStop(r, returnUrl, host).GetURL()
```

### Decrypt

Returns the decrypted SWAN data from the base 64 encoded encrypted data 
provided.

```
swanPairs := connection.Decrypt(encrypted)
```

### DecryptRaw

Returns the decrypted raw SWAN data from the base 64 encoded encrypted data 
provided. Must only be used by User Interface Providers.

```
rawPairs := connection.DecryptRaw(encrypted)
```

### CreateSWID

Returns a new SWID in OWID from from the SWAN Operator. 

```
swid := connection.CreateSWID()
```

### HomeNode

Returns the domain of the home node for the web browser associated with the 
request.

```
homeNode := connection.HomeNode(request)
```

