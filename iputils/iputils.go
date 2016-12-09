// Copyright (c) 2016 turbobytes
//
// This file is part of geoipdb, a library of GeoIP related helper functions
// for TurboBytes stack.
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package iputils

import (
	"net"
)

func init() {
	// Make sure nonGlobalIPv*CIDRs are parseable.
	cidrs := append(nonGlobalIPv4CIDRs, nonGlobalIPv6CIDRs...)
	for _, cidr := range cidrs {
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("unparseable CIDR '" + cidr + "': " + err.Error())
		}
	}
}

// nonGlobalIPv4CIDRs contains IANA IPv4 Special-Purpose Address Registry,
// where 'Global' flag is false.
//
// http://www.iana.org/assignments/iana-ipv4-special-registry/
var nonGlobalIPv4CIDRs = []string{
	"127.0.0.0/8",        // Loopback, RFC1122
	"192.168.0.0/16",     // Private-Use, RFC1918
	"10.0.0.0/8",         // Private-Use, RFC1918
	"172.16.0.0/12",      // Private-Use, RFC1918
	"0.0.0.0/8",          // "This host on this network", RFC1122 section 3.2.1.3
	"100.64.0.0/10",      // Shared Address Space, RFC6598
	"169.254.0.0/16",     // Link Local, RFC3927
	"192.0.0.0/24",       // IETF Protocol Assignments, RFC6890
	"192.0.2.0/24",       // Documentation (TEST-NET-1), RFC5737
	"198.18.0.0/15",      // Benchmarking, RFC2544
	"198.51.100.0/24",    // Documentation (TEST-NET-2), RFC5737
	"203.0.113.0/24",     // Documentation (TEST-NET-3), RFC5737
	"240.0.0.0/4",        // Reserved, RFC1112
	"255.255.255.255/32", // Limited Broadcast, RFC919
}

// nonGlobalIPv6CIDRs contains IANA IPv6 Special-Purpose Address Registry,
// where 'Global' flag is false.
//
// http://www.iana.org/assignments/iana-ipv6-special-registry/
var nonGlobalIPv6CIDRs = []string{
	"::1/128",       // Loopback Address, RFC4291
	"fc00::/7",      // Unique-Local, RFC4193
	"::ffff:0:0/96", // IPv4-mapped Address, RFC4291
	"2001::/23",     // IETF Protocol Assignments, RFC2928
	"fe80::/10",     // Linked-Scoped Unicast, RFC4291
	"2001:db8::/32", // Documentation, RFC3849
	"2001:2::/48",   // Benchmarking, RFC5180
	"2001::/32",     // TEREDO, RFC4380
	"100::/64",      // Discard-Only Address Block, RFC6666
	"::/128",        // Unspecified Address, RFC4291
}

// IsLocalIP tells if an IP address is not forwardable beyond a network.
//
// Returns if the given IP address is local.
func IsLocalIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	ip4 := ip.To4()
	if ip4 != nil {
		for _, cidr := range nonGlobalIPv4CIDRs {
			_, inet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			if inet.Contains(ip4) {
				return true
			}
		}
		return false
	}
	ip6 := ip.To16()
	if ip6 != nil {
		for _, cidr := range nonGlobalIPv6CIDRs {
			_, inet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			if inet.Contains(ip6) {
				return true
			}
		}
	}
	return false
}
