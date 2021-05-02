# Ideas

## Current

I currently get issues to even start a pipline.
Try to make the trigger-binding as basic as possible. Aka send the entire body and no other incoming params.
Lets echo that output in the hello task/ubuntu image.
In short get a basic eventbinding to work agian.

### TODO

- Create a super simple go app that works with the incoming data.
  - probably just a single body for now, then discuss with falco community.
- Create a Dockerfile
- Push dockerfile to quay.io
- Verify e2e workflow

## Output example

This is how I can easily test my go application.
This output is from a command that is not okay in falco.

export BODY='{"output":"14:49:49.264147779: Notice A shell was spawned in a container with an attached terminal (user=root user_loginuid=-1 k8s.ns=default k8s.pod=alpine container=a15057582acc shell=sh parent=runc cmdline=sh -c uptime terminal=34816 container_id=a15057582acc image=alpine) k8s.ns=default k8s.pod=alpine container=a15057582acc k8s.ns=default k8s.pod=alpine container=a15057582acc","priority":"Notice","rule":"Terminal shell in container","time":"2021-05-01T14:49:49.264147779Z", "output_fields": {"container.id":"a15057582acc","container.image.repository":"alpine","evt.time":1619880589264147779,"k8s.ns.name":"default","k8s.pod.name":"alpine","proc.cmdline":"sh -c uptime","proc.name":"sh","proc.pname":"runc","proc.tty":34816,"user.loginuid":-1,"user.name":"root"}}'

## Potential issue

This is a log from falco. Try to find which rule generates this log.

{"output":"15:29:29.219916520: Notice Unexpected connection to K8s API Server from container (command=eventlistenersi --el-name=falco-listener --el-namespace=falcoresponse --port=8000 --readtimeout=5 --writetimeout=40 --idletimeout=120 --timeouthandler=30 --is-multi-ns=false --tls-cert= --tls-key= k8s.ns=<NA> k8s.pod=<NA> container=9c48dd873ab9 image=<NA>:<NA> connection=172.17.0.6:49936->10.96.0.1:443) k8s.ns=<NA> k8s.pod=<NA> container=9c48dd873ab9 k8s.ns=<NA> k8s.pod=<NA> container=9c48dd873ab9","priority":"Notice","rule":"Contact K8S API Server From Container","time":"2021-05-01T15:29:29.219916520Z", "output_fields": {"container.id":"9c48dd873ab9","container.image.repository":null,"container.image.tag":null,"evt.time":1619882969219916520,"fd.name":"172.17.0.6:49936->10.96.0.1:443","k8s.ns.name":null,"k8s.pod.name":null,"proc.cmdline":"eventlistenersi --el-name=falco-listener --el-namespace=falcoresponse --port=8000 --readtimeout=5 --writetimeout=40 --idletimeout=120 --timeouthandler=30 --is-multi-ns=false --tls-cert= --tls-key="}}
