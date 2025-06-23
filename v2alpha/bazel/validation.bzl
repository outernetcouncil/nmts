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

"""Provides a simple nmts_validation_test rule for ensuring an NMTS graph is
valid."""

_script_template = """\
#!/bin/bash
{} validate {}
"""

def _nmts_validation_test_impl(ctx):
    # Emit the executable shell script.
    script = ctx.actions.declare_file("%s-script" % ctx.label.name)
    script_content = _script_template.format(
        ctx.executable._nmtscli.short_path,
        " ".join([f.short_path for f in ctx.files.srcs]),
    )
    ctx.actions.write(script, script_content, is_executable = True)

    runfiles = ctx.runfiles([ctx.executable._nmtscli] + ctx.files.srcs)
    return DefaultInfo(
        executable = script,
        runfiles = runfiles,
    )

_nmts_validation_test = rule(
    implementation = _nmts_validation_test_impl,
    test = True,
    attrs = {
        "srcs": attr.label_list(
            allow_files = [".txtpb"],
            doc = "A list of txtpb files to be validated as a whole",
        ),
        "_nmtscli": attr.label(
            default = "//v2alpha/cmd/nmtscli",
            cfg = "exec",
            executable = True,
        ),
    },
)

nmts_validation_test = _nmts_validation_test
