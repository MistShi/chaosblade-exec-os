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
	"fmt"
	"path"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

type KillProcessActionCommandSpec struct {
	spec.BaseExpActionCommandSpec
}

func NewKillProcessActionCommandSpec() spec.ExpActionCommandSpec {
	return &KillProcessActionCommandSpec{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name: "process",
					Desc: "Process name",
				},
				&spec.ExpFlag{
					Name: "process-cmd",
					Desc: "Process name in command",
				},
			},
			ActionFlags:    []spec.ExpFlagSpec{},
			ActionExecutor: &KillProcessExecutor{},
		},
	}
}

func (*KillProcessActionCommandSpec) Name() string {
	return "kill"
}

func (*KillProcessActionCommandSpec) Aliases() []string {
	return []string{"k"}
}

func (*KillProcessActionCommandSpec) ShortDesc() string {
	return "Kill process"
}

func (*KillProcessActionCommandSpec) LongDesc() string {
	return "Kill process by process id or process name"
}

type KillProcessExecutor struct {
	channel spec.Channel
}

func (kpe *KillProcessExecutor) Name() string {
	return "kill"
}

var killProcessBin = "chaos_killprocess"

func (kpe *KillProcessExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	if kpe.channel == nil {
		return spec.ReturnFail(spec.Code[spec.ServerError], "channel is nil")
	}
	if _, ok := spec.IsDestroy(ctx); ok {
		return spec.ReturnSuccess(uid)
	}
	process := model.ActionFlags["process"]
	processCmd := model.ActionFlags["process-cmd"]
	if process == "" && processCmd == "" {
		return spec.ReturnFail(spec.Code[spec.IllegalParameters], "less process matcher")
	}
	flags := ""
	if process != "" {
		flags = fmt.Sprintf("--process %s", process)
	} else if processCmd != "" {
		flags = fmt.Sprintf("--process-cmd %s", processCmd)
	}
	return kpe.channel.Run(ctx, path.Join(kpe.channel.GetScriptPath(), killProcessBin), flags)
}

func (kpe *KillProcessExecutor) SetChannel(channel spec.Channel) {
	kpe.channel = channel
}
