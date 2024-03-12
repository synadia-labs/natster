---
title: Initializing Natster
weight: 10
geekdocToC: 3
---

In the **Getting Started** section, you went through the process of signing up for a Synadia Cloud account if you didn't already have one. The next thing we need to do is configure the Natster CLI so that it can communicate with Synadia Cloud on our behalf. 

<!--more-->

## Creating a Personal Access Token
A personal access token allows applications to utilize the Synadia Cloud API on your behalf. This is exactly what Natster will do for nearly all of its core functions. To do this, we use a _personal access token_ or **PAT** for short. If you've created SSH keys in GitHub or performed OAuth authorizations then this concept should be pretty familiar.

First, log into your [Synadia Cloud](https://cloud.synadia.com) account. Once you're on the home page, click the person icon in the top right and choose "profile". Once on your profile page, the left navigation bar should contain two categories: _General_ and _Access Tokens_.

Click on _Access Tokens_.

Once on the Personal Access Tokens page, click the **Create Token** button. You'll be prompted for a name and expiration date to be used for the token. It's a good idea to include a reference to Natster somewhere in the token name so you can keep activity separate from other tokens. The expiration date is up to you, but keep in mind that you'll have to re-initialize Natster if your access token expires.

As soon as you create the token, it will be displayed to you. You will **never** see this token again, so you'll want to copy it and store it in a password manager or in some other safe place. You will need this token for the next step.

Here is a sample showing an active access token:

![access token](/userguide/personal_access_token.png)

## Creating a Natster User
The Natster CLI, both for ad hoc commands and for running a catalog server, needs to connect to NGS (Synadia Cloud's NATS infrastructure) with a specific user credential. This user will belong to your account and have whatever privileges you decide. Prior to initializing the CLI, you'll want to create this new user. To do so, use the following steps:

1. Click the NGS account that you want to use for your Natster activity. You can create a separate account just for use with Natster if you like.
2. On the top tab section, click the **Users** tab
3. Click the **Create User** button and fill in the appropriate fields. 

![new user](/userguide/create_natster_user.png)

With the new user created, you're able to initialize the Natster CLI. Don't worry about copying the credentials, the Natster initialization process will take care of that for you.

## Initializing the CLI
With your Synadia Cloud token read on the clipboard and a NATS user ready, you can run the `init` command of the Natster CLI

```
natster init --token <SYNADIA CLOUD TOKEN>
```
You will be prompted for a number of different things when you run `init`. It will ask you to identify the team, account, and user to establish a new Natster context. Here's an example of the initialization command right before finishing:

```
? Select a team: The Team
? Select a system: NGS
? Select an account: The Account
✅ Catalog service export is configured
✅ Media stream export is configured
✅ Natster global service import is configured
✅ Natster global events import is configured
? Select a user for NATS authentication:  [Use arrows to move, type to filter]
> kevin
```

## Verifying Your Identity
After you've initialized and you've seen no error messages and it all looks like you're ready to go, you can verify the identity that the Natster CLI has for you. To check your identity, you can run the following command (we've obfuscated the identifiers but your output will have the full values):

```
$ natster whoami
╭─────────────────────────────────────────────────────────────────────────────────╮
│                              Accounty McAccountFace                             │
├──────────────────────┬──────────────────────────────────────────────────────────┤
│ Account              │ ACTZW5NQGNUQHWDFNPBPROEC4HJO76K7H3QRARK4DWA2P2KBJEHRCUT7 │
│ Initialized At       │ 2024-03-08 12:52:52                                      │
│ Synadia Cloud Team   │ 2bN3..................1MVlQ                              │
│ Synadia Cloud System │ 2bN3..................1MVN1                              │
│ Synadia Cloud User   │ 2bN3..................1MV5N                              │
│ Credentials          │ /home/kevin/.natster/accounty.creds                      │
│ Natster.io Login     │ (unlinked)                                               │
╰──────────────────────┴──────────────────────────────────────────────────────────╯

```

Congratulations! Now that you're all set up and running, it's time to start using Natster.