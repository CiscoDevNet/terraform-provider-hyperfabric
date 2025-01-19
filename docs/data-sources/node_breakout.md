---
subcategory: "Blueprint"
layout: "hyperfabric"
page_title: "Nexus Hyperfabric: hyperfabric_node_breakout"
sidebar_current: "docs-hyperfabric-data-source-hyperfabric_node_breakout"
description: |-
  Data source for a Breakout of a Node in a Nexus Hyperfabric Fabric
---

# hyperfabric_node_breakout

Data source for a Breakout of a Node in a Nexus Hyperfabric Fabric

A Breakout defines the mode in which high-speed, channelized Node Ports are each broken down into multiple low-speed Breakout Ports, each connecting to separate network elements. For example, a switch with 400G ports can be connected to four 100G ports.

## API Paths ##

* `/fabrics/{fabricId|fabricName}/nodes/{nodeId|nodeName}/breakouts/{breakoutId|name}` `GET`

## GUI Information ##

* Location: `> Fabrics > {fabric} > Nodes > {node} > Configure > Port configuration`

## Example Usage ##

```hcl
data "hyperfabric_node_breakout" "example_node_breakout" {
  node_id = hyperfabric_node.example_node.id
  name = "exampleBreakout"
}
```

## Schema ##

### Required ###
* `node_id` - (string) The unique identifier (id) of a Node in a Fabric. Use the id attribute of the [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/resources/node) resource or [hyperfabric_node](https://registry.terraform.io/providers/cisco-open/hyperfabric/latest/docs/data-sources/node) data source.
* `name` - (string) The name of the Breakout of the Node.

### Read-Only ###

* `id` - (string) The unique identifier (id) of the Breakout of the Node in the Fabric.
* `breakout_id` - (string) The unique identifier (id) of a Breakout.
* `description` - (string) The description is a user defined field to store notes about the Breakout of the Node.
* `enabled` - (bool) The enabled state of the Breakout of the Node.
* `mode` - (string) The mode used to Breakout the Ports.
* `ports` - (list of strings) A list of Node Ports names to be broken down into Breakout Ports based on the Breakout `mode`.
* `breakouts` - (list of strings) A list of Breakout Port names of the Breakout of the Node.
* `pluggable` - (string) The type of pluggable used for the Breakout.
* `metadata` - (map) A map of the Metadata of the Node Breakout:
  * `created_at` - (string) The timestamp when this object was created in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `created_by` - (string) The user that created this object.
  * `modified_at` - (string) The timestamp when this object was last modified in [RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format.
  * `modified_by` - (string) The user that modified this object last.
  * `revision_id` - (string) An integer that represent the current revision of the object.
* `labels` - (list of strings) A list of user-defined labels that can be used for grouping and filtering objects.
* `annotations` - (list of maps) A list of key-value annotations to store user-defined data including complex data such as JSON.
  * `name` - (string) The name used to uniquely identify the annotation.
  * `value` - (string) The value of the annotation.
  * `data_type` - (string) The type of data stored in the value of the annotation.
      - Possible Values: `STRING`, `INT32`, `UINT32`, `INT64`, `UINT64`, `BOOL`, `TIME`, `UUID`, `DURATION`, `JSON`.
