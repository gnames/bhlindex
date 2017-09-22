#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

function wait_for_db {
  sleep 3

  while [[ $(pg_isready -h "${POSTGRES_HOST}" \
           -U "${POSTGRES_USER}") = "no response" ]]; do
    echo "Waiting for postgresql to start..."
    sleep 1
  done
}

function touch_bhlindex_db {
  psql -U ${POSTGRES_USER} -h ${POSTGRES_HOST} -tc "SELECT 1 FROM pg_database WHERE datname = 'bhlindex'" | grep -q 1 || psql -U ${POSTGRES_USER} -h ${POSTGRES_HOST} -c "CREATE DATABASE bhlindex"
}

function development {
  touch_bhlindex_db
  migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/bhlindex?sslmode=disable -path db drop
  migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/bhlindex?sslmode=disable -path db up
  ginkgo watch
}


if [[ ! ${POSTGRES_HOST:?Requires POSTGRES_HOST} \
   || ! ${POSTGRES_USER:?Requires POSTGRES_USER} \
   ]]; then
  exit 1
fi

wait_for_db

development
