# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details

FROM alpine as builder
ARG TARGETPLATFORM

WORKDIR /

RUN --mount=target=/build tar xf /build/dist/helm-mirror_*_$(echo ${TARGETPLATFORM} | tr '/' '_' | sed -e 's/arm_/arm/').tar.gz
RUN cp helm-mirror /usr/bin/helm-mirror

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /usr/bin/helm-mirror /usr/bin/helm-mirror
USER 65532:65532

ENTRYPOINT ["/usr/bin/helm-mirror"]