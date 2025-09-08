# Golang Architecture Maker

https://github.com/user-attachments/assets/d996b362-8f06-46a3-851d-5b8c8559afb4

## What is that?

Binary which create a ready project to run a server with preinstalled dependencies/tools:

- ✅ `gin` / `fiber`
- ✅ `viper`
- ✅ `pgxpool`/`mysql`/`go-sqlite3`
- ✅ `jwt`
- ✅ `docker`
- ✅ `linters`
- ✅ `Makefile`
- ✅ `best practice architecture`
- ❌ `postman collection`

## Why?

Every time when you need to create a project, need to:

- create an architecture
- intstall dependencies, like: `gin`, `viper`, `pgx`, `jwt`, etc...
- create an instancies of all dependencies: `app.Init()`, `db.Init()`, `config.Init()`, etc...
- create `ping` route
- ...

and other things. It's really hard to do it every time and it's very annoying. So this project
will help you to prevent routine

## Contributions

We will happy to get a help from you

- Just fork it
- Modify code
- Create a pull request
