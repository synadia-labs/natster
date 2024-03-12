---
title: Creating a Media Catalog
weight: 20
---

A media catalog is a collection of logically related media files. Their association with each other could be very specific, or you could have a simple catalog that contains all of the media you wish to share. Since catalogs are the unit of sharing, you'll want to use your sharing plans to influence how many catalogs you have and what they contain.


## Locate or Place Media Files
Every catalog has a root directory and all files contained in the catalog must be somewhere (however many levels deep) under that root. Before you create your catalog, you'll want to place all of the files you want in that catalog in its root directory. You can set up the directory structure however you like.

## Create the Catalog
Creating a catalog is a simple process. Decide on a name (alphanumeric only) for the catalog that will serve as its accessible identifier. Then, decide on its root directory and description and then you can use `natster catalog new` to create the catalog:

```
natster catalog new sample "sample catalog" /home/kevin/medialibrary
New catalog created: sample
```

Natster has gone and ingested the contents of that root directory and prepopulated the catalog's JSON file. For example, the preceding catalog can be found at `~/.natster/sample.json` and contains the following contents (yours will vary):

```json
{
  "name": "sample",
  "root_dir": "/home/kevin/medialibrary",
  "description": "sample catalog",
  "last_modified": 1710251244,
  "entries": [
    {
      "path": "/README.txt",
      "description": "Auto-imported library entry",
      "mime_type": "text/plain; charset=utf-8",
      "hash": "df2e9d5d77459d97e0784928f7d30b894732a3e83e460f80254dd3a693d227b9",
      "byte_size": 40
    },
    {
      "path": "/bookvideos/Programming WebAssembly with Rust - Trailer Video.mp4",
      "description": "Auto-imported library entry",
      "mime_type": "video/mp4",
      "hash": "d1b1ca5d0bebd22266e9078977390dec28addf854ff12a90609a72f2e108cb67",
      "byte_size": 26922231
    },
    {
      "path": "/c/leveltwo.txt",
      "description": "Auto-imported library entry",
      "mime_type": "text/plain; charset=utf-8",
      "hash": "e6379468fa0990e3be979b6df4aaeea8b88c5d763e793c910fb87a0dc3efc4dd",
      "byte_size": 22
    }
  ]
}
```

Here you can see that Natster has automatically produced the hash of each of the files (this will be important later for downloads), stored the byte size, and even made a best-effort attempt at determining the mime type of the file. As you'll see later in this guide, files with a streamable mime type can viewed directly from the [natster.io](https://natster.io) application. The JSON file contains a flat list of files, even if those files have a hierarchical structure.

{{< hint type=important title="Security Context" >}}
It's important to remember that Natster is decentralized and leverages NATS' decentralized security mechanism. Further, the contents of your catalog are **not** made available to Synadia nor anyone else until you explicitly share them. The catalog contents are never cached and the moment you stop your catalog server, no metadata or files will be available.

{{< /hint >}}

With our shiny new catalog in hand, we can start sharing it with our friends.



