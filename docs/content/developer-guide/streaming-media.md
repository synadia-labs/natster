---
title: Streaming Media
weight: 100
---

When we originally designed Natster, streaming media was a stretch goal. We planned on supporting single downloads via file chunking and that was all. But because Natster leveraged so much underlying NATS technology, we ended up _ahead_ of schedule, and so had time to explore streaming media.

The [natster.io](https://natster.io) website supports streaming `video/mp4` content from catalogs to which you have access. This means that you can log into Natster on your phone and watch the latest videos and tutorials from Synadia.

We could have used a third party media server or encoder or a media transcoding proxy, but we didn't need any of those. We already had in place the ability to securely download _encrypted_ chunks of data from the catalog source.

When you stream media on **natster.io**, the site's server creates a subscription to the media download topic hosted by the source catalog. As it receives chunks on this subject, it decrypts them and then encodes them in a format that is compatible with JavaScript's in-browser media element.

We were able to put together an experience that looks and feels like a high-end video streaming product using pieces of our architecture that we'd already developed. The most difficult part of implementing streaming media was the buffering and encoding work that had to be done in the browser JavaScript. The back-end portion of this feature was relatively simple.