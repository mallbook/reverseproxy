# reverseproxy

## Installation
```
go get -u github.com/mallbook/reverseproxy
```

## Quick Start
Support for reverse proxy services in the gorestful service, just import the package `github.com/mallbook/reverseproxy` and write a configuration file, as follows:  

### 1. Import package  
```go
import (
    _ "github.com/mallbook/reverseproxy"
)
```

### 2. The configuration file  
The configuration file is `etc/conf/reverse_proxy.json`  
  
The file contents are as follows:
```json
{
    "reverseProxy": {
        "tokens": {
            "rootPath": "/iamservice/v1/auth/tokens",
            "targetPath": "/v3/auth/tokens",
            "proxyPass": ["https://127.0.0.1:5000"],
            "routes": [
                {
                    "subPath": "/{tokenID}",
                    "httpMethod": "GET"
                }
            ]
        }
    }
}
```
|Key|Description|
|:--|:--|
|reverseProxy|Fixed key, represent reverse proxy|
|rootPath|The root path of the reverse proxy|
|targetPath|The target path of the reverse proxy|
|proxyPass|The all backends of the reverse proxy|
|routes|The all routes of the reverse proxy|
|subPath|The sub path of the reverse proxy|
|httpMethod|The http method of the reverse proxy, must use all uppercase, such as `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`|
