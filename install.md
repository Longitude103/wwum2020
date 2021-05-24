# GO Installation

To install a program, use the binary file.

Install the binary file to `/usr/bin` and it will work without change to path since
`/usr/bin` is in the path.

If you put the binary file in another location you can add that location to the
path by running `export PATH=$PATH:/path/to/your/install/directory` to add it.

You can see what is in the path by running `echo $PATH`

In "GO" you can get the binary by running `go build` and if you want a different name
you can use `go build -o <name>` for a different file name.

There are several potential paths that can be used, but this is good for Linux, might be others for Windows.