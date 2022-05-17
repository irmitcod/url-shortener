# Url-Shortener
Let's design a URL shortening service like TinyURL. This service will provide short aliases redirecting to long URLs.


> **WARNING** \
> This project is completely unstable, and it's under the development and may change alot

## Features
- [x] ✔️ Echo
- [x] ✔️ Mongodb
- [x] ✔️ Redis
- [x] ✔️ Redis cache
- [x] ✔️ LRU cache
- [x] ✔️ LFU cache
- [x] ✔️ Clean architecture
- [ ] ❌     Mongodb sharding
- [ ] ❌      ...

## We'll cover the following
- [x] ✔️ System APIs
- [x] ✔️ Database design (MongoDB)
- [x] ✔️ Encoding url with base58
- [x] ✔️ Cache with LRU and Redis
- [x] ✔️ Load Balancer with nginx to server
- [x] ✔️ Telemetry
- [ ] ❌ Load Balancer for MongoDB
- [ ] ❌  Load Balancer for Redis
- [ ] ❌  Load Balancer for Redis
## Table of Contents

## Built-in Components


## Quick Start
* This project using `docker-compose` as package manager.
* Run `docker-compose up --build` 

## API request
* Create User POST http://localhost/user ->      params {Username:test,Password:test,Name:test}
* Login request POST  http://localhost/login -> params {Username:test,Password:test}
* set url POST  http://localhost/UrlShotener -> params {OriginalURL:http://test.com}
* redirect url GET  http://localhost/{key} -> e.g  http://localhost/hgu4sdss
