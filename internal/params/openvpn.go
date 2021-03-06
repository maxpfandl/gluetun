package params

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/models"
	libparams "github.com/qdm12/golibs/params"
)

// GetUser obtains the user to use to connect to the VPN servers.
// It first tries to use the OPENVPN_USER environment variable (easier for the end user)
// and then tries to read from the secret file openvpn_user if nothing was found.
func (r *reader) GetUser() (user string, err error) {
	const compulsory = true
	return r.getFromEnvOrSecretFile("OPENVPN_USER", compulsory, []string{"USER"})
}

// GetPassword obtains the password to use to connect to the VPN servers.
// It first tries to use the OPENVPN_PASSWORD environment variable (easier for the end user)
// and then tries to read from the secret file openvpn_password if nothing was found.
func (r *reader) GetPassword() (s string, err error) {
	const compulsory = true
	return r.getFromEnvOrSecretFile("OPENVPN_PASSWORD", compulsory, []string{"PASSWORD"})
}

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL.
func (r *reader) GetNetworkProtocol() (protocol models.NetworkProtocol, err error) {
	s, err := r.env.Inside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return models.NetworkProtocol(s), err
}

// GetOpenVPNVerbosity obtains the verbosity level for verbosity between 0 and 6
// from the environment variable OPENVPN_VERBOSITY.
func (r *reader) GetOpenVPNVerbosity() (verbosity int, err error) {
	return r.env.IntRange("OPENVPN_VERBOSITY", 0, 6, libparams.Default("1"))
}

// GetOpenVPNRoot obtains if openvpn should be run as root
// from the environment variable OPENVPN_ROOT.
func (r *reader) GetOpenVPNRoot() (root bool, err error) {
	return r.env.YesNo("OPENVPN_ROOT", libparams.Default("no"))
}

// GetTargetIP obtains the IP address to override over the list of IP addresses filtered
// from the environment variable OPENVPN_TARGET_IP.
func (r *reader) GetTargetIP() (ip net.IP, err error) {
	s, err := r.env.Get("OPENVPN_TARGET_IP")
	if len(s) == 0 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("target IP address %q is not valid", s)
	}
	return ip, nil
}

// GetOpenVPNCipher obtains a custom cipher to use with OpenVPN
// from the environment variable OPENVPN_CIPHER.
func (r *reader) GetOpenVPNCipher() (cipher string, err error) {
	return r.env.Get("OPENVPN_CIPHER")
}

// GetOpenVPNAuth obtains a custom auth algorithm to use with OpenVPN
// from the environment variable OPENVPN_AUTH.
func (r *reader) GetOpenVPNAuth() (auth string, err error) {
	return r.env.Get("OPENVPN_AUTH")
}

// GetOpenVPNIPv6 obtains if ipv6 should be tunneled through the
// openvpn tunnel from the environment variable OPENVPN_IPV6.
func (r *reader) GetOpenVPNIPv6() (ipv6 bool, err error) {
	return r.env.OnOff("OPENVPN_IPV6", libparams.Default("off"))
}

func (r *reader) GetOpenVPNMSSFix() (mssFix uint16, err error) {
	n, err := r.env.IntRange("OPENVPN_MSSFIX", 0, 10000, libparams.Default("0"))
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
