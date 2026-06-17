---
name: debugger
version: 1.0.0
description: Use when a runtime defect requires stepping through live execution state — variables, call stack, evaluated expressions — and reading code or log lines is not enough to locate the root cause. Also use when a test failure has no obvious cause from its output alone. Do NOT use for trivial issues resolvable by reading code or a single log line. In sdd-apply, diagnose AND fix in the same context. In sdd-verify, diagnose only — defer fixes to a subsequent sdd-apply. Requires the per-language debug toolchain to be present (run the preflight below before opening any session). Complements the `mcp__debugger__*` MCP: the MCP is the DAP engine; this skill is how to drive it.
license: MIT
metadata: {"hermes":{"tags":["debugger","debugging","breakpoint","dap","runtime","step-through","stack-trace"],"category":"development","related_skills":["drawio"]},"author":"matecito-ai","version":"1.0.0"}
---

# Debugger

**On-demand only — never automatic. The `mcp__debugger__*` MCP is the DAP engine; this skill is how to drive it well.**

This skill knows the MCP **exists**; it does not hardcode what the MCP **exposes**. The tool names, signatures, and per-tool descriptions are owned by the MCP server and may change over time — you learn them at runtime when its tools load. So this skill describes the **judgment, order, and gotchas** by role, and leaves the concrete tool selection to the live `mcp__debugger__*` surface (its own descriptions are authoritative). If those tools are not yet loaded when you need them, load them first, then proceed.

## When to reach for the debugger

**Use it when:**

- A runtime defect is present and reading the code or log output does not reveal where it happens.
- A test fails and the failure message does not point to a clear line.
- You need to inspect the live value of a variable or expression mid-execution (not what the code says it should be, but what it actually is at runtime).

**Do NOT use it when:**

- The issue is obvious from reading the code or a single log line — open it, read it, fix it.
- There is no runnable program in scope (e.g., a pure documentation or config change).
- The per-language debug toolchain is absent — run the preflight; if the binary is missing, install it first (see Install Helper below).

## Phase split: apply vs. verify

| Phase | Allowed actions |
|-------|----------------|
| `sdd-apply` | Diagnose the root cause AND apply a fix in the same context. Primary home for the debugger. |
| `sdd-verify` | Diagnose only — understand why a test or scenario fails. MUST NOT apply fixes here. Any fix found belongs in a subsequent `sdd-apply` invocation. |

## Preflight (run this before every debug session)

The preflight is **procedural** — it resolves the real toolchain for THIS project's language at use-time. It is NOT a lookup table to copy from.

### Step 1 — Detect the project's language

Read the project's manifest or source files to determine the primary language. Common signals: `go.mod` → Go, `package.json` → JavaScript/TypeScript, `requirements.txt` or `pyproject.toml` → Python, `Cargo.toml` → Rust.

### Step 2 — Ask the MCP which languages it supports

Use the MCP's language-listing capability. A language appearing there means the **adapter** is present — it does NOT mean the debug binary for that language is installed on the machine.

### Step 3 — Verify the real toolchain binary

**Critical distinction: adapter supported ≠ toolchain binary present.**

The MCP adapter is the bridge code that speaks DAP; the toolchain binary is the language-specific executable the adapter launches (e.g., `dlv` for Go, `debugpy` for Python, `node --inspect` for Node.js). A language can be reported as supported while its binary is completely absent from the machine. The MCP will error or hang when it tries to launch a missing binary.

Verify the binary is actually installed for the detected language. Example for Go:

```bash
which dlv          # or: dlv version
```

If the command succeeds and prints a path or version, the toolchain is present. If it fails (`not found`, exit 1), the toolchain is **absent** — do NOT proceed to a debug session; go to the Install Helper below.

Apply the same verification for whatever language your project uses: look up the language's standard debug binary name, run `which <binary>` (or `<binary> --version`), and treat a non-zero exit as absent.

### Step 4 — Proceed or install

- **Binary present** → continue to the Debug Loop.
- **Binary absent** → follow the Install Helper, then re-run Step 3 before opening a session.

## Install Helper

When the preflight finds the binary absent, do not fail silently or error out. Surface the missing toolchain to the user, give them the exact install command for the detected language, and offer to run it.

The install command is specific to that language's ecosystem. Do not ship an exhaustive table here — resolve it for the language at hand. Example for Go:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Then confirm with `which dlv` or `dlv version` before proceeding.

For other languages, apply the same pattern: identify the language's standard debug toolchain, find its install command from the language's official documentation or package manager, run it, and verify the binary is present before reopening the preflight.

## Debug Loop

Once the preflight passes, drive the session in this order. Each step is a **role** — map it to whichever `mcp__debugger__*` tool currently provides that capability (consult the MCP's own tool list and descriptions; do not assume names):

1. **Open a session** for the detected language. Keep the returned session identifier for every later call in the loop.
2. **Set a breakpoint** on an **executable line** — a statement, function call, assignment, or return.
   - **Gotcha — non-executable lines.** The MCP accepts a breakpoint on any line number, including closing braces (`}`), blank lines, and comments, without erroring — but execution then steps oddly: it may land on the wrong line, skip the breakpoint entirely, or stop at an unexpected position. Never target a structural delimiter or comment. When unsure, pick the first statement inside the block you want to inspect.
3. **Start execution.** It runs and pauses at the first breakpoint hit.
4. **Inspect the live state.** Read the local variables in the current frame, a broader scope, or the full call stack — to see what the values *actually are* at runtime, not what the code says they should be.
5. **Step through execution** as the situation requires: execute the current line and stay in the frame, descend into a call, step out to the caller, or resume to the next breakpoint (or program end).
6. **Evaluate an expression** in the current frame to confirm a hypothesis or inspect a computed value without adding another breakpoint.
7. **Close the session** when done — always, even after an error or an early exit. Do not leave dangling DAP sessions open.

When you are unsure which tool performs one of these roles, the live `mcp__debugger__*` tool list and its descriptions are the source of truth — this skill intentionally does not duplicate them, so it never goes stale when the MCP evolves.
