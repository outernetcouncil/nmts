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

load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_proto_grpc//cpp:defs.bzl", "cpp_proto_library")
load("@rules_proto_grpc//java:defs.bzl", "java_proto_library")
load("@rules_proto_grpc//python:defs.bzl", "python_proto_library")

package(
    default_visibility = ["//visibility:public"],
)

##
# Per-proto file library rules.
##
proto_library(
    name = "azimuth_elevation_mask_proto",
    srcs = ["azimuth_elevation_mask.proto"],
    deps = [
        ":coordinates_proto",
    ],
)

proto_library(
    name = "coordinates_proto",
    srcs = ["coordinates.proto"],
    deps = ["@com_google_googleapis//google/type:quaternion_proto"],
)

proto_library(
    name = "geodesy_proto",
    srcs = ["geodesy.proto"],
)

proto_library(
    name = "motion_proto",
    srcs = ["motion.proto"],
    deps = [
        ":coordinates_proto",
        ":geodesy_proto",
        ":orbital_proto",
        "@com_google_googleapis//google/type:interval_proto",
        "@com_google_googleapis//google/type:quaternion_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

proto_library(
    name = "orbital_proto",
    srcs = ["orbital.proto"],
    deps = [
        ":coordinates_proto",
        "@com_google_protobuf//:timestamp_proto",
    ],
)

cpp_proto_library(
    name = "geophys_cpp_proto",
    protos = [
        ":azimuth_elevation_mask_proto",
        ":coordinates_proto",
        ":geodesy_proto",
        ":motion_proto",
        ":orbital_proto",
    ],
    deps = [
        "@com_google_googleapis//google/type:interval_cc_proto",
        "@com_google_googleapis//google/type:quaternion_cc_proto",
    ],
)

go_proto_library(
    name = "geophys_go_proto",
    importpath = "outernetcouncil.org/nmts/proto/types/geophys",
    protos = [
        ":azimuth_elevation_mask_proto",
        ":coordinates_proto",
        ":geodesy_proto",
        ":motion_proto",
        ":orbital_proto",
    ],
    deps = [
        "@org_golang_google_genproto//googleapis/type/interval",
        "@org_golang_google_genproto//googleapis/type/quaternion",
    ],
)

java_proto_library(
    name = "geophys_java_proto",
    protos = [
        ":azimuth_elevation_mask_proto",
        ":coordinates_proto",
        ":geodesy_proto",
        ":motion_proto",
        ":orbital_proto",
    ],
    deps = [
        "@com_google_googleapis//google/type:type_java_proto",
    ],
)

python_proto_library(
    name = "geophys_python_proto",
    protos = [
        ":azimuth_elevation_mask_proto",
        ":coordinates_proto",
        ":geodesy_proto",
        ":motion_proto",
        ":orbital_proto",
    ],
    deps = [
        "@com_google_googleapis//google/type:interval_py_proto",
        "@com_google_googleapis//google/type:quaternion_py_proto",
    ],
)
