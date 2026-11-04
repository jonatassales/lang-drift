# The event log of the microservices: Redis with append-only persistence, so
# the stream survives a restart. The HEALTHCHECK is what `ex` turns into the
# component's readiness probe (spec 0032).
#
#   docker build -t langdrift-redis:local apps/microservices/redis

FROM redis:8.10.3-alpine
EXPOSE 6379
HEALTHCHECK --interval=5s --timeout=2s --retries=12 \
	CMD redis-cli ping | grep -q PONG || exit 1
CMD ["redis-server", "--appendonly", "yes", "--save", ""]
