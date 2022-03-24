/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestValidateNameserverIpAndPort(t *testing.T) {
	for _, tc := range []struct {
		wantErr bool
		ns      string
		ip      string
		port    string
	}{
		{wantErr: true},
		{ns: "1.2.3.4", ip: "1.2.3.4", port: "53"},
		{ns: "1.2.3.4:53", ip: "1.2.3.4", port: "53"},
		{wantErr: true, ns: "1.1.1.1:abc"},
		{wantErr: true, ns: "1.1.1.1:123456789"},
		{wantErr: true, ns: "1.1.1.1:"},
		{wantErr: true, ns: "invalidip"},
		{wantErr: true, ns: "invalidip:80"},
	} {
		ip, port, err := ValidateNameserverIpAndPort(tc.ns)
		gotErr := err != nil
		if gotErr != tc.wantErr {
			t.Errorf("ValidateNameserverIpAndPort(%q) = %q, %q, %v; gotErr = %t, want %t", tc.ns, ip, port, err, gotErr, tc.wantErr)
		}
		if ip != tc.ip || port != tc.port {
			t.Errorf("ValidateNameserverIpAndPort(%q) = %q, %q, nil; want %q, %q, nil", tc.ns, ip, port, tc.ip, tc.port)
		}
	}
}

func TestExtractIP(t *testing.T) {
	for _, tc := range []struct {
		testName string
		ptr      string
		wantIP   string
		wantOk   bool
	}{
		{
			testName: "valid IPv4 ptr",
			ptr:      "255.2.0.192.in-addr.arpa.",
			wantIP:   "192.0.2.255",
			wantOk:   true,
		},
		{
			testName: "valid IPv6 ptr",
			ptr:      "b.a.9.8.7.6.5.0.4.0.0.0.3.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.1.2.3.4.ip6.arpa.",
			wantIP:   "4321::2:3:4:567:89ab",
			wantOk:   true,
		},
		{
			testName: "valid IPv6 ptr has :0: instead of ::",
			ptr:      "b.a.9.8.7.6.5.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.0.0.0.0.1.2.3.4.ip6.arpa.",
			wantIP:   "4321:0:1:2:3:4:567:89ab",
			wantOk:   true,
		},
		{
			testName: "empty ptr",
			wantOk:   false,
		},
		{
			testName: "IPv4 ptr with incorrect suffix",
			ptr:      "255.2.0.192.ip6.arpa.",
			wantOk:   false,
		},
		{
			testName: "IPv6 ptr with incorrect suffix",
			ptr:      "b.a.9.8.7.6.5.0.4.0.0.0.3.0.0.0.2.0.0.0.0.0.0.0.0.0.0.0.1.2.3.4.in-addr.arpa",
			wantOk:   false,
		},
		{
			testName: "large number of nibbles in ipv6 ptr",
			ptr:      "a.b.a.9.8.7.6.5.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.0.0.0.0.1.2.3.4.ip6.arpa.",
			wantOk:   false,
		},
		{
			testName: "unexpected char",
			ptr:      "z.a.9.8.7.6.5.0.4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.0.0.0.0.1.2.3.4.ip6.arpa.",
			wantOk:   false,
		},
		{
			testName: "custom text",
			ptr:      "custom text",
			wantOk:   false,
		},
	} {
		ip, ok := ExtractIP(tc.ptr)
		if ok != tc.wantOk || ip != tc.wantIP {
			t.Errorf("Test %q: ExtractIP(%q) = %q, %t; want %q, %t", tc.testName, tc.ptr, ip, ok, tc.wantIP, tc.wantOk)
		}
	}
}

func TestGetClusterIPs(t *testing.T) {
	for _, tc := range []struct {
		service *v1.Service
		wantIPs []string
	}{
		{
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					ClusterIP:  "2001:db8:0:0:aaaa::1",
					ClusterIPs: []string{"2001:db8:0:0:aaaa::1"},
				},
			},
			wantIPs: []string{"2001:db8::aaaa:0:0:1"},
		},
		{
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					ClusterIP: "2001:db8::aaaa:0:0:1",
				},
			},
			wantIPs: []string{"2001:db8::aaaa:0:0:1"},
		},
		{
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					ClusterIP:  "2001:db8:0::aaaa:0:0:1",
					ClusterIPs: []string{"2001:db8:0::aaaa:0:0:1", "255.255.255.0"},
				},
			},
			wantIPs: []string{"2001:db8::aaaa:0:0:1", "255.255.255.0"},
		},
	} {
		assert.ElementsMatch(t, tc.wantIPs, GetClusterIPs(tc.service))
	}
}
