// Copyright 2018 Paul Greenberg (greenpau@outlook.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ovsdb

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
)

// OvnLogicalRouter holds basic information about a logical router.
type OvnLogicalRouter struct {
	UUID        string `json:"uuid" yaml:"uuid"`
	Name        string `json:"name" yaml:"name"`
	TunnelKey   uint64 `json:"tunnel_key" yaml:"tunnel_key"`
	DatapathID  string
	ExternalIDs map[string]string
	Ports       []string `json:"ports" yaml:"ports"`
}

// GetLogicalRouters returns a list of OVN logical routers.
func (cli *OvnClient) GetLogicalRouters() ([]*OvnLogicalRouter, error) {
	routers := []*OvnLogicalRouter{}
	// First, get basic information about OVN logical routers.
	query := "SELECT _uuid, external_ids, name, ports, policies, FROM Logical_Router"
	result, err := cli.Database.Northbound.Client.Transact(cli.Database.Northbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Northbound.Name, "Logical_Router", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no router found", cli.Database.Northbound.Name)
	}
	for _, row := range result.Rows {
		rt := &OvnLogicalRouter{}
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			rt.UUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("name", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			rt.Name = r.(string)
		}
		if r, dt, err := row.GetColumnValue("ports", result.Columns); err != nil {
			continue
		} else {
			switch dt {
			case "string":
				rt.Ports = append(rt.Ports, r.(string))
			case "[]string":
				rt.Ports = r.([]string)
			default:
				continue
			}
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err != nil {
			rt.ExternalIDs = make(map[string]string)
		} else {
			if dt == "map[string]string" {
				rt.ExternalIDs = r.(map[string]string)
			}
		}
		routers = append(routers, rt)
	}
	// !TODO(ccope): What does this do for switches? Do we want an equivalent function for routers?
	// Next, obtain a tunnel key for the datapath associated with the switch.
	/*query = "SELECT _uuid, external_ids, tunnel_key FROM Datapath_Binding"
	result, err = cli.Database.Southbound.Client.Transact(cli.Database.Southbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Southbound.Name, "Datapath_Binding", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no datapath binding found", cli.Database.Southbound.Name)
	}
	for _, row := range result.Rows {
		var bindUUID string
		var bindExternalIDs map[string]string
		var bindTunnelKey uint64
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			bindUUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("tunnel_key", result.Columns); err != nil {
			continue
		} else {
			if dt != "integer" {
				continue
			}
			bindTunnelKey = uint64(r.(int64))
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err != nil {
			continue
		} else {
			if dt != "map[string]string" {
				continue
			}
			bindExternalIDs = r.(map[string]string)
		}
		if len(bindExternalIDs) < 1 {
			continue
		}
		if _, exists := bindExternalIDs["logical-switch"]; !exists {
			continue
		}
		for _, sw := range switches {
			if bindExternalIDs["logical-switch"] == sw.UUID {
				sw.TunnelKey = bindTunnelKey
				sw.DatapathID = bindUUID
				break
			}
		}
	}*/
	return routers, nil
}

/*
!TODO(ccope): this function seems unused in the ovs_logical_switch.go file, not going to try to adapt it for routers right now
// MapPortToRouter update logical router ports with the entries from the
// logical routers associated with the ports.
func (cli *OvnClient) MapPortToRouter(logicalRouters []*OvnLogicalRouter, logicalRouterPorts []*OvnLogicalRouterPort) {
	portRef := make(map[string]string)
	portMap := make(map[string]*OvnLogicalRouter)
	for _, logicalRouter := range logicalRouters {
		for _, port := range logicalRouter.Ports {
			portRef[port] = logicalRouter.UUID
			portMap[port] = logicalRouter
		}
	}
	for _, logicalRouterPort := range logicalRouterPorts {
		if _, exists := portRef[logicalRouterPort.UUID]; !exists {
			continue
		}
		logicalRouterPort.LogicalRouterUUID = portMap[logicalRouterPort.UUID].UUID
		logicalRouterPort.LogicalRouterName = portMap[logicalRouterPort.UUID].Name
		for k, v := range portMap[logicalRouterPort.UUID].ExternalIDs {
			logicalRouterPort.ExternalIDs[k] = v
		}
	}
}
*/
