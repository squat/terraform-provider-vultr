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

// validateIPAddress ensures that the string value is a valid IPv4 or IPv6
// address and returns an error otherwise.
func validateIPAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ip := net.ParseIP(value)
	if ip == nil {
		errors = append(errors, fmt.Errorf("%q must contain a valid IP address", k))
		return
	}
	return
}

// validateReservedIPType ensures that the string value is either "v4" or "v6".
func validateReservedIPType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "v6" && value != "v4" {
		errors = append(errors, fmt.Errorf("%q must be either 'v4' or 'v6'", k))
		return
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

// validateStartupScriptType ensures that the string value is a valid
// startup script type and returns an error otherwise.
func validateStartupScriptType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validProtocols := map[string]struct{}{
		"boot": {},
		"pxe":  {},
	}
	if _, ok := validProtocols[value]; !ok {
		errors = append(errors, fmt.Errorf("%q contains an invalid startup script type %q; valid types are: %q and %q", k, value, "boot", "pxe"))
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
