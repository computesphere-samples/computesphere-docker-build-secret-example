# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -o /out/server .

# The build secret is mounted only for this one RUN step, at
# /run/secrets/SMOKE_SECRET. It never becomes part of a layer: only its
# SHA-256 is written to the image, so a deployed app can show which value
# the build saw without revealing it. The build fails if the secret is missing.
RUN --mount=type=secret,id=SMOKE_SECRET,required=true \
    test -s /run/secrets/SMOKE_SECRET || { echo 'the SMOKE_SECRET build secret is missing' >&2; exit 1; }; \
    sha256sum /run/secrets/SMOKE_SECRET | cut -d' ' -f1 > /out/secret.sha256

FROM scratch
COPY --from=build /out/server /server
COPY --from=build /out/secret.sha256 /secret.sha256
EXPOSE 8080
ENTRYPOINT ["/server"]
