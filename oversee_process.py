import subprocess
import sys
from pathlib import Path

PROJECT_DIR = Path(__file__).resolve().parent


def run_command(command: list[str], description: str) -> bool:
    """Run one workflow step and return True only when it succeeds."""
    print(f"===== {description.upper()} =====")

    try:
        subprocess.run(
            command,
            cwd=PROJECT_DIR,
            check=True,
        )
    except FileNotFoundError as error:
        print(f"Error: command not found: {error.filename}")
        return False
    except subprocess.CalledProcessError as error:
        print(
            f"Error: {description} failed "
            f"with exit code {error.returncode}."
        )
        return False

    print(f"{description} completed successfully.")
    return True


def main() -> int:
    steps = [
        (
            ["/usr/local/go/bin/go", "run", "."],
            "Fetching roster data",
        ),
        (
            [sys.executable, "sheet_injection.py"],
            "Updating Google Sheets",
        ),
    ]

    for command, description in steps:
        if not run_command(command, description):
            return 1

    return 0


if __name__ == "__main__":
    raise SystemExit(main())