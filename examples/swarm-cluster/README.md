#### Vultr credentials

Set an environment variable containing the Vultr API key:

    export VULTR_API_KEY=<your-vultr-api-key>

    or

Save vultr api key in a file `~/.creds/vultr_api_token`

Note: as an alternative, the API key can be specified in configuration as shown below.

    // Configure the Vultr provider.

    provider "vultr" {
      api_key = "<your-vultr-api-key>"
    }

#### `jq` is needed on your local computer.

    sudo apt install jq

or

    brew install jq



the same on digital ocean

https://knpw.rs/blog/docker-swarm-terraform
