#!/usr/bin/env node
/** Smoke the publisher helper: health → publish 2 editions → list → buy. */
import { createHash, randomUUID } from "node:crypto";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { spawn } from "node:child_process";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "../..");
const HELPER = process.env.HELPER_URL || "http://127.0.0.1:8787";
const PORT = Number(process.env.HELPER_PORT || 8787);

function sleep(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

async function main() {
  const child = spawn("node", [path.join(__dirname, "server.mjs")], {
    env: { ...process.env, HELPER_PORT: String(PORT) },
    cwd: ROOT,
    stdio: ["ignore", "pipe", "pipe"],
  });
  let boot = "";
  child.stdout.on("data", (d) => {
    boot += d;
  });
  child.stderr.on("data", (d) => {
    boot += d;
  });

  try {
    for (let i = 0; i < 40; i++) {
      try {
        const r = await fetch(`${HELPER}/api/health`);
        if (r.ok) break;
      } catch {
        /* retry */
      }
      await sleep(250);
    }
    const health = await (await fetch(`${HELPER}/api/health`)).json();
    if (!health.ok) throw new Error(`helper unhealthy: ${JSON.stringify(health)}`);
    if (!health.lcdOk) throw new Error("LCD not reachable — start bookstored first");
    if (!health.keys?.publisher) throw new Error("alice key missing — make node-init");

    const work = path.join(os.tmpdir(), `bs-smoke-${randomUUID()}`);
    await fs.mkdir(work);
    const book = path.join(work, "novella.txt");
    await fs.writeFile(book, "Once upon a publisher-friendly chain\n");

    const fd = new FormData();
    const blob = new Blob([await fs.readFile(book)], { type: "text/plain" });
    fd.set("file", blob, "novella.txt");
    fd.set("title", "Novella");
    fd.set("editions", "2");
    fd.set("list", "1");
    fd.set("price", "1000000");
    fd.set("from", "alice");

    const pubRes = await fetch(`${HELPER}/api/publish`, { method: "POST", body: fd });
    const pub = await pubRes.json();
    if (!pubRes.ok || !pub.ok) throw new Error(pub.error || "publish failed");
    if (pub.results?.length !== 2) throw new Error(`expected 2 editions, got ${pub.results?.length}`);
    const hashes = new Set(pub.results.map((r) => r.contentHash));
    if (hashes.size !== 2) throw new Error("edition content hashes must differ");
    for (const r of pub.results) {
      if (!r.listingId) throw new Error(`edition ${r.edition} missing listing`);
    }

    if (!health.keys?.buyer) {
      console.log("bob missing — skip buy; publish OK");
    } else {
      const listingId = pub.results[0].listingId;
      const buyRes = await fetch(`${HELPER}/api/buy`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ listingId, from: "bob", seller: "alice" }),
      });
      const buy = await buyRes.json();
      if (!buyRes.ok || !buy.ok) throw new Error(buy.error || "buy failed");
      console.log("buy ok", { buyTx: buy.buyTx, rewrapOk: buy.rewrapOk });
    }

    console.log("PUBLISHER SMOKE PASSED");
    console.log(
      pub.results
        .map((r) => `edition ${r.edition}: original #${r.originalId} listing #${r.listingId} hash=${r.contentHash.slice(0, 12)}…`)
        .join("\n"),
    );
    // silence unused
    void createHash;
  } finally {
    child.kill("SIGTERM");
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
