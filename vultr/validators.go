package vultr

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	// stringSlashIntFormatErrTemplate is the default error template for parsing a string/int resource ID.
	stringSlashIntFormatErrTemplate = "Error parsing %s: should be of form <%s>/<%s>, where <%s> is an integer; got %q"
	// stringSlashStringFormatErrTemplate is the default error template for parsing a string/string resource ID.
	stringSlashStringFormatErrTemplate = "Error parsing %s: should be of form <%s>/<%s>; got %q"
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

// parseStringSlashInt parses a commond ID format <string>/<int> into its components and returns
// the provided error otherwise.
func parseStringSlashInt(id, idType, first, second string) (string, int, error) {
	baseErr := fmt.Errorf(stringSlashIntFormatErrTemplate, idType, first, second, second, id)
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		return "", 0, baseErr
	}
	s := idParts[0]
	i, err := strconv.Atoi(idParts[1])
	if err != nil {
		return "", 0, fmt.Errorf("%v: %v", baseErr, err)
	}
	return s, i, nil
}

// parseStringSlashString parses a commond ID format <string>/<string> into its components and returns
// the provided error otherwise.
func parseStringSlashString(id, idType, first, second string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf(stringSlashStringFormatErrTemplate, idType, first, second, id)
	}
	return parts[0], parts[1], nil
}
