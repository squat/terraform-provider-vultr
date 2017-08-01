package vultr

import (
	"fmt"
	"net"
	"regexp"
)

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address and returns an error otherwise.
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must contain a valid CIDR, got error parsing: %v", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf("%q must contain a valid network CIDR, got %q", k, value))
	}
	return
}

// validateFirewallRuleProtocol ensures that the string value is a valid
// firewall rule protocol and returns an error otherwise.
func validateFirewallRuleProtocol(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validProtocols := map[string]struct{}{
		"gre":  {},
		"icmp": {},
		"tcp":  {},
		"udp":  {},
	}
	if _, ok := validProtocols[value]; !ok {
		errors = append(errors, fmt.Errorf("%q contains an invalid firewall rule protocol %q; valid types are: %q, %q, %q, and %q", k, value, "gre", "icmp", "tcp", "udp"))
	}
	return
}

// validateRegex ensures that the string is a valid regular expression.
func validateRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid regular expression: %v", k, err))
	}
	return
}
