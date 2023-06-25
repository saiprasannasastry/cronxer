# Cron Parser
This command-line application parses a cron string and expands each field to show the times at which it will run. It takes a cron string as input and outputs a formatted table with the expanded schedule.

# Installation
Clone the repository: git clone https://github.com/saiprasannasastry/cron-parser.git
Navigate to the project directory: cd cron-parser

```
go build.
```

# Usage
Run the application with the following command:

```
./cron-parser "cron_string"
```

Replace "cron_string" with the actual cron string you want to parse. The cron string should follow the standard cron format with five time fields (minute, hour, day of month, month, and day of week) plus a command.

# Example
For example, if you want to parse the cron string `"*/15 0 1,15 * 1-5 /usr/bin/find"`, you can run the following command:

```
./cron-parser "*/15 0 1,15 * 1-5 /usr/bin/find"
```
The output will be formatted as a table with the field name taking the first 14 columns and the times as a space-separated list following it:

```
minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```

# Benchmarking
```

BenchmarkParse-12    	  186549	      6144 ns/op	    2574 B/op	      88 allocs/op
BenchmarkValidateCronString-12    	 2420242	       453.6 ns/op	     224 B/op	       6 allocs/op


```
# Notes
- This application only supports the standard cron format with five time fields and a command. Special time strings such as "@yearly" are not supported.
- The cron string should be provided as a single argument enclosed in quotes.
- The application does not rely on existing cron parser libraries but implements its own logic to parse and expand the cron schedule.
