# Building
Instructions for how to build the software for BookPi.

## System Requirements
In order to build the project, it is recommended you use a Raspberry Pi.
Any version should work, but to get it done in a reasonable amount of time, I'd recommend using a Raspberry Pi 4.
If you do not need to build the display manager, then you should be able to build all the components on any machine.

## Build Structure
Each component has its own `Makefile` for building, installing dependencies, and cleaning generated files.
There is also a project-level `Makefile` which orchestrates the operations on each project.

### Build
This exists in each of the components as well as the top-level `Makefile`.
As the name states, it is used for building/compiling the components or project.
The biggest different when running at the project-level is that it automatically copies the built frontend to the server folder.
This eliminates some friction for compiling the server due to the embedding of the frontend in the compiled binary.

### Clean
This exists in each of the components as well as the top-level `Makefile`.
It is used for removing all the files generated by the build command.
This will irrevocably delete any dependencies or generated files, so use with caution.

### Dist
This command only exists in the top-level `Makefile` as it depends on the building of each component.
After each component gets built, all the generated and static files are copied to a folder to be compressed into a tarball.
This makes it easy to use for distribution.

## Inter-component Dependencies
The server component requires the frontend to be built before the server is compiled.
This is due to the static file embedding of the frontend within the server.
The frontend does not need to be built every time you compile the server, the `build` directory simply needs to be copied into the `server` folder.
Then you must run `go generate .` to create all the statically embedded files.

This manual process is only required if you choose not to use the top-level `Makefile`. 
