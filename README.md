# gogurt
Downloads and builds static packages from source.
Where possible, packages found already built (e.g. `go`) are downloaded
instead.

The successor to `install-from-source`, packaged in one, neat binary.

For now, this will ignore musl, and installing gcc to get libstdc++.a.
To this end, we use a Docker container to build these applications.
