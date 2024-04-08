<div align="center">

# Calcuroller
![code coverage badge](https://github.com/daneofmanythings/calcuroller/actions/workflows/tests.yml/badge.svg)

## Summary
A gRPC service intended to calculate the results of dice rolls

</div>

## Description
This project has a few parts to it. The heart of it is an interpreter
written in go (based on [Writing an interpreter in Go](https://interpreterbook.com/)),
which recognizes 'dice strings' and returns the associated metadata.
More on the specifics of those in the [Usage](#usage) section. A cli 
REPL is provided which can be used to familiarize yourself with the syntax.

There is also a gRPC server which uses the parser and can be run inside
a lightweight, secure docker container. This can easily be pinged from [grpcurl](https://github.com/fullstorydev/grpcurl),
a barebones supplied client which can be ran locally, or any service you like!

## Setup

- If you plan to run anything locally, make sure [go](https://go.dev/) is installed.
- install [docker](https://docs.docker.com/get-docker/).
- ensure the `make` utility is installed.

## Quick Start

- `make run-server-docker` to build and run the server inside a docker container
This should be all you need to get the server up and running.

- *optional* `make run-client-local` to run the client.
This will let you send requests to the server through the client using the command line.


## Usage

#### Makefile
All commands here should be preceeded by `make`:
- `gen_grpc`: Re-run protoc to generate the protobuf files.
- `build-repl`: builds the binary for the repl only.
- `run-repl`: runs `build-repl` and starts the application locally.
- `build-server-local`: Builds the binary for the server.
- `run-server-local`: Runs `build-server-local` and starts the application locally.
- `build-server-docker-multistage`: Builds the server inside a lightweight docker container.
- `run-server-docker`: Runs `build-server-docker-multistage` and starts the docker container on the "host" network.
- `test`: run tests for the whole project.
- `ping`: send a request through grpcurl to Roller.Roll on port 8080.
- `clean`: remove the temporary directory holding the built binaries.

#### Interpreter
The interpreter is a calculator that recognizes a new type of primitive that I've been referring to as a 'dice string.'
There are currently 5 modifiers implemented that will alter the calculations done on a dice expression. Those are:

- `qu#`: 'Quantity' The number of dice to be rolled.
- `mi#`: 'Minimum' The floor that any of the dice rolled in the expression can be.
- `ma#`: 'Maximum' The ceiling that any of the dice rolled in the expression can be.
- `kl#`: 'Keep Lowest' Keeps the # lowest dice rolled in the expression.
- `kh#`: 'Keep Highest' Keeps the # highest dice rolled in the expression.
- `[tag]`: The tag modifier. Has no influence on the roll, but it is tracked and returned with the associated expression in the metadata (covered in server usage).

Here is an example using some of the modifiers: `d12qu4mi2kh2[cold]`
This dice expression will roll 4 d12 dice. Lets say, for example, the rolls ended up being [8, 1, 6, 2].
The `1` will be turned into a `2` from the minimum modifier. Then the `8` and `6` will be selected from the
'keep highest' modifier, returning 14.

A dice string is made up of multiple dice expressions, such as `d20 + d10 - 5 + d4`. 
The interpreter can handle all standard arithematic operations: addition, subtraction, multiplication, integer division, modulo, and exponentiation.
Note that any exponent < 1 will return 1.

Limitations:
- The standard syntax `xdy` to represent rolling 'x' dice of size 'y' is not currently supported.
- Cannot use arithmatic expressions to determine dice or dice modifier size. ex: `d(1 + d4)`

Both of these will be implemented soon so you can roll some dice to see how many dice you roll inline with the dice roll!


#### Server API
There is currently a single service implemented in the gRPC, Roller, with two procedures, Ping and Roll.

Ping takes no arguements and returns the struct: `{ ping: pong }`

Roll takes data of the shape 
```json
{ dice_string: <string>, caller_id: <string> }
```
.
Upon recieving a request resulting in an error, it will return:
```json
{
    message: {
        status: {
            code: int32,
            message: <string>,
        }
    }
}

A successful request will return:
```json
{
    message: {
        data: {
            data: {
                literal: <string>,
                metadata: <json>, -- See below for the shape of this field
            }
            caller_id: <string>,
        },
    },
}
```
The metadata json is a map with <string> pointing to <DiceData>. Dice data has the following shape:
```json
{
    Literal: <string>,
    Tags: []<string>
    RawRolls: []<int>,
    FinalRolls: []<int>,
    Value: <int>,
}
```


## Licensing
This project is Liscensed under the MiT Liscence.
