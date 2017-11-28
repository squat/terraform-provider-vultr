package vultr

import (
	"log"
	"time"

	"github.com/JamesClonk/vultr/lib"
)

// Config is the configuration structure used to instantiate the Vultr
// provider.
type Config struct {
	APIKey string
}

// Client wraps a JamesClonk/vultr/lib.
type Client struct {
	*lib.Client
}

// Client configures and returns a fully initialized Vultr Client.
func (c *Config) Client() (interface{}, error) {
	client := Client{lib.NewClient(c.APIKey, &lib.Options{RateLimitation: 500 * time.Millisecond})}
	log.Printf("[INFO] Vultr Client configured for URL: %s", client.Endpoint)
	return &client, nil
}
