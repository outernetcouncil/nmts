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

proto_library(
    name = "service_proto",
    srcs = ["service.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//proto/types/ietf:inet_proto",
        "@googleapis//google/type:interval_proto",
        "@protobuf//:duration_proto",
    ],
)

go_proto_library(
    name = "service_go_proto",
    importpath = "outernetcouncil.org/nmts/proto/service",
    proto = ":service_proto",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/types/ietf:ietf_go_proto",
        "@org_golang_google_genproto//googleapis/type/interval",
    ],
)
