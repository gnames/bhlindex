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

function development {
  migrate -database postgres://postgres@pg:5432/bhlindex?sslmode=disable -path db drop
  migrate -database postgres://postgres@pg:5432/bhlindex?sslmode=disable -path db up
  ginkgo watch
}


if [[ ! ${POSTGRES_HOST:?Requires POSTGRES_HOST} \
   || ! ${POSTGRES_USER:?Requires POSTGRES_USER} \
   ]]; then
  exit 1
fi

wait_for_db

development
