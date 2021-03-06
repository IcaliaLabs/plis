[![Code Climate](https://codeclimate.com/repos/5713cd028daddd007c000a55/badges/23a412d4ad98ff7d79c6/gpa.svg)](https://codeclimate.com/repos/5713cd028daddd007c000a55/feed)
[![Test Coverage](https://codeclimate.com/repos/5713cd028daddd007c000a55/badges/23a412d4ad98ff7d79c6/coverage.svg)](https://codeclimate.com/repos/5713cd028daddd007c000a55/coverage)
[![Issue Count](https://codeclimate.com/repos/5713cd028daddd007c000a55/badges/23a412d4ad98ff7d79c6/issue_count.svg)](https://codeclimate.com/repos/5713cd028daddd007c000a55/feed)
![Made with Love by Icalia Labs](https://img.shields.io/badge/With%20love%20by-Icalia%20Labs-ff3434.svg)

# Plis

Helps your development process with Docker Compose by asking nicely :)

## Install

On macOS, install it via Homebrew:

```
brew tap icalialabs/formulae
brew install plis
```

On other systems you can:

 * Download the executable for your system from the
[Releases Page](https://github.com/IcaliaLabs/plis/releases)
 * Place it on any of the paths reachable in $PATH

## Special behaviors:
* `plis start [services-optional]`: Starts a docker-compose project, with the following extra functionality:
  * If some or all of the requested project's containers are missing, issues a `docker-compose up -d` command.
  * If all of the requested project's containers are present, issues a `docker-compose start` command.
* `plis attach [service_name]`: It figures out the given service's container, and attaches the console to it.
* `plis run [service_name] [command]`: It runs the given command:
  * If there's a running container for the given service, it executes it issuing a `docker exec -ti` command.
  * If there are no running containers for the given service, it executes it issuing a `docker-compose run --rm` command.

```bash

# Start a docker-compose project:
plis start

# Restart a service:
plis restart web

# Attach the console to a service:
plis attach web

# Run a command on an existing or new container:
plis run web rails c

# Stop a service:
plis stop web

```

## TODO's:
- [x] `build` command to invoke the docker-compose build command.
- [x] Split up the big `plis.go` file.
- [x] `check context` command to list the files that will pass to the Docker build context.
- [ ] `start` command with just one service should attach to the container immediately. (i.e.: `plis start web` starts a rails web container and attaches to it, mimicking the behavior of running `rails server` on the host)
- [ ] Change the `run` command to use `docker-compose exec` instead of `docker exec` whenever a running container is already available.
- [ ] Copy (from existing templates/examples) or generate blank dotenv files referenced in the Compose file.
- [ ] Install Docker (for Mac/Windows or native for Linux) if it is missing.
- [ ] Make `plis start github.com/some_org/some_dockerized_app` clone the project and run it.
- [ ] `upgrade` command that upgrades `plis` to the newest version.
- [ ] `prune` command to invoke `docker system prune`.
