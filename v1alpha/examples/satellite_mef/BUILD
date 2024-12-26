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

load("@org_outernetcouncil_nmts//v1alpha/bazel:rule.bzl", "nmts_graph")
load("@org_outernetcouncil_nmts//v1alpha/bazel:validation.bzl", "nmts_validation_test")

package(
    default_testonly = 1,
    default_visibility = ["//visibility:public"],
)

filegroup(
    name = "txtpbs",
    srcs = [
        "//satellite_mef/geo_satellites:txtpbs",
        "//satellite_mef/ground_stations:txtpbs",
        "//satellite_mef/pops:txtpbs",
        "//satellite_mef/sdn:txtpbs",
        "//satellite_mef/user_terminals:txtpbs",
    ],
)

nmts_validation_test(
    name = "txtpbs_test",
    srcs = [":txtpbs"],
)

nmts_graph(
    name = "satellite_mef",
    srcs = [":txtpbs"],
    html = ":satellite_mef.html",
    rankdir = "TB",
    svg = ":satellite_mef.svg",
)
