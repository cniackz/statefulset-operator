### Description:

An operator that delivers a statefulset

<img width="433" alt="Screenshot 2024-08-27 at 3 44 27 PM" src="https://github.com/user-attachments/assets/e819c97f-1511-4de1-a89e-e468d7b4b7f9">

### Acknowledgment

https://github.com/cniackz/basic-k8s-operator

### Instructions:

### Instructions:

1. Create an empty cluster

```shell
createcluster
```

2. Clone the repository

```shell
rm -rf ~/statefulset-operator/
cd; git clone git@github.com:cniackz/statefulset-operator.git
cd ~/statefulset-operator
make docker-build IMG=radical-123
kind load docker-image docker.io/library/radical-123
```
   
3. Deploy the Operator

```shell
kubectl apply -f ~/statefulset-operator/config/manager
kubectl apply -k ~/statefulset-operator/config/rbac
kubectl apply -k ~/statefulset-operator/config/crd
kubectl apply -f ~/statefulset-operator/cr.yaml
```

### Development:

After each change run:

```shell
make generate
make manifests
make docker-build IMG=radical-123
kind load docker-image docker.io/library/radical-123
```

<img width="1840" alt="Screenshot 2024-08-30 at 11 07 55 AM" src="https://github.com/user-attachments/assets/b86b1d5c-8c5d-4039-8c76-cb84badda2a6">

<img width="1840" alt="Screenshot 2024-08-30 at 11 09 41 AM" src="https://github.com/user-attachments/assets/f2029302-e2fd-48d4-b66b-1a75874441b1">


<img width="1840" alt="Screenshot 2024-08-30 at 11 09 06 AM" src="https://github.com/user-attachments/assets/3898342a-d76e-4a5e-a504-8f9f5d40a229">


