# QATARINA-CLI

This is a command line interface for the Qatarina (https://github.com/golang-malawi/qatarina) platform. It allows developers and testers to interact with the Qatarina API directly from the terminal - Managing projects, testcases, modules and many more.


## Tech Stack

The qatarina-cli is built with these wonderful technologies:

- Go 1.25
- Cobra for CLI
- Bubble tea for TUI 

## Building and running

In order to build and run use the following commands

```sh
$ go build

# place the binary in your $PATH e.g. sudo mv ./qatarina-cli /usr/local/bin/

$ qatarina-cli
```

## USAGE

# QATARINA CLI Developer Docs

## Authentication

Before using commands, you must log in. By default this uses the instance running on `http://localhost:4597`.
You can use the environment variable `QATARINA_HOST` to customize the server that the CLI interacts with.

```sh
$ qatarina-cli login --email user@example.com --password secret123


# OR Customize the host using environment variable

$ export QATARINA_HOST="https://qatarina.yourdomain.com/"
$ qatarina-cli login --email user@example.com --password secret123
```

Replace user@example.com and secret123 with our actual email and password
To log out:
```
$ qatarina-cli logout
```

Read the full doumentation [here](./docs/developer.md)

## Licence
MIT License
