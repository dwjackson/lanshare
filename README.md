# LANshare

It's weirdly difficult to get files from one computer to another, even when
they're on the same network. LANshare tries to mitigate that difficulty by
allowing you to serve a directory over HTTP. This is like Python's http.server
(formerly SimpleHTTPServer) or deno's file\_server.ts, but it's written in Go
because I wanted a single executable file that was trivially cross-compilable.

## Usage

In the directory which you want to allow files to be downloaded from, just run
the following:

```sh
$ lanshare
listening on :8080...
```

To set the maximum upload file size, pass the `-m` flag with an argument having
the form `[0-9]+(KiB|MiB|GiB)?`.

To set the port, use the `-p [PORT]` flag.

## License

LANshare is licensed under the MIT license.
