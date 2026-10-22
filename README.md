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
- The app stores only the secret's SHA-256 and serves it at `/`. You can check which value the build used without exposing it.

## Deploy it

1. Create a web service from this repository. Set **Build method** to **Dockerfile**, and the port to `8080`.
2. Open the service's **Settings** tab and go to **Variables**. Choose **Edit variables**, open the **Secrets** tab and use **Add variable** to add `BUILD_SECRET`. Set its scope to **Build only** (or **Build & runtime**). Secrets are **Runtime only** by default, and a runtime-only secret never reaches the build.
3. Deploy. Then compare the hash the service returns with your value's hash:

   ```sh
   printf %s "$BUILD_SECRET" | sha256sum
   ```

A secret whose scope includes the build is passed to a Dockerfile build as a BuildKit secret, with the variable's name as its `id`. Read it with `RUN --mount=type=secret,id=<NAME>`.

## Run it locally

```sh
printf %s "my-test-value" > secret.txt
docker build --secret id=BUILD_SECRET,src=secret.txt -t build-secret-example .
docker run --rm -p 8080:8080 build-secret-example
curl localhost:8080
```
