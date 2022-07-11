// Copyright 2019 PingCAP, Inc.
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
	"context"
	"strconv"
	"testing"

	"github.com/pingcap/tidb/plugin"
	"github.com/stretchr/testify/require"
)

func TestLoadPlugin(t *testing.T) {
	ctx := context.Background()
	pluginName := "s3audit"
	pluginVersion := uint16(1)
	pluginSign := pluginName + "-" + strconv.Itoa(int(pluginVersion))

	cfg := plugin.Config{
		Plugins:    []string{pluginSign},
		PluginDir:  "",
		EnvVersion: map[string]uint16{"go": 1112},
	}

	// setup load test hook.
	loadOne := func(p *plugin.Plugin, dir string, pluginID plugin.ID) (manifest func() *plugin.Manifest, err error) {
		return func() *plugin.Manifest {
			m := &plugin.AuditManifest{
				Manifest: plugin.Manifest{
					Kind:       plugin.Audit,
					Name:       pluginName,
					Version:    pluginVersion,
					OnInit:     OnInit,
					OnShutdown: OnShutdown,
					Validate:   Validate,
				},
				OnGeneralEvent:    OnGeneralEvent,
				OnConnectionEvent: OnConnectionEvent,
			}
			return plugin.ExportManifest(m)
		}, nil
	}
	plugin.SetTestHook(loadOne)

	// trigger load.
	err := plugin.Load(ctx, cfg)
	require.NoErrorf(t, err, "load plugin [%s] fail, error [%s]\n", pluginSign, err)

	err = plugin.Init(ctx, cfg)
	require.NoErrorf(t, err, "init plugin [%s] fail, error [%s]\n", pluginSign, err)

	err = plugin.ForeachPlugin(plugin.Audit, func(auditPlugin *plugin.Plugin) error {
		plugin.DeclareAuditManifest(auditPlugin.Manifest).OnGeneralEvent(context.Background(), nil, plugin.Completed, "QUERY")
		return nil
	})
	require.NoErrorf(t, err, "query event fail, error [%s]\n", err)
}
