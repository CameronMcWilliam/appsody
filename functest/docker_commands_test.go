// Copyright © 2019 IBM Corporation and others.
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
package functest

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	cmd "github.com/appsody/appsody/cmd"
	"github.com/appsody/appsody/cmd/cmdtest"
)

var invalidDockerCmdsTest = []struct {
	file     string
	expected string
}{
	{"imageName", "invalid reference format: repository name must be lowercase"},
	{"imagename", "No such image: imagename"},
}

func TestDockerInspect(t *testing.T) {

	for _, testData := range invalidDockerCmdsTest {
		// need to set testData to a new variable scoped under the for loop
		// otherwise tests run in parallel may get the wrong testData
		// because the for loop reassigns it before the func runs
		test := testData

		t.Run(fmt.Sprintf("Test Invalid DockerInspect"), func(t *testing.T) {
			var outBuffer bytes.Buffer
			config := &cmd.LoggingConfig{}
			config.InitLogging(&outBuffer, &outBuffer)

			out, err := cmd.RunDockerInspect(config, test.file)
			t.Log(outBuffer.String())

			if err == nil {
				t.Error("Expected an error from '", test.file, "' name but it did not return one.")
			} else if !strings.Contains(out, test.expected) {
				t.Error("Expected the stdout to contain '" + test.expected + "'. It actually contains: " + out)
			}
		})
	}
}

var buildahBuild = []struct {
	config   *cmd.RootCommandConfig
	expected string
}{
	{&cmd.RootCommandConfig{Dryrun: false}, "[Buildah] Writing manifest to image destination"},
	{&cmd.RootCommandConfig{Dryrun: true}, "Dryrun complete"},
}

func TestBuildahBuild(t *testing.T) {

	for _, testData := range buildahBuild {
		// need to set testData to a new variable scoped under the for loop
		// otherwise tests run in parallel may get the wrong testData
		// because the for loop reassigns it before the func runs
		test := testData

		config := test.config

		t.Run(fmt.Sprintf("Test Buildah Build"), func(t *testing.T) {
			if runtime.GOOS != "linux" {
				t.Skip()
			}

			sandbox, cleanup := cmdtest.TestSetupWithSandbox(t, true)
			defer cleanup()

			var outBuffer bytes.Buffer
			log := &cmd.LoggingConfig{}
			log.InitLogging(&outBuffer, &outBuffer)

			// Because the 'starter' folder has been copied, the stack.yaml file will be in the 'starter'
			// folder within the temp directory that has been generated for sandboxing purposes, rather than
			// the usual core temp directory
			sandbox.ProjectDir = filepath.Join(sandbox.TestDataPath, "starter")

			args := []string{"build", "--buildah", "--buildah-options", "--format=docker"}
			output, err := cmdtest.RunAppsody(sandbox, args...)
			if err != nil {
				t.Fatalf("Test failed unexpectedly when dryrun value: %v with error: %v", config.Dryrun, err)
			} else {
				if !strings.Contains(output, test.expected) {
					t.Error("String ", test.expected, " not found in output")
				}
			}
		})
	}
}
