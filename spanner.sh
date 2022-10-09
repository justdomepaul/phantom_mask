#!/usr/bin/env bash

export SPANNER_EMULATOR_HOST="localhost:9010"

PROJECT_DEFAULT=test-project
INSTANCE_DEFAULT=test-instance
DATABASE_DEFAULT=test-database

printf "enter the spanner project name (default: test-project): "
read PROJECT
PROJECT=${PROJECT:-$PROJECT_DEFAULT}

printf "enter the spanner instance name (default: test-instance): "
read INSTANCE
INSTANCE=${INSTANCE:-$INSTANCE_DEFAULT}

printf "enter the spanner database name (default: test-database): "
read DATABASE
DATABASE=${DATABASE:-$DATABASE_DEFAULT}

docker run -ti --network="host" --rm justdomepaul/spanner-cli -p $PROJECT -i $INSTANCE -d $DATABASE
