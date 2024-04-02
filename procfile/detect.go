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
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	l := bard.NewLogger(os.Stdout)
	// Create Procfile from source path or binding, if both exist, merge into one. The binding takes precedence on duplicate name/command pairs.
	p, err := NewProcfileFromEnvironmentOrPathOrBinding(context.Application.Path, context.Platform.Bindings)
	if err != nil {
		return libcnb.DetectResult{}, err
	}

	if len(p) == 0 {
		l.Logger.Info("SKIPPED: No procfile found from environment, source path, or binding.")
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "procfile"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "procfile", Metadata: p},
				},
			},
		},
	}, nil
}
