# ntpclient

a small little CLI thing that grabs the current time from any NTP server.

## install

needs:

- Go (version 1.11+)

```
go install github.com/nkcmr/ntpclient
```

## usage

`ntpclient` will simply grab the time from an NTP and report it alongside the local time. it defaults to using `time.google.com`, but can be told to use other servers. it can also output in JSON. here are some examples:

```
> ntpclient --help
USAGE: ntpclient [NTP Host] [options...]

OPTIONS:
        --json  Enable JSON Output

> ntpclient --json
{"local_time":"2018-12-15T20:08:46.916173Z","network_time":"2018-12-15T20:08:46.870379693Z","port":123,"server":"time.google.com"}

> ntpclient time.nist.gov
(server: time.nist.gov, port: 123)
local time:   2018-12-15T20:10:42.525361Z
network time: 2018-12-15T20:10:42.49277417Z

> ntpclient pool.ntp.org --json
{"local_time":"2018-12-15T20:12:56.598717Z","network_time":"2018-12-15T20:12:56.552408584Z","port":123,"server":"pool.ntp.org"}
```

## license

MIT