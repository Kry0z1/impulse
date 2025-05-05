# System prototype for biathlon competitions
## What is it
This app is build for parsing, transforming and aggregating information of triathlon events.
Input, output, log and config formats are described in [desc.md](./desc.md).

## Usage
```shell
git clone https://github.com/Kry0z1/impulse.git
cd impulse
go run -c {path_to_config} -o {path_to_output_file} -i {path_to_input_file} -l {path_to_log_file}
```

Paths to log and output files might be omitted - stdout will be used instead.

## Testing
```shell
cd tests
go test -v
```