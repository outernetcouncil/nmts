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

"""Provides an nmts_gen_txtpbs rule for using a golang template with JSON
files to generate txtpb output files."""

def _nmts_gen_txtpbs_impl(ctx):
    outputs = []
    for f in ctx.files.srcs:
        infile = f
        outfile = ctx.actions.declare_file(infile.path.replace(".json", ".txtpb"))
        outputs.extend([outfile])

        inputs = [infile, ctx.file.template]
        arguments = [
            "--output",
            outfile.path,
            "--input",
            infile.path,
            "--tmpl_filename",
            ctx.file.template.short_path,
        ]
        for template in ctx.files.templates:
            inputs.append(template)
            arguments += ["--tmpl_filename", template.short_path]

        ctx.actions.run(
            outputs = [outfile],
            inputs = inputs,
            arguments = arguments,
            executable = ctx.executable._template2txtpb,
        )

    return DefaultInfo(
        files = depset(direct = outputs),
    )

_nmts_gen_txtpbs = rule(
    implementation = _nmts_gen_txtpbs_impl,
    attrs = {
        "template": attr.label(
            allow_single_file = [".tmpl"],
            doc = "The primary entry point template for the txtpb to be used for each src -> dst transformation",
            mandatory = True,
        ),
        "templates": attr.label_list(
            allow_files = [".tmpl"],
            doc = "A list of supplemental templates for the txtpb to be used for each src -> dst transformation",
            mandatory = False,
        ),
        "srcs": attr.label_list(
            allow_files = [".json"],
            doc = "A list of JSON files to be used as input",
            mandatory = True,
        ),
        "_template2txtpb": attr.label(
            default = "//v0/cmd/template2txtpb",
            cfg = "exec",
            executable = True,
        ),
    },
)

nmts_gen_txtpbs = _nmts_gen_txtpbs
