# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	industrial-key-rotation/cmd/keyrotator	[no test files]
?   	industrial-key-rotation/internal/audit	[no test files]
ok  	industrial-key-rotation/internal/controller	0.028s
ok  	industrial-key-rotation/internal/crypto	0.001s
?   	industrial-key-rotation/internal/model	[no test files]
ok  	industrial-key-rotation/internal/persistence	0.025s
?   	industrial-key-rotation/internal/policy	[no test files]
ok  	industrial-key-rotation/internal/report	0.002s
--- FAIL: TestKeyRotationRejectsBadEnvelope (0.01s)
    workflow_test.go:33: expected malformed sensor digest to reject rotation
    workflow_test.go:40: previous active key was not preserved: [{ID:sensor-a Name:Line Sensor ControllerID:controller-a Algorithm:AES-GCM KeyLength:24 Status:active CreatedAt:1700000000 UpdatedAt:1700000000 ActiveVersion:1}]
FAIL
FAIL	industrial-key-rotation/internal/rotation	0.024s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/keyrotator): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/keyrotator): exit `0`
