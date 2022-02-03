// Copyright © 2021 Obol Technologies Inc.
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
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/obolnetwork/charon/app"
	"github.com/obolnetwork/charon/discovery"
	"github.com/obolnetwork/charon/p2p"
)

func TestCmdFlags(t *testing.T) {
	tests := []struct {
		Name            string
		Args            []string
		VersionConfig   *versionConfig
		BootstrapConfig *bootstrapConfig
		appConfig       *app.Config
	}{
		{
			Name:          "version verbose",
			Args:          slice("version", "--verbose"),
			VersionConfig: &versionConfig{Verbose: true},
		},
		{
			Name:          "version no verbose",
			Args:          slice("version", "--verbose=false"),
			VersionConfig: &versionConfig{Verbose: false},
		},
		{
			Name: "bootstrap flags",
			Args: slice("bootstrap"),
			BootstrapConfig: &bootstrapConfig{
				Out:          "./keys",
				Shares:       4,
				PasswordFile: "",
				Bootnodes:    nil,
			},
		},
		{
			Name: "bootstrap with flags",
			Args: slice("bootstrap", "-o=./gen_keys", "-n=6", "--password-file=./pass", `--bootnodes=hello,world`),
			BootstrapConfig: &bootstrapConfig{
				Out:          "./gen_keys",
				Shares:       6,
				PasswordFile: "./pass",
				Bootnodes:    []string{"hello", "world"},
			},
		},
		{
			Name: "run command",
			Args: slice("run"),
			appConfig: &app.Config{
				Discovery: discovery.Config{ListenAddr: "127.0.0.1:30309", DBPath: ""},
				P2P: p2p.Config{
					Addrs:     []string{"127.0.0.1:13900"},
					Allowlist: "",
					Denylist:  "",
				},
				ClusterDir:       "./charon/manifest.json",
				DataDir:          "./charon/data",
				MonitoringAddr:   "127.0.0.1:8088",
				ValidatorAPIAddr: "127.0.0.1:3500",
				BeaconNodeAddr:   "http://localhost/",
				JaegerAddr:       "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			root := newRootCmd(
				newVersionCmd(func(_ io.Writer, config versionConfig) {
					require.NotNil(t, test.VersionConfig)
					require.Equal(t, *test.VersionConfig, config)
				}),
				newBootstrapCmd(func(_ io.Writer, config bootstrapConfig) error {
					require.NotNil(t, test.BootstrapConfig)
					require.Equal(t, *test.BootstrapConfig, config)

					return nil
				}),
				newRunCmd(func(_ context.Context, config app.Config) error {
					require.NotNil(t, test.appConfig)
					require.Equal(t, *test.appConfig, config)

					return nil
				}),
			)

			root.SetArgs(test.Args)
			require.NoError(t, root.Execute())
		})
	}
}

// slice is a convenience function for creating string slice literals.
func slice(strs ...string) []string {
	return strs
}