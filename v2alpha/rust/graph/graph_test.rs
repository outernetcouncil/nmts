// Copyright (c) Outernet Council and Contributors.
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

use nmts_proto::nmts::v2alpha::{Entity,Fragment,Relationship,Rk};
use nmts_proto::nmts::v2alpha::entity::Kind;
use nmts_proto::platform_proto::nmts::v2alpha::ek::physical::Platform;
use nmts_proto::network_node_proto::nmts::v2alpha::ek::logical::NetworkNode;

#[test]
fn basic_protobuf_types_instantiation() {
    // Create EK_PLATFORM
    let mut ek_platform = Platform::default();
    ek_platform.name = "platform name".to_string();
    let mut platform = Entity::default();
    platform.id = "platform".to_string();
    platform.kind = Some(Kind::EkPlatform(ek_platform));

    // Create EK_NETWORK_NODE
    let mut ek_network_node = NetworkNode::default();
    ek_network_node.name = "network_node name".to_string();
    let mut network_node = Entity::default();
    network_node.id = "network_node".to_string();
    network_node.kind = Some(Kind::EkNetworkNode(ek_network_node));

    // Create RK_CONTAINS between them.
    let rk_contains = Relationship{
        kind: Rk::Contains.into(),
        a: platform.id.clone(),
        z: network_node.id.clone(),
    };

    let mut fragment = Fragment::default();
    fragment.entity.push(platform);
    fragment.entity.push(network_node);
    fragment.relationship.push(rk_contains);

    let mut fragments = Vec::<Fragment>::new();
    fragments.push(fragment);
}
