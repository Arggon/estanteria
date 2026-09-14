#!/usr/bin/env python3
"""Persistent MCP client for the arggon stdio server.

The arggon MCP server (`arggon mcp`) speaks JSON-RPC over stdio with
newline-delimited messages. ZCode shell calls don't share processes, so the
server session lives in a small daemon that keeps the subprocess open and
forwards requests over a Unix socket. Every tracker operation in this repo's
implementation phase goes through this one session (MCP-first experiment).

Usage:
  mcp_client.py start [--server CMD]      spawn daemon + arggon mcp + handshake
  mcp_client.py call <tool> '<json-args>' tools/call; prints the response JSON
  mcp_client.py status                    health check (tools/list round trip)
  mcp_client.py stop                      shutdown + reap the server

State: .mcp-session/ (socket, pid, log). Response log: .mcp-session/responses.log
(one JSON line per tools/call) kept as evidence of the session's health.
"""

import json
import os
import signal
import socket
import subprocess
import sys
import threading
import time

HERE = os.path.dirname(os.path.abspath(__file__))
STATE_DIR = os.path.join(HERE, "..", ".mcp-session")
SOCK_PATH = os.path.join(STATE_DIR, "mcp.sock")
PID_PATH = os.path.join(STATE_DIR, "daemon.pid")
LOG_PATH = os.path.join(STATE_DIR, "responses.log")


def recv_line(stream):
    line = stream.readline()
    if not line:
        raise ConnectionError("server closed the stream")
    return json.loads(line)


class Bridge:
    """One client connection: read request lines, forward to the server."""

    def __init__(self, conn, proc, lock, log):
        self.conn = conn
        self.proc = proc
        self.lock = lock
        self.log = log

    def run(self):
        try:
            with self.conn:
                buf = self.conn.makefile("r")
                for raw in buf:
                    raw = raw.strip()
                    if not raw:
                        continue
                    resp = self.forward(json.loads(raw))
                    self.conn.sendall((json.dumps(resp) + "\n").encode())
        except (ConnectionError, json.JSONDecodeError) as exc:
            print(f"bridge closed: {exc}", file=sys.stderr)

    def forward(self, req):
        with self.lock:
            self.proc.stdin.write((json.dumps(req) + "\n").encode())
            self.proc.stdin.flush()
            while True:
                msg = recv_line(self.proc.stdout)
                if msg.get("id") == req.get("id") or ("id" not in msg and "method" in msg):
                    if "method" in msg and "id" not in msg:
                        continue  # server-pushed notification; keep waiting
                    break
            if req.get("method") == "tools/call":
                self.log.write(json.dumps(msg, ensure_ascii=False) + "\n")
                self.log.flush()
            return msg


def start(server_cmd):
    os.makedirs(STATE_DIR, exist_ok=True)
    if os.path.exists(SOCK_PATH):
        try:
            with socket.socket(socket.AF_UNIX) as s:
                s.settimeout(1)
                s.connect(SOCK_PATH)
            print("already running")
            return 0
        except OSError:
            os.unlink(SOCK_PATH)  # stale socket

    proc = subprocess.Popen(
        server_cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    req_id = 0

    def rpc(method, params=None):
        nonlocal req_id
        req_id += 1
        msg = {"jsonrpc": "2.0", "id": req_id, "method": method}
        if params is not None:
            msg["params"] = params
        proc.stdin.write((json.dumps(msg) + "\n").encode())
        proc.stdin.flush()
        while True:
            resp = recv_line(proc.stdout)
            if resp.get("id") == req_id:
                return resp
            # ignore notifications while handshaking

    init = rpc("initialize", {
        "protocolVersion": "2024-11-05",
        "capabilities": {},
        "clientInfo": {"name": "estanteria-mcp-client", "version": "0.1.0"},
    })
    server_name = init["result"]["serverInfo"]["name"]
    proc.stdin.write(b'{"jsonrpc":"2.0","method":"notifications/initialized"}\n')
    proc.stdin.flush()

    tools = rpc("tools/list")
    names = [t["name"] for t in tools["result"]["tools"]]
    print(f"connected: {server_name} | tools: {', '.join(names)}")

    log = open(LOG_PATH, "a")
    lock = threading.Lock()
    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.bind(SOCK_PATH)
    server.listen(4)
    with open(PID_PATH, "w") as f:
        f.write(str(os.getpid()))

    def reap(_sig, _frame):
        proc.terminate()
        sys.exit(0)

    signal.signal(signal.SIGTERM, reap)
    signal.signal(signal.SIGINT, reap)
    try:
        while True:
            conn, _ = server.accept()
            Bridge(conn, proc, lock, log).run()
    except KeyboardInterrupt:
        proc.terminate()


def daemon_pid():
    try:
        return int(open(PID_PATH).read())
    except (OSError, ValueError):
        return None


def call(tool, args_json):
    req = {"jsonrpc": "2.0", "id": 1, "method": "tools/call",
           "params": {"name": tool, "arguments": json.loads(args_json or "{}")}}
    resp = _roundtrip(req)
    print(json.dumps(resp, ensure_ascii=False, indent=1))
    return 0 if "error" not in resp and not resp.get("result", {}).get("isError") else 1


def status():
    t0 = time.time()
    resp = _roundtrip({"jsonrpc": "2.0", "id": 1, "method": "tools/list"})
    ms = (time.time() - t0) * 1000
    names = [t["name"] for t in resp["result"]["tools"]]
    print(json.dumps({"ok": "error" not in resp, "latency_ms": round(ms, 1),
                      "tools": names, "daemon_pid": daemon_pid()}))
    return 0 if "error" not in resp else 1


def _roundtrip(req):
    if not os.path.exists(SOCK_PATH):
        sys.exit("no live session — run: tools/mcp_client.py start")
    with socket.socket(socket.AF_UNIX) as s:
        s.settimeout(30)
        s.connect(SOCK_PATH)
        s.sendall((json.dumps(req) + "\n").encode())
        buf = b""
        while not buf.endswith(b"\n"):
            chunk = s.recv(65536)
            if not chunk:
                break
            buf += chunk
    return json.loads(buf)


def stop():
    pid = daemon_pid()
    if pid:
        os.kill(pid, signal.SIGTERM)
        print("stopped")
    else:
        print("not running")
    return 0


def main():
    cmd = sys.argv[1] if len(sys.argv) > 1 else "status"
    if cmd == "start":
        server = sys.argv[sys.argv.index("--server") + 1] if "--server" in sys.argv else "arggon mcp"
        sys.exit(start(server.split()))
    elif cmd == "call":
        sys.exit(call(sys.argv[2], sys.argv[3] if len(sys.argv) > 3 else "{}"))
    elif cmd == "stop":
        sys.exit(stop())
    else:
        sys.exit(status())


if __name__ == "__main__":
    main()
