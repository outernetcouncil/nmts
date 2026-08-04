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
    name = "graph",
    srcs = [
        "graph.go",
        "traverse.go",
    ],
    importpath = "outernetcouncil.org/nmts/v2alpha/lib/graph",
    deps = [
        "//v2alpha/lib/entityrelationship",
        "//v2alpha/proto:nmts_go_proto",
        "@com_github_deckarep_golang_set_v2//:golang-set",
        "@com_github_samber_lo//:lo",
    ],
)

go_test(
    name = "graph_test",
    srcs = [
        "graph_test.go",
        "traverse_test.go",
    ],
    embed = [":graph"],
    deps = [
        "//v2alpha/lib/entityrelationship",
        "//v2alpha/proto:nmts_go_proto",
        "//v2alpha/proto/ek/logical:logical_go_proto",
        "@com_github_deckarep_golang_set_v2//:golang-set",
        "@com_github_google_go_cmp//cmp",
        "@com_github_samber_lo//:lo",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)

# Benchmarks-only sibling of :graph_test, routed to the Pool=bench Buildbarn
# executors and published to Bencher. Written out longhand, not via
# //bazel/bench:defs.bzl because nmts is a separate Bazel module
# (org_outernetcouncil_nmts) and cannot load from the parent workspace. Keep in
# sync with bazel/bench/defs.bzl there. Gazelle manages go_test rules by name
# and only generates ":graph_test", so it leaves this target alone -- but it
# also never updates these deps. Keep them in sync with :graph_test above.
go_test(
    name = "graph_bench",
    timeout = "eternal",
    srcs = [
        "graph_test.go",
        "traverse_test.go",
    ],
    args = [
        "-test.bench=.",
        "-test.run=^$$",
        "-test.benchtime=2s",
        "-test.count=1",
        "-test.benchmem",
    ],
    embed = [":graph"],
    exec_properties = {"test.Pool": "bench"},
    msan = "off",
    race = "off",
    tags = [
        "bench",
        "bench_mr",
        "manual",
        "no-cache",
        "no_ci_pipeline",
    ],
    deps = [
        "//v2alpha/lib/entityrelationship",
        "//v2alpha/proto:nmts_go_proto",
        "//v2alpha/proto/ek/logical:logical_go_proto",
        "@com_github_deckarep_golang_set_v2//:golang-set",
        "@com_github_google_go_cmp//cmp",
        "@com_github_samber_lo//:lo",
        "@org_golang_google_protobuf//encoding/prototext",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//testing/protocmp",
    ],
)
