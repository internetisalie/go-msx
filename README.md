# go-msx

go-msx is a Go library for microservices and tools interacting with MSX. 

## Versioning

Currently this library and tools are in a pre-alpha state.  They are subject to
backwards-incompatible changes at *any* time.  After reaching the first stable release (v1.0.0),
[SemVer](https://semver.org/) will be used per industry and golang best practices.     

## Requirements

- Go 1.16+

    - Ensure your GOPATH is correctly set and referenced in your PATH.  For example:
        ```bash
        export GOPATH=~/go
        export PATH=$PATH:$GOPATH/bin
        ```

    - Be sure to set your Go proxy settings correctly.  For example:
        ```bash
        go env -w GOPRIVATE=cto-github.cisco.com/NFV-BU
        ```

- Git SSH configuration for `cto-github.cisco.com`

    - Ensure you have a registered SSH key referenced in your `~/.ssh/config`:
    
        ```
        Host cto-github.cisco.com
              HostName cto-github.cisco.com
              User git
              IdentityFile ~/.ssh/github.key
        ```
      
      Note that this key must be registered via the [Github UI](https://cto-github.cisco.com/settings/keys).

    - Ensure you have SSH protocol override for git HTTPS urls to our github in your `~/.gitconfig`:
    
      ```
      [url "ssh://git@cto-github.cisco.com/"]
              insteadOf = https://cto-github.cisco.com/
      ```

- Skel tool for code generation

    - Check out go-msx into your local workspace:
        
        ```bash
        mkdir -p $HOME/msx && cd $HOME/msx
        git clone git@cto-github.cisco.com:NFV-BU/go-msx.git
        cd go-msx
        go mod download
        ```

    - Install `skel`:
  
        ```bash
        make install-skel
        ```

## Quick Start

- To continue working on an existing go-msx project:

    - Return to the original project README instructions
      and continue.

- To add go-msx to an existing module-enabled go project:

    ```bash
    go get -u cto-github.cisco.com/NFV-BU/go-msx
    ```

- To create a new go-msx microservice skeleton project:
    
    ```bash
    cd $HOME/msx
    skel
    ```

## Contributing

- Ensure you create a meaningfully named topic branch for your code:

    `feature/sql-transactions`
    
    `bugfix/populate-error-handling`
    
- Make your code changes

- Run `make precommit` to regenerate and reformat.  You will likely need to
  install the `staticfiles` package the first time:
  
    `go get bou.ke/staticfiles`

- Commit your code to your topic branch

- Rebase your topic branch onto master (do not reverse merge master into your branch)

- Ensure your commits are cohesive, or just squash them

- Create a Pull Request with a meaningful title similar to your topic branch name

## Documentation

### Cross-Cutting Concerns
* [Logging](log/README.md)
* [Configuration](config/README.md)
* [Lifecycle](app/README.md)
* [Dependencies](app/context.md)
* [Stats](stats/README.md)
* [Tracing](trace/README.md)

### Application Components
* Web Service
    * [Controller](webservice/controller.md)
    * [Filter](#)

* Persistence
  
    * [Repository](sqldb/repository.md)
    * [Migration](#)
    
* Communication
    * [Integration](#)
    * [Streaming](#)

### Utilities

* [Audit Events](#)
* [Auditable Models](#)
* [Cache](cache/lru/README.md)
* [Certificates and TLS](certificate/README.md)
* [Executing Commands](#)
* [Health Checks](#)
* [Http Client](#)
* [Leader Election](#)
* [Pagination](#)
* [Resources](resource/README.md)
* [Retry](#)
* [Sanitization](sanitize/README.md)
* [Scheduled Tasks](scheduled/README.md)
* [Transit Encryption](transit/README.md)
* [Validation](#)

## License

Copyright (C) 2020 Cisco Systems, Inc.  All rights reserved.
