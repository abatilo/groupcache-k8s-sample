# groupcache-k8s-sample

[mailgun/groupcache](https://github.com/mailgun/groupcache) is a loading cache
that was inspired by the
[golang/groupcache](https://github.com/golang/groupcache) but with some
additional capabilities like `context.Context` support.

`groupcache` has the ability to register a list of peers and will hash request
keys and fetch the value from a peer. This way, you have an in-memory cache
that should be consistent and spread across all of your replicas.

The official repository doesn't have any examples or documentation for doing
peer discovery in a Kubernetes environment, so I did a quick sample showing a
way that you can do that with a headless `kind: Service`.

Run `./setup.sh` which will install required tools with
[asdf](https://asdf-vm.com/), then will create a local kubernetes cluster for
you, and build and deploy the local sample application.

In a different terminal window, you can run `./test.sh` to send a bunch of
requests and the response bodies will include cache hit statistics and which
peer the loading cache sourced the value from.

P.S. Hi Tim's co-workers! :wave: Hi Ryan. I'm a Casey Muratori fan too :D

## Video demo

Click
[here](https://raw.githubusercontent.com/abatilo/groupcache-k8s-sample/master/demo.svg)
to watch a demo svg.
