---
title: Import a Catalog
weight: 60
---

In Natster, a catalog import is really not much more than a simple NATS account import. For every catalog you import, Natster will create an NGS/NATS subject import in your account.

{{<hint type=warning title="Account Quotas" >}}
Depending on the type of Synadia Cloud team you have and the quotas set on the account you're using for Natster, you might see catalog imports fail. Every catalog import consumes **2** (one for the service, one for media) NATS subject imports. If you don't have enough subject imports remaining in your account, the catalog import will fail.

Note that every Natster account has only one export for catalogs and this will not grow no matter how many catalogs you share.
{{</hint >}}

## Locating the Source
Just as you needed the account public key of the target when you shared a catalog, the recipients of that share need your account public key in order to import. You've already seen how to obtain this information by using the Natster `whoami` command. It might seem a little awkward to have to use these opaque keys, but remember Natster is a decentralized system and these keys ensure that Synadia doesn't need to be involved in facilitating shares.

With the public key of the sharer in your clipboard, you can run the `import` command, as shown below:

```
natster catalog import kevbuzz ACTZW5NQGNUQHWDFNPBPROVV4HJO76K7H3QRARK4DWA2P2KBJEHRCUT7
```

If all goes well and you haven't exceeded your account import quota, you should now have a fully functioning one-way, private exchange between two accounts.

## Viewing Catalog Shares
There's a very handy Natster command that will list all catalogs _known to you_. Essentially this is a combination of all of the catalogs you have shared with others and all of the catalogs others have shared with you. 

```
$ natster catalog ls
╭──────────────────────────────────────────────────────────────────────────────────────╮
│                                    Shared Catalogs                                   │
├────┬───────────────┬────────────────────────────────┬────────────────────────────────┤
│    │ Catalog       │ From                           │ To                             │
├────┼───────────────┼────────────────────────────────┼────────────────────────────────┤
│ 🟢 │ kevbuzz       │ me                             │ ACSH...PRQOLE3A2YOR7YALKXCKTPA │   
│ 🟢 │ jordans_stuff │ ACSH...PRQOLE3A2YOR7YALKXCKTPA │ me                             │
│ 🟢 │ synadiahub    │ AC5V...KAAANER6P6DHKBYGVGJHSNC │ me                             │
│ 🟢 │ kt            │ AATM...5NCGVLMDBBBBCE2XLBQKZEE │ me                             │
╰────┴───────────────┴─────────────────────────────────────────────────────────────────╯
```

The green dot in the preceding output indicates whether that catalog is _online_. An online catalog is one that has a running catalog server that has submitted a heartbeat to the Natster global service recently. If a catalog server explicitly shuts down or loses network connectivity, its status will switch to offline.