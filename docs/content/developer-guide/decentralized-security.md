---
title: Decentralized Auth
weight: 20
---

Rather than build its own security and user account tier like we have to do with most applications, Natster utilizes NATS' [decentralized authn/authz](https://docs.nats.io/running-a-nats-service/configuration/securing_nats/auth_intro/jwt) directly. A Natster account is an NGS account (which is a NATS account). Natster authenticates to Synadia Cloud using JSON Web Tokens (JWT) and seeds.

Traditionally, when we get asked to build an application out of a suite of microservices, we might add services to deal with authentication and authorization. There might be a central database where all of the users are stored and maybe hashes of their passwords. Because nearly every user activity needs to be secured, it's remarkably easy for the security services to become a central point of failure or a performance bottleneck.

With NATS, we don't have to create or maintain any extra layers.

If a user's identity and permissions (_claims_) can be represented as a JWT, then our services don't need to consult a central database to provide security middleware. The impact of this fact might not be immediately realized, but when we start running our applications in multiple regions on multiple clouds with complicated network topologies, the confidence we gain from being able to verify user authorization and identity _in isolation_ without needing to contact a central authority is immeasurable.

## Natster.IO
The [natster.io](https://natster.io) website is hosted by a Go binary. This binary is given a set of credentials to log in as a user within the **Natster IO** Synadia Cloud account. This is used kind of like a bootstrap so that when an authenticated user logs in, the website can communicate with Synadia Cloud on their behalf, using their credentials.

## Catalog Server
The Natster catalog server authenticates to NGS using credentials maintained in the natster "context", which can be found in the `~/.natster` directory after a user has run `init`.

## Global Service
The global service is another Go binary, also running in our (Synadia) infrastructure. This service authenticates to NGS using its own set of unprivileged credentials. Consumers communicate with the global service using [secure imports](../secure-sharing).

## Attack and Loss Mitigation
A pretty popular attack vector for malicious actors is the message bus, whatever product or technology that may be. The idea behind that is if you can put a tap on message traffic, eventually you might be able to capture enough secrets to be able to escalate privilege and do whatever damage you like.

This is where both NATS and decentralized authentication truly shine. When a user authenticates with NATS via JWT and seed, _no keys are ever transmitted on the wire_. A NATS server _does not store any secrets_ on behalf of connected users. You could literally obtain a memory dump of the machine on which a NATS server resides and still not find a single secret.
