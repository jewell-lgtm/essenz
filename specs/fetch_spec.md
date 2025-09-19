# Fetch Command Executable Specification

## SPEC: Fetch Command with Local File Support

GIVEN a local HTML file exists
WHEN the user runs `sz fetch /path/to/file.html`
THEN the output should display the file contents

## SPEC: Fetch Command with HTTP URL Support

GIVEN a valid HTTP URL
WHEN the user runs `sz fetch http://example.com`
THEN the output should display the HTTP response content

## SPEC: Fetch Command with HTTPS URL Support

GIVEN a valid HTTPS URL
WHEN the user runs `sz fetch https://example.com`
THEN the output should display the HTTPS response content

## SPEC: Fetch Command Error Handling - File Not Found

GIVEN a file path that does not exist
WHEN the user runs `sz fetch /nonexistent/file.html`
THEN the command should exit with error and show helpful message

## SPEC: Fetch Command Error Handling - Invalid URL

GIVEN an invalid URL
WHEN the user runs `sz fetch invalid-url`
THEN the command should exit with error and show helpful message

## SPEC: Fetch Command Usage Help

GIVEN the user needs help with the fetch command
WHEN the user runs `sz fetch --help`
THEN the output should show usage examples for both files and URLs