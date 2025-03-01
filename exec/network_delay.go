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
	"context"
	"path"
	"fmt"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
	"github.com/chaosblade-io/chaosblade-spec-go/util"
)

type DelayActionSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewDelayActionSpec() spec.ExpActionCommandSpec {
	return &DelayActionSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: commFlags,
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:     "time",
					Desc:     "Delay time, ms",
					Required: true,
				},
				&spec.ExpFlag{
					Name: "offset",
					Desc: "Delay offset time, ms",
				},
			},
			ActionExecutor: &NetworkDelayExecutor{},
		},
	}
}

func (*DelayActionSpec) Name() string {
	return "delay"
}

func (*DelayActionSpec) Aliases() []string {
	return []string{}
}

func (*DelayActionSpec) ShortDesc() string {
	return "Delay experiment"
}

func (*DelayActionSpec) LongDesc() string {
	return "Delay experiment"
}

type NetworkDelayExecutor struct {
	channel spec.Channel
}

func (de *NetworkDelayExecutor) Name() string {
	return "delay"
}

func (de *NetworkDelayExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	if de.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}
	netInterface := model.ActionFlags["interface"]
	if netInterface == "" {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less interface parameter")
	}
	time := model.ActionFlags["time"]
	if time == "" {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less time flag")
	}
	offset := model.ActionFlags["offset"]
	if offset == "" {
		offset = "10"
	}
	localPort := model.ActionFlags["local-port"]
	remotePort := model.ActionFlags["remote-port"]
	excludePort := model.ActionFlags["exclude-port"]
	destIp := model.ActionFlags["destination-ip"]
	if _, ok := spec.IsDestroy(ctx); ok {
		return de.stop(netInterface, ctx)
	} else {
		return de.start(localPort, remotePort, excludePort, destIp, time, offset, netInterface, ctx)
	}
}

func (de *NetworkDelayExecutor) start(localPort, remotePort, excludePort, destIp, time, offset, netInterface string,
	ctx context.Context) *spec.Response {
	args := fmt.Sprintf("--start --interface %s --time %s --offset %s --debug=%t", netInterface, time, offset, util.Debug)
	args, err := getCommArgs(localPort, remotePort, excludePort, destIp, args)
	if err != nil {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], err.Error())
	}
	return de.channel.Run(ctx, path.Join(de.channel.GetScriptPath(), dlNetworkBin), args)
}

func (de *NetworkDelayExecutor) stop(netInterface string, ctx context.Context) *spec.Response {
	return de.channel.Run(ctx, path.Join(de.channel.GetScriptPath(), dlNetworkBin),
		fmt.Sprintf("--stop --interface %s --debug=%t", netInterface, util.Debug))
}

func (de *NetworkDelayExecutor) SetChannel(channel spec.Channel) {
	de.channel = channel
}
