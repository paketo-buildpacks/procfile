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

	. "github.com/onsi/gomega"
	"github.com/paketoio/procfile/procfile"
	"github.com/sclevine/spec"
)

func testProcfile(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path string
	)

	it.Before(func() {
		var err error
		path, err = ioutil.TempDir("", "procfile")
		Expect(err).NotTo(HaveOccurred())
	})

	it("returns an empty Procfile when file does not exist", func() {
		Expect(procfile.NewProcfileFromPath(path)).To(HaveLen(0))
	})

	it("returns a parsed Profile", func() {
		Expect(ioutil.WriteFile(filepath.Join(path, "Procfile"), []byte("test-type: test-command"), 0644)).To(Succeed())

		Expect(procfile.NewProcfileFromPath(path)).To(Equal(procfile.Procfile{"test-type": "test-command"}))
	})
}
