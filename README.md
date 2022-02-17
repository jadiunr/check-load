[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/jadiunr/check-load)
![Go Test](https://github.com/jadiunr/check-load/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/jadiunr/check-load/workflows/goreleaser/badge.svg)

# Check Load

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu load average check is a [Sensu Check][6] that provides alerting and metrics for load averages. Metrics are provided in [nagios_perfdata](https://docs.sensu.io/sensu-go/latest/observability-pipeline/observe-schedule/collect-metrics-with-checks/#supported-output-metric-formats) format.

## Usage examples

```
Check load averages and provide metrics

Usage:
  check-load [flags]
  check-load [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -c, --critical string   Critical threshold for load averages (default "0.85,0.8,0.75")
  -h, --help              help for check-load
  -m, --metricsonly       Outputs only the metrics without checking the threshold.
  -r, --percpu            Divide the load averages by the number of CPUs
  -w, --warning string    Warning threshold for load averages (default "0.75,0.7,0.65")

Use "check-load [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add jadiunr/check-load
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/jadiunr/check-load].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-load
  namespace: default
spec:
  command: check-load -r -w 0.8,0.7,0.6 -c 0.9,0.8,0.7
  output_metric_format: nagios_perfdata
  output_metric_handlers:
  - influxdb
  subscriptions:
  - system
  runtime_assets:
  - jadiunr/check-load
  interval: 60
  publish: true
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the check-load repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/check-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/check-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu-community/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
