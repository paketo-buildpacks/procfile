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

package procfile_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/procfile/procfile"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		build procfile.Build
		ctx   libcnb.BuildContext
	)

	it("does nothing without plan", func() {
		Expect(build.Build(ctx)).To(Equal(libcnb.BuildResult{}))
	})

	it("adds metadata to result", func() {
		ctx.Plan = libcnb.BuildpackPlan{
			Entries: []libcnb.BuildpackPlanEntry{
				{
					Name: "procfile",
					Metadata: map[string]interface{}{
						"test-type-1": "test-command-1",
						"test-type-2": "test-command-2",
					},
				},
			},
		}

		result := libcnb.NewBuildResult()
		result.Processes = append(result.Processes,
			libcnb.Process{Type: "test-type-1", Command: "test-command-1"},
			libcnb.Process{Type: "test-type-2", Command: "test-command-2"},
		)

		Expect(build.Build(ctx)).To(Equal(result))
	})
}
