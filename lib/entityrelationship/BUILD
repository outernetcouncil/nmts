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
    importpath = "outernetcouncil.org/nmts/lib/entityrelationship",
    deps = [
        "//proto:nmts_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

go_test(
    name = "entityrelationship_test",
    srcs = ["entity_test.go"],
    deps = [
        ":entityrelationship",
        "//proto:nmts_go_proto",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)
