// Copyright 2019 Istio Authors
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

package cmd

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"istio.io/istio/istioctl/pkg/auth"

	"istio.io/istio/pilot/test/util"
)

func runCommandAndCheckGoldenFile(name, command, golden string, t *testing.T) {
	out := runCommand(name, command, t)
	util.CompareContent(out.Bytes(), golden, t)
}

func runCommandAndCheckExpectedString(name, command, expected string, t *testing.T) {
	out := runCommand(name, command, t)
	if !reflect.DeepEqual(out.String(), expected) {
		t.Errorf("test %q failed. \nExpected\n%s\nGot%s\n", name, expected, out.String())
	}
}

func runCommand(name, command string, t *testing.T) bytes.Buffer {
	t.Helper()
	var out bytes.Buffer
	rootCmd := GetRootCmd(strings.Split(command, " "))
	rootCmd.SetOutput(&out)

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("%s: unexpected error: %s", name, err)
	}
	return out
}

func TestAuthCheck(t *testing.T) {
	testCases := []struct {
		name   string
		in     string
		golden string
	}{
		{
			name:   "listeners and clusters",
			in:     "testdata/auth/productpage_config_dump.json",
			golden: "testdata/auth/productpage.golden",
		},
	}

	for _, c := range testCases {
		command := fmt.Sprintf("experimental auth check -f %s", c.in)
		runCommandAndCheckGoldenFile(c.name, command, c.golden, t)
	}
}

func TestAuthUpgrade(t *testing.T) {
	testCases := []struct {
		name     string
		in       string
		services []string
		golden   string
	}{
		{
			name:     "v1 policies",
			in:       "testdata/auth/authz-policy.yaml",
			services: []string{"testdata/auth/svc-other.yaml", "testdata/auth/svc-bookinfo.yaml"},
			golden:   "testdata/auth/authz-policy.golden",
		},
	}

	for _, c := range testCases {
		command := fmt.Sprintf("experimental auth upgrade -f %s --service %s", c.in, strings.Join(c.services, ","))
		runCommandAndCheckGoldenFile(c.name, command, c.golden, t)
	}
}

func TestAuthValidator(t *testing.T) {
	testCases := []struct {
		name     string
		in       []string
		expected string
	}{
		{
			name:     "good policy",
			in:       []string{"testdata/auth/authz-policy.yaml"},
			expected: "",
		},
		{
			name: "bad policy",
			in:   []string{"../pkg/auth/testdata/validator/unused-role.yaml", "../pkg/auth/testdata/validator/notfound-role-in-binding.yaml"},
			expected: fmt.Sprintf("%s%s",
				fmt.Sprintf(auth.RoleNotFound, "some-role", "bind-service-viewer", "default"),
				fmt.Sprintf(auth.RoleNotUsed, "unused-role", "default")),
		},
	}
	for _, c := range testCases {
		command := fmt.Sprintf("experimental auth validate -f %s", strings.Join(c.in, ","))
		runCommandAndCheckExpectedString(c.name, command, c.expected, t)
	}
}
