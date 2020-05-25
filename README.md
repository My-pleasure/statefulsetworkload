# StatefulSetWorkload
Run OAM StatefulSetWorkload on a Kubernetes cluster, like ContainerizedWorkload.

# Progress
Implemented a StatefulSetWorkload Definition, including some fields: Name, Image, Env, Config, Port.

# How to use it
- Install OAM Application Controller and OAM Core workload and trait controller. (You can also follow [addon-oam-kubernetes-local](https://github.com/crossplane/addon-oam-kubernetes-local))
```
kubectl create namespace crossplane-system

helm repo add crossplane-alpha https://charts.crossplane.io/alpha

helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane

git clone git@github.com:crossplane/addon-oam-kubernetes-local.git

kubectl create namespace oam-system

helm install controller -n oam-system ./charts/oam-core-resources/ 
```
- Run statefulsetworkload controller.
```
git clone https://github.com/My-pleasure/statefulsetworkload.git

cd $GOPATH/src/statefulsetworkload

make install

make run
```
- Apply the sample application config
```
kubectl apply -f config/samples/statefulsetworkload-test.yaml
```
- Verify it you should see a statefulset looking like below
```
kubectl get statefulset
NAME                                     READY   AGE
example-appconfig-statefulset-workload   1/1     13s
```
