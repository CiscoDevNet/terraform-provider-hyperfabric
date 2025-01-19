// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNodeBreakoutResource(t *testing.T) {
	fabricName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	breakoutName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify provided and default Hyperfabric values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - Create with minimum config and verify provided and default Hyperfabric values.")
				},
				Config:             testNodeBreakoutResourceHclConfig(fabricName, breakoutName, "minimal"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "name", breakoutName),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.0", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.#", "4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.0", "Ethernet1_1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.1", "Ethernet1_1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.2", "Ethernet1_1_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.3", "Ethernet1_1_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "mode", "4x25G(4)"),
				),
			},
			// Update with all config and verify provided values.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - Update with all config and verify provided values.")
				},
				Config:             testNodeBreakoutResourceHclConfig(fabricName, breakoutName, "full"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "name", breakoutName),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "description", "Ports to be configured as Breakout"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.0", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.1", "Ethernet1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.#", "8"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.0", "Ethernet1_1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.1", "Ethernet1_1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.2", "Ethernet1_1_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.3", "Ethernet1_1_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.4", "Ethernet1_2_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.5", "Ethernet1_2_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.6", "Ethernet1_2_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.7", "Ethernet1_2_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "mode", "4x25G(4)"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "pluggable", "QSFP-4SFP25G-CU1M"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "annotations.#", "2"),
				),
			},
			// Update with minimum config and verify config is unchanged.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - Update with minimum config and verify config is unchanged.")
				},
				Config:             testNodeBreakoutResourceHclConfig(fabricName, breakoutName, "minimal+"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "name", breakoutName),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "description", "Ports to be configured as Breakout"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.0", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.1", "Ethernet1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.#", "8"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.0", "Ethernet1_1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.1", "Ethernet1_1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.2", "Ethernet1_1_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.3", "Ethernet1_1_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.4", "Ethernet1_2_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.5", "Ethernet1_2_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.6", "Ethernet1_2_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.7", "Ethernet1_2_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "mode", "4x25G(4)"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "pluggable", "QSFP-4SFP25G-CU1M"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "labels.#", "2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "annotations.#", "2"),
				),
			},
			// ImportState testing with pre-existing Id.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - ImportState testing with pre-existing Id.")
				},
				ResourceName:      "hyperfabric_node_breakout.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// ImportState testing with fabric and node name.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - ImportState testing with fabric, node and breakout name.")
				},
				ResourceName:      "hyperfabric_node_breakout.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fabricName + "/nodes/node1/breakouts/" + breakoutName,
			},
			// Update with config containing all optional attributes with empty values and verify config is cleared.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - Update with config containing all optional attributes with empty values and verify config is cleared.")
				},
				Config:             testNodeBreakoutResourceHclConfig(fabricName, breakoutName, "clear"),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "name", breakoutName),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.0", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.#", "4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.0", "Ethernet1_1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.1", "Ethernet1_1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.2", "Ethernet1_1_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.3", "Ethernet1_1_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "mode", "4x25G(4)"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "pluggable", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "annotations.#", "0"),
				),
			},
			// Run Plan Only with minimal config and check that plan is empty.
			{
				PreConfig: func() {
					fmt.Println("= RUNNING: Node Breakout - Run Plan Only with minimal config and check that plan is empty.")
				},
				Config:             testNodeBreakoutResourceHclConfig(fabricName, breakoutName, "minimal"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "name", breakoutName),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "description", ""),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "enabled", "true"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "ports.0", "Ethernet1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.#", "4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.0", "Ethernet1_1_1"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.1", "Ethernet1_1_2"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.2", "Ethernet1_1_3"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "breakouts.3", "Ethernet1_1_4"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "mode", "4x25G(4)"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hyperfabric_node_breakout.test", "annotations.#", "0"),
				),
			},
		},
	})
}

func testNodeBreakoutResourceHclConfig(fabricName string, breakoutName string, configType string) string {
	if configType == "full" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
}

resource "hyperfabric_node_breakout" "test" {
	node_id     = hyperfabric_node.test.id
	name        = "%[2]s"
	description = "Ports to be configured as Breakout"
	ports       = ["Ethernet1_1", "Ethernet1_2"]
	mode        = "4x25G(4)"
	pluggable   = "QSFP-4SFP25G-CU1M"
	labels      = [
		"sj01-1-101-AAA01",
		"blue"
	]
	annotations = [
		{
			name      = "color"
			value     = "blue"
		},
		{
			data_type = "UINT32"
			name      = "rack"
			value     = "1"
		}
	]
}
`, fabricName, breakoutName)
	} else if configType == "clear" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
	fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
}

resource "hyperfabric_node_breakout" "test" {
	node_id     = hyperfabric_node.test.id
	name        = "%[2]s"
	mode        = "4x25G(4)"
	description = ""
	ports       = ["Ethernet1_1"]
	pluggable   = ""
	labels      = []
	annotations = []
}
`, fabricName, breakoutName)
	} else if configType == "minimal+" {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
	fabric_id   = hyperfabric_fabric.test.id
	name        = "node1"
	model_name  = "HF6100-32D"
}

resource "hyperfabric_node_breakout" "test" {
	node_id     = hyperfabric_node.test.id
	name        = "%[2]s"
	ports       = ["Ethernet1_1", "Ethernet1_2"]
	mode        = "4x25G(4)"
}
`, fabricName, breakoutName)
	} else {
		return fmt.Sprintf(`
resource "hyperfabric_fabric" "test" {
	name = "%[1]s"
}

resource "hyperfabric_node" "test" {
    fabric_id  = hyperfabric_fabric.test.id
	name       = "node1"
	model_name = "HF6100-32D"
}

resource "hyperfabric_node_breakout" "test" {
    node_id = hyperfabric_node.test.id
	name    = "%[2]s"
	ports   = ["Ethernet1_1"]
	mode    = "4x25G(4)"
}
`, fabricName, breakoutName)
	}
}
