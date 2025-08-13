# Go

## Demonstration

> **Only for Golang backend developers**

## What is that?

Binary which create a ready project to run a server with preinstalled dependencies/tools:

- [x] `gin`
- [x] `viper`
- [x] `pgxpool`/`mysql`/`go-sqlite3`
- [x] `jwt`
- [x] `docker`
- [x] `linters`
- [x] `Makefile`
- [x] `best practice architecture`
- [ ] `postman collection`

## Why?

Every time when you need to create a project, need to:

- create an architecture
- intstall dependencies, like: `gin`, `viper`, `pgx`, `jwt`, etc...
- create an instancies of all dependencies: `app.Init()`, `db.Init()`, `config.Init()`, etc...
- create `ping` route
- ...

and other things. It's really hard to do it every time and it's very annoying. So this project
will help you to prevent routen

## Contributions

We will happy to get a help from you

- Just fork it
- Modify code
- Create pull request
