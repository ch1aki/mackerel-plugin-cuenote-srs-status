mackerel-plugin-cuenote-srs-status
===

Cuenote SR-S custom metrics plugin for mackerel.io agent.

## Synopsis

```sh
mackerel-plugin-cuenote-srs -H <host> -u <username> -p <password>
```

## Example of mackerel-agent.conf

```
[plugin.metrics.cuenote-srs-status]
command = "/path/to/mackerel-plugin-cuenote-srs-status -H srsXXXX.cuenote.jp -u xxxx -p xxxxxxxx"
```

## cuenote-srs.Queue

- cuenote-srs.Queue.delivering
- cuenote-srs.Queue.undelivered
- cuenote-srs.Queue.resend
