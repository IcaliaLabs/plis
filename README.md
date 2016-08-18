[![Code Climate](https://codeclimate.com/repos/5713cd028daddd007c000a55/badges/23a412d4ad98ff7d79c6/gpa.svg)](https://codeclimate.com/repos/5713cd028daddd007c000a55/feed)
[![Test Coverage](https://codeclimate.com/repos/5713cd028daddd007c000a55/badges/23a412d4ad98ff7d79c6/coverage.svg)](https://codeclimate.com/repos/5713cd028daddd007c000a55/coverage)
[![Issue Count](https://codeclimate.com/repos/5713cd028daddd007c000a55/badges/23a412d4ad98ff7d79c6/issue_count.svg)](https://codeclimate.com/repos/5713cd028daddd007c000a55/feed)

# Plis

Helps your development process with Docker Compose by asking nicely :)

## Install
See ([Releases Page](https://github.com/IcaliaLabs/plis/releases))

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
- [ ] `build` command to invoke the docker-compose build command.
- [ ] Change the `run` command to use `docker-compose exec` instead of `docker exec` whenever a running container is already available.
- [ ] Copy (from existing templates/examples) or generate blank dotenv files referenced in the Compose file.
- [ ] Install Docker (for Mac/Windows or native for Linux) if it is missing.
- [ ] Make `plis start github.com/some_org/some_dockerized_app` clone the project and run it.
- [ ] `upgrade` command that upgrades `plis` to the newest version.
