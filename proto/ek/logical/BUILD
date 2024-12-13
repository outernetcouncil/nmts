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
    name = "bp_agent_fn_proto",
    srcs = ["bp_agent_fn.proto"],
    import_prefix = "nmts",
)

proto_library(
    name = "interface_proto",
    srcs = ["interface.proto"],
    import_prefix = "nmts",
    deps = [
        "//proto/types/ieee:ieee802_proto",
        "//proto/types/ietf:inet_proto",
        "@protobuf//:empty_proto",
    ],
)

proto_library(
    name = "logical_packet_link_proto",
    srcs = ["logical_packet_link.proto"],
    import_prefix = "nmts",
    deps = [
        "//proto/types/ietf:routing_proto",
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
        "//proto/types/ietf:routing_proto",
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

proto_library(
    name = "uni_proto",
    srcs = ["uni.proto"],
    import_prefix = "nmts",
    deps = ["//proto/types/mef:uni_proto"],
)

cpp_proto_library(
    name = "logical_cpp_proto",
    protos = [
        ":bp_agent_fn_proto",
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
        ":uni_proto",
    ],
    deps = [
        "//proto/types/ieee:ieee_cpp_proto",
        "//proto/types/ietf:ietf_cpp_proto",
        "//proto/types/mef:mef_cpp_proto",
    ],
)

go_proto_library(
    name = "logical_go_proto",
    importpath = "outernetcouncil.org/nmts/proto/ek/logical",
    protos = [
        ":bp_agent_fn_proto",
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
        ":uni_proto",
    ],
    deps = [
        "//proto/types/ieee:ieee_go_proto",
        "//proto/types/ietf:ietf_go_proto",
        "//proto/types/mef:mef_go_proto",
    ],
)

java_proto_library(
    name = "logical_java_proto",
    protos = [
        ":bp_agent_fn_proto",
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
        ":uni_proto",
    ],
    deps = [
        "//proto/types/ieee:ieee_java_proto",
        "//proto/types/ietf:ietf_java_proto",
        "//proto/types/mef:mef_java_proto",
    ],
)

python_proto_library(
    name = "logical_python_proto",
    protos = [
        ":bp_agent_fn_proto",
        ":interface_proto",
        ":logical_packet_link_proto",
        ":network_node_proto",
        ":route_fn_proto",
        ":sdn_agent_proto",
        ":switch_fn_proto",
        ":uni_proto",
    ],
    deps = [
        "//proto/types/ieee:ieee_python_proto",
        "//proto/types/ietf:ietf_python_proto",
        "//proto/types/mef:mef_python_proto",
    ],
)
