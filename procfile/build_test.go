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
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/procfile/v5/procfile"
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
						"test-type-2": "test-command-2 argument",
					},
				},
			},
		}

		result := libcnb.NewBuildResult()
		result.Processes = append(result.Processes,
			libcnb.Process{
				Type:    "test-type-1",
				Command: "test-command-1",
			},
			libcnb.Process{
				Type:    "test-type-2",
				Command: "test-command-2 argument",
			},
		)

		Expect(build.Build(ctx)).To(Equal(result))
	})

	context("given BP_DIRECT_PROCESS=true", func() {
		it.Before(func() {
			t.Setenv("BP_DIRECT_PROCESS", "true")
		})

		it("uses a process with direct=true", func() {
			ctx.Plan = libcnb.BuildpackPlan{
				Entries: []libcnb.BuildpackPlanEntry{
					{
						Name: "procfile",
						Metadata: map[string]interface{}{
							"test-type": "test-command arg",
						},
					},
				},
			}

			result := libcnb.NewBuildResult()
			result.Processes = append(result.Processes,
				libcnb.Process{
					Type:      "test-type",
					Command:   "test-command",
					Arguments: []string{"arg"},
					Direct:    true,
				},
			)

			Expect(build.Build(ctx)).To(Equal(result))
		})
	})

	context("given a special process name", func() {
		var assertMarkedAsDefault = func(name string) {
			ctx.Plan = libcnb.BuildpackPlan{
				Entries: []libcnb.BuildpackPlanEntry{
					{
						Name: "procfile",
						Metadata: map[string]interface{}{
							"test-type-1": "test-command-1",
							name:          "test-command-2 argument",
						},
					},
				},
			}

			result := libcnb.NewBuildResult()
			result.Processes = append(result.Processes,
				libcnb.Process{
					Type:    "test-type-1",
					Command: "test-command-1",
				},
				libcnb.Process{
					Type:    name,
					Command: "test-command-2 argument",
					Default: true,
				},
			)

			Expect(build.Build(ctx)).To(Equal(result))
		}

		it("adds metadata to result, marks web process as default", func() {
			assertMarkedAsDefault("web")
		})

		it("adds metadata to result, marks worker process as default", func() {
			assertMarkedAsDefault("worker")
		})

		it("has only one process marked as default", func() {
			{
				ctx.Plan = libcnb.BuildpackPlan{
					Entries: []libcnb.BuildpackPlanEntry{
						{
							Name: "procfile",
							Metadata: map[string]interface{}{
								"web":    "test-command-1",
								"worker": "test-command-2 argument",
							},
						},
					},
				}

				result := libcnb.NewBuildResult()
				result.Processes = append(result.Processes,
					libcnb.Process{
						Type:    "web",
						Command: "test-command-1",
						Default: true,
					},
					libcnb.Process{
						Type:    "worker",
						Command: "test-command-2 argument",
					},
				)

				Expect(build.Build(ctx)).To(Equal(result))
			}
		})
	})

	context("bionic tiny stack", func() {
		it.Before(func() {
			ctx.StackID = libpak.BionicTinyStackID
		})

		it("adds metadata to result", func() {
			ctx.Plan = libcnb.BuildpackPlan{
				Entries: []libcnb.BuildpackPlanEntry{
					{
						Name: "procfile",
						Metadata: map[string]interface{}{
							"test-type-1": "test-command-1",
							"test-type-2": "test-command-2 argument",
						},
					},
				},
			}

			result := libcnb.NewBuildResult()
			result.Processes = append(result.Processes,
				libcnb.Process{
					Type:      "test-type-1",
					Command:   "test-command-1",
					Arguments: []string{},
					Direct:    true,
				},
				libcnb.Process{
					Type:      "test-type-2",
					Command:   "test-command-2",
					Arguments: []string{"argument"},
					Direct:    true,
				},
			)

			Expect(build.Build(ctx)).To(Equal(result))
		})

	})

	context("jammy tiny stack", func() {
		it.Before(func() {
			ctx.StackID = libpak.JammyTinyStackID
		})

		it("adds metadata to result", func() {
			ctx.Plan = libcnb.BuildpackPlan{
				Entries: []libcnb.BuildpackPlanEntry{
					{
						Name: "procfile",
						Metadata: map[string]interface{}{
							"test-type-1": "test-command-1",
							"test-type-2": "test-command-2 argument",
						},
					},
				},
			}

			result := libcnb.NewBuildResult()
			result.Processes = append(result.Processes,
				libcnb.Process{
					Type:      "test-type-1",
					Command:   "test-command-1",
					Arguments: []string{},
					Direct:    true,
				},
				libcnb.Process{
					Type:      "test-type-2",
					Command:   "test-command-2",
					Arguments: []string{"argument"},
					Direct:    true,
				},
			)

			Expect(build.Build(ctx)).To(Equal(result))
		})
	})
}
