---
title: Downloading Files
weight: 65
---

Downloading files can be done directly from the Natster CLI. First you need to query a given catalog's contents via the `natster catalog contents`
command. As you've seen already this will contain the unique hash of each file in that catalog. To download a file,
you need to know the catalog in which it resides, and the hash.

The following command will download a video from the **synadiahub** catalog:

```
$ natster catalog download synadiahub \ 
   dbf499dc63b7f990762d578208b8ac5ee9d74d193573219e2ad2f3077841e769 \ 
   ./s1ep1.mp3
File download request acknowledged: 33656581 bytes (8217 chunks of 4096 bytes each.) 
  from XAQM376L7EBVX6DEKS5BAFLSUJEGE5GUBOEDUFRVFU4GSM6GDGHF72GA
Received chunk 0 (4096 bytes)
Received chunk 1 (4096 bytes)
Received chunk 2 (4096 bytes)
Received chunk 3 (4096 bytes)
Received chunk 4 (4096 bytes)
...
Received chunk 8216 (3845 bytes)
```

Your output may differ as the CLI experience improves. You should now have a file called **s1ep1.mp3** which is Season 1, Episode 1, of the [NATS.fm](https://nats.fm) podcast.

```
$ ls s1ep1.mp3
-rw-rw-r-- 1 kevin kevin 434176 Mar 12 11:11 s1ep1.mp3
```

{{<hint type=tip title="Encrypted Data Transfer" >}}
Every time you request the download or streaming of a file from a remote catalog server, a brand new set of encryption keys will be generated. These keys (xkeys) 
are directional, and so only the downloader's private key can be used to decrypt the contents transmitted by the sender. Not only does this ensure that all
of your shared data remains private, but the same download key can't be used twice, providing even more security.
{{</ hint>}}