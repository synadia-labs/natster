---
title: Global Service
weight: 40
---
The global service does not manage any of your sensitive information other than holding onto the last personal access token you used so that it can make Synadia Cloud API requests on your behalf. We don't want any personal information and nearly all data stored is just reference keys and aggregate statistics to power our dashboard.

Your catalog contents are never queried, stored, or cached by this service.

The Natster global service is consumed by the `natster` CLI and nearly every function outlined in the API below can be invoked from the CLI.

All of the following API functions are prefixed with `natster.global`, so the `whoami` function is actually on `natster.global.whoami`.


| Subject | Description |
| :-- | :-- |
| `events.put` | The payload is a cloud event. If valid, the event is written to the appropriate log |
| `heartbeats.put` | Submits a heartbeat so the backend can keep track of which catalogs are online |
| `stats` | An empty request to this subject will return global summary stats, such as how many catalogs are online |
| `my.shares` | Requests the list of all catalogs known to the caller. This includes catalogs shared _by_ the caller and shared _to_ the caller. |
| `otc.generate` | Generates a one-time code that can be claimed by an authenticated web user to perform [context binding](../context-binding). These tokens will expire in 5 minutes, by default |
| `otc.claim` | Called by the Natster website to associate a valid OAuth ID with the context previously published on the `otc.generate` subject |
| `whoami` | Called by a valid user to obtain information about who Natster thinks they are. Very useful in troubleshooting. |
| `context.get` | Retrieves a bound context based on the OAuth ID attempting to log in. Only usable by the Natster.IO site account |
| `catalogs.validatename` | Used by the CLI to vet the names used when creating and sharing catalogs. Rejects names of catalogs that have already been created and shared, for example |