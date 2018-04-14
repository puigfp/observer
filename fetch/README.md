# observer/fetch

This sub-package contains all the functions related to the first concern: *fetching raw metrics by polling websites*

## How does it work ?

This package exports a CLI command object. When this command is launched:

- the configuration file is read and parsed

- a buffer used to temporarily store the metrics in memory is created

- a channel through which the metrics will be sent back is created

- a polling goroutine is launched for each website: each goroutine will poll the website and send back raw metrics through the channel

- a buffer emptying goroutine is launched: this goroutine will empty the raw metrics buffer and send its content to influxDB every 5 seconds

- the 'main' goroutine receives the metrics by listening to the channel and store those metrics in the buffer

## Folder structure

    fetch
    ├── README.md   // this file
    ├── cmd.go      // "main"
    ├── poll.go     // website polling related functions
    ├── store.go    // influxDB related functions
    └── type.go     // type declarations
