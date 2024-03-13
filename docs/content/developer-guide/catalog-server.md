---
title: Catalog Server
weight: 60
---

The catalog server is a daemon process that you can run anywhere that has direct access to the files that make up a given catalog. This means you could run it on a Raspberry Pi connected to network storage or you could run it in the cloud pointing at an EBS volume, or, more commonly, you could run the service in your own infrastructure to ensure your files are never stored outside your own environment.

The catalog server is started via the Natster CLI with the `catalog serve` command. It takes the name of an existing catalog as a parameter. Once running, it utilizes the imports and exports we've previously discussed to expose the catalog as a service.

## Catalog Service API
The catalog server exposes a fairly simple API for accessing the contents of the hosted catalog. As with all APIs in the Natster ecosystem, this is a core NATS API. The functions listed below are also used by the `natster` CLI so you won't need to access them directly unless you're troubleshooting.

Unless otherwise indicated, all service endpoints described below are assumed to be prefixed with `natster.catalog.`, so the `catalog.*.get` function is really `natster.catalog.{catalog}.get`.

| Subject | Description |
| :-- | :-- |
| `{catalog}.get` | Obtains the contents of the given catalog. If the catalog is not running, the caller will receive the typical "no responders" reply. If the caller does not have permission (the catalog has not been shared with them), this will return an envelope indicating unauthorize |
| `{catalog}.download` | Requests the download of a file indicated by its hash. If the download is approved, the caller may download the encrypted file chunks on `natster.media.{catalog}.{hash}` |
| `natster.local.inbox` | A function only exposed _inside_ the account running the catalog server. Used to query a list of pending shares, e.g. catalogs that have been shared with the caller but not yet imported. This is a queue subscription so only one catalog server will handle this request |