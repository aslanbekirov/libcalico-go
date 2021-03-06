// Copyright (c) 2017 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

import (
	"fmt"
	"strings"

	"github.com/projectcalico/libcalico-go/lib/net"

	log "github.com/Sirupsen/logrus"
)

// This file contains various name conversion methods that can be used to convert
// between Calico key types and resource names.

// IPToResourceName converts an IP address to a name used for a k8s resource.
func IPToResourceName(ip net.IP) string {
	name := strings.Replace(ip.String(), ".", "-", 3)
	name = strings.Replace(name, ":", "-", 7)

	log.WithFields(log.Fields{
		"Name": name,
		"IP":   ip.String(),
	}).Debug("Converting IP to resource name")

	return name
}

// ResourceNameToIP converts a name used for a k8s resource to an IP address.
func ResourceNameToIP(name string) (*net.IP, error) {
	ip := net.ParseIP(resourceNameToIPString(name))
	if ip == nil {
		return nil, fmt.Errorf("invalid resource name %s: does not follow Calico IP name format", name)
	}
	return ip, nil
}

// IPNetToResourceName converts the given IPNet into a name used for a k8s resource.
func IPNetToResourceName(net net.IPNet) string {
	name := strings.Replace(net.String(), ".", "-", 3)
	name = strings.Replace(name, ":", "-", 7)
	name = strings.Replace(name, "/", "-", 1)

	log.WithFields(log.Fields{
		"Name":  name,
		"IPNet": net.String(),
	}).Debug("Converting IPNet to resource name")

	return name
}

// ResourceNameToIPNet converts a name used for a k8s resource to an IPNet.
func ResourceNameToIPNet(name string) (*net.IPNet, error) {
	// The last dash should be replaced by a "/"
	idx := strings.LastIndex(name, "-")
	if idx == -1 {
		return nil, fmt.Errorf("invalid resource name: %s: does not follow Calico IPNet name format", name)
	}
	ipstr := resourceNameToIPString(name[:idx])
	size := name[idx+1:]

	_, cidr, err := net.ParseCIDR(ipstr + "/" + size)
	if err != nil {
		return nil, fmt.Errorf("invalid resource name %s: does not follow Calico IPNet name format", name)
	}
	return cidr, nil
}

// resourceNameToIPString converts a name used for a k8s resource to an IP address string.
// This function does not check the validity of the result - it merely reverses the
// character conversion used to convert an IP address to a k8s compatible name.
func resourceNameToIPString(name string) string {
	// The IP address is stored in the name with periods and colons replaced
	// by dashes.  To determine if this is IPv4 or IPv6 count the dashes.  If
	// either of the following are true, it's IPv6:
	// -  There is a "--"
	// -  The number of "-" is greater than 3.
	var ipstr string
	if strings.Contains(name, "--") || strings.Count(name, "-") > 3 {
		// IPv6:  replace - with :
		ipstr = strings.Replace(name, "-", ":", 7)
	} else {
		// IPv4:  replace - with .
		ipstr = strings.Replace(name, "-", ".", 3)
	}

	log.WithFields(log.Fields{
		"Name": name,
		"IP":   ipstr,
	}).Debug("Converting resource name to IP String")
	return ipstr
}
