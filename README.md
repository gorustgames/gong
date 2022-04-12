# Gong

Classic pong game written in Go. Powered by [Ebiten](https://ebiten.org/) 2D game library for Go. Fun project to do some Go coding and to learn thing or two.

## How to play

Two games modes are available:

### Singleplayer

Use **A** and **Z** or **Arrow Up** & **Arrow Down**

### Multiplayer

* First player
  * Use **A** and **Z** or **Arrow Up** & **Arrow Down**
* Second player
  * Use **K** and **M**

## Where to get binaries


| Version     | Binary          | OS     |
|-------------|:---------------:|--------|
| 0.1.0       | [Download here](https://github.com/gorustgames/gong/releases/tag/v0.1.0) | Windows |
 
## Actor model
All games objects are modeled as actors. Actor is component that holds some internal state (e.g. position). Actor can basically do two things:

* update it's state
* draw it's state

Main game loop then basically iterates over all actors, updates them and draws them. Actor communicate with the game and with other actors using in-memory pub/sub bus (see next section). 

## Game bus
Game contains internal game bus for communication between actors. Inspired by:
* https://blog.logrocket.com/building-pub-sub-service-go/
* https://github.com/krazygaurav/pubsub-go
* https://eli.thegreenplace.net/2020/pubsub-using-channels-in-go/