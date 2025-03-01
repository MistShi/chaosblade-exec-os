/*
 * Copyright 1999-2019 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package exec

import (
	"fmt"
	"strings"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type NetworkCommandSpec struct {
	spec.BaseExpModelCommandSpec
}

func NewNetworkCommandSpec() spec.ExpModelCommandSpec {
	return &NetworkCommandSpec{
		spec.BaseExpModelCommandSpec{
			ExpActions: []spec.ExpActionCommandSpec{
				NewDelayActionSpec(),
				NewDropActionSpec(),
				NewDnsActionSpec(),
				NewLossActionSpec(),
			},
			ExpFlags: []spec.ExpFlagSpec{},
		},
	}
}

func (*NetworkCommandSpec) Name() string {
	return "network"
}

func (*NetworkCommandSpec) ShortDesc() string {
	return "Network experiment"
}

func (*NetworkCommandSpec) LongDesc() string {
	return "Network experiment"
}

func (*NetworkCommandSpec) Example() string {
	return `network delay --interface eth0 --time 3000

# You can execute "blade query network interface" command to query the interfaces`
}

// dlNetworkBin for delay and loss experiments
var dlNetworkBin = "chaos_dlnetwork"

var commFlags = []spec.ExpFlagSpec{
	&spec.ExpFlag{
		Name: "local-port",
		Desc: "Ports for local service. Support for configuring multiple ports, separated by commas or connector representing ranges, for example: 80,8000-8080",
	},
	&spec.ExpFlag{
		Name: "remote-port",
		Desc: "Ports for remote service. Support for configuring multiple ports, separated by commas or connector representing ranges, for example: 80,8000-8080",
	},
	&spec.ExpFlag{
		Name: "exclude-port",
		Desc: "Exclude local ports. Support for configuring multiple ports, separated by commas or connector representing ranges, for example: 22,8000. This flag is invalid when --local-port or --remote-port is specified",
	},
	&spec.ExpFlag{
		Name: "destination-ip",
		Desc: "destination ip. Support for using mask to specify the ip range, for example, 192.168.1.0/24. You can also use 192.168.1.1 or 192.168.1.1/32 to specify it.",
	},
	&spec.ExpFlag{
		Name:     "interface",
		Desc:     "Network interface, for example, eth0",
		Required: true,
	},
}

func getCommArgs(localPort, remotePort, excludePort, destinationIp string, args string) (string, error) {
	if localPort != "" {
		localPorts, err := util.ParseIntegerListToStringSlice(localPort)
		if err != nil {
			return "", err
		}
		args = fmt.Sprintf("%s --local-port %s", args, strings.Join(localPorts, ","))
	}
	if remotePort != "" {
		remotePorts, err := util.ParseIntegerListToStringSlice(remotePort)
		if err != nil {
			return "", err
		}
		args = fmt.Sprintf("%s --remote-port %s", args, strings.Join(remotePorts, ","))
	}
	if excludePort != "" {
		excludePorts, err := util.ParseIntegerListToStringSlice(excludePort)
		if err != nil {
			return "", err
		}
		args = fmt.Sprintf("%s --exclude-port %s", args, strings.Join(excludePorts, ","))
	}
	if destinationIp != "" {
		args = fmt.Sprintf("%s --destination-ip %s", args, destinationIp)
	}
	return args, nil
}
