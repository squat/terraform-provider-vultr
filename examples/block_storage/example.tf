// Configure the Vultr provider.
// Alternatively, export the API key as an environment variable: `export VULTR_API_KEY=<your-vultr-api-key>`.
provider "vultr" {
  api_key = "<your-vultr-api-key>"
}

// Find the ID of the Silicon Valley region.
data "vultr_region" "has_block_storage" {
  filter {
    name   = "block_storage"
    values = ["true"]
  }
}

// Create block storage.
resource "vultr_block_storage" "example" {
  name      = "example"
  region_id = "${data.vultr_region.has_block_storage.id}"
  size      = 50
}
