#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

dir=$(dirname ${BASH_SOURCE[0]})

function wait_for_db {
  sleep 3

  while [[ $(pg_isready -h "${POSTGRES_HOST}" \
           -U "${POSTGRES_USER}") = "no response" ]]; do
    echo "Waiting for postgresql to start..."
    sleep 1
  done
}

function touch_bhlindex_db {
  psql -U ${POSTGRES_USER} -h ${POSTGRES_HOST} -tc "SELECT 1 FROM pg_database WHERE datname = '${POSTGRES_DB}'" | grep -q 1 || psql -U ${POSTGRES_USER} -h ${POSTGRES_HOST} -c "CREATE DATABASE ${POSTGRES_DB}"
}

function development {
  touch_bhlindex_db
  ${dir}/migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db drop
  ${dir}/migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db up
  ginkgo watch -v
}


if [[ ! ${POSTGRES_HOST:?Requires POSTGRES_HOST} \
   || ! ${POSTGRES_USER:?Requires POSTGRES_USER} \
   || ! ${POSTGRES_DB:?Requires POSTGRES_DB} \
   ]]; then
  exit 1
fi

wait_for_db

development
