# gogurt

[![Build Status](https://travis-ci.org/AlexandreCarlton/gogurt.svg?branch=master)](https://travis-ci.org/AlexandreCarlton/gogurt)

Downloads and builds static packages from source.

The successor to [`install-from-source`](https://github.com/AlexandreCarlton/install-from-source),
packaged in one, neat binary.

To install it:

```bash
  go install github.com/alexandrecarlton/gogurt/cmd/gogurt
```

## Caveats
For now, this will ignore `musl`, and install `gcc` to get `libstdc++.a`.
To this end, we use a Docker container to build these applications.

## Tips to making

`cmake -LAH` will list all cached variables. Useful when trying to override
settings to inject custom libraries.

## Related projects

 - [Sabotage](http://sabotage.tech/)
 - [Morpheus](http://git.2f30.org/ports/)
 - [minos](http://s.minos.io/)

# Todo:
 - Replicate stow functionality? Or just let stow do it.
   - Or add helper script to do it. Add a 'gogurt config' to get the values.
 - Add 'package' method to make tarball to download instead of build
   - Keep record of what can be downloaded from an archive, and what cannot (e.g. vim)
   - Or not, and just warn that depending on how you set it, things may break
     if you package them to other users.

## Testing
Look into `go test`, use jq for PoC since it is small.
We want to be able to check that we installed the correct things.

## Merging build and install
This was kept separate so we could elevate privileges when necessary.
However, we probably don't need to do this, since we choose where to install
it.
It might pay to optionally elevate if necessary, though.
Yep, that's why we'll keep it that way.

## Archives
All source releases should be extracted - they (SourceArchive) have a url that
needs to be downloaded, cached and extracted.
Note that, when extracting it, we'll need to strip out the first folder.

## Logging
All output from the configure/cmake/make steps should be redirected into a
logging file.
We'll configure this inside the install() funcs - more work, but more control,
too (can have 'steps', like 'configure', 'cmake', etc.)


## Server (ambitious)
We have a server which can communicate with other servers. When one receives a
multiple jobs it can farm them out to others. Great for running in kubernetes.

Or, instead, we have multiple agents (in containers), and we have a client that
knows about these servers. Then the client can coordinate instead, and forward
built archives to other agents if they require it.

Use raft protocol?
