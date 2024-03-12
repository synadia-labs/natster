---
title: Installation
weight: 0
---

Natster installation is quite simple. The only thing you need after you've set up your Synadia Cloud account is the Natster CLI.

<!-- more -->

You can download the CLI using the following shell command:

```
curl -sSf https://natster.io/install.sh | sh
```

After following the prompts and installing the CLI, you can verify your installation with the following command:


```
natster --version

v1.0.0 [none] | BuiltOn: 12 Mar 24 08:10 EDT
```

Your version may vary from what's shown above, but if you can run the CLI then you should be able to move on to the next step.

{{< hint type=note >}}
All of the metadata about your catalogs and your Natster context will be stored in the `.natster` directory in the root of your home directory, e.g. `~/.natster` on Mac or Linux. You may find yourself editing files in this directory.
{{</hint>}}

