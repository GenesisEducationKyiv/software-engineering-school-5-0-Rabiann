# Selecting programming language for project

Status: Accepted<br>
Date: Jun 10 2025<br>
Author: Pavlo Novosolov (Rabiann)

## Context and Problem Statement
There is a need to select good programming language which is able to complete task. The task is building simple web-api application to allow users to subscribe on periodical weather information updates.


<!-- This is an optional element. Feel free to remove. -->
## Decision Drivers

* Application must be backend
* Language must be simple and performant
* Language must have good ecosystem but not being to dependent on it (for learning purposes developer should understand how each system component really works)
* Language must support lightweight concurrency out of the box

## Considered Options

* Golang

Golang is modern and powerful procedural system language which allows to build high-performant backend applications. It allows users to write simple code to solve simple tasks.

** Pros **
- Simple and easy to use
- Have lots of docs and learning resources
- Could be easily used for development in any text editor
- Has built-in coroutine mechanism, very easy to use
- Has a lot of web-frameworks, and standard-library `net/http` module is also good

** Cons **
- Doesn't have so large ecosystem
- Has pretty primitive syntax and external tools, many things must be implemented from scratch

* C#

C# is a language originally developed by Microsoft, now is developed by community. Has one most-popular open source web-framework .NET Core. Has full support of OOP.

** Pros **
 - Due to full OOP support, it's easy to build advanced architecture.
 - Has excellent support in IDEs.
 - Has extensive ecosystem.
 - Has a lot of docs.
 - Supports async

** Cons **
 - Not so performant
 - Pretty difficult and overbloated
 - Not very good choice for learning purposes, as many internal things are hidden away.

* Rust

Rust is a modern high-performant system language, is alternative for C++ in some industries. Pretty immature.

** Pros **
 - Very high performant
 - A lot of low-level abilities

** Cons **
 - Pretty hard, has many of its own concepts
 - Has immature ecosystem
 - Has small amount of docs
 - Overkill for simple web-appications in general

## Decision Outcome

Choosen option `Golang` because it's pretty simple and powerful. Also it has almost all features needed in solving task.

## Outcomes

 - Using built-in coroutines for email scheduler
 - Using other concurrency constructuins in cache f.e.
 - Difficulty of suporting codebase will grow as service will grow.
 - Many things will be implemented from the ground-up, instead of using automatic tools.
