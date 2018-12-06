#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

dir=$(dirname ${BASH_SOURCE[0]})

function drop_data {
  read -p "Delete all data from ${POSTGRES_DB} on ${POSTGRES_HOST} (y/N)" -n 1 -r
  echo    # (optional) move to a new line
  if [[ $REPLY =~ ^[Yy]$ ]]
  then
    echo "Removing old data from the database"
    ${dir}/migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db drop
  fi

}

function production {
  drop_data
  ${dir}/migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db up
  echo "You updated the schema for ${POSTGRES_DB}"
  time ${dir}/bhlindex find -w 15
  time ${dir}/bhlindex verify -w 5
}


if [[ ! ${POSTGRES_HOST:?Requires POSTGRES_HOST} \
   || ! ${POSTGRES_USER:?Requires POSTGRES_USER} \
   || ! ${POSTGRES_DB:?Requires POSTGRES_DB} \
   ]]; then
  exit 1
fi

production
