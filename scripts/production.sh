#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

dir=$(dirname ${BASH_SOURCE[0]})

function touch_bhlindex_db {
  psql -U ${POSTGRES_USER} -h ${POSTGRES_HOST} -tc "SELECT 1 FROM pg_database WHERE datname = '${POSTGRES_DB}'" | grep -q 1 || psql -U ${POSTGRES_USER} -h ${POSTGRES_HOST} -c "CREATE DATABASE ${POSTGRES_DB}"
}

function drop_data {
  read -p "Delete all data on ${GNINDEX_HOST} (y/N)" -n 1 -r
  echo    # (optional) move to a new line
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
    echo "Removing old data from the database"
    migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db drop
  fi

}

function production {
  touch_bhlindex_db
  drop_data
  migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db up
  echo "You updated the schema for ${POSTGRES_DB}"
}


if [[ ! ${POSTGRES_HOST:?Requires POSTGRES_HOST} \
   || ! ${POSTGRES_USER:?Requires POSTGRES_USER} \
   || ! ${POSTGRES_DB:?Requires POSTGRES_DB} \
   ]]; then
  exit 1
fi

production
