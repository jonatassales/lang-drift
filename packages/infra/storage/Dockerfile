# The archive of the products: MinIO, built from its AGPL source at a
# pinned release, since MinIO no longer publishes community images or
# binaries. The HEALTHCHECK is what `dft` turns into the component's readiness
# probe (spec 0032). Root credentials arrive at run time as MINIO_ROOT_USER
# and MINIO_ROOT_PASSWORD, never here.
#
#   docker build -t langdrift-storage:local packages/infra/storage

FROM golang:1.25-alpine AS build
ARG MINIO_RELEASE=RELEASE.2025-10-15T17-29-55Z
RUN apk add --no-cache git
WORKDIR /src
RUN git clone --depth 1 --branch "${MINIO_RELEASE}" https://github.com/minio/minio.git . \
	&& CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X github.com/minio/minio/cmd.ReleaseTag=${MINIO_RELEASE}" -o /minio .

FROM alpine:3.22
COPY --from=build /minio /usr/bin/minio
RUN mkdir -p /data && chown nobody:nobody /data
USER nobody
VOLUME /data
EXPOSE 9000 9001
HEALTHCHECK --interval=5s --timeout=3s --retries=12 \
	CMD wget -qO- http://127.0.0.1:9000/minio/health/live || exit 1
CMD ["minio", "server", "/data", "--console-address", ":9001"]
