# psusage ![gha build](https://github.com/karantan/psusage/workflows/Go/badge.svg)

Tool for monitoring CPU usage for a program (and it's forks).

Inspired by https://github.com/struCoder/pidusage but I've decided to only get the
information from the [`ps`](https://man7.org/linux/man-pages/man1/ps.1.html) tool and
not parse `/proc/<pid>/stat` file, because this is already done by the `ps` tool.

> /proc/[pid]/stat
>    Status information about the process.  This is used by
>    [ps(1)](https://man7.org/linux/man-pages/man1/ps.1.html).

Ref: https://man7.org/linux/man-pages/man5/proc.5.html

With this change the code should be simpler and easiler to hack.


## Key concepts and definitions

We use the same logic as AWS (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/burstable-credits-baseline-concepts.html).

### CPU credit

A unit of vCPU-time.

Examples:

1 CPU credit = 1 vCPU * 100% utilization * 60 seconds.
1 CPU credit = 6000 vCPU utilisation seconds.

So one server has 8_640_000 vCPU utilisation seconds (i.e. CPU credits) per day per CPU.
Around this number one must balance out processes and their CPU usage on a server.

## How it works

Basically run `ps -o pcpu=,time=,pid=,user:32=,comm= $(pidof <program>)` and add process this
information. Example of an output for program `php-fpm`

```
 0.1 00:00:00 3477510 root myprogram
36.5 00:01:23 3477549 foo_com myprogram
10.1 00:00:21 3477680 bar_com myprogram
 0.2 00:00:00 3477884 baz_com myprogram
```

We update cpu utilization of the process every second and when the process no longer
exists (i.e. it stopped) we send CPU credit used to the InfluxDB.

From there we can then do all sort of aggregation because we have all the information we
need:

1. Average CPU utilisation (in %)
2. Duration (in seconds)
3. Program name (e.g. `php-fpm`)
4. Effective user (e.g. `foo_com`)


## (Known) Limitations
TODO

## Testing

Use `stress` tool. It imposes a configurable amount of CPU, memory, I/O, and disk stress
on the system.

## Build the package with nix

Run the following command:

```
$ nix-build -E "with import <nixpkgs> {}; callPackage ./default.nix {}"
```

The version of the program will be passed during the build process via buildFlagsArray
in the `default.nix`.

## Usage
TODO

## Installation
TODO

## Development
TODO


## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request
