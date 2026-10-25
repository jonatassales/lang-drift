# syntax=docker/dockerfile:1
# Builds the website from the repository root: docker build -f apps/website/Dockerfile .
# node:24-alpine
FROM node:24-alpine@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53

WORKDIR /repo

ENV NEXT_TELEMETRY_DISABLED=1
ENV PNPM_HOME=/pnpm
ENV PATH="$PNPM_HOME:$PATH"

# pnpm comes from its npm tarball, verified against the registry's published hash.
ARG PNPM_VERSION=12.4.2
ARG PNPM_SHA512=08adc6613180275c7c9edada39dcf08c9c61ad4e7eaf330a4f3461f102b0f907423454d117f98e72d47fef0616070644d7bffc973a6a57f5090a6d7c368b07c9
RUN wget -qO /tmp/pnpm.tgz "https://registry.npmjs.org/pnpm/-/pnpm-${PNPM_VERSION}.tgz" \
  && echo "${PNPM_SHA512}  /tmp/pnpm.tgz" | sha512sum -c - \
  && mkdir -p /pnpm/lib \
  && tar -xzf /tmp/pnpm.tgz -C /pnpm/lib --strip-components=1 \
  && ln -s /pnpm/lib/bin/pnpm.cjs /pnpm/pnpm \
  && rm /tmp/pnpm.tgz

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json ./
COPY apps/website ./apps/website
COPY packages/tokens ./packages/tokens
COPY packages/web ./packages/web

# The repository's install scripts set up Git hooks and the drift CLI, which an image does not need.
RUN pnpm install --frozen-lockfile --ignore-scripts

# The Docs zone upstream; next.config.ts refuses to build without it.
ARG DOCS_ORIGIN
ENV DOCS_ORIGIN=$DOCS_ORIGIN
ARG WEBSITE_PREVIEW_SIGNUP_DISABLED=true
RUN pnpm exec turbo run build --filter=website...

ENV NODE_ENV=production
WORKDIR /repo/apps/website

EXPOSE 3000

CMD ["pnpm", "start"]
