# StatefulSetWorkload
Run OAM StatefulSetWorkload on a Kubernetes cluster， like ContainerizedWorkload.

# Progress
Implemented a StatefulSetWorkload Definition, including some fields: Name, Image, Env, Config, Port.

# How to use it
- You should follow this firstly: [addon-oam-kubernetes-local](https://github.com/crossplane/addon-oam-kubernetes-local/tree/659c1331b0c734f4567c3fbc042ef0dea37f0624). To install cert manager, OAM Application Controller and OAM Core workload and trait controllers.
- Git clone this project and put it on $GOPATH/src
- Install and run it
```
cd $GOPATH/src/statefulsetworkload
make install
make run
```
- Apply the sample application config
```
kubectl apply -f config/samples/statefulsetworkload-test.yaml
```
