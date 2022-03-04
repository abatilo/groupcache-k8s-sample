allow_k8s_contexts("kind-groupcache")

docker_build(
  ref="sample",
  context="."
)

k8s_yaml("./k8s/sample.yaml")
k8s_resource(workload="sample", port_forwards=["8000:8000"])
