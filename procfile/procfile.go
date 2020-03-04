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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Procfile is a map between a logical name and a command.
type Procfile map[string]interface{}

// NewProcfile creates a Procfile by reading Procfile from path if it exists.  If it does not exist, returns an
// empty Procfile.
func NewProcfileFromPath(path string) (Procfile, error) {
	pat := regexp.MustCompile("^([A-Za-z0-9_-]+):\\s*(.+)$")

	f := filepath.Join(path, "Procfile")
	file, err := os.OpenFile(f, os.O_RDONLY, 0644)
	if err != nil && os.IsNotExist(err) {
		return Procfile{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to open Procfile %s: %w", f, err)
	}
	defer file.Close()

	p := Procfile{}

	s := bufio.NewScanner(file)
	for s.Scan() {
		parts := pat.FindStringSubmatch(s.Text())
		if len(parts) > 0 {
			p[parts[1]] = parts[2]
		}
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("unable to parse Procfile %s: %w", f, err)
	}

	return p, nil
}
