package utils

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache is a global cache instance that can be used across your application.
var Cache *gocache.Cache

// DefaultExpiration holds the default expiration time for cached items.
var DefaultExpiration = gocache.DefaultExpiration

func init() {
	// Initialize the cache with a default expiration of 5 minutes and a cleanup interval of 10 minutes.
	Cache = gocache.New(5*time.Minute, 10*time.Minute)
}
