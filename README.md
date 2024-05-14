# Admission Controller 101
This is a hello world Kubernetes Golang based Validating Admission Controller 101 project.

This validating admission controller watches (only) the creation of Pods and disallows
to create Pods named `pod-with-an-invalid-name`. All other Pods names are allowed.

It is inspired by this project and maintains the same Apache License 2.0:
https://github.com/stackrox/admission-controller-webhook-demo

See more information below on the why section and the differences between the projects.

Target audience:
- Developers getting started with implementing Validating Admission Controllers in Kubernetes
- Developers keen to quickly debug and isolate issues related to Validating Admission Controllers
- Curious souls

## Getting started
Build and deploy the project on your cluster:
```shell
git clone https://github.com/ofirc/admission-controller.git
cd admission-controller
task build package gen-certs
# (Optional): task push-image
#             and replace `ofirc` with your container registry prefix.
task deploy
```

Verify it was deployed successfully:
```bash
$ kubectl get pod -n webhook-demo
NAME                              READY   STATUS    RESTARTS   AGE
webhook-server-68496b6dff-nv95v   1/1     Running   0          9s
$
```

And then try to create a Pod that is not allowed:
```bash
$ kubectl apply -f examples/pod-with-an-invalid-name.yaml
Error from server: error when creating "examples/pod-with-an-invalid-name.yaml": admission webhook "webhook-server.webhook-demo.svc" denied the request: Pod with an invalid name is not allowed.
$
```

And create a regular Pod with a different name, should be allowed:
```bash
$ kubectl apply -f examples/pod-with-a-valid-name.yaml
pod/pod-with-a-valid-name created
$
```

To uninstall:
```bash
task uninstall
```

and clean the pods that were created:
```bash
kubectl delete -f examples/pod-with-a-valid-name.yaml --force
```

For more information refer to the source code and the help:
```bash
task help
```

## Requirements
This project uses a Taskfile to simplify building and deploying the webhook server.

Prerequisites:
- bash
- Go 1.22.3+
- Docker
- Kubernetes cluster
- Task ([Installation instructions](https://taskfile.dev/installation/))

## Why this project?
I wanted to tinker around with admission controllers for debugging an issue related
to how the Kubernetes API Server and CA and TLS certificates behave, and then I realized
that I don't have a minimalistic Go project I can build, package, release, deploy
and then uninstall. So this project was born.

## How is it different than admission-controller-webhook-demo?
It is inspired by and builds upon the mutating webhook blueprint from:
https://github.com/stackrox/admission-controller-webhook-demo

But it differs from it in a few different ways:
1. Validating Admission Controller - instead of a mutating admission webhook
2. Idempotency - it allows one to deploy the manifest an infinite amount of times
                 without failing (i.e. no imperative kubectl create commands are involved).
                 Hint: `task deploy`
3. Uninstall - it allows to uninstall the deployment on the cluster.
               Hint: `task uninstall`
4. CI/CD friendly - it features a clear separation of a typical CI/CD pipeline:
                    build -> package -> release -> deploy
                    and can be used as a reference for integration.
5. Multiple container images - both `admission-controller` the distroless build,
                               as well as `admission-controller:dev`
                               which has a package manager and a shell and is useful
                               for hacking and debugging.

## What are good resources to learn about admission controller?
You may refer to the following resource:
https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/

And the official Kubernetes documentation:
https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/

Note, however, that there are many nuances and sometimes the best way to learn
is to tinker and hack your way around.

## Live debugging of the webhook server inside the Pod
First, execute into a shell:
```bash
$ kubectl exec -it $(kubectl get pod -n webhook-demo -lapp=webhook-server -oname) -n webhook-demo -- sh
```

Trigger the root endpoint in the webhook server:
```bash
~ $ curl -k https://localhost:8443/
Hello from root endpoint!
~ $
```

## FAQ
Q: Is this a project I can run in my production Kubernetes clusters?
A: This is definitely not advised, it's meant as an educational tool and not
   as a robust policy engine that enforces security or misconfiguration policices.


Q: The default container image has a shell and a package manager, why?

A: This project actually builds two container images:
   ofirc/admission-controller:dev
   ofirc/admission-controller:latest

   The dev variant contains a package manager (apk), a shell and curl.
   This is useful when one wants to simulate the webhook server endpoints
   when the webhook is deployed in-cluster.

   The latest variant is a distroless container image with no shell or a package manager
   and has a minimal footprint (and thus attack surface).
