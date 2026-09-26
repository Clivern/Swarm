# Copyright 2026 Swarm. All rights reserved.
# License can be found in the LICENSE file.
#
# Headless Pi on a mounted git repo (published as clivern/swarm on Docker Hub).
FROM node:24-bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g --ignore-scripts @earendil-works/pi-coding-agent

WORKDIR /repo

COPY entrypoint.sh /usr/local/bin/swarm-pi-entrypoint
RUN chmod +x /usr/local/bin/swarm-pi-entrypoint

ENTRYPOINT ["/usr/local/bin/swarm-pi-entrypoint"]
