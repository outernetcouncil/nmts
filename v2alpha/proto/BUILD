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
load("@rules_rust_prost//:defs.bzl", "rust_prost_library")

# gazelle:proto disable_global

package(
    default_visibility = ["//visibility:public"],
)

proto_library(
    name = "nmts_proto",
    srcs = ["nmts.proto"],
    import_prefix = "nmts",
    deps = [
        "//v2alpha/proto/ek/logical:bp_agent_fn_proto",
        "//v2alpha/proto/ek/logical:interface_proto",
        "//v2alpha/proto/ek/logical:internal_fabric_proto",
        "//v2alpha/proto/ek/logical:logical_packet_link_proto",
        "//v2alpha/proto/ek/logical:network_node_proto",
        "//v2alpha/proto/ek/logical:route_fn_proto",
        "//v2alpha/proto/ek/logical:sdn_agent_proto",
        "//v2alpha/proto/ek/logical:switch_fn_proto",
        "//v2alpha/proto/ek/logical:uni_proto",
        "//v2alpha/proto/ek/physical:antenna_proto",
        "//v2alpha/proto/ek/physical:carrier_configuration_proto",
        "//v2alpha/proto/ek/physical:modem_proto",
        "//v2alpha/proto/ek/physical:physical_medium_link_proto",
        "//v2alpha/proto/ek/physical:platform_proto",
        "//v2alpha/proto/ek/physical:port_proto",
        "//v2alpha/proto/ek/physical:signal_processing_chain_proto",
        "//v2alpha/proto/ek/physical:transceiver_proto",
    ],
)

cpp_proto_library(
    name = "nmts_cpp_proto",
    protos = [":nmts_proto"],
    deps = [
        "//v2alpha/proto/ek/logical:logical_cpp_proto",
        "//v2alpha/proto/ek/physical:physical_cpp_proto",
    ],
)

# gazelle:prefix outernetcouncil.org/nmts/v2alpha/proto/
go_proto_library(
    name = "nmts_go_proto",
    importpath = "outernetcouncil.org/nmts/v2alpha/proto",
    proto = ":nmts_proto",
    deps = [
        "//v2alpha/proto/ek/logical:logical_go_proto",
        "//v2alpha/proto/ek/physical:physical_go_proto",
    ],
)

java_proto_library(
    name = "nmts_java_proto",
    protos = [":nmts_proto"],
    deps = [
        "//v2alpha/proto/ek/logical:logical_java_proto",
        "//v2alpha/proto/ek/physical:physical_java_proto",
    ],
)

python_proto_library(
    name = "nmts_python_proto",
    protos = [":nmts_proto"],
    deps = [
        "//v2alpha/proto/ek/logical:logical_python_proto",
        "//v2alpha/proto/ek/physical:physical_python_proto",
    ],
)

rust_prost_library(
    name = "nmts_rust_proto",
    proto = ":nmts_proto",
)
