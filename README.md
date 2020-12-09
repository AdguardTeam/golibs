[![Code Coverage](https://img.shields.io/codecov/c/github/AdguardTeam/golibs/master.svg)](https://codecov.io/github/AdguardTeam/golibs?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdguardTeam/golibs)](https://goreportcard.com/report/AdguardTeam/golibs)
[![GolangCI](https://golangci.com/badges/github.com/AdguardTeam/golibs.svg)](https://golangci.com/r/github.com/AdguardTeam/golibs)
[![Go Doc](https://godoc.org/github.com/AdguardTeam/golibs?status.svg)](https://godoc.org/github.com/AdguardTeam/golibs)

# golibs

This repository contains several useful functions and interfaces for Go:

* Cache - in-memory cache with LRU, limits and statistics
* Log - logger with configurable log-level on top of standard "log"
* File:
    * safe file writing
* JSON:
    * JSON format helper functions
* Utils:
    * hostname validator


## Cache

A quick example:

    conf := cache.Config{}
	conf.EnableLRU = true
	c := cache.New(conf)
    c.Set([]byte("key"), []byte("value"))
    val := c.Get([]byte("key"))
