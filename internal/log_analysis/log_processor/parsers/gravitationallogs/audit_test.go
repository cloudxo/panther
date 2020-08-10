package gravitationallogs

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"fmt"
	"testing"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/testutil"
)

var logTypeTeleportAudit = TypeTeleportAudit.Describe().Name

func TestTeleportAudit(t *testing.T) {
	type testCase struct {
		Name    string
		Input   string
		Expect  []string
		LogType string
	}
	for _, tc := range []testCase{
		{
			Name:    "session.start",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "addr.local": "127.0.0.1:3022",
			  "addr.remote": "1.1.1.1:63558",
			  "code": "T2000I",
			  "ei": 0,
			  "event": "session.start",
			  "login": "root",
			  "namespace": "default",
			  "server_hostname": "ip-172-31-14-137.us-west-2.compute.internal",
			  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
			  "server_labels": {
				"arch": "x86_64",
				"env": "demo",
				"hostname": "ip-172-31-14-137.us-west-2.compute.internal",
				"role": "test-cluster"
			  },
			  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
			  "size": "80:25",
			  "time": "2020-08-07T07:52:09.821Z",
			  "uid": "34bb01d5-8cef-4925-875a-783b2dbee3b6",
			  "user": "kostaspap"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "addr.local": "127.0.0.1:3022",
				  "addr.remote": "1.1.1.1:63558",
				  "code": "T2000I",
				  "ei": 0,
				  "event": "session.start",
				  "login": "root",
				  "namespace": "default",
				  "server_hostname": "ip-172-31-14-137.us-west-2.compute.internal",
				  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
				  "server_labels": {
					"arch": "x86_64",
					"env": "demo",
					"hostname": "ip-172-31-14-137.us-west-2.compute.internal",
					"role": "test-cluster"
				  },
				  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
				  "size": "80:25",
				  "time": "2020-08-07T07:52:09.821Z",
				  "uid": "34bb01d5-8cef-4925-875a-783b2dbee3b6",
				  "user": "kostaspap",
				  "p_event_time": "2020-08-07T07:52:09.821Z",
				  "p_any_ip_addresses": ["127.0.0.1", "1.1.1.1"],
				  "p_any_domain_names": ["ip-172-31-14-137.us-west-2.compute.internal"],
				  "p_any_trace_ids": ["e527ab2a-d882-11ea-9f82-0a588c28e4c2"],
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "session.end",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "code": "T2004I",
			  "ei": 22,
			  "enhanced_recording": true,
			  "event": "session.end",
			  "interactive": true,
			  "namespace": "default",
			  "participants": [
				"kostaspap"
			  ],
			  "server_addr": "[::]:3022",
			  "server_hostname": "ip-172-31-14-137.us-west-2.compute.internal",
			  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
			  "session_start": "2020-08-07T07:52:09.817638134Z",
			  "session_stop": "2020-08-07T07:52:44.607207415Z",
			  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
			  "time": "2020-08-07T07:52:44.607Z",
			  "uid": "2e3b58c6-e135-4ca4-b79b-ee7acdafa8ab",
			  "user": "kostaspap"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "code": "T2004I",
				  "ei": 22,
				  "enhanced_recording": true,
				  "event": "session.end",
				  "interactive": true,
				  "namespace": "default",
				  "participants": [
					"kostaspap"
				  ],
				  "server_addr": "[::]:3022",
				  "server_hostname": "ip-172-31-14-137.us-west-2.compute.internal",
				  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
				  "session_start": "2020-08-07T07:52:09.817638134Z",
				  "session_stop": "2020-08-07T07:52:44.607207415Z",
				  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
				  "time": "2020-08-07T07:52:44.607Z",
				  "uid": "2e3b58c6-e135-4ca4-b79b-ee7acdafa8ab",
				  "user": "kostaspap",
				  "p_event_time": "2020-08-07T07:52:44.607Z",
				  "p_any_ip_addresses": ["::"],
				  "p_any_domain_names": ["ip-172-31-14-137.us-west-2.compute.internal"],
				  "p_any_trace_ids": ["e527ab2a-d882-11ea-9f82-0a588c28e4c2"],
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "session.data",
			LogType: logTypeTeleportAudit,
			Input: ` {
				"addr.local": "127.0.0.1:3022",
				"addr.remote": "1.1.1.1:63558",
				"code": "T2006I",
				"ei": 2147483646,
				"event": "session.data",
				"login": "root",
				"rx": 5286,
				"server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
				"sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
				"time": "2020-08-07T07:52:25Z",
				"tx": 5848,
				"uid": "9f2bc778-e87e-4536-9c3f-0d4a57955fd0",
				"user": "kostaspap"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
					"addr.local": "127.0.0.1:3022",
					"addr.remote": "1.1.1.1:63558",
					"code": "T2006I",
					"ei": 2147483646,
					"event": "session.data",
					"login": "root",
					"rx": 5286,
					"server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
					"sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
					"time": "2020-08-07T07:52:25Z",
					"tx": 5848,
					"uid": "9f2bc778-e87e-4536-9c3f-0d4a57955fd0",
					"user": "kostaspap",
					"p_event_time": "2020-08-07T07:52:25Z",
					"p_any_ip_addresses": ["127.0.0.1", "1.1.1.1"],
					"p_any_trace_ids": ["e527ab2a-d882-11ea-9f82-0a588c28e4c2"],
					"p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "session.command",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "argv": [
				"-u"
			  ],
			  "cgroup_id": 4294967537,
			  "code": "T4000I",
			  "ei": 11,
			  "event": "session.command",
			  "login": "root",
			  "namespace": "default",
			  "path": "/usr/bin/id",
			  "pid": 11371,
			  "ppid": 11370,
			  "program": "id",
			  "return_code": 0,
			  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
			  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
			  "time": "2020-08-07T07:52:09.932Z",
			  "uid": "d5f13a69-539b-4978-b6d3-e73cf72f1024",
			  "user": "kostaspap"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "argv": [
					"-u"
				  ],
				  "cgroup_id": 4294967537,
				  "code": "T4000I",
				  "ei": 11,
				  "event": "session.command",
				  "login": "root",
				  "namespace": "default",
				  "path": "/usr/bin/id",
				  "pid": 11371,
				  "ppid": 11370,
				  "program": "id",
				  "return_code": 0,
				  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
				  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
				  "time": "2020-08-07T07:52:09.932Z",
				  "uid": "d5f13a69-539b-4978-b6d3-e73cf72f1024",
				  "user": "kostaspap",
				  "p_event_time": "2020-08-07T07:52:09.932Z",
				  "p_any_trace_ids": ["e527ab2a-d882-11ea-9f82-0a588c28e4c2"],
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "user.create",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "code": "T1002I",
			  "connector": "github",
			  "event": "user.create",
			  "expires": "2020-08-08T13:39:41.895085658Z",
			  "name": "kostaspap",
			  "roles": [ "admin" ],
			  "time": "2020-08-07T07:39:42Z",
			  "uid": "14551101-4d8e-4f35-b40b-a1b1ead65d43",
			  "user": "system"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "code": "T1002I",
				  "connector": "github",
				  "event": "user.create",
				  "expires": "2020-08-08T13:39:41.895085658Z",
				  "name": "kostaspap",
				  "roles": [ "admin" ],
				  "time": "2020-08-07T07:39:42Z",
				  "uid": "14551101-4d8e-4f35-b40b-a1b1ead65d43",
				  "user": "system",
				  "p_event_time": "2020-08-07T07:39:42Z",
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "user.login",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "code": "T1001W",
			  "error": "list of user teams is empty, did you grant access?",
			  "event": "user.login",
			  "method": "github",
			  "success": false,
			  "time": "2020-08-06T20:43:13Z",
			  "uid": "8750c3dd-8fcb-4f1e-8cc4-634473d1e8bc"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "code": "T1001W",
				  "error": "list of user teams is empty, did you grant access?",
				  "event": "user.login",
				  "method": "github",
				  "success": false,
				  "time": "2020-08-06T20:43:13Z",
				  "uid": "8750c3dd-8fcb-4f1e-8cc4-634473d1e8bc",
				  "p_event_time": "2020-08-06T20:43:13Z",
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "github.created",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "code": "T8000I",
			  "event": "github.created",
			  "name": "github",
			  "time": "2020-08-06T20:34:17Z",
			  "uid": "0212fa21-4011-4831-8934-a29653b78bb0",
			  "user": "411b9b66-b686-471a-b2c6-f6dc6c745f93.aws"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "code": "T8000I",
				  "event": "github.created",
				  "name": "github",
				  "time": "2020-08-06T20:34:17Z",
				  "uid": "0212fa21-4011-4831-8934-a29653b78bb0",
				  "user": "411b9b66-b686-471a-b2c6-f6dc6c745f93.aws",
				  "p_event_time": "2020-08-06T20:34:17Z",
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "session.leave",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "code": "T2003I",
			  "ei": 34,
			  "event": "session.leave",
			  "namespace": "default",
			  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
			  "sid": "49aa4466-d824-11ea-a94f-0a588c28e4c2",
			  "time": "2020-08-06T20:41:24.042Z",
			  "uid": "9d13e8a3-b8ed-4749-b09f-f1c96edb822d",
			  "user": "benarent"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "code": "T2003I",
				  "ei": 34,
				  "event": "session.leave",
				  "namespace": "default",
				  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
				  "sid": "49aa4466-d824-11ea-a94f-0a588c28e4c2",
				  "time": "2020-08-06T20:41:24.042Z",
				  "uid": "9d13e8a3-b8ed-4749-b09f-f1c96edb822d",
				  "user": "benarent",
				  "p_event_time": "2020-08-06T20:41:24.042Z",
				  "p_any_trace_ids": ["49aa4466-d824-11ea-a94f-0a588c28e4c2"],
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
		{
			Name:    "resize",
			LogType: logTypeTeleportAudit,
			Input: `{
			  "code": "T2002I",
			  "ei": 1,
			  "event": "resize",
			  "login": "root",
			  "namespace": "default",
			  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
			  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
			  "size": "197:55",
			  "time": "2020-08-07T07:52:09.839Z",
			  "uid": "cce1ada6-34cf-40e7-b527-eb2bf7f14e0c",
			  "user": "kostaspap"
			}`,
			Expect: []string{
				fmt.Sprintf(`{
				  "code": "T2002I",
				  "ei": 1,
				  "event": "resize",
				  "login": "root",
				  "namespace": "default",
				  "server_id": "411b9b66-b686-471a-b2c6-f6dc6c745f93",
				  "sid": "e527ab2a-d882-11ea-9f82-0a588c28e4c2",
				  "size": "197:55",
				  "time": "2020-08-07T07:52:09.839Z",
				  "uid": "cce1ada6-34cf-40e7-b527-eb2bf7f14e0c",
				  "user": "kostaspap",
				  "p_event_time": "2020-08-07T07:52:09.839Z",
				  "p_any_trace_ids": ["e527ab2a-d882-11ea-9f82-0a588c28e4c2"],
				  "p_log_type": "%s"
				}`, logTypeTeleportAudit),
			},
		},
	} {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			testutil.CheckRegisteredParser(t, tc.LogType, tc.Input, tc.Expect...)
		})
	}
}
