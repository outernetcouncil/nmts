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
    name = "entityrelationship",
    srcs = [
        "collection.go",
        "collection_builder.go",
        "entity.go",
        "fragments.go",
        "relationship.go",
    ],
    importpath = "outernetcouncil.org/nmts/v1/lib/entityrelationship",
    deps = [
        "//v1/proto:nmts_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

go_test(
    name = "entityrelationship_test",
    srcs = ["entity_test.go"],
    deps = [
        ":entityrelationship",
        "//v1/proto:nmts_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

# Benchmarks-only sibling of :entityrelationship_test, routed to the Pool=bench
# Buildbarn executors and published to Bencher. Written out longhand rather
# than via //bazel/bench:defs.bzl because nmts is a separate Bazel module
# (org_outernetcouncil_nmts) and cannot load from the parent workspace. Keep in
# sync with bazel/bench/defs.bzl there. Gazelle manages go_test rules by name
# and only generates ":entityrelationship_test", so it leaves this target alone
# -- but it also never updates these deps. Keep them in sync with
# :entityrelationship_test above.
go_test(
    name = "entityrelationship_bench",
    timeout = "eternal",
    srcs = ["entity_test.go"],
    args = [
        "-test.bench=.",
        "-test.run=^$$",
        "-test.benchtime=2s",
        "-test.count=1",
        "-test.benchmem",
    ],
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
        ":entityrelationship",
        "//v1/proto:nmts_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)
