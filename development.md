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

# Docker Implementation
Docker is the primary way the application is debugged and compiled for use by others. This method insures the latest tools are
being deployed to build as well as there are no outside package differences from the developer's machine to compile againts.
There is a docker file included in the application to all this to be run as a container to enable it to be isolated and to run on different systems with minimal issues. To build the container use the following command: `docker build -t wwum2020:<version> .` the "version" can be removed if you are just building it locally.

To run a container from an image then use the following command `docker run -it -v $(pwd)/bin/CropSimOutput:/app/CropSimOutput -v $(pwd)/bin/OutputFiles:/app/OutputFiles wwum2020:1.2.4`. This command does the following things:

- runs the container in "-it" interactive mode so you have the command line from bash
- adds a bind mount to the local system for the CropSimOutput to the location within the container that is required
- adds another bind mount for the OutputFiles so when they are created those files are already on the host system for futher processes

If you are using podman for RHEL, Fedora, or another distro, use the `--privledged` to allow the container to access the bind mount local system files. For mac I use: `docker run -it -v /Users/hkuntz/Documents/OutputFiles:/app/OutputFiles -v /Users/hkuntz/Documents/CropSimOutput:/app/CropSimOutput wwum2020`

# Local Compiling
If you wish to compile and run the program locally. You can do this through the terminal using a few of the following commands.

These commands look like:
- `go run main.go runModel --CSDir "<path>"`
    - "--CSDir" is the path to the monthly output CropSim .txt files, you must qualify it in "<path>" if the path contains spaces
        - My example is `--CSDir "<path>/WWUMM2020/CropSim/Run005_WWUM2020/Output"`
- `go run main.go runModel --CSDir "<path>" --debug`