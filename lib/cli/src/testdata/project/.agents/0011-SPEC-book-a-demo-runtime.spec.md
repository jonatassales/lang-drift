# 0011: Book a Demo Runtime

**Status:** ready-for-agent

## Components

1. services/signup-api:8080
2. services/notify-api
3. platform/nats
4. platform/redis
5. apps/signup

## Acceptance criteria

- [ ] The five components start inside the envelope without port collisions.
- [ ] The runtime serves book-a-demo requests across signup, notification, queue and cache.
- [ ] Stopping the runtime releases all resources.
