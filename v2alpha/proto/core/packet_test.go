// Copyright (c) Outernet Council and Contributors.
//
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

package core_test

import (
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"outernetcouncil.org/nmts/v2alpha/proto/core"
)

func TestTextParsingExamples(t *testing.T) {
	testCases := []struct {
		name string
		text string
	}{
		{
			name: "MPLS label stack",
			text: `
field: { mpls { label: 1 } }
field: { mpls { label: 2 } }
field: { mpls { label: 3 s:true } }`,
		},
		{
			name: "Seven layer bean dip, with IPv6 flowlabel for entropy",
			text: `
field: { eth_type: { eth: ETH_CTAG } }
field: { vlan { vid { u12: 100 } } }
field: { eth_type: { eth: ETH_MPLS } }
field: { mpls { label: 1 } }
field: { mpls { label: 2 } }
field: { mpls { label: 3 } }
field: { ipv6: { src: { str: "2001:db8:1:1::1" }
                 dst: { str: "2001:db8:2:2::2/128" }
                 flow_label: { u20: 17 }
                 ip_proto: { well_known: IPP_RH6 }
               } }
field: { srh: { sid: { str: "5f00:1:1:1::" }
                sid: { str: "5f00:2:2:2::" }
                ip_proto: { well_known: IPP_UDP }
              } }
# See RFC 7510, MPLS-in-UDP
field: { udp: { dst_port: { u16: 6635 } } }
field: { mpls { label: 100 } }
field: { mpls { label: 200 s:true } }
field: { ipv4: { src: { str: "192.0.2.1/32" }
                 dst: { str: "192.0.2.2" }
                 ip_proto: { well_known: IPP_GRE }
               } }
field: { gre: { proto: { eth: ETH_ARP }}}
`,
		},
	}

	for _, testCase := range testCases {
		parsed := new(core.PacketDescription)
		err := prototext.Unmarshal([]byte(testCase.text), parsed)
		if err != nil {
			t.Fatalf("%s: %q", testCase.name, err)
		}
	}
}
