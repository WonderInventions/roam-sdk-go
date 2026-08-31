# Releasing `ro.am/roamhq`

Go modules are git tags. Merging a regeneration PR does not publish. Pushing
a `vX.Y.Z` tag does.

## One-time setup — vanity import

`go get ro.am/roamhq` resolves through a `go-import` meta tag that must be
served, permanently, from the production apex:

```
GET https://ro.am/roamhq?go-get=1

<meta name="go-import"
      content="ro.am/roamhq git https://github.com/WonderInventions/roam-sdk-go">
```

The path is a prefix: `ro.am/roamhq/client`, `ro.am/roamhq/option`, and
`ro.am/roamhq/webhooks` all resolve from the same tag. Once proxy.golang.org
caches it, dropping that route breaks every `go get`.

This route lives in the roam product, not this repository. The first tagged
release should not happen until it is live; without it, `go get ro.am/roamhq`
fails even if the tag exists.

Until then, a replace directive works for development:

```
go mod edit -replace=ro.am/roamhq=github.com/WonderInventions/roam-sdk-go@master
```

## Cutting a release

1. Merge the "Regenerate SDK from OpenAPI spec" PR (or confirm master is the
   tree you want).
2. Pick the bump:
   - new optional field or reworded description → patch
   - new endpoint, method, or optional parameter → minor
   - renamed or removed method, newly required field → major
3. Tag and push:

   ```bash
   git checkout master
   git pull
   git tag v0.1.0
   git push origin v0.1.0
   ```

4. pkg.go.dev indexes on its own. Force a refresh with:

   ```bash
   curl 'https://proxy.golang.org/ro.am/roamhq/@v/v0.1.0.info'
   ```

There is no nested `webhooks/` module. Signature verification is a package
inside this module (`ro.am/roamhq/webhooks`), so it versions with the parent
tag. Do not push a `webhooks/vX.Y.Z` tag.
