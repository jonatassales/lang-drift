# The demo-requests service, built from the repository root so Turborepo can
# prune the workspace to this app and the packages it depends on.
#
#   docker build -f apps/microservices/demo-requests/Dockerfile -t langdrift-demo-requests:local .

FROM node:24-alpine AS base
ARG PNPM_VERSION=12.4.2
ARG TARGETARCH
# pnpm as a checksum-pinned static binary: no npm, no corepack.
RUN set -eu; \
	case "$TARGETARCH" in \
		amd64) arch=x64; sha256=5ac22f6d1da756482490a10b3be1b929a62fa1022a5b034b37856fab609cbd6c ;; \
		arm64) arch=arm64; sha256=bf7597812c814030e759c20691dab88a1b271f495576fdbc4dae568458367105 ;; \
		*) echo "unsupported architecture: $TARGETARCH" >&2; exit 1 ;; \
	esac; \
	wget -qO /tmp/pnpm.tar.gz "https://github.com/pnpm/pnpm/releases/download/v${PNPM_VERSION}/pnpm-linux-${arch}-musl.tar.gz"; \
	echo "${sha256}  /tmp/pnpm.tar.gz" | sha256sum -c -; \
	tar -xzf /tmp/pnpm.tar.gz -C /usr/local/bin pnpm; \
	rm /tmp/pnpm.tar.gz
WORKDIR /repo

FROM base AS prune
COPY . .
RUN pnpm dlx turbo@2.11.7 prune demo-requests --docker

FROM base AS build
COPY --from=prune /repo/out/json/ .
RUN pnpm install --frozen-lockfile --ignore-scripts
COPY --from=prune /repo/out/full/ .
RUN pnpm turbo run build --filter=demo-requests
RUN pnpm --filter=demo-requests deploy --prod --legacy /out

FROM node:24-alpine
ENV NODE_ENV=production
ENV PORT=3402
WORKDIR /app
COPY --from=build --chown=node:node /out /app
USER node
EXPOSE 3402
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=6 \
	CMD wget -qO- http://127.0.0.1:3402/health || exit 1
CMD ["node", "dist/main.js"]
