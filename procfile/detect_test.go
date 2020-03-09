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
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/procfile/procfile"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect procfile.Detect
		path   string
	)

	it.Before(func() {
		var err error
		path, err = ioutil.TempDir("", "procfile")
		Expect(err).NotTo(HaveOccurred())

		ctx.Application.Path = path
	})

	it("fails without Procfile", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
	})

	it("fails with empty Procfile", func() {
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte(""), 0644))

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
	})

	it("passes with Procfile", func() {
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte(`test-type-1: test-command-1
test-type-2: test-command-2`), 0644))

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "procfile"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "procfile", Metadata: procfile.Procfile{
							"test-type-1": "test-command-1",
							"test-type-2": "test-command-2",
						}},
					},
				},
			},
		}))
	})
}
