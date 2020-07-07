#!/bin/sh

# setup
docker run --name mapcleaner_db --rm \
 -e POSTGRES_PASSWORD=enter \
 -p 5432:5432 \
 postgres
