---
title: Nothing but NATS
weight: 80
---

_"Nothing but NATS"_ is a philosophy of building applications where we rely on primitives available in NATS for everything we need, eschewing traditional costly, difficult, and complicated cloud products. Some of us treat this philosophy more like a way of life (the 道 Dao of NATS?)

## Security
Probably the largest return on investment from this philosophy for Natster is in the area of security. By using NATS decentralized authentication and authorization as a core part of the Natster sharing mechanics, we got a tremendous amount of features and functionality for free. 

Secure [imports and exports](../secure-sharing) took a feature that could have taken us _months_ to build and gave it to us at no cost with low complexity.

The details of this are covered in more detail in the [decentralized security](../decentralized-security) section.

## One Time Codes
Without having to use a single external library to manage one-time codes, we simply created a key-value bucket and set a short time-to-live on the values. This gave us expiring codes for free, as well as a means of attaching meaningful context to each code to ensure that codes can't be claimed by man-in-the-middle or snooping attackers.

## Streams and Consumers
The application needs to store some information in streams and dynamically consume that information. The streams need to be persistent and, in some cases, the status of each pull consumer also needs to be persistent.

We got all of this for free with NATS JetStream. No third party products were harmed (or used) in the production of this application.

## Event Log
The Natster [event log](../global-event-log) is an append-only JetStream stream. There are hundreds of libraries and products and servers and expensive things out in the market that supply the functionality we get automatically.

## Projections
In an event-sourced world, [projections](https://event-driven.io/en/projections_and_read_models_in_event_driven_architecture/) are the data that is read by your code and presented to your users. Without the use of any other products or libraries, we were able to combine JetStream streams and consumers with key-value buckets to create a robust and predictable projection system.



