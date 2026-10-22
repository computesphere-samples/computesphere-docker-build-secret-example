# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -o /out/server .

# The build secret is mounted only for this one RUN step, at
# /run/secrets/BUILD_SECRET. It never becomes part of a layer: only its
# SHA-256 is written to the image. At runtime the app compares it with
# EXPECTED_SHA256 and answers only matched / not matched; it never shows the
# secret or the hash. The build fails if the secret is missing.
RUN --mount=type=secret,id=BUILD_SECRET,required=true \
    test -s /run/secrets/BUILD_SECRET || { echo 'the BUILD_SECRET build secret is missing' >&2; exit 1; }; \
    sha256sum /run/secrets/BUILD_SECRET | cut -d' ' -f1 > /out/secret.sha256

FROM scratch
COPY --from=build /out/server /server
COPY --from=build /out/secret.sha256 /secret.sha256
EXPOSE 8080
ENTRYPOINT ["/server"]
