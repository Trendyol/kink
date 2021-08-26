# Webhook server executing configured commands

This is based on https://github.com/adnanh/webhook

The container image contains `kubectl` in it.
As a result, one can easily create a webhook that can trigger updates to the Kubernetes API server.
This can be very useful when combining with utilities like [configmap-reload][configmap-reload] that watches for Kubernetes API object changes.

## Simple usage

```
docker run --rm webhook [WEBHOOK ARGS]
```

See more details about [webhook][webhook] in its documentation.

[configmap-reload]: https://github.com/jimmidyson/configmap-reload
[webhook]: https://github.com/adnanh/webhook
