# psusage ![gha build](https://github.com/karantan/psusage/workflows/Go/badge.svg)
Cross-platform process cpu % and memory usage of a program.

Inspired by https://github.com/struCoder/pidusage but I've decided to only get the
information from the [`ps`](https://man7.org/linux/man-pages/man1/ps.1.html) tool and
not parse `/proc/<pid>/stat` file, because this is already done by the `ps` tool.

> /proc/[pid]/stat
>    Status information about the process.  This is used by
>    [ps(1)](https://man7.org/linux/man-pages/man1/ps.1.html).

Ref: https://man7.org/linux/man-pages/man5/proc.5.html

With this change the code should be simpler and easiler to hack.

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
