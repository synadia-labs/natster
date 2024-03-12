---
title: Web UI
weight: 70
---

Using Natster's web UI is just a matter of pointing your browser at [natster.io](https://natster.io) and logging in. Natster supports both Google (e.g. gmail) and Github authentication.

Before going to the website, however, you should bind your web identity to your Natster account/CLI.

## Binding your Web Context
The Natster web application not only needs to know who you are via Google or Github authentication, but it also needs to know which Natster account context to use whenever you log in.

This association between your web identity and your Natster account/CLI credentials is easily established by running the following command:

```
natster login
```

This will generate a one-time code (OTC) and open your browser to the website. Once you're on the website, you can log in with Google or Github and then, after it links your web identity to your Natster account, you will see a screen that looks similar to the following:

![natster UI](/assets/natster_screen.png)

The web application is easy to use, but if you need contextual help while using it, the application provides tooltips and additional help.
