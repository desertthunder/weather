"""Coverage is a utility module used to take Golang coverage reports and convert them to a markdown table."""

import datetime
import json
import os
import enum
import jinja2
import playwright
import playwright.sync_api


class Stream(enum.StrEnum):
    """Streams for the coverage report."""

    Markdown = "md"
    Html = "html"
    StdOut = "stdout"


def get_module_name():
    """Get the go module name."""
    with open("go.mod", "r") as f:
        lines = f.readlines()
        f.close()

    mod_line = lines[0]

    return mod_line.split(" ")[1].strip()


def get_coverage(cov_file: str = ".cov/coverage.txt"):
    """Open the coverage file and return the coverage data."""
    module_name = get_module_name()
    with open(cov_file, "r") as f:
        lines = f.readlines()
        f.close()

    coverage = {}

    for line in lines:
        if line.startswith(module_name):
            file_path, func_name, cov_pct = line.replace("\t\t", "\t").split("\t")

            filename, line_number, _ = file_path.split(":")

            if filename not in coverage:
                coverage[filename] = []

            coverage[filename] += [
                [
                    v.strip()
                    for v in [
                        line_number,
                        func_name,
                        cov_pct,
                    ]
                ]
            ]

    return coverage


def render_html(coverage: dict):
    """Render the coverage data to an HTML file."""
    j2_env = jinja2.Environment(
        loader=jinja2.FileSystemLoader("templates"),
        autoescape=True,
    )

    timestamp = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    template = j2_env.get_template("coverage.j2")
    html = template.render(coverage=coverage, timestamp=timestamp)

    with open("coverage.html", "w") as f:
        f.write(html)
        f.close()


def export_coverage(coverage: dict, stream: Stream = Stream.StdOut):
    """Export the coverage data to a markdown table."""
    if stream == Stream.Html:
        render_html(coverage)

        return json.dumps(coverage, indent=4)

    mdown = [
        "## Coverage Report",
    ]

    for k in coverage:
        mdown.append(f"### {k}\n")
        mdown.append("| Line | Function | Coverage |")
        mdown.append("|---|---|---|")

        for line in coverage[k]:
            ln, fn, cov = line

            mdown.append(f"| {ln} | {fn} | {cov} |")

        mdown.append("\n")

    if stream == "md":
        with open("coverage.md", "w") as f:
            f.write("".join(mdown))
            f.close()

    return "\n".join(mdown)


def screenshot():
    """Take a screenshot of the coverage report."""

    with playwright.sync_api.sync_playwright() as p:
        browser = p.chromium.launch()
        page = browser.new_page()
        root_path = f"file://{os.getcwd()}"
        page.goto(
            f"{root_path}/coverage.html",
        )

        page.locator("main").screenshot(path="assets/coverage.png")
        page.close()

        os.remove("coverage.html")


if __name__ == "__main__":
    import sys

    cov_file = sys.argv[1] if len(sys.argv) > 1 else ".cov/coverage.txt"
    coverage = get_coverage(cov_file)

    cov = export_coverage(coverage, stream=Stream.Html)

    screenshot()
