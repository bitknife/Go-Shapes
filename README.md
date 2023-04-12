# We The Forsaken (WTF!)

## Introduction
A super-massive PvP dungeon brawler, all copy-righted to Swedish company Bitknife AB.

The game consists of different client and server applications.

## Servers

### Main Go server (DGS)
The dedicated game-server is written in the language Go. Server developers should continue with [README.md](server/README.md)
for that sub-project.

For now it is a single process with an upper game and lower socket layer. We also have a "common/core"-layer in between
just to keep the two main layers apart. The core layer is responsible for distributing messages between the layers.

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
The [./shared](./common) directory contains the shared message model.