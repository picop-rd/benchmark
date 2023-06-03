## Common Steps
Please prepare Kubernetes Cluster deployed Proxy Controller and Host Controller DB.
You can see [usage](https://github.com/picop-rd/demo#kubernetes) in demo.
You do not need to deploy proxies and services because they are deployed in several benchmark steps.

```
# In the demo directory
cd kubernetes/manifests
vim proxy-controller.yaml # Please replace the IP address 192.168.0.4 with that of the Host Controller DB.
kubectl apply namespace.yaml
kubectl apply proxy-controller.yaml
```
