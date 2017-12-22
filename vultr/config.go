package vultr

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/JamesClonk/vultr/lib"
	"github.com/hashicorp/terraform/helper/logging"
)

const logReqMsg = `Vultr API Request Details:
---[ REQUEST ]---------------------------------------
%s
-----------------------------------------------------`

const logRespMsg = `Vultr API Response Details:
---[ RESPONSE ]--------------------------------------
%s
-----------------------------------------------------`

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

	if logging.IsDebugOrHigher() {
		client.OnRequestCompleted(logRequestAndResponse)
	}

	log.Printf("[INFO] Vultr Client configured for URL: %s", client.Endpoint)

	return &client, nil
}

func logRequestAndResponse(req *http.Request, resp *http.Response) {
	reqData, err := httputil.DumpRequest(req, true)
	if err == nil {
		log.Printf("[DEBUG] "+logReqMsg, string(reqData))
	} else {
		log.Printf("[ERROR] Vultr API Request error: %#v", err)
	}

	respData, err := httputil.DumpResponse(resp, true)
	if err == nil {
		log.Printf("[DEBUG] "+logRespMsg, string(respData))
	} else {
		log.Printf("[ERROR] Vultr API Response error: %#v", err)
	}
}
