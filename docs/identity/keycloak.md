# Keycloak group refresh on CLI login

The Infra server refreshes OIDC group membership synchronously when approving
`infra login` through a browser session. It queries the provider's UserInfo
endpoint before approving the device flow, even if that session was synchronized
within the normal 120-minute interval. A failed synchronization prevents approval.
New OIDC browser logins also synchronize groups as before.

This handles both newly granted and removed group membership. Existing CLI
sessions retain the normal refresh behavior; run `infra login` again after changing
groups to explicitly synchronize. The connector does not need a code change.

## Keycloak configuration

Configure a Group Membership mapper on the Infra OIDC client's client scope:

- Token Claim Name: `groups`.
- Add to UserInfo: enabled.
- Ensure the client scope is assigned to the Infra client (as a default scope,
  or as an optional `groups` scope supported by discovery).
- Match Full group path to Infra group names: `/cluster-edit` and `cluster-edit`
  are different names.

Infra reads groups from UserInfo, so putting groups only in an access or ID token
is insufficient. The client must also return an email and support refresh tokens.

## Container publication

The `cd / containers` workflow tests the server and builds the existing core
Dockerfile for Linux amd64 and arm64. Publication also requires the SSH connector
acceptance suite on Ubuntu, Debian and Red Hat UBI, including grant and revocation
checks. The UI image is published separately to `ghcr.io/muroed/infra-ui`. Pushes to `main`, `v*` tags and manual runs
publish to `ghcr.io/muroed/infra` using the built-in `GITHUB_TOKEN`; no Docker Hub
credentials are required. Pull requests test and build without publishing.
Tags include `sha-<full commit SHA>`, `latest` on main and the Git tag on releases.
Prefer the SHA tag or image digest for deployment. The image starts `infra server`
by default and also contains the connector command, as in the upstream image.

For unauthenticated pulls, set the GHCR package visibility to public after the
first publication, or configure an image pull secret for a private package.
The workflow publishes an image; it does not deploy to a cluster.

## Verify after deploying the server image

1. Use an existing browser session and an existing user with no cluster access.
2. Add that user to the Keycloak group mapped to Kubernetes `edit` access.
3. Run `infra login`, approve in the existing browser session, then `infra list`.
   The cluster should appear without waiting for the periodic sync interval.
4. Remove the membership and repeat login; access should disappear.
5. Verify a revoked Keycloak session cannot approve a new CLI login.

If membership is still missing, inspect the UserInfo response for that client
without logging or sharing tokens. Confirm that `groups` contains the exact Infra
group name and that the group has a grant for the expected cluster.
