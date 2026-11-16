#!/usr/bin/env node
/**
 * Local publisher helper — thin HTTP API over bookstored for the reference UI.
 * Binds 127.0.0.1 only. Uses the local test keyring; never custodial remote keys.
 */
import http from "node:http";
import { spawn } from "node:child_process";
import { createHash, randomUUID } from "node:crypto";
import fs from "node:fs/promises";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "../..");

const PORT = Number(process.env.HELPER_PORT || 8787);
const HOME = process.env.BOOKSTORE_HOME || path.join(os.homedir(), ".bookstore");
const BIN =
  process.env.BOOKSTORED_BIN ||
  path.join(ROOT, ".bin", "bookstored");
const CHAIN_ID = process.env.CHAIN_ID || "bookstore-1";
const NODE = process.env.BOOKSTORE_NODE || "tcp://127.0.0.1:26657";
const DENOM = process.env.BOOKSTORE_DENOM || "stake";
const KEYRING = process.env.KEYRING_BACKEND || "test";
const LCD = process.env.BOOKSTORE_LCD || "http://127.0.0.1:1317";

function sha256Hex(buf) {
  return createHash("sha256").update(buf).digest("hex");
}

function run(args, opts = {}) {
  return new Promise((resolve, reject) => {
    const child = spawn(BIN, args, {
      env: { ...process.env },
      ...opts,
    });
    let stdout = "";
    let stderr = "";
    child.stdout.on("data", (d) => {
      stdout += d;
    });
    child.stderr.on("data", (d) => {
      stderr += d;
    });
    child.on("error", reject);
    child.on("close", (code) => {
      if (code !== 0) {
        reject(new Error(`${args.join(" ")}\n${stderr || stdout}`.trim()));
        return;
      }
      resolve({ stdout, stderr });
    });
  });
}

function txFlags(from) {
  return [
    "--from",
    from,
    "--home",
    HOME,
    "--keyring-backend",
    KEYRING,
    "--chain-id",
    CHAIN_ID,
    "--node",
    NODE,
    "--yes",
    "--gas",
    "auto",
    "--gas-adjustment",
    "1.5",
    "--fees",
    `200${DENOM}`,
    "-o",
    "json",
  ];
}

async function waitTx(txhash) {
  for (let i = 0; i < 40; i++) {
    try {
      const { stdout } = await run([
        "q",
        "tx",
        txhash,
        "--home",
        HOME,
        "--node",
        NODE,
        "-o",
        "json",
      ]);
      const j = JSON.parse(stdout);
      if (j.code === 0) return j;
      throw new Error(`tx failed code=${j.code}: ${j.raw_log || ""}`);
    } catch (err) {
      if (String(err.message).includes("tx failed")) throw err;
    }
    await new Promise((r) => setTimeout(r, 500));
  }
  throw new Error(`timeout waiting for tx ${txhash}`);
}

async function broadcast(args) {
  const { stdout } = await run(args);
  // gas estimate may print before JSON
  const line = stdout
    .trim()
    .split("\n")
    .reverse()
    .find((l) => l.trim().startsWith("{"));
  if (!line) throw new Error(`no json from: ${args.join(" ")}\n${stdout}`);
  const j = JSON.parse(line);
  if (!j.txhash) throw new Error(`no txhash: ${stdout}`);
  await waitTx(j.txhash);
  return j.txhash;
}

async function keyAddr(name) {
  const { stdout } = await run([
    "keys",
    "show",
    name,
    "-a",
    "--home",
    HOME,
    "--keyring-backend",
    KEYRING,
  ]);
  return stdout.trim();
}

async function contentPut(from, filePath) {
  const { stdout } = await run([
    "content",
    "put",
    "--from",
    from,
    "--file",
    filePath,
    "--home",
    HOME,
    "--keyring-backend",
    KEYRING,
  ]);
  const j = JSON.parse(stdout);
  return j.content_hash;
}

async function findOriginalIdByHash(contentHash) {
  const res = await fetch(
    `${LCD}/ksedleni/bookstore/bookstore/original?pagination.limit=200`,
  );
  if (!res.ok) throw new Error(`LCD originals ${res.status}`);
  const data = await res.json();
  const rows = data.Original || data.original || [];
  const hit = rows.find((o) => (o.contentHash || o.content_hash) === contentHash);
  if (!hit) throw new Error(`original not found for hash ${contentHash}`);
  return String(hit.id);
}

async function findListingId(originalId) {
  const res = await fetch(
    `${LCD}/ksedleni/bookstore/marketplace/listing?pagination.limit=200`,
  );
  if (!res.ok) throw new Error(`LCD listings ${res.status}`);
  const data = await res.json();
  const rows = data.Listing || data.listing || [];
  const hit = rows.find(
    (l) =>
      l.active &&
      String(l.originalId || l.original_id) === String(originalId),
  );
  return hit ? String(hit.id) : null;
}

function cors(res) {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader("Access-Control-Allow-Methods", "GET,POST,OPTIONS");
  res.setHeader("Access-Control-Allow-Headers", "Content-Type");
}

function sendJson(res, status, body) {
  cors(res);
  res.writeHead(status, { "Content-Type": "application/json" });
  res.end(JSON.stringify(body, null, 2));
}

async function readBody(req) {
  const chunks = [];
  for await (const c of req) chunks.push(c);
  return Buffer.concat(chunks);
}

/** Minimal multipart parser for one file field + text fields. */
function parseMultipart(buf, boundary) {
  const parts = {};
  const files = {};
  const sep = Buffer.from(`--${boundary}`);
  let start = buf.indexOf(sep) + sep.length;
  while (start < buf.length) {
    if (buf[start] === 45 && buf[start + 1] === 45) break; // --
    if (buf[start] === 13 && buf[start + 1] === 10) start += 2;
    const next = buf.indexOf(sep, start);
    const chunk = buf.subarray(start, next === -1 ? buf.length : next - 2);
    const headerEnd = chunk.indexOf("\r\n\r\n");
    if (headerEnd === -1) break;
    const header = chunk.subarray(0, headerEnd).toString("utf8");
    const body = chunk.subarray(headerEnd + 4);
    const nameMatch = /name="([^"]+)"/.exec(header);
    const fileMatch = /filename="([^"]*)"/.exec(header);
    if (nameMatch) {
      const name = nameMatch[1];
      if (fileMatch) {
        files[name] = { filename: fileMatch[1] || "upload.bin", data: body };
      } else {
        parts[name] = body.toString("utf8");
      }
    }
    if (next === -1) break;
    start = next + sep.length;
  }
  return { parts, files };
}

async function handleStatus() {
  let publisher = null;
  let buyer = null;
  try {
    publisher = await keyAddr("alice");
  } catch {
    /* missing */
  }
  try {
    buyer = await keyAddr("bob");
  } catch {
    /* missing */
  }
  let lcdOk = false;
  try {
    const r = await fetch(`${LCD}/cosmos/base/tendermint/v1beta1/node_info`);
    lcdOk = r.ok;
  } catch {
    lcdOk = false;
  }
  return {
    ok: true,
    helper: true,
    home: HOME,
    chainId: CHAIN_ID,
    denom: DENOM,
    lcd: LCD,
    lcdOk,
    bin: BIN,
    keys: { publisher: publisher ? { name: "alice", address: publisher } : null, buyer: buyer ? { name: "bob", address: buyer } : null },
    note: "Local helper signs with the test keyring only. Keplr can replace this later.",
  };
}

async function handlePublish(parts, files) {
  const file = files.file;
  if (!file?.data?.length) throw new Error("file required");
  const from = parts.from || "alice";
  const title = parts.title || file.filename || "Untitled";
  const editions = Math.max(1, Math.min(20, Number(parts.editions || 1)));
  const price = parts.price ? String(parts.price).trim() : "";
  const list = parts.list === "1" || parts.list === "true" || Boolean(price);
  if (list && !price) throw new Error("price required when listing");

  const work = path.join(os.tmpdir(), `bookstore-pub-${randomUUID()}`);
  await fs.mkdir(work, { recursive: true });
  const results = [];

  try {
    for (let i = 1; i <= editions; i++) {
      const editionBytes =
        editions === 1
          ? file.data
          : Buffer.concat([
              file.data,
              Buffer.from(
                `\n\n---\nBookStore edition ${i} of ${editions}\ntitle: ${title}\n`,
                "utf8",
              ),
            ]);
      const editionName =
        editions === 1
          ? file.filename || "book.bin"
          : `${path.parse(file.filename || "book").name}-edition-${i}${path.extname(file.filename || "") || ".bin"}`;
      const filePath = path.join(work, editionName);
      await fs.writeFile(filePath, editionBytes);

      const contentHash = await contentPut(from, filePath);
      const metaObj = {
        title,
        edition: i,
        editions,
        filename: editionName,
        content_hash: contentHash,
      };
      const metaRaw = Buffer.from(JSON.stringify(metaObj), "utf8");
      const metaHash = sha256Hex(metaRaw);
      const metaUri = `inline:edition-${i}`;

      const mintTx = await broadcast([
        "tx",
        "bookstore",
        "mint-original",
        contentHash,
        metaHash,
        metaUri,
        ...txFlags(from),
      ]);
      const originalId = await findOriginalIdByHash(contentHash);

      let listingId = null;
      let listTx = null;
      if (list) {
        listTx = await broadcast([
          "tx",
          "marketplace",
          "create-listing",
          originalId,
          price,
          DENOM,
          ...txFlags(from),
        ]);
        listingId = await findListingId(originalId);
      }

      results.push({
        edition: i,
        filename: editionName,
        contentHash,
        originalId,
        mintTx,
        listed: Boolean(list),
        price: list ? `${price}${DENOM}` : null,
        listingId,
        listTx,
      });
    }
  } finally {
    await fs.rm(work, { recursive: true, force: true });
  }

  return {
    ok: true,
    title,
    editions,
    results,
    feeNote: "Protocol fee 0%. Paid sales include a fixed 0.1% creator royalty.",
  };
}

async function handleBuy(body) {
  const listingId = String(body.listingId || "");
  const buyer = body.from || "bob";
  const seller = body.seller || "alice";
  if (!listingId) throw new Error("listingId required");

  const buyTx = await broadcast([
    "tx",
    "marketplace",
    "buy-listing",
    listingId,
    ...txFlags(buyer),
  ]);

  // Resolve content hash from listing → original
  const listRes = await fetch(
    `${LCD}/ksedleni/bookstore/marketplace/listing/${listingId}`,
  );
  if (!listRes.ok) throw new Error(`listing query ${listRes.status}`);
  const listData = await listRes.json();
  const listing = listData.Listing || listData.listing;
  const originalId = String(listing.originalId || listing.original_id);
  const origRes = await fetch(
    `${LCD}/ksedleni/bookstore/bookstore/original/${originalId}`,
  );
  const origData = await origRes.json();
  const original = origData.Original || origData.original;
  const contentHash = original.contentHash || original.content_hash;

  let rewrapOk = false;
  let fetchPath = null;
  try {
    await run([
      "content",
      "rewrap",
      "--from",
      seller,
      "--to",
      buyer,
      "--content-hash",
      contentHash,
      "--home",
      HOME,
      "--keyring-backend",
      KEYRING,
    ]);
    rewrapOk = true;
    const out = path.join(os.tmpdir(), `bookstore-buy-${contentHash.slice(0, 12)}.bin`);
    await run([
      "content",
      "get",
      "--from",
      buyer,
      "--content-hash",
      contentHash,
      "--out",
      out,
      "--home",
      HOME,
      "--keyring-backend",
      KEYRING,
    ]);
    fetchPath = out;
  } catch (err) {
    return {
      ok: true,
      buyTx,
      originalId,
      contentHash,
      rewrapOk: false,
      rewrapError: String(err.message || err),
      note: "Purchase settled on-ledger. Ask the seller to rewrap content for your key, then fetch.",
    };
  }

  return {
    ok: true,
    buyTx,
    originalId,
    contentHash,
    rewrapOk,
    fetchPath,
    feeNote: "0% protocol fee; 0.1% royalty applied on paid sale.",
  };
}

const server = http.createServer(async (req, res) => {
  cors(res);
  if (req.method === "OPTIONS") {
    res.writeHead(204);
    res.end();
    return;
  }

  try {
    const url = new URL(req.url || "/", `http://127.0.0.1:${PORT}`);

    if (req.method === "GET" && url.pathname === "/api/health") {
      sendJson(res, 200, await handleStatus());
      return;
    }

    if (req.method === "POST" && url.pathname === "/api/publish") {
      const ctype = req.headers["content-type"] || "";
      const m = /boundary=(?:"([^"]+)"|([^;]+))/i.exec(ctype);
      if (!m) {
        sendJson(res, 400, { ok: false, error: "multipart required" });
        return;
      }
      const boundary = m[1] || m[2];
      const buf = await readBody(req);
      const { parts, files } = parseMultipart(buf, boundary);
      const out = await handlePublish(parts, files);
      sendJson(res, 200, out);
      return;
    }

    if (req.method === "POST" && url.pathname === "/api/buy") {
      const buf = await readBody(req);
      const body = JSON.parse(buf.toString("utf8") || "{}");
      const out = await handleBuy(body);
      sendJson(res, 200, out);
      return;
    }

    sendJson(res, 404, { ok: false, error: "not found" });
  } catch (err) {
    sendJson(res, 500, { ok: false, error: String(err.message || err) });
  }
});

server.listen(PORT, "127.0.0.1", () => {
  console.log(`BookStore helper on http://127.0.0.1:${PORT}`);
  console.log(`  home=${HOME} bin=${BIN} lcd=${LCD}`);
});
