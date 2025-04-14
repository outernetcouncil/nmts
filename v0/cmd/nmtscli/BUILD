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
    name = "nmtscli_lib",
    srcs = [
        "d2.go",
        "dot.go",
        "html.go",
        "main.go",
        "nquads.go",
        "prolog.go",
        "validate.go",
    ],
    importpath = "outernetcouncil.org/nmts/v0/cmd/nmtscli",
    deps = [
        "//v0/lib/entityrelationship",
        "//v0/lib/validation",
        "//v0/proto:nmts_go_proto",
        "//v0/proto/ek/logical:logical_go_proto",
        "@com_github_ichiban_prolog//:prolog",
        "@com_github_urfave_cli_v2//:cli",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

go_binary(
    name = "nmtscli",
    embed = [":nmtscli_lib"],
    visibility = ["//visibility:public"],
)
