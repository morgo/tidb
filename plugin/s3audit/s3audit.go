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
	"fmt"
	"strings"

	"github.com/pingcap/tidb/plugin"
	"github.com/pingcap/tidb/sessionctx/variable"
)

// Validate implements TiDB plugin's Validate SPI.
// It is called before OnInit
// nolint: unused, deadcode
func Validate(ctx context.Context, m *plugin.Manifest) error {
	return nil
}

// OnInit implements TiDB plugin's OnInit SPI.
// nolint: unused, deadcode
func OnInit(ctx context.Context, manifest *plugin.Manifest) error {
	return nil
}

// OnShutdown implements TiDB plugin's OnShutdown SPI.
// nolint: unused, deadcode
func OnShutdown(ctx context.Context, manifest *plugin.Manifest) error {
	return nil
}

// OnGeneralEvent implements TiDB Audit plugin's OnGeneralEvent SPI.
// nolint: unused, deadcode
func OnGeneralEvent(ctx context.Context, sctx *variable.SessionVars, event plugin.GeneralEvent, cmd string) {
	// Is this connection auditable?
	// i.e. is it a DBA or user or an application
	if !isAuditable(sctx.User) {
		return
	}
	// The connection is auditable.
	// Fetch as much context as we can.
	if event == plugin.Starting || event == plugin.Error {
		// Statement and Digest are not accurate in the event of a parse error.
		// See: https://github.com/pingcap/tidb/issues/36109
		// We'll fix it once this bug closes.
		return
	}
	// It's Completed event (or a new event that is unknown)
	info := make(map[string]string)
	if sctx != nil {
		info["statement"] = sctx.StmtCtx.OriginalSQL
		digest, _ := sctx.StmtCtx.SQLDigest()
		info["digest"] = digest
		// Fetch involved tables
		if len(sctx.StmtCtx.Tables) > 0 {
			var tables []string
			for _, t := range sctx.StmtCtx.Tables {
				tables = append(tables, fmt.Sprintf("%s.%s", t.DB, t.Table))
			}
			info["tables"] = strings.Join(tables, ",")
		}
		// TODO: Fetch statistics about the statement, such as how many rows it modified.
		// This is not conclusive, since the statement itself might not be committed.
	}
	log(sctx.User, event.String(), info)
}

// OnConnectionEvent implements TiDB Audit plugin's OnConnectionEvent SPI.
// nolint: unused, deadcode
func OnConnectionEvent(ctx context.Context, event plugin.ConnectionEvent, info *variable.ConnectionInfo) error {
	// The OnConnectionEvent is fired twice during connection establishment.
	// TODO: For which 'event' do we log, or do we log both the same?
	return nil
	/*
		var reason string
		if r := ctx.Value(plugin.RejectReasonCtxValue{}); r != nil {
			reason = r.(string)
		}
	*/
}
