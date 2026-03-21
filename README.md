# mealie-markdown-exporter

A CLI tool that exports recipes from a [Mealie](https://mealie.io) instance as markdown files with YAML front matter.

## Installation

```sh
go install github.com/iamkirkbater/mealie-markdown-exporter@latest
```

Or build from source:

```sh
go build -o dist/mealie-markdown-exporter .
```

## Usage

```sh
mealie-markdown-exporter export --base-url https://your-mealie-instance --api-token your-token
```

### Flags

- `--base-url` (required) - Base URL of your Mealie instance
- `--api-token` (required) - API token for authentication
- `--output-dir` - Output directory (default: `mealie-markdown-export`)
- `--log-level` - Log level: `debug`, `info`, `warn`, `error` (default: `info`)

### Environment Variables

All flags can be set via environment variables with the `MME_` prefix. Flags take precedence over environment variables.

- `MME_BASE_URL`
- `MME_API_TOKEN`
- `MME_OUTPUT_DIR`
- `MME_LOG_LEVEL`

### API Token from File

To avoid exposing your API token in shell history, you can load it from a file:

```sh
mealie-markdown-exporter export --base-url https://your-mealie-instance --api-token file:///path/to/token
```

Or use the environment variable:

```sh
export MME_API_TOKEN=file:///path/to/token
```

## Running Tests

```sh
go test ./...
```
