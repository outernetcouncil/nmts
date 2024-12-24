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
load("@rules_go//go:def.bzl", "go_test")
load("@rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto_grpc_cpp//:defs.bzl", "cpp_proto_library")
load("@rules_proto_grpc_java//:defs.bzl", "java_proto_library")
load("@rules_proto_grpc_python//:defs.bzl", "python_proto_library")

package(
    default_visibility = ["//visibility:public"],
)

proto_library(
    name = "classifier_proto",
    srcs = ["classifier.proto"],
    import_prefix = "nmts",
    deps = [
        ":packet_proto",
        "@protobuf//:field_mask_proto",
    ],
)

proto_library(
    name = "packet_proto",
    srcs = ["packet.proto"],
    import_prefix = "nmts",
    deps = [
        "//v0/proto/types/iana:ipproto_proto",
        "//v0/proto/types/ieee:ethertype_proto",
        "//v0/proto/types/ieee:ieee802_proto",
        "//v0/proto/types/ietf:common_proto",
        "//v0/proto/types/ietf:inet_proto",
    ],
)

cpp_proto_library(
    name = "core_cpp_proto",
    protos = [
        ":classifier_proto",
        ":packet_proto",
    ],
    deps = [
        "//v0/proto/types/iana:iana_cpp_proto",
        "//v0/proto/types/ieee:ieee_cpp_proto",
        "//v0/proto/types/ietf:ietf_cpp_proto",
    ],
)

go_proto_library(
    name = "core_go_proto",
    importpath = "outernetcouncil.org/nmts/v0/proto/core",
    protos = [
        ":classifier_proto",
        ":packet_proto",
    ],
    deps = [
        "//v0/proto/types/iana:iana_go_proto",
        "//v0/proto/types/ieee:ieee_go_proto",
        "//v0/proto/types/ietf:ietf_go_proto",
    ],
)

java_proto_library(
    name = "core_java_proto",
    protos = [
        ":classifier_proto",
        ":packet_proto",
    ],
    deps = [
        "//v0/proto/types/iana:iana_java_proto",
        "//v0/proto/types/ieee:ieee_java_proto",
        "//v0/proto/types/ietf:ietf_java_proto",
    ],
)

python_proto_library(
    name = "core_python_proto",
    protos = [
        ":classifier_proto",
        ":packet_proto",
    ],
    deps = [
        "//v0/proto/types/iana:iana_python_proto",
        "//v0/proto/types/ieee:ieee_python_proto",
        "//v0/proto/types/ietf:ietf_python_proto",
    ],
)

go_test(
    name = "core_test",
    srcs = ["packet_test.go"],
    deps = [
        ":core_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)
