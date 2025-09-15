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

package(default_visibility = ["//visibility:public"])

##
# Per-proto file library rules.
##
proto_library(
    name = "antenna_proto",
    srcs = ["antenna.proto"],
    import_prefix = "nmts",
    deps = [
        "//v1/proto/types/geophys:azimuth_elevation_mask_proto",
        "//v1/proto/types/geophys:coordinates_proto",
        "//v1/proto/types/physical:antenna_pattern_proto",
    ],
)

proto_library(
    name = "common_proto",
    srcs = ["common.proto"],
    import_prefix = "nmts",
    deps = [
        "//v1/proto/types/physical:polarization_proto",
        "//v1/proto/types/physical:waveform_proto",
    ],
)

proto_library(
    name = "modem_proto",
    srcs = ["modem.proto"],
    import_prefix = "nmts",
    deps = [
        "//v1/proto/types/physical:waveform_proto",
    ],
)

proto_library(
    name = "physical_medium_link_proto",
    srcs = ["physical_medium_link.proto"],
    import_prefix = "nmts",
    deps = [
        ":common_proto",
        "@protobuf//:duration_proto",
    ],
)

proto_library(
    name = "platform_proto",
    srcs = ["platform.proto"],
    import_prefix = "nmts",
    deps = [
        "//v1/proto/types/geophys:motion_proto",
    ],
)

proto_library(
    name = "port_proto",
    srcs = ["port.proto"],
    import_prefix = "nmts",
    deps = [
        "//v1/proto/types/ietf:inet_proto",
    ],
)

proto_library(
    name = "signal_processing_chain_proto",
    srcs = ["signal_processing_chain.proto"],
    import_prefix = "nmts",
    deps = [
        "@protobuf//:duration_proto",
    ],
)

proto_library(
    name = "transceiver_proto",
    srcs = ["transceiver.proto"],
    import_prefix = "nmts",
    deps = [
        ":common_proto",
    ],
)

cpp_proto_library(
    name = "physical_cpp_proto",
    protos = [
        ":antenna_proto",
        ":common_proto",
        ":modem_proto",
        ":physical_medium_link_proto",
        ":platform_proto",
        ":port_proto",
        ":signal_processing_chain_proto",
        ":transceiver_proto",
    ],
    deps = [
        "//v1/proto/types/geophys:geophys_cpp_proto",
        "//v1/proto/types/ietf:ietf_cpp_proto",
        "//v1/proto/types/physical:physical_cpp_proto",
    ],
)

go_proto_library(
    name = "physical_go_proto",
    importpath = "outernetcouncil.org/nmts/v1/proto/ek/physical",
    protos = [
        ":antenna_proto",
        ":common_proto",
        ":modem_proto",
        ":physical_medium_link_proto",
        ":platform_proto",
        ":port_proto",
        ":signal_processing_chain_proto",
        ":transceiver_proto",
    ],
    deps = [
        "//v1/proto/types/geophys:geophys_go_proto",
        "//v1/proto/types/ietf:ietf_go_proto",
        "//v1/proto/types/physical:physical_go_proto",
    ],
)

java_proto_library(
    name = "physical_java_proto",
    protos = [
        ":antenna_proto",
        ":common_proto",
        ":modem_proto",
        ":physical_medium_link_proto",
        ":platform_proto",
        ":port_proto",
        ":signal_processing_chain_proto",
        ":transceiver_proto",
    ],
    deps = [
        "//v1/proto/types/geophys:geophys_java_proto",
        "//v1/proto/types/ietf:ietf_java_proto",
        "//v1/proto/types/physical:physical_java_proto",
    ],
)

python_proto_library(
    name = "physical_python_proto",
    protos = [
        ":antenna_proto",
        ":common_proto",
        ":modem_proto",
        ":physical_medium_link_proto",
        ":platform_proto",
        ":port_proto",
        ":signal_processing_chain_proto",
        ":transceiver_proto",
    ],
    deps = [
        "//v1/proto/types/geophys:geophys_python_proto",
        "//v1/proto/types/ietf:ietf_python_proto",
        "//v1/proto/types/physical:physical_python_proto",
    ],
)
