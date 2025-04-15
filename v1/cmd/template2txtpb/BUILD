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

load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "template2txtpb_lib",
    srcs = ["main.go"],
    importpath = "outernetcouncil.org/nmts/v1/cmd/template2txtpb",
    visibility = ["//visibility:private"],
    deps = [
        "//v1/proto:nmts_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

go_binary(
    name = "template2txtpb",
    embed = [":template2txtpb_lib"],
    visibility = ["//visibility:public"],
)
