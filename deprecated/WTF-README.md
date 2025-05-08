# We The Forsaken (WTF!)

## Deprecation note.
Project WTF not continued on this code-base.

The working bits of the proto-game was refactored/renamed to Go-Shapes-Server-Client. See new main README.md in root.

## Introduction
WTF is an attempt to make an online multiplayer dungeon brawler / roguelike.

## Go Game Source - wtf-main
The [](../wtf-main) directory contains both the go-based server and the client(s).

The first version of the client was writen using Ebitengine for Go.

The game is split into three main packages

- client
- server
- shared

We try to share as much logic as possible between the client and server as not to duplicate similar
concepts of the game.

### Client in Ebitengine
Links to Ebitengine and a list of very useful additions to it.

https://github.com/hajimehoshi/ebiten
https://github.com/sedyh/awesome-ebitengine

better camera?
https://github.com/MelonFunction/ebiten-camera

#### Suggestions
Based on nice presentation from experienced Ebitengine developer Quasilyte
https://speakerdeck.com/quasilyte/ebitengine-ecosystem-overview?slide=60
https://github.com/quasilyte/

Also note his "itch" page as well. 

Sprites/Anim:
https://github.com/yohamta/ganim8
https://github.com/tanema/gween

Map editor / Tiles etc.
- https://www.mapeditor.org/
- https://github.com/lafriks/go-tiled

Collision detection
- https://github.com/SolarLune/resolv
- https://github.com/ByteArena/box2d

Scripting (Go in Go)
Makes for dynamically loading of functions runtime. For Plugins or scripting objects etc.
- https://github.com/traefik/yaegi
- https://github.com/open2b/scriggo
See also Lua, Lisp etc for scripting.

Pathfinding, Quasilyte evaulated three of them but could recommend any of them.

2D Math
- https://github.com/quartercastle/vector (float64? slow)
- https://github.com/quasilyte/gmath (his own lib, compat. w ebitengine-input)

Input
He has good points on why to add another layer:
- https://github.com/quasilyte/ebitengine-input
Good practice to keep a virtual layer in between, ie Input -> Action, better let game depend on Action.

Signals/Slots
Good to reduce objects coupling, adds event listening to Go
- https://github.com/quasilyte/gsignal

Resource management
- https://github.com/quasilyte/ebitengine-resource

A game made by qasilyte, good to look for structure of larger project etc.
https://github.com/quasilyte/roboden-game
Also found his own ge "util" repo:
https://github.com/quasilyte/ge

Ebiten game template: https://github.com/sinisterstuf/ebitengine-game-template nice!
also contains github workflows for integrating with itch.io.

### Shared
The [./shared](./common) directory contains the shared message model.

### Server


## Ideas for future

### Message systems 
If we go for a distributed server solution, we need a hyper-fast messaging system to separate the game layer from the
socket-layer by a fast distributed messaging system (either a broker or broker-less):

Serverless go messaging library looks interesting, roll your own I guess?

- https://github.com/nanomsg/mangos

Full messaging platforms

- ZeroMQ, fast and lightweight
- RabbitMQ, have experience with this, hard to cluster
- Apache Pulsar: https://pulsar.apache.org/
- Apache Kafka, monster
- Redis, simple, but single threaded, said to be slower with higher loads

### Database
We will need to persist stuff at some point. I imagine an SQL-server will do just fine for this. We will let the Django
web server be the master for this. The game server will just load a few tables.

- Player accounts (username, password)
- High-scores.
- Player stuff (skins, gifts etc.)

### Web server
Not started, but I imagine the game having a web-server for user-registration/management and possibly a web-based game
client as well (built from the Unit project). A Django-based solution will be perfect for this.

Will share the SQL database with the game-server.

### Util servers
Expect to have log-server, metric-server etc.

## Other Clients (Maybe Deprecated!)

Outside of the wtf-main directory is the [](./clients) directory in where we may build additional game clients in 
other languages than Go.

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
