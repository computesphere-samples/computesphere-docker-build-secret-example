# Use a build secret in a Dockerfile on ComputeSphere

This sample shows how a Dockerfile reads a secret at build time without the secret ending up in the image.

## How it works

The Dockerfile mounts the secret for a single `RUN` step:

```dockerfile
RUN --mount=type=secret,id=SMOKE_SECRET,required=true \
    test -s /run/secrets/SMOKE_SECRET || { echo 'the SMOKE_SECRET build secret is missing' >&2; exit 1; }; \
    sha256sum /run/secrets/SMOKE_SECRET | cut -d' ' -f1 > /out/secret.sha256
```

- The secret is available at `/run/secrets/SMOKE_SECRET` during that step only. It is never written to a layer.
- The `test -s` check (with `required=true`) makes the build fail if the secret isn't provided or is empty.
- The app stores only the secret's SHA-256 and serves it at `/`. You can check which value the build used without exposing it.

## Deploy it

1. Create a web service from this repository with the **Docker** build mode, on port `8080`.
2. Add a secret variable named `SMOKE_SECRET` and make it available at **build** time.
3. Deploy. Then compare the hash the service returns with your value's hash:

   ```sh
   printf %s "$SMOKE_SECRET" | sha256sum
   ```

A secret variable that's available at build time is passed to a Docker build as a BuildKit secret, with the variable's name as its `id`.

## Run it locally

```sh
printf %s "my-test-value" > secret.txt
docker build --secret id=SMOKE_SECRET,src=secret.txt -t build-secret-example .
docker run --rm -p 8080:8080 build-secret-example
curl localhost:8080
```
