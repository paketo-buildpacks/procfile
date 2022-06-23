/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package procfile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/mattn/go-shellwords"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	r := libpak.PlanEntryResolver{Plan: context.Plan}
	e, ok, err := r.Resolve("procfile")
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve buildpack plan entry procfile\n%w", err)
	} else if !ok {
		return libcnb.BuildResult{}, nil
	}

	for k, v := range e.Metadata {
		process := libcnb.Process{Type: k}

		if (context.StackID == libpak.BionicTinyStackID) || (context.StackID == libpak.JammyTinyStackID) {
			s, err := shellwords.Parse(v.(string))
			if err != nil {
				return libcnb.BuildResult{}, fmt.Errorf("unable to parse %s\n%w", s, err)
			}

			process.Command = s[0]
			process.Arguments = s[1:]
			process.Direct = true
		} else {
			process.Command = v.(string)
		}

		result.Processes = append(result.Processes, process)
	}

	markDefaultProcess(result)

	sort.Slice(result.Processes, func(i int, j int) bool {
		return result.Processes[i].Type < result.Processes[j].Type
	})

	return result, nil
}

func markDefaultProcess(result libcnb.BuildResult) {
	for _, magicType := range []string{"web", "worker"} {
		for i, proc := range result.Processes {
			if strings.EqualFold(magicType, proc.Type) {
				result.Processes[i].Default = true
				return
			}
		}
	}
}
