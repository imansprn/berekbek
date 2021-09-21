# berekbek

A Golang project.

## Table of contents
* [How to contribute](#how-to-contribute)
* [Requirement](#requirement)
* [How to configure on your local machine](#how-to-configure-on-your-local-machine)
* [How to run migration](#how-to-run-migration)
* [Need more help?](#need-more-help)


## How to contribute

Read [CONTRIBUTING.md](https://github.com/gobliggg/berekbek/blob/master/CONTRIBUTING.md) to know more about how to contribute to this repo and how to deploy service. If you are new, it is mandatory to read this file first.

## Requirement

1. Go 1.16 or above (download [here](https://golang.org/dl/)).
2. Git (download [here](https://git-scm.com/downloads)).

## How to configure on your local machine

1. Clone this repository to your local.
   ```bash
   $ git clone git@github.com:gobliggg/berekbek.git
   ```
   
2. Add this command in the `~/.bashrc`

   ```bash
   $ export GOPRIVATE="github.com/valbury-repos"
   ```

3. Change working directory to `berekbek` folder.
   ```bash
   $ cd berekbek
   ```

4. Create configuration files.
   ```bash
   $ cp params/berekbek.toml.sample params/berekbek.toml
   ```

5. Edit configuration values in `params/berekbek.toml` according to your setting.

6. Running it then should be as simple as:
    ```bash
    $ make build
    $ ./bin/berekbek
    ```

## How to run migration

This migration can do these actions:

1. Migration up

   This command will migrate the database to the most recent version available. Migration files can be seen in this folder `migrations/sql/`.
   ```bash
   $ go run main.go migrate up
   ```

2. Migration down

   This command will undo/rollback database migration.
   ```bash
   $ go run main.go migrate down
   ```

3. Migration new

   This command will generate new migration file based on unix timestamp format
   ```bash
   $ go run main.go migrate new create_A_table
   ```

To get any help about these command, add `--help` at the end of the command.
```bash
$ go run main.go --help
$ go run main.go migrate up --help
$ go run main.go migrate down --help
$ go run main.go migrate new --help
$ go run main.go migrate force --help
```

## Need more help

If you need more help or anything else, please ask Valbury backend engineer team. We would be happy to help you.
