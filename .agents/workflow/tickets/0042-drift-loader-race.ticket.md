# 0042: Stop the drift loader's delayed writer with the loader

**Spec:** None

**What to build:** a drift command that finishes within a second leaves no error behind. Today the loader's one-second delay runs as an untracked background job; when the command ends first, the loader deletes its temp directory and the job then fails writing `ready` into it, printing `output.sh: line 93: …/ready: No such file or directory` after the prompt returns.

**Blocked by:** None (can start immediately).

**Status:** ready-for-agent

- [ ] The delayed `ready` writer's process is tracked, and stopped and reaped before the loader restarts or deletes its directory.
- [ ] A command slower than a second still shows the spinner after the one-second delay.
- [ ] No `drift-loader.*` temp directory is left behind.

### Verification

- `DRIFT_FORCE_LOADER=1 sh -c 'drift env validate; sleep 1.6'` prints no `No such file or directory`, three runs in a row.
- Sourcing `output.sh` and holding the loader for 1.5 s draws spinner frames; restarting it within the second and stopping it leaves no error.
