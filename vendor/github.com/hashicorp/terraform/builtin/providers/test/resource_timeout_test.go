package test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceTimeout_create(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_timeout" "foo" {
	create_delay = "2s"
	timeouts {
		create = "1s"
	}
}
				`),
				ExpectError: regexp.MustCompile("timeout while creating resource"),
			},
		},
	})
}
func TestResourceTimeout_update(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_timeout" "foo" {
	update_delay = "1s"
	timeouts {
		update = "1s"
	}
}
				`),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_timeout" "foo" {
	update_delay = "2s"
	timeouts {
		update = "1s"
	}
}
				`),
				ExpectError: regexp.MustCompile("timeout while updating resource"),
			},
		},
	})
}

func TestResourceTimeout_read(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckResourceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_timeout" "foo" {
}
				`),
			},
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_timeout" "foo" {
	read_delay = "30m"
}
				`),
				ExpectError: regexp.MustCompile("timeout while reading resource"),
			},
			// we need to remove the read_delay so that the resource can be
			// destroyed in the final step, but expect an error here from the
			// pre-existing delay.
			resource.TestStep{
				Config: strings.TrimSpace(`
resource "test_resource_timeout" "foo" {
}
				`),
				ExpectError: regexp.MustCompile("timeout while reading resource"),
			},
		},
	})
}
