---
title: Serving a Catalog
weight: 50
---

Serving a catalog involves starting the Natster catalog service. You can run this service anywhere that it has physical access to the contents of the media catalog. Thanks to the global connectivity offered by NGS, you can run this server on your laptop, in the cloud, or even on a Raspberry Pi if the mood suits.

## Starting a Catalog Server
Starting the server is easy.

```
natster catalog serve myvideos
```

Here **myvideos** is the name of a previously created (and hopefully shared) catalog. You'll see output similar to the following once it starts up:

```
2024/03/07 09:10:14 INFO Established Natster NATS connection servers=tls://connect.ngs.global
2024/03/07 09:10:14 INFO Opened Media Catalog name=myvideos rootpath=/home/user/myvideos/
2024/03/07 09:10:14 INFO Natster Media Catalog Server Started
2024/03/07 09:10:14 INFO Local (private) services are available on 'natster.local.>'
```

With this server running, only those people to whom you have shared this catalog will be able to access it. No one else can
query its contents or the files within.

## Interacting with the Catalog Service
When a catalog server starts, it subscribes to a set of topics that are then exposed via account exports in your NGS account. This is how other users with their own separate Synadia Cloud accounts can have secure access to the catalog.

In addition to making the service available globally, it is also available locally, allowing you to interact with it. As you'll see in the [web UI](../website) section, you can even securely see the contents of your catalog from the [natster.io](https://natster.io) website.

One incredibly handy thing you might want to do is query the contents of your catalog, as seen by those who have access. To do this, simply use the Natster CLI to query your catalog's contents:

```
natster catalog contents sample
╭──────────────────────────────────────────────────────────────────╮
│ Items in Catalog kevbuzz                                         │
├──────────────────────────────────────────────────────────────────┤
│ Path                │ Hash                │ Mime Type            │
├──────────────────────────────────────────────────────────────────┤
│ bookvideos/book.mp4 │ d1b1ca5d0beb...cb67 │ video/mp4            │
│ README.txt          │ df2e9d5d7745...27b9 │ text/plain           │
╰──────────────────────────────────────────────────────────────────╯
```
Here the output has been condensed a bit to make it easily readable in the documentation, but you can see that the catalog contents are identical to what is defined in the catalog's JSON file.

Now that you've shared a catalog with someone and the catalog server is running, it's time to take a look at how to securely import shared catalogs.