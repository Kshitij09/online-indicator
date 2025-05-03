## Design Document - v0

At this stage, the goal was to have minimal working components that build an
end-to-end working model of online indicator. The system consists of three
components:
* backend
* pulsesim
* dashboard

![architecture](v0-architecture.svg)

### backend

* An HTTP server written in Golang
* Uses in-memory storage for user and session details, protected by RWMutex
* Follows a domain-driven architecture with these layers:
  * domain - models, interfaces, and business logic
  * inmem - in-memory implementation of domain's DAO interfaces
  * transport - DTOs, minimal wrapper on net/http module
    to facilitate better handlers and middlewares

### pulsesim

* A Golang HTTP client that simulates user activity
* Registers configurable number of fake users (`-n`) and creates active sessions
* Launches goroutines to periodically ping the backend, creating a mix of online/offline statuses

### dashboard

* A web frontend built with Svelte + Vite
* Displays users with their online/offline status
* Features include name search and reactive statistics (total/online/offline users)
* Polls the backend's batch API every 5 seconds for status updates
