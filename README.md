eggo is an opinionated ACME client for dynamic SSL certificate management.

Based on a Redis pubsub queue, eggo will submit requests for an SSL certificate from the configured ACME server, and upon success will store it into a backend storage system for retrieval by other tools.

## setup

eggo relies on environmental variables for its configuration, either directly from the environment or from a `.env` file.
See `.env.example` for a current list of environmental variables.

## limitations

eggo is in early development, and currently has very limited flexibility. eggo is currently built to specifically work with the ACME flow for **ZeroSSL** and currently has hardcoded storage protocol support.

currently only new certificate issuance is supported; renewals and revocation will be added in the future.
