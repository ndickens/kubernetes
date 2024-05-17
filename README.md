This is a small project to demonstrate a basic understanding of Kubernetes.

**Prompt:**

 - Create a hello world application and deploy it within Kubernetes.
 - The application should accept HTTP requests from clients and simply return “Hello, world!”.
 - The application should be reachable from a browser/client outside of the cluster.
 - The application should keep an audit log of IP addresses which sent a request to the application.
 - Be prepared to talk about how your application works, paying particular attention to the path the HTTP request takes from the browser to the application.

**Requirements:**

 - The application can be written by you, or you can use one from the internet.
 - The application can be written in any language. Using Golang is a stretch goal.
 - The Kubernetes cluster can be any distribution. For simplicity, kind is recommended.
 - The audit log should persist through container restarts. Strive to persist the audit log permanently.
 - There are no limitations on the technologies used for traffic ingress or storage.
 - Accepting HTTPS requests is not a requirement.

**[Github Project](https://github.com/ndickens/kubernetes)**

I created the application in go; see `main.go` in the app folder.  The app exposes two resources:

`/hello` returns the “Hello, world!” string and logs the ip address.  
`/logs` returns the contents of the log file for test verification purposes. 

As suggested I used kind for the distribution and created the kind-config.yml file in the 
kubernetes folder. I also created the kube-config.yml to define all the kubernetes components 
(deployment, service, persistent volume, etc.).

With everything created, I then ran through the below steps to test and verify the requirements.

**Note, these steps assume docker desktop, go, kubectl, and kind are installed**

1. Build main.go.
```
% go build
% go mod tidy
  ```

2. Build the image.
```
% docker build -t hello-app $KUBE_HOME/app
```

3. Create the cluster and load the image into the worker. Add the nginx controller and wait for it to be ready.
```
% kind create cluster --name=hello-cluster --config=$KUBE_HOME/kubernetes/kind-config.yml
% kind load docker-image hello-app --name hello-cluster --nodes hello-cluster-worker
% kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
% kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

4. Apply the kubernetes config file.
```
% kubectl apply -f $KUBE_HOME/kubernetes/kube-config.yml
```

5. Verify accessibility and cause a log event.
```
% curl http://localhost/hello 
  OR 
view at http://localhost/hello
```
6. Verify the log was created in the logs directory.
```
% curl http://localhost/logs 
  OR 
view at http://localhost/logs
```

7. Restart the pod by scaling instances down to 0 and back up.
```
% kubectl scale deployment hello --replicas=0
% kubectl get pods // ensure the pod is gone
% kubectl scale deployment hello --replicas=1
```

8. Add another log event.
```
% curl http://localhost/hello 
  OR 
view at http://localhost/hello
```

9. Verify the log has two entries.
```
% curl http://localhost/logs 
  OR 
view at http://localhost/logs
```

10. Clean up.
```
% kind delete cluster --name=hello-cluster
% docker image rm hello-app
```

References
- [Docker Go Image Building](https://docs.docker.com/language/golang/build-images/)
- [Kind Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Persistent Volumes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/)
- [Restarting Pods](https://spacelift.io/blog/restart-kubernetes-pods-with-kubectl)
- [Kind Ingress](https://kind.sigs.k8s.io/docs/user/ingress/)