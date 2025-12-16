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

load("@rules_go//go:def.bzl", "go_library", "go_test")

package(
    default_visibility = ["//visibility:public"],
)

go_library(
    name = "utilities",
    srcs = ["utilities.go"],
    importpath = "outernetcouncil.org/nmts/v1/lib/utilities",
    deps = [
        "//v1/lib/graph",
        "//v1/proto:nmts_go_proto",
        "@com_github_deckarep_golang_set_v2//:golang-set",
        "@com_github_samber_lo//:lo",
    ],
)

go_test(
    name = "utilities_test",
    srcs = ["utilities_test.go"],
    deps = [
        ":utilities",
        "//v1/lib/entityrelationship",
        "//v1/lib/graph",
        "//v1/lib/utilities/testing",
        "//v1/proto:nmts_go_proto",
        "//v1/proto/ek/logical:logical_go_proto",
        "//v1/proto/types/ietf:ietf_go_proto",
        "@com_github_deckarep_golang_set_v2//:golang-set",
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_samber_lo//:lo",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
