from pathlib import Path
import re
import subprocess


_STABLE_TAG_PATTERN = re.compile(r"^v\d+\.\d+\.\d+$")


def _latest_stable_tag() -> str:
    result = subprocess.run(
        [
            "git",
            "tag",
            "--merged",
            "HEAD",
            "--list",
            "v[0-9]*",
            "--sort=-version:refname",
        ],
        cwd=Path(__file__).resolve().parent.parent,
        check=True,
        capture_output=True,
        text=True,
    )
    for line in result.stdout.splitlines():
        tag = line.strip()
        if _STABLE_TAG_PATTERN.fullmatch(tag):
            return tag
    raise RuntimeError("no stable git tag found")


def define_env(env) -> None:
    env.variables["latest_stable_tag"] = _latest_stable_tag()
