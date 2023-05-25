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

type OvnLogicalRouterPort struct {
	UUID string //`ovsdb:"_uuid"`
	//Enabled        []bool            //`ovsdb:"enabled"`
	ExternalIDs    map[string]string //`ovsdb:"external_ids"`
	GatewayChassis []string          //`ovsdb:"gateway_chassis"`
	HaChassisGroup []string          //`ovsdb:"ha_chassis_group"`
	//Ipv6Prefix     []string          //`ovsdb:"ipv6_prefix"`
	//Ipv6RaConfigs  map[string]string //`ovsdb:"ipv6_ra_configs"`
	MAC      string            //`ovsdb:"mac"`
	Name     string            //`ovsdb:"name"`
	Networks []string          //`ovsdb:"networks"`
	Options  map[string]string //`ovsdb:"options"`
	Peer     []string          //`ovsdb:"peer"`
}

// GetLogicalRouterPorts returns a list of OVN logical router ports.
func (cli *OvnClient) GetLogicalRouterPorts() ([]*OvnLogicalRouterPort, error) {
	// First, fetch logical router ports.
	ports := []*OvnLogicalRouterPort{}
	query := "SELECT _uuid, external_ids, gateway_chassis, ha_chassis_group, mac, name, networks, options, peer FROM Logical_Router_Port"
	result, err := cli.Database.Northbound.Client.Transact(cli.Database.Northbound.Name, query)
	if err != nil {
		return nil, fmt.Errorf("%s: '%s' table error: %s", cli.Database.Northbound.Name, "Logical_Router_Port", err)
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("%s: no logical router port found", cli.Database.Northbound.Name)
	}
	for _, row := range result.Rows {
		port := OvnLogicalRouterPort{}
		if r, dt, err := row.GetColumnValue("_uuid", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			port.UUID = r.(string)
		}
		if r, dt, err := row.GetColumnValue("name", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			port.Name = r.(string)
		}
		if r, dt, err := row.GetColumnValue("mac", result.Columns); err != nil {
			continue
		} else {
			if dt != "string" {
				continue
			}
			port.MAC = r.(string)
		}
		if r, dt, err := row.GetColumnValue("networks", result.Columns); err != nil {
			continue
		} else {
			if dt != "[]string" {
				continue
			}
			port.Networks = r.([]string)
		}
		if r, dt, err := row.GetColumnValue("external_ids", result.Columns); err == nil {
			if dt == "map[string]string" {
				port.ExternalIDs = r.(map[string]string)
			}
		} else {
			port.ExternalIDs = make(map[string]string)
		}
		ports = append(ports, &port)
	}
	return ports, nil
}
