# 0001-SPEC: Run the console and extub as one microruntime

**Status:** ready
**PRD:** 0001-PRD-console-preview.prd.md

## Solution

The console and extub ship as one image with two processes: extub on `127.0.0.1:3401` and the console on `0.0.0.0:3001`. Both apps accept only loopback hosts, and `ex` runs one container per component without injected variables, so sharing one container keeps every guard as it is. A person opens the running console through a port-forward, which keeps the Host loopback.

## Decisions

- One component instead of two. Two would need configurable bind addresses in extub, a non-loopback `EXTUB_URL` in the console, and environment injection in `ex`.
- Readiness on `/icon.svg`, a static route outside the console's middleware, because the kubelet probes with the pod address as Host and every data route refuses it.
- The image build is pre-warmed before a start, because each agent operation is bounded to 120 seconds.

## Components

```langdrift-plan
{
	"components": [
		{
			"path": ".",
			"image": "langdrift-console:local",
			"port": 3001,
			"readinessPath": "/icon.svg",
			"build": {
				"executable": "platform/microruntime/console/build-image.sh",
				"args": []
			}
		}
	]
}
```

## Acceptance criteria

- [ ] The runtime reaches ready, and `/en/home` answers through the port-forward.
- [ ] CPU and memory are observed against the envelope, with their timestamps.
- [ ] Stopping removes every owned resource and keeps the evidence.
