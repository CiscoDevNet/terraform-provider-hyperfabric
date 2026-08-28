// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

// API response fields use pointers so refresh can distinguish an omitted field
// from a legitimate zero value. Collection fields remain in their decoded JSON
// representation because the Terraform model conversion helpers consume that
// shape.

type fabricAPIResponse struct {
	FabricId    *string                `json:"fabricId"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Topology    *string                `json:"topology"`
	Location    *string                `json:"location"`
	Address     *string                `json:"address"`
	City        *string                `json:"city"`
	Country     *string                `json:"country"`
	Metadata    map[string]interface{} `json:"metadata"`
	Labels      []interface{}          `json:"labels"`
	Annotations []interface{}          `json:"annotations"`
}

type nodeAPIResponse struct {
	NodeId       *string                `json:"nodeId"`
	Name         *string                `json:"name"`
	Description  *string                `json:"description"`
	Enabled      *bool                  `json:"enabled"`
	Location     *string                `json:"location"`
	ModelName    *string                `json:"modelName"`
	SerialNumber *string                `json:"serialNumber"`
	DeviceId     *string                `json:"deviceId"`
	Roles        []interface{}          `json:"roles"`
	Metadata     map[string]interface{} `json:"metadata"`
	Labels       []interface{}          `json:"labels"`
	Annotations  []interface{}          `json:"annotations"`
}

type vrfAPIResponse struct {
	FabricId    *string                `json:"fabricId"`
	Id          *string                `json:"id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Enabled     *bool                  `json:"enabled"`
	IsDefault   *bool                  `json:"isDefault"`
	ASN         *float64               `json:"asn"`
	VNI         *float64               `json:"vni"`
	RouteTarget *string                `json:"routeTarget"`
	Metadata    map[string]interface{} `json:"metadata"`
	Labels      []interface{}          `json:"labels"`
	Annotations []interface{}          `json:"annotations"`
}

type vniAPIResponse struct {
	FabricId    *string                `json:"fabricId"`
	Id          *string                `json:"id"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Enabled     *bool                  `json:"enabled"`
	IsDefault   *bool                  `json:"isDefault"`
	VrfID       *string                `json:"vrfId"`
	VNI         *float64               `json:"vni"`
	MTU         *float64               `json:"mtu"`
	Members     []interface{}          `json:"members"`
	SVIs        []interface{}          `json:"svis"`
	Metadata    map[string]interface{} `json:"metadata"`
	Labels      []interface{}          `json:"labels"`
	Annotations []interface{}          `json:"annotations"`
}

type nodeLoopbackAPIResponse struct {
	Id          *string                `json:"id"`
	FabricId    *string                `json:"fabricId"`
	NodeId      *string                `json:"nodeId"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	IPv4Address *string                `json:"ipv4Address"`
	IPv6Address *string                `json:"ipv6Address"`
	VrfID       *string                `json:"vrfId"`
	Metadata    map[string]interface{} `json:"metadata"`
	Labels      []interface{}          `json:"labels"`
	Annotations []interface{}          `json:"annotations"`
}

type nodeSubInterfaceAPIResponse struct {
	Id            *string                `json:"id"`
	FabricId      *string                `json:"fabricId"`
	NodeId        *string                `json:"nodeId"`
	Name          *string                `json:"name"`
	Description   *string                `json:"description"`
	Enabled       *bool                  `json:"enabled"`
	IPv4Addresses []interface{}          `json:"ipv4Addresses"`
	IPv6Addresses []interface{}          `json:"ipv6Addresses"`
	VlanId        *float64               `json:"vlanId"`
	VrfID         *string                `json:"vrfId"`
	Parent        *string                `json:"parent"`
	Metadata      map[string]interface{} `json:"metadata"`
	Labels        []interface{}          `json:"labels"`
	Annotations   []interface{}          `json:"annotations"`
}

type nodeBreakoutAPIResponse struct {
	Id          *string                `json:"id"`
	FabricId    *string                `json:"fabricId"`
	NodeId      *string                `json:"nodeId"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Enabled     *bool                  `json:"enabled"`
	Breakouts   []interface{}          `json:"breakouts"`
	Ports       []interface{}          `json:"ports"`
	Mode        *string                `json:"mode"`
	Pluggable   *string                `json:"pluggable"`
	Metadata    map[string]interface{} `json:"metadata"`
	Labels      []interface{}          `json:"labels"`
	Annotations []interface{}          `json:"annotations"`
}

type nodePortAPIResponse struct {
	Id                 *string                `json:"id"`
	FabricId           *string                `json:"fabricId"`
	NodeId             *string                `json:"nodeId"`
	Name               *string                `json:"name"`
	Description        *string                `json:"description"`
	Enabled            *bool                  `json:"enabled"`
	Index              *float64               `json:"index"`
	IPv4Addresses      []interface{}          `json:"ipv4Addresses"`
	IPv6Addresses      []interface{}          `json:"ipv6Addresses"`
	Linecard           *float64               `json:"linecard"`
	PreventForwarding  *bool                  `json:"linkDown"`
	LLDPHost           *string                `json:"lldpHost"`
	LLDPInfo           *string                `json:"lldpInfo"`
	LLDPPort           *string                `json:"lldpPort"`
	MaxSpeed           *string                `json:"maxSpeed"`
	MTU                *float64               `json:"mtu"`
	Roles              []interface{}          `json:"roles"`
	Speed              *string                `json:"speed"`
	SubInterfacesCount *float64               `json:"subInfCount"`
	VLANIDs            []interface{}          `json:"vlanIds"`
	VNIs               []interface{}          `json:"vnis"`
	VrfID              *string                `json:"vrfId"`
	Metadata           map[string]interface{} `json:"metadata"`
	Labels             []interface{}          `json:"labels"`
	Annotations        []interface{}          `json:"annotations"`
}

type nodeManagementPortAPIResponse struct {
	Id                *string                `json:"id"`
	Name              *string                `json:"name"`
	Description       *string                `json:"description"`
	Enabled           *bool                  `json:"enabled"`
	CloudURLs         []interface{}          `json:"cloudUrls"`
	IPv4ConfigType    *string                `json:"ipv4ConfigType"`
	IPv4Address       *string                `json:"ipv4Address"`
	IPv4Gateway       *string                `json:"ipv4Gateway"`
	IPv6ConfigType    *string                `json:"ipv6ConfigType"`
	IPv6Address       *string                `json:"ipv6Address"`
	IPv6Gateway       *string                `json:"ipv6Gateway"`
	DNSAddresses      []interface{}          `json:"dnsAddresses"`
	NTPAddresses      []interface{}          `json:"ntpAddresses"`
	NoProxy           []interface{}          `json:"noProxy"`
	ProxyAddress      *string                `json:"proxyAddress"`
	ProxyCredentialId *string                `json:"proxyCredentialId"`
	ProxyUsername     *string                `json:"proxyUsername"`
	ConnectedState    *string                `json:"connectedState"`
	ConfigOrigin      *string                `json:"configOrigin"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type nodeManagementPortsAPIResponse struct {
	Ports []nodeManagementPortAPIResponse `json:"ports"`
}

type connectionAPIResponse struct {
	Id           *string                `json:"id"`
	FabricId     *string                `json:"fabricId"`
	Description  *string                `json:"description"`
	Pluggable    *string                `json:"pluggable"`
	Local        map[string]interface{} `json:"local"`
	Remote       map[string]interface{} `json:"remote"`
	OSType       *string                `json:"osType"`
	Unrecognized *bool                  `json:"unrecognized"`
}

type userAPIResponse struct {
	Id          *string                `json:"id"`
	Email       *string                `json:"email"`
	Provider    *string                `json:"provider"`
	LastLogin   *string                `json:"lastLogin"`
	Enabled     *bool                  `json:"enabled"`
	Role        *string                `json:"role"`
	Metadata    map[string]interface{} `json:"metadata"`
	Labels      []interface{}          `json:"labels"`
	Annotations []interface{}          `json:"annotations"`
}

type bearerTokenAPIResponse struct {
	TokenId     *string                `json:"tokenId"`
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	NotBefore   *string                `json:"notBefore"`
	NotAfter    *string                `json:"notAfter"`
	Scope       *string                `json:"scope"`
	Token       *string                `json:"token"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type deviceAPIResponse struct {
	DeviceId     *string       `json:"deviceId"`
	FabricId     *string       `json:"fabricId"`
	ModelName    *string       `json:"modelName"`
	NodeId       *string       `json:"nodeId"`
	OSType       *string       `json:"osType"`
	RackId       *string       `json:"rackId"`
	Roles        []interface{} `json:"roles"`
	SerialNumber *string       `json:"serialNumber"`
}

type devicesAPIResponse struct {
	Devices []deviceAPIResponse `json:"devices"`
}

type fabricsAPIResponse struct {
	Fabrics []fabricAPIResponse `json:"fabrics"`
}

type nodesAPIResponse struct {
	Nodes []nodeAPIResponse `json:"nodes"`
}

type vrfsAPIResponse struct {
	VRFs []vrfAPIResponse `json:"vrfs"`
}

type vnisAPIResponse struct {
	VNIs []vniAPIResponse `json:"vnis"`
}

type nodeLoopbacksAPIResponse struct {
	Loopbacks []nodeLoopbackAPIResponse `json:"loopbacks"`
}

type nodeSubInterfacesAPIResponse struct {
	SubInterfaces []nodeSubInterfaceAPIResponse `json:"subInterfaces"`
}

type nodeBreakoutsAPIResponse struct {
	Breakouts []nodeBreakoutAPIResponse `json:"breakouts"`
}

type connectionsAPIResponse struct {
	Connections []connectionAPIResponse `json:"connections"`
}

type usersAPIResponse struct {
	Users []userAPIResponse `json:"users"`
}

type bearerTokensAPIResponse struct {
	Tokens []bearerTokenAPIResponse `json:"tokens"`
	Token  *string                  `json:"token"`
}
