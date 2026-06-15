# Cookie Analyser

Cookie Analyser is a simple CLI tool that analyses cookie usage. Currently, it identifies the most frequently occurring cookie(s) for a given day from a CSV file.

## Installation

### Requirements

- Go 1.20+

### Clone & Run

Clone the repo and cd into the root directory.

```bash
git clone https://github.com/tayy-b/cookies.git
go run ./cmd -d <date> -f <file>
```

### Example

```bash
go run ./cmd -d 2018-12-09 -f ./data/input.csv
```

### Available Flags:

```
  -f, --file string   Path to the CSV file containing cookie logs (required)
  -d, --date string   Date in YYYY-MM-DD format (UTC, required)
  -h, --help          Help
```

## Input Format

The tool expects a CSV file with the following structure:

```csv
cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00
```

- **timestamp**: ISO 8601 format (RFC3339)

---

### Example Output

```text
AtY0laUfhglK3lC7
```

If multiple cookies share the highest frequency:

```text
cookieA
cookieB
```

If no cookies are found for the given date:

```text
no cookies found on given day
```

---

## Assumptions

- The CSV file is sorted by timestamp in descending order (most recent first)
- The `-d` flag specifies a date in UTC
- The CSV file fits into memory
- Timestamps are in RFC3339 format

---

## Implementation Details

- **Language**: Go
- **CLI Framework**: Cobra
- **Logging**: Logrus

**Approach**:

- Load and parse all records from the CSV file
- Normalise timestamps to date-only (YYYY-MM-DD)
- Count occurrences of each cookie on the target date
- Stop early when records fall outside the target date
- Return all cookies with the highest frequency

The tool handles:

- Invalid date format (must be `YYYY-MM-DD`)
- Invalid timestamp format in CSV rows
- Missing required flags
- File read errors
- Malformed CSV rows
