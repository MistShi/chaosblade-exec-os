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
)

type LossActionSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewLossActionSpec() spec.ExpActionCommandSpec {
	return &LossActionSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: commFlags,
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:     "percent",
					Desc:     "loss percent, [0, 100]",
					Required: true,
				},
			},
			ActionExecutor: &NetworkLossExecutor{},
		},
	}
}

func (*LossActionSpec) Name() string {
	return "loss"
}

func (*LossActionSpec) Aliases() []string {
	return []string{}
}

func (*LossActionSpec) ShortDesc() string {
	return "Loss network package"
}

func (*LossActionSpec) LongDesc() string {
	return "Loss network package"
}

type NetworkLossExecutor struct {
	channel spec.Channel
}

func (*NetworkLossExecutor) Name() string {
	return "loss"
}

func (nle *NetworkLossExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	if nle.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}
	var dev = ""
	if netInterface, ok := model.ActionFlags["interface"]; ok {
		if netInterface == "" {
			return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less interface flag")
		}
		dev = netInterface
	}
	percent := model.ActionFlags["percent"]
	if percent == "" {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less percent flag")
	}
	if _, ok := spec.IsDestroy(ctx); ok {
		return nle.stop(dev, ctx)
	}
	localPort := model.ActionFlags["local-port"]
	remotePort := model.ActionFlags["remote-port"]
	excludePort := model.ActionFlags["exclude-port"]
	destIp := model.ActionFlags["destination-ip"]
	return nle.start(dev, localPort, remotePort, excludePort, destIp, percent, ctx)
}

func (nle *NetworkLossExecutor) start(netInterface, localPort, remotePort, excludePort, destIp, percent string,
	ctx context.Context) *spec.Response {
	args := fmt.Sprintf("--start --interface %s --percent %s", netInterface, percent)
	args, err := getCommArgs(localPort, remotePort, excludePort, destIp, args)
	if err != nil {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], err.Error())
	}
	return nle.channel.Run(ctx, path.Join(nle.channel.GetScriptPath(), dlNetworkBin), args)
}

func (nle *NetworkLossExecutor) stop(netInterface string, ctx context.Context) *spec.Response {
	return nle.channel.Run(ctx, path.Join(nle.channel.GetScriptPath(), dlNetworkBin),
		fmt.Sprintf("--stop --interface %s", netInterface))
}

func (nle *NetworkLossExecutor) SetChannel(channel spec.Channel) {
	nle.channel = channel
}
