# Copyright (c) Outernet Council and Contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@protobuf//bazel:proto_library.bzl", "proto_library")
load("@rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto_grpc_cpp//:defs.bzl", "cpp_proto_library")
load("@rules_proto_grpc_java//:defs.bzl", "java_proto_library")
load("@rules_proto_grpc_python//:defs.bzl", "python_proto_library")

package(
    default_visibility = ["//visibility:public"],
)

##
# Per-proto file library rules.
##
proto_library(
    name = "interface_proto",
    srcs = ["interface.proto"],
    import_prefix = "nmts",
    deps = [
        "//v0/proto/types/ieee:ieee802_proto",
        "//v0/proto/types/ietf:inet_proto",
        "@protobuf//:empty_proto",
    ],
)

proto_library(
    name = "logical_packet_link_proto",
    srcs = ["logical_packet_link.proto"],
    import_prefix = "nmts",
    deps = [
        "//v0/proto/types/ietf:routing_proto",
    ],
)

proto_library(
    name = "network_node_proto",
    srcs = ["network_node.proto"],
    import_prefix = "nmts",
)

proto_library(
    name = "route_fn_proto",
    srcs = ["route_fn.proto"],
    import_prefix = "nmts",
    deps = [
        "//v0/proto/types/ietf:routing_proto",
    ],
)

proto_library(
    name = "sdn_agent_proto",
    srcs = ["sdn_agent.proto"],
    import_prefix = "nmts",
)

proto_library(
    name = "switch_fn_proto",
    srcs = ["switch_fn.proto"],
    import_prefix = "nmts",
)

cpp_proto_library(
    name = "logical_cpp_proto",
    protos = [
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
    ],
    deps = [
        "//v0/proto/types/ieee:ieee_cpp_proto",
        "//v0/proto/types/ietf:ietf_cpp_proto",
    ],
)

go_proto_library(
    name = "logical_go_proto",
    importpath = "outernetcouncil.org/nmts/v0/proto/ek/logical",
    protos = [
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
    ],
    deps = [
        "//v0/proto/types/ieee:ieee_go_proto",
        "//v0/proto/types/ietf:ietf_go_proto",
    ],
)

java_proto_library(
    name = "logical_java_proto",
    protos = [
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
    ],
    deps = [
        "//v0/proto/types/ieee:ieee_java_proto",
        "//v0/proto/types/ietf:ietf_java_proto",
    ],
)

python_proto_library(
    name = "logical_python_proto",
    protos = [
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
    ],
    deps = [
        "//v0/proto/types/ieee:ieee_python_proto",
        "//v0/proto/types/ietf:ietf_python_proto",
    ],
)
