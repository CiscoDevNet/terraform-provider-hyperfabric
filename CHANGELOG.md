# Terraform Provider Cisco Nexus Hyperfabric - Changelog

All notable changes to this project will be documented in this file.

## 0.1.0 (August 28, 2026)

IMPROVEMENTS:
- Add resources and data-sources for hyperfabric_fabric, hyperfabric_node, hyperfabric_node_port, hyperfabric_node_management_port, hyperfabric_device, hyperfabric_connection, hyperfabric_bind_to_node, hyperfabric_user, hyperfabric_vni, hyperfabric_vrf
- Add resources and data-sources for hyperfabric_node_loopback
- Add topology attribute back to hyperfabric_fabric
- Add resources and data-sources for hyperfabric_node_sub_interface
- Add resources and data-sources for hyperfabric_node_breakout
- Add hyperfabric_bearer_token resource and datasource, fix documentation typo and format example

BUG FIXES:
- Fix a series of bug across hyperfabric_bearer_token, hyperfabric_node_port and hyperfabric_node_management_port
- Fix a series of bug across hyperfabric_bearer_token, hyperfabric_node_port, hyperfabric_node, hyperfabric_connection, hyperfabric_user, hyperfabric_vrf and hyperfabric_vni
- Add support to ignore 404 when deleting
- Extract uuid from full id when provided and a specific resource uuid is expected
- Fix hyperfabric_node_port issue on destroy when roles set to fabric_port and fix test to ignore vrf_id bug issue for now
- Make credentials and tokens sensitive attributes and add regression check for it.

<!--
## 1.0.0 (Unreleased)
BREAKING CHANGE:
- Migration to Terraform Provider SDK v2. Remove support for Terraform v0.11.x or below
- Fix and update netflow monitor relation in aci_leaf_access_port_policy_group and aci_leaf_access_bundle_policy_group
- Fix tcp_rules from string to list in aci_filter_entry

IMPROVEMENTS:
- Add `aci_netflow_record_policy` resource and data-source. (#1220)
- Add `aci_l3out_node_sid_profile` resource and data-source
- Add `aci_netflow_monitor_policy` and `aci_relation_to_netflow_exporter` resources and data-sources (#1208)
- Add `aci_l3out_provider_label` resource and data-source (#1200)
- Add `aci_relation_to_fallback_route_group` resource and data-source (#1195)
- Add attributes `pc_tag` and `scope` to `aci_vrf` (#1238)
- Allow dn based filtering for `aci_client_end_point` data-source
- Remove `flood_on_encap` and `prio` attributes and change the non required attributes to read-only in `aci_endpoint_security` data-source
- Enable toggling of escaping of HTML characters with escape_html attribute in `aci_rest_managed` payloads (#1199)

BUG FIXES:
- Add error handling in try login function for `aaa_user`
- Prevent error by setting `flood_on_encap` and `prio` for `aci_endpoint_security_group`
- Fix to avoid known after applies for children when they are not provided and not configured on APIC
- Fix import functionality for `aci_rest_managed` when brackets are present in DN
-->
## 0.1.0 (Unreleased)

- Initial Release

FEATURES:
