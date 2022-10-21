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
    docker-compose run --rm postgresql-migrate
}

function migrate() { #cmd migrate
    docker-compose run --rm postgresql-migrate $@
}

function upDefault() { #cmd upDefault
  PROJECT=test-project
  INSTANCE=test-instance
  DATABASE=test-database

	migrate -path=/migrations/ -database "postgres://phantom:phantom@postgresql:5432/phantom?sslmode=disable" up
}

function up() { #cmd up
	migrate -path=/migrations/ -database "postgres://phantom:phantom@postgresql:5432/phantom?sslmode=disable" up
}

function down() { #cmd up
		migrate -path=/migrations/ -database "postgres://phantom:phantom@postgresql:5432/phantom?sslmode=disable" down
}

function version() { #cmd version
 	migrate -path=/migrations/ -database "postgres://phantom:phantom@postgresql:5432/phantom" version
}

function goto() { #cmd goto
   printf "enter the postgresql goto version:"
  read VERSION
  if [ -z "$VERSION" ]
  then
      echo "postgresql goto version not found"
      exit 2
  fi
  printf "\n"
	migrate -path=/migrations/ -database "postgres://phantom:phantom@postgresql:5432/phantom" goto $VERSION
}

function force() { #cmd force
  printf "enter the postgresql force version:"
  read VERSION
  if [ -z "$VERSION" ]
  then
      echo "postgresql force version not found"
      exit 2
  fi
  printf "\n"
	migrate -path=/migrations/ -database "postgres://phantom:phantom@postgresql:5432/phantom" force $VERSION
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
