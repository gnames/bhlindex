#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

function wait_for_db {
  sleep 3

  while [[ $(pg_isready -h "${RACKAPP_DB_HOST}" \
           -U "${RACKAPP_DB_USERNAME}") = "no response" ]]; do
    echo "Waiting for postgresql to start..."
    sleep 1
  done
}

function development {
  cd /app

  bundle exec rake db:drop
  bundle exec rake db:create
  bundle exec rake db:migrate
  bundle exec rake db:migrate RACK_ENV=test
  ASSET_HOST="webpack:8080" puma config.ru -b "tcp://0.0.0.0:9292"
}

function production {
  cd /app && npm run build && mkdir -p public
  cp -r ./dist/* ./public && rm -R ./dist/*

  bundle exec rake db:migrate
  RACK_ENV=production puma -C config/docker/puma.rb
}


if [[ ! ${RACKAPP_DB_HOST:?Requires RACKAPP_DB_HOST} \
   || ! ${RACKAPP_DB_USERNAME:?Requires RACKAPP_DB_USERNAME} \
   || ! ${RACK_ENV:?Requires RACK_ENV} ]]; then
  exit 1
fi

wait_for_db

if [[ ${RACK_ENV} = "production" ]]; then production; else development; fi
