#!/usr/bin/env bash

set -e
cd $( cd $(dirname $0) ; pwd -P )

function usage () {
	cat <<EOS
To execute pre-defined commands with Docker.
Usage:
	$(basename $0) <Command> [args...]
Command:
EOS
	egrep -o "^\s*function.*#cmd.*" $(basename $0) | sed "s/^[ \t]*function//" | sed "s/[ \(\)\{\}]*#cmd//" \
	    | awk '{CMD=$1; $1=""; printf "\t%-16s%s\n", CMD, $0}'
}

function help() { #cmd help
    docker-compose run --rm spanner-migrate
}

function migrate() { #cmd migrate
    docker-compose run --rm spanner-migrate $@
}

function upDefault() { #cmd upDefault
  PROJECT=test-project
  INSTANCE=test-instance
  DATABASE=test-database

	migrate -path=/migrations/ -database "spanner://projects/$PROJECT/instances/$INSTANCE/databases/$DATABASE?x-clean-statements=true" up
}

function up() { #cmd up
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

	migrate -path=/migrations/ -database "spanner://projects/$PROJECT/instances/$INSTANCE/databases/$DATABASE?x-clean-statements=true" up
}

function down() { #cmd up
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

	migrate -path=/migrations/ -database "spanner://projects/$PROJECT/instances/$INSTANCE/databases/$DATABASE?x-clean-statements=true" down
}

function version() { #cmd version
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

	migrate -path=/migrations/ -database "spanner://projects/$PROJECT/instances/$INSTANCE/databases/$DATABASE?x-clean-statements=true" version
}

function goto() { #cmd goto
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

  printf "enter the spanner goto version:"
  read VERSION
  if [ -z "$VERSION" ]
  then
      echo "spanner goto version not found"
      exit 2
  fi
  printf "\n"
	migrate -path=/migrations/ -database "spanner://projects/$PROJECT/instances/$INSTANCE/databases/$DATABASE?x-clean-statements=true" goto $VERSION
}

function force() { #cmd force
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

  printf "enter the spanner force version:"
  read VERSION
  if [ -z "$VERSION" ]
  then
      echo "spanner force version not found"
      exit 2
  fi
  printf "\n"
	migrate -path=/migrations/ -database "spanner://projects/$PROJECT/instances/$INSTANCE/databases/$DATABASE?x-clean-statements=true" force $VERSION
}

function generate() { #cmd generate $@
    printf "enter the generate table name:"
	  read TABLE
    migrate create -ext sql -dir migrations $TABLE
}

if [ $# -eq 0 ] ; then
	usage
else
	export COMPOSE_HTTP_TIMEOUT=600
	CMD=$1
	shift
	$CMD $@
fi
