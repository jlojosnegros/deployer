/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2021 Red Hat, Inc.
 */

package flagcodec

import (
	"reflect"
	"testing"
)

func TestParseStringRoundTrip(t *testing.T) {
	type testCase struct {
		name     string
		command  string
		args     []string
		expected []string
	}

	testCases := []testCase{
		{
			name: "nil",
		},
		{
			name: "empty",
			args: []string{},
		},
		{
			name:    "only command",
			command: "/bin/true",
			expected: []string{
				"/bin/true",
			},
		},
		{
			name:    "simple",
			command: "/bin/resource-topology-exporter",
			args: []string{
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
			},
			expected: []string{
				"/bin/resource-topology-exporter",
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fl := ParseArgvKeyValueWithCommand(tc.command, tc.args)
			got := fl.Argv()
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected %v got %v", tc.expected, got)
			}
		})
	}
}

func TestAddFlags(t *testing.T) {
	type testOpt struct {
		name  string
		value string
	}

	type testCase struct {
		name     string
		command  string
		args     []string
		options  []testOpt
		expected []string
	}

	testCases := []testCase{
		{
			name:    "add-mixed",
			command: "/bin/resource-topology-exporter",
			args: []string{
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
			},
			options: []testOpt{
				{
					name:  "--hostname",
					value: "host.test.net",
				},
				{
					name:  "--v",
					value: "2",
				},
			},
			expected: []string{
				"/bin/resource-topology-exporter",
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
				"--hostname=host.test.net",
				"--v=2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fl := ParseArgvKeyValueWithCommand(tc.command, tc.args)
			for _, opt := range tc.options {
				fl.SetOption(opt.name, opt.value)
			}
			got := fl.Argv()
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected %v got %v", tc.expected, got)
			}
		})
	}
}

func TestDeleteFlags(t *testing.T) {

	type testCase struct {
		name     string
		command  string
		args     []string
		options  []string
		expected []string
	}

	testCases := []testCase{
		{
			name:    "remove-option",
			command: "/bin/resource-topology-exporter",
			args: []string{
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
			},
			options: []string{
				"--sleep-interval",
			},
			expected: []string{
				"/bin/resource-topology-exporter",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
			},
		},
		{
			name:    "remove-toggle",
			command: "/bin/resource-topology-exporter",
			args: []string{
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
				"--pods-fingerprint",
			},
			options: []string{
				"--pods-fingerprint",
			},
			expected: []string{
				"/bin/resource-topology-exporter",
				"--sleep-interval=10s",
				"--sysfs=/host-sys",
				"--kubelet-state-dir=/host-var/lib/kubelet",
				"--podresources-socket=unix:///host-var/lib/kubelet/pod-resources/kubelet.sock",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fl := ParseArgvKeyValueWithCommand(tc.command, tc.args)
			for _, opt := range tc.options {
				fl.Delete(opt)
			}
			got := fl.Argv()
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("expected %v got %v", tc.expected, got)
			}
		})
	}
}
