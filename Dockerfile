# Copyright 2026 Swarm. All rights reserved.
# License can be found in the LICENSE file.
#
# Headless Pi on a mounted git repo.
#
#   docker build -t swarm:v0.8.0 .
#
#   docker run --rm \
#     -e OPENROUTER_API_KEY \
#     -e PROMPT='your task' \
#     -e PI_MODEL=openrouter/anthropic/claude-sonnet-4.5 \
#     -v /path/to/git-clone:/repo \
#     -v /path/to/output:/out \
#     swarm:v0.8.0
FROM node:24-bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g --ignore-scripts @earendil-works/pi-coding-agent

WORKDIR /repo

COPY entrypoint.sh /usr/local/bin/swarm-pi-entrypoint
RUN chmod +x /usr/local/bin/swarm-pi-entrypoint

ENTRYPOINT ["/usr/local/bin/swarm-pi-entrypoint"]
