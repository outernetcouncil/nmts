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

"""Rules for rendering NMTS graphs."""

def _nmts_to_dot_impl(ctx):
    dot_file = ctx.actions.declare_file(ctx.label.name)

    to_dot_args = ctx.actions.args()
    to_dot_args.add(ctx.executable._nmtscli)
    to_dot_args.add(ctx.attr.rankdir)
    to_dot_args.add(dot_file)
    to_dot_args.add_all(ctx.files.srcs)
    ctx.actions.run_shell(
        outputs = [dot_file],
        inputs = ctx.files.srcs,
        tools = [ctx.executable._nmtscli],
        progress_message = "Generating %{output} from an NMTS graph (%{input})",
        arguments = [to_dot_args],
        command = '''
        tool="$1"; shift
        rankdir="$1"; shift
        dst="$1"; shift

        "$tool" export dot --rankdir "$rankdir" "$@" > "$dst"
    ''',
    )

    return DefaultInfo(files = depset([dot_file]))

_nmts_to_dot = rule(
    implementation = _nmts_to_dot_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "A list of txtpb files to use as the graph definition",
        ),
        "rankdir": attr.string(
            doc = "rankdir attribute for the graph",
            default = "LR",
        ),
        "_nmtscli": attr.label(
            default = "//v0/cmd/nmtscli",
            cfg = "exec",
            executable = True,
        ),
    },
)

def _nmts_to_d2_impl(ctx):
    d2_file = ctx.actions.declare_file(ctx.label.name)

    to_d2_args = ctx.actions.args()
    to_d2_args.add(ctx.executable._nmtscli)
    to_d2_args.add(d2_file)
    to_d2_args.add_all(ctx.files.srcs)
    ctx.actions.run_shell(
        outputs = [d2_file],
        inputs = ctx.files.srcs,
        tools = [ctx.executable._nmtscli],
        progress_message = "Generating %{output} from an NMTS graph (%{input})",
        arguments = [to_d2_args],
        command = '''
        tool="$1"; shift
        dst="$1"; shift
        "$tool" export d2 "$@" > "$dst"
    ''',
    )

    return DefaultInfo(files = depset([d2_file]))

_nmts_to_d2 = rule(
    implementation = _nmts_to_d2_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "A list of txtpb files to use as the graph definition",
        ),
        "_nmtscli": attr.label(
            default = "//v0/cmd/nmtscli",
            cfg = "exec",
            executable = True,
        ),
    },
)

def _nmts_to_prolog_impl(ctx):
    pl_file = ctx.actions.declare_file(ctx.label.name)

    to_pl_args = ctx.actions.args()
    to_pl_args.add(ctx.executable._nmtscli)
    to_pl_args.add(pl_file)
    to_pl_args.add_all(ctx.files.srcs)
    ctx.actions.run_shell(
        outputs = [pl_file],
        inputs = ctx.files.srcs,
        tools = [ctx.executable._nmtscli],
        progress_message = "Generating %{output} from an NMTS graph (%{input})",
        arguments = [to_pl_args],
        command = '''
        tool="$1"; shift
        dst="$1"; shift
        "$tool" export prolog "$@" > "$dst"
    ''',
    )

    return DefaultInfo(files = depset([pl_file]))

_nmts_to_prolog = rule(
    implementation = _nmts_to_prolog_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "A list of txtpb files to use as the graph definition",
        ),
        "_nmtscli": attr.label(
            default = "//v0/cmd/nmtscli",
            cfg = "exec",
            executable = True,
        ),
    },
)

def _nmts_to_nquads_impl(ctx):
    nquads_file = ctx.actions.declare_file(ctx.label.name)

    to_nquads_args = ctx.actions.args()
    to_nquads_args.add(ctx.executable._nmtscli)
    to_nquads_args.add(nquads_file)
    to_nquads_args.add_all(ctx.files.srcs)
    ctx.actions.run_shell(
        outputs = [nquads_file],
        inputs = ctx.files.srcs,
        tools = [ctx.executable._nmtscli],
        progress_message = "Generating %{output} from an NMTS graph (%{input})",
        arguments = [to_nquads_args],
        command = '''
        tool="$1"; shift
        dst="$1"; shift
        "$tool" export nquads "$@" > "$dst"
    ''',
    )

    return DefaultInfo(files = depset([nquads_file]))

_nmts_to_nquads = rule(
    implementation = _nmts_to_nquads_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "A list of txtpb files to use as the graph definition",
        ),
        "_nmtscli": attr.label(
            default = "//v0/cmd/nmtscli",
            cfg = "exec",
            executable = True,
        ),
    },
)

def _nmts_to_svg_impl(ctx):
    svg_file = ctx.actions.declare_file(ctx.label.name)
    to_svg_args = ctx.actions.args()
    to_svg_args.add(ctx.file.src)
    to_svg_args.add(svg_file)
    to_svg_args.add(ctx.executable._embedstyles)

    ctx.actions.run_shell(
        inputs = [ctx.file.src],
        outputs = [svg_file],
        tools = [ctx.executable._embedstyles],
        arguments = [to_svg_args],
        progress_message = "Running dot on %{input} to generate %{output}",
        # TODO: figure out how to get graphviz into our dependency tree
        command = '''
    dot -Tsvg "$1" | "$3" > "$2"
    ''',
    )

    run_tmpl = ctx.actions.declare_file(ctx.label.name + "run_tmpl.txt")
    ctx.actions.write(run_tmpl, "cat \"@@FILE@@\"", is_executable = True)

    run_script = ctx.actions.declare_file(ctx.label.name + "_run.sh")
    ctx.actions.expand_template(template = run_tmpl, output = run_script, substitutions = {"@@FILE@@": svg_file.short_path}, is_executable = True)

    return DefaultInfo(files = depset([svg_file]), executable = run_script, runfiles = ctx.runfiles(files = [svg_file]))

_dot_to_svg = rule(
    implementation = _nmts_to_svg_impl,
    attrs = {
        "src": attr.label(
            allow_single_file = [".dot"],
            mandatory = True,
            doc = "A list of txtpb files to use as the graph definition",
        ),
        "_embedstyles": attr.label(
            default = "//v0/cmd/embedstyles",
            cfg = "exec",
            executable = True,
        ),
    },
    executable = True,
)

def _nmts_to_html_impl(ctx):
    out_file = ctx.actions.declare_file(ctx.label.name)

    args = ctx.actions.args()
    args.add(ctx.executable._nmtscli)
    args.add(out_file)
    args.add_all(ctx.files.srcs)
    ctx.actions.run_shell(
        outputs = [out_file],
        inputs = ctx.files.srcs,
        tools = [ctx.executable._nmtscli],
        progress_message = "Generating %{output} from an NMTS graph (%{input})",
        arguments = [args],
        command = '''
        tool="$1"; shift
        dst="$1"; shift
        "$tool" export html "$@" > "$dst"
    ''',
    )

    run_tmpl = ctx.actions.declare_file(ctx.label.name + "run_tmpl.txt")
    ctx.actions.write(run_tmpl, "cat \"@@FILE@@\"", is_executable = True)

    run_script = ctx.actions.declare_file(ctx.label.name + "_run.sh")
    ctx.actions.expand_template(template = run_tmpl, output = run_script, substitutions = {"@@FILE@@": out_file.short_path}, is_executable = True)

    return DefaultInfo(files = depset([out_file]), executable = run_script, runfiles = ctx.runfiles(files = [out_file]))

_nmts_to_html = rule(
    implementation = _nmts_to_html_impl,
    executable = True,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "A list of txtpb files to use as the graph definition",
        ),
        "_nmtscli": attr.label(
            default = "//v0/cmd/nmtscli",
            cfg = "exec",
            executable = True,
        ),
    },
)

def _check_output_test_impl(ctx):
    diff_template = ctx.actions.declare_file("%s_diff_tmpl.txt" % ctx.label.name)
    ctx.actions.write(diff_template, """#!/usr/bin/env bash
main() {
  local target="@@TARGET@@"
  local want_file="@@WANT_FILE@@"
  local got_file="@@GOT_FILE@@"

  if ! diff "$want_file" "$got_file";
  then
    echo >&2 "$got_file doesn't match the generated output from the nmtscli tool."
    echo >&2 "Run the following to regenerate the in-tree output:"
    echo >&2 "bazel run $target > \\$PWD/$got_file"
    return 1
  fi
}

main "$@"
""")

    target_parts = ctx.file.want.short_path.rsplit("/", 1)

    diff_script = ctx.actions.declare_file("%s_diff.sh" % ctx.label.name)
    ctx.actions.expand_template(template = diff_template, output = diff_script, substitutions = {
        "@@WANT_FILE@@": ctx.file.want.short_path,
        "@@GOT_FILE@@": ctx.file.got.short_path,
        "@@TARGET@@": "//%s:%s" % (target_parts[0], target_parts[1]),
    }, is_executable = True)

    diff_output = ctx.actions.declare_file("%s_diff.txt" % ctx.label.name)
    ctx.actions.run(outputs = [diff_output], inputs = [ctx.file.want, ctx.file.got], executable = diff_script)

    return DefaultInfo(executable = diff_script, runfiles = ctx.runfiles(files = [ctx.file.want, ctx.file.got]))

_check_output_test = rule(
    implementation = _check_output_test_impl,
    test = True,
    attrs = {
        "want": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "got": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
    },
)

def nmts_graph(name, srcs, svg = None, html = None, **kwargs):
    """Defines a number of rules for processing a given NMTS graph into a variety of different output formats.

    Args:
        name: The basename to use in the rules. This will be the name of the rendered SVG, but will also be used as the basename of the output files (i.e. "${NAME}.dot")
        srcs: The .txtpb files that make up the NMTS graph.
        **kwargs: Extra arguments to pass to the underlying rules.
    """
    rankdir = kwargs.pop("rankdir", "LR")

    _nmts_to_nquads(name = "%s.nquads" % name, srcs = srcs, **kwargs)
    _nmts_to_dot(name = "%s.dot" % name, srcs = srcs, rankdir = rankdir, **kwargs)
    _nmts_to_d2(name = "%s.d2" % name, srcs = srcs, **kwargs)
    _nmts_to_prolog(name = "%s.pl" % name, srcs = srcs, **kwargs)

    rendered_svg = "%s_rendered.svg" % name
    rendered_html = "%s_rendered.html" % name
    _dot_to_svg(name = rendered_svg, src = "%s.dot" % name)
    _nmts_to_html(name = rendered_html, srcs = srcs, **kwargs)

    if svg:
        _check_output_test(
            name = "%s_svg_test" % name,
            want = rendered_svg,
            got = svg,
        )

    if html:
        _check_output_test(
            name = "%s_html_test" % name,
            want = rendered_html,
            got = html,
        )

    # we want the rule to be buildable but the filename to have the .svg
    # extension, so alias the rule name to the svg target
    native.alias(name = name, actual = rendered_svg)
