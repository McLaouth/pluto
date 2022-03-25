// Copyright 2020 Fairwinds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "embed"
	"os"

	"github.com/fairwindsops/pluto/v3/cmd"
	"k8s.io/klog/v2"
)

var (
	// version is set during build
	version = "development"
	// commit is set during build
	commit = "n/a"

	//go:embed versions.yaml
	versionsFile []byte
)

func main() {
	command, err := cmd.NewRootCommand(version, commit, versionsFile)
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
	if err := command.Execute(); err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
