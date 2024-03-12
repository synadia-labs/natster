---
title: Sharing a Catalog
weight: 30
---

Catalog sharing in Natster is a secure, one-direction agreement between two accounts. This share is directly between you, the sharer, and the target or recipient of your share. Synadia doesn't get any inside information about this share. 

## Getting a Target Account Key
Part of the decentralized appeal of this system (and, of course, NATS itself) is that you can do things without requiring central administration or facilitation. This is true of the Natster share system, which builds on NATS' account exports. To share your catalog with someone else, you will need their account public key.

This information is easily obtained. Have the intended recipient of your catalog share run the `natster whoami` command. The first line of data output from that command is a 56-character, all capital letter ID beginning with the letter **A**. If you've used NATS decentralized security before, you'll recognize this as an account's public key. Natster doesn't hide this from you and instead leverages NATS accounts for application security.

## Sharing a Catalog
Once you have the recipient's public key, you can now share the catalog with a command like so:

```
natster catalog share sample ACTZW5NQGNUQHWDFNPVPROEC4HJO76K7H3QRARK4DWA2P2KBJEHRCUT7
```

This shares the `sample` catalog with the account `ACTZW5NQGNUQHWDFNPVPROEC4HJO76K7H3QRARK4DWA2P2KBJEHRCUT7`. Note that there is no secure or sensitive information exposed in this action. 

{{< hint type=caution title="Don't Use this Account" >}}
The account **ACTZW5NQGNUQ...UT7** is a fake account key and sharing to this account will fail.
{{< /hint >}}

One the other side of this share, your friend should [import](../import-catalog) this catalog. As you'll see, the `inbox` command can be used to list off shared catalogs that have not yet been imported.

## Revoking Catalog Access
If at any time you no longer want the recipient to have access to your catalog, you just need to use the `catalog unshare` command with the same arguments you originally passed to the `share` command. There's no guarantee that the recipient's quota will be freed for that catalog, but they will no longer be able to access the catalog.


Once your friends have imported your catalog, you're ready to start your first catalog server.