# We The Forsaken (WTF!)

## Introduction
A super-massive PvP dungeon brawler, all copy-righted to Swedish company Bitknife AB.

The game consists of different client and server applications.

## Servers

### Main Go server
The main game-server is written in the language Go. Server developers should continue with [README.md](server/README.md)
for that sub-project.

### Web server
[TODO]

Not started, but I imagine the game having a web-server for user-registration/management and possibly a web-based game
client as well (built from the Unit project).

## Clients
### Main C# Unity game client
This is the main graphical game client, read about how to get started in the [README.md](clients/unity-client/README.md).

### Admin Python clients
This directory contains a set of "scripts" made for internal use mainly for:

- Performance testing.
- Integration testing.
- Remote admin and monitoring.

It uses the same set of message system used by the Unity client, but from a more controlled and scripted context for
automatic testing of the game-server, regarding performance etc.

May be simple console/prompt based applications, or something else, continue here: [README.md](clients/py-client/README.md).

## Common code
The [./common](./common) directory contains the shared message model.