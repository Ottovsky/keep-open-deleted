### PoC for https://github.com/kubernetes/kubernetes/issues/83107
Simple program to check if the issue in question can cause a k8s worker to fail.

Compile with 
```bash
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o keepopendeleted  main.go
```

Docker image available at:
```bash
docker pull ottovsky/keepopendeleted
```


Run in k8s with:
```bash
kubectl create deployment keepopendeleted --image=ottovsky/keepopendeleted
```

**Run only in the sandbox environment!**