// Configure the Vultr provider. 
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Create a DNS domain.
resource "vultr_dns_domain" "example" {
  domain = "example.com"
  ip     = "10.0.0.1"
}

// Create a new DNS record.
resource "vultr_dns_record" "example_web" {
  domain = "${vultr_dns_domain.example.id}"
  name   = "www"
  type   = "A"
  data   = "${vultr_dns_domain.example.ip}"
  ttl    = 300
}
