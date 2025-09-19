# CLI Executable Specification

## SPEC: Version Command

GIVEN the sz command line tool is available
WHEN the user runs `sz version`
THEN the output should display the current version number in the format "sz version X.Y.Z"

## SPEC: Help Command

GIVEN the sz command line tool is available
WHEN the user runs `sz help`
THEN the output should display help information including:
- The tool description "Distill the web into semantic markdown"
- Usage information
- Available commands

## SPEC: Default Behavior

GIVEN the sz command line tool is available
WHEN the user runs `sz` without any arguments
THEN the output should display the help information by default