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

load("@org_outernetcouncil_nmts//v2alpha/bazel:gen_txtpbs.bzl", "nmts_gen_txtpbs")

package(
    default_testonly = 1,
    default_visibility = ["//visibility:public"],
)

nmts_gen_txtpbs(
    name = "txtpbs",
    srcs = glob(["*.json"]),
    template = "user_terminal_gen1.txtpb.tmpl",
)
