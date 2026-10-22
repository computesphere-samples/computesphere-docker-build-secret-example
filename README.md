# Use a build secret in a Dockerfile on ComputeSphere

This sample shows how a Dockerfile reads a secret at build time without the secret ending up in the image.

## How it works

The Dockerfile mounts the secret for a single `RUN` step:

```dockerfile
RUN --mount=type=secret,id=BUILD_SECRET,required=true \
    test -s /run/secrets/BUILD_SECRET || { echo 'the BUILD_SECRET build secret is missing' >&2; exit 1; }; \
    sha256sum /run/secrets/BUILD_SECRET | cut -d' ' -f1 > /out/secret.sha256
```

- The secret is available at `/run/secrets/BUILD_SECRET` during that step only. It is never written to a layer.
- The `test -s` check (with `required=true`) makes the build fail if the secret isn't provided or is empty.
- The image keeps only the secret's SHA-256. At runtime the app compares it with the `EXPECTED_SHA256` variable you set and answers `Build secret check: matched` or `not matched`. It never shows the secret or its hash.

## Deploy it

Use a throwaway test value for this sample, never a real secret.

1. Create a web service from this repository. Set **Build method** to **Dockerfile**, and the port to `8080`. Its first build fails until `BUILD_SECRET` exists: that is the sample's own check.
2. Open the service's **Settings** tab and go to **Variables**. Choose **Edit**.
3. On the **Secrets** tab, use **Add secret** to add `BUILD_SECRET` with your test value, and set its scope to **Build only** (or **Build & runtime**). Secrets are **Runtime only** by default, and a runtime-only secret never reaches the build.
4. On the **General** tab, use **Add variable** to add `EXPECTED_SHA256` with the SHA-256 of the same test value, scoped **Runtime only**:

   ```sh
   printf %s "my-test-value" | sha256sum
   ```

5. Save with **Save & rebuild**. A change to a build-scoped variable needs a new build. The service then answers `Build secret check: matched`.

A secret whose scope includes the build is passed to a Dockerfile build as a BuildKit secret, with the variable's name as its `id`. Read it with `RUN --mount=type=secret,id=<NAME>`.

## Run it locally

```sh
printf %s "my-test-value" > secret.txt   # listed in .dockerignore, never part of the build context
docker build --secret id=BUILD_SECRET,src=secret.txt -t build-secret-example .
docker run --rm -p 8080:8080 -e EXPECTED_SHA256="$(printf %s "my-test-value" | sha256sum | cut -d' ' -f1)" build-secret-example
curl localhost:8080
```
