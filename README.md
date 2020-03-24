mackerel-plugin-cuenote-srs-status
===

![](https://github.com/ch1aki/mackerel-plugin-cuenote-srs-status/workflows/test/badge.svg)
![](https://github.com/ch1aki/mackerel-plugin-cuenote-srs-status/workflows/Release/badge.svg)

Cuenote SR-S custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-cuenote-srs --host=<host> --user=<username> --password=<password> [--group-stats] [--prefix=<prefix>] [--template=<tempfile>]
```

Options:

- `--host`: Cuenote SR-S hostname (e.g. `srsXXXX.cuenote.jp`)
- `--user`: Cuenote SR-S username
- `--password`: Cuenote SR-S password
- `--prefix`: metric key prefix (default: `cuenote-srs-stat`)
- `--group-stats`: Enable Grouped status (default: `false`)
- `--tempfile=`: Override tempfile path (default: mackerel default)

## Install

```shell
mkr plugin install ch1aki/mackerel-plugin-cuenote-srs-status@v0.0.1
```

## Example of mackerel-agent.conf

```toml
[plugin.metrics.cuenote-srs-status]
command = "/path/to/mackerel-plugin-cuenote-srs-status -H srsXXXX.cuenote.jp -u xxxx -p xxxxxxxx"
```

```toml
[plugin.metrics.cuenote-srs-status]
command = "/path/to/mackerel-plugin-cuenote-srs-status -H srsXXXX.cuenote.jp -u xxxx -p xxxxxxxx --group-stats"
```

## cuenote-srs-stat.queue_total

- cuenote-srs-stat.queue_total.delivering
- cuenote-srs-stat.queue_total.undelivered
- cuenote-srs-stat.queue_total.resend

## cuenote-srs.queue_group

- cuenote-srs-stat.queue_group.delivering.*
- cuenote-srs-stat.queue_group.undelivered.*
- cuenote-srs-stat.queue_group.resend.*
