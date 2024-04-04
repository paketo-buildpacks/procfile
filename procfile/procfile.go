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

	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"

	"github.com/buildpacks/libcnb"
)

// Procfile is a map between a logical name and a command.
type Procfile map[string]interface{}

const (
	BindingType = "Procfile" // BindingType is used to resolve a binding containing a Procfile
)

// NewProcfileFromEnvironment creates a Procfile by reading environment variable BP_PROCFILE_DEFAULT_PROCESS if it exists.
// If it does not exist, returns an empty Procfile.
func NewProcfileFromEnvironment() (Procfile, error) {
	if process, isSet := os.LookupEnv("BP_PROCFILE_DEFAULT_PROCESS"); isSet {
		if process != "" {
			return Procfile{"web": process}, nil
		}
	}
	return nil, nil
}

// NewProcfileFromPath creates a Procfile by reading Procfile from path if it exists.  If it does not exist, returns an
// empty Procfile.
func NewProcfileFromPath(path string) (Procfile, error) {
	pat := regexp.MustCompile(`^([A-Za-z0-9_-]+):\s*(.+)$`)

	f := filepath.Join(path, "Procfile")
	file, err := os.OpenFile(f, os.O_RDONLY, 0644)
	if err != nil && os.IsNotExist(err) {
		return Procfile{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("unable to open Procfile %s\n%w", f, err)
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
		return nil, fmt.Errorf("unable to parse Procfile %s\n%w", f, err)
	}

	return p, nil
}

// NewProcfileFromBinding creates a Procfile by reading Procfile from bindings if it exists.  If it does not exist, returns an
// empty Procfile.
func NewProcfileFromBinding(binds libcnb.Bindings) (Procfile, error) {

	p := Procfile{}
	if binding, ok, err := bindings.ResolveOne(binds, bindings.OfType(BindingType)); err != nil {
		return nil, fmt.Errorf("unable to resolve binding\n%w", err)
	} else if ok {
		if path, ok := binding.SecretFilePath(BindingType); ok {
			if p, err = NewProcfileFromPath(filepath.Dir(path)); err != nil {
				return nil, err
			}
			return p, nil
		} else {
			return nil, fmt.Errorf("unable to find Procfile from binding")
		}
	} else {
		return p, nil
	}
}

// NewProcfileFromEnvironmentOrPathOrBinding attempts to create a merged Procfile from environment and/or given path and bindings.
// If none can be created, returns an empty Procfile.
func NewProcfileFromEnvironmentOrPathOrBinding(path string, binds libcnb.Bindings) (Procfile, error) {
	procEnv, err := NewProcfileFromEnvironment()
	if err != nil {
		return nil, err
	}
	procPath, err := NewProcfileFromPath(path)
	if err != nil {
		return nil, err
	}
	procBind, err := NewProcfileFromBinding(binds)
	if err != nil {
		return nil, err
	}
	if len(procEnv) > 0 && len(procPath)+len(procBind) > 0 {
		l := bard.NewLogger(os.Stdout)
		l.Logger.Info("A Procfile exists and BP_PROCFILE_DEFAULT_PROCESS is set, entries in Procfile take precedence")
	}

	procBind = mergeProcfiles(procEnv, procPath, procBind)
	return procBind, nil
}

// merge procfiles from binding + path, overwriting duplicate keys - binding takes precedence
func mergeProcfiles(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
