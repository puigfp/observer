# observer/process

This sub-package contains all the functions related to the second concern: *processing raw metrics and computing more interesting metrics and alerts*

## How does it work ?

This package exports a CLI command object. When this command is launched:

- the configuration file is read and parsed

- several goroutines are launched to process, the metrics one for each (window length, refresh rate) combination we want to support

  - the supported (window length, refresh rate) combinations are

    - (2m, 5s), metrics are stored in the `metrics_2m` influxDB table

    - (10m, 10s), metrics are stored in the `metrics_10m` influxDB table

    - (1h, 1m), metrics are stored in the `metrics_1h` influxDB table

  - basic metrics, such as mean, max, min and percentile are computed using `SELECT [...] INTO [...]` influxDB queries

  - `success=true` and `success=false` counts are also computed using `SELECT [...] INTO [...]` influxDB queries

  - more complicated metrics, such as status strings count, are computed manually, and require several queries (the code can get quite long because of the way the influxDB golang client returns the data: a lot of preprocessing and validation needs to be done manually)

- the main goroutine performs edge detection on the metrics computed over 2m long timeframes to detect when a website become available or unavailable

## Folder structure

    process
    ├── README.md // this file
    ├── alerts.go // alerts related functions
    ├── cmd.go    // "main" file
    ├── loop.go   // metrics/alerts computation loops
    ├── simple.go // "simple" (avg/min/max/count) metrics related functions
    ├── status.go // status strings count computation related functions
    ├── type.go   // type declarations
    └── util.go   // misc functions
