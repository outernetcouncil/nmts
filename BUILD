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

load("@gazelle//:def.bzl", "gazelle", "gazelle_test")

# For go module source files
# use the files as is to prevent getting out of sync
filegroup(
    name = "go_module_source",
    srcs = [
        "go.mod",
        "go.sum",
    ],
    visibility = ["//visibility:public"],
)

# gazelle:go_naming_convention_external import
gazelle(name = "gazelle")

gazelle_test(
    name = "gazelle_test",
    workspace = "//:BUILD",
)
