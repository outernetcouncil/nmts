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
    name = "common_proto",
    srcs = ["common.proto"],
    import_prefix = "nmts",
)

proto_library(
    name = "sdwan_proto",
    srcs = ["sdwan.proto"],
    import_prefix = "nmts",
)

proto_library(
    name = "subscriber_proto",
    srcs = ["subscriber.proto"],
    import_prefix = "nmts",
    deps = [":common_proto"],
)

proto_library(
    name = "uni_proto",
    srcs = ["uni.proto"],
    import_prefix = "nmts",
    deps = [
        ":common_proto",
        "//proto/types/ieee:ieee802_proto",
    ],
)

cpp_proto_library(
    name = "mef_cpp_proto",
    protos = [
        ":common_proto",
        ":sdwan_proto",
        ":subscriber_proto",
        ":uni_proto",
    ],
    deps = ["//proto/types/ieee:ieee_cpp_proto"],
)

go_proto_library(
    name = "mef_go_proto",
    importpath = "outernetcouncil.org/nmts/proto/types/mef",
    protos = [
        ":common_proto",
        ":sdwan_proto",
        ":subscriber_proto",
        ":uni_proto",
    ],
    deps = ["//proto/types/ieee:ieee_go_proto"],
)

java_proto_library(
    name = "mef_java_proto",
    protos = [
        ":common_proto",
        ":sdwan_proto",
        ":subscriber_proto",
        ":uni_proto",
    ],
    deps = ["//proto/types/ieee:ieee_java_proto"],
)

python_proto_library(
    name = "mef_python_proto",
    protos = [
        ":common_proto",
        ":sdwan_proto",
        ":subscriber_proto",
        ":uni_proto",
    ],
    deps = ["//proto/types/ieee:ieee_python_proto"],
)
