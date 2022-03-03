allow_k8s_contexts("kind-groupcache")

docker_build(
  ref="sample",
  context="."
)

k8s_yaml("./k8s/sample.yaml")
