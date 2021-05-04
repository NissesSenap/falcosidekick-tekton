# Ideas

## TODO

- Create a super simple go app that works with the incoming data.
  - probably just a single body for now, then discuss with falco community.
- Create a Dockerfile
- Push dockerfile to quay.io
- Verify e2e workflow

## Output example

This is how I can easily test my go application.
This output is from a command that is not okay in falco.

export BODY='{"output":"17:30:51.428027354: Notice A shell was spawned in a container with an attached terminal (user=root user_loginuid=-1 k8s.ns=falcoresponse k8s.pod=alpine container=ac7c5359a04b shell=sh parent=runc cmdline=sh -c cat /etc/resolv.conf terminal=34816 container_id=ac7c5359a04b image=alpine) k8s.ns=falcoresponse k8s.pod=alpine container=ac7c5359a04b k8s.ns=falcoresponse k8s.pod=alpine container=ac7c5359a04b","priority":"Notice","rule":"Terminal shell in container","time":"2021-05-02T17:30:51.428027354Z", "output_fields": {"container.id":"ac7c5359a04b","container.image.repository":"alpine","evt.time":1619976651428027354,"k8s.ns.name":"falcoresponse","k8s.pod.name":"alpine","proc.cmdline":"sh -c cat /etc/resolv.conf","proc.name":"sh","proc.pname":"runc","proc.tty":34816,"user.loginuid":-1,"user.name":"root"}}'
