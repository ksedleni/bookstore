import type { Listing, Original } from "./types";

const DEFAULT_LCD = "/lcd";

let lcdOverride: string | null = null;

export function lcdBase(): string {
  const raw = lcdOverride ?? import.meta.env.VITE_LCD_URL?.trim();
  return (raw && raw.length > 0 ? raw : DEFAULT_LCD).replace(/\/$/, "");
}

export function setLcdBase(url: string) {
  lcdOverride = url.trim().replace(/\/$/, "") || DEFAULT_LCD;
}

export function rpcUrl(): string {
  return (import.meta.env.VITE_RPC_URL?.trim() || "http://127.0.0.1:26657").replace(
    /\/$/,
    "",
  );
}

function pickArray<T>(payload: Record<string, unknown>, keys: string[]): T[] {
  for (const key of keys) {
    const v = payload[key];
    if (Array.isArray(v)) return v as T[];
  }
  return [];
}

function normalizeOriginal(raw: Record<string, unknown>): Original {
  return {
    id: String(raw.id ?? raw.Id ?? ""),
    contentHash: String(raw.contentHash ?? raw.content_hash ?? ""),
    metadataHash: String(raw.metadataHash ?? raw.metadata_hash ?? ""),
    metadataUri: String(raw.metadataUri ?? raw.metadata_uri ?? ""),
    owner: String(raw.owner ?? ""),
    royaltyPayee: String(raw.royaltyPayee ?? raw.royalty_payee ?? ""),
    creator: String(raw.creator ?? ""),
  };
}

function normalizeListing(raw: Record<string, unknown>): Listing {
  return {
    id: String(raw.id ?? ""),
    originalId: String(raw.originalId ?? raw.original_id ?? ""),
    seller: String(raw.seller ?? ""),
    priceAmount: String(raw.priceAmount ?? raw.price_amount ?? ""),
    denom: String(raw.denom ?? ""),
    active: Boolean(raw.active ?? false),
  };
}

async function getJson(path: string): Promise<Record<string, unknown>> {
  const url = `${lcdBase()}${path}`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`${res.status} ${res.statusText} fetching ${url}`);
  }
  return (await res.json()) as Record<string, unknown>;
}

export async function fetchOriginals(): Promise<Original[]> {
  const data = await getJson("/ksedleni/bookstore/bookstore/original?pagination.limit=100");
  const rows = pickArray<Record<string, unknown>>(data, ["Original", "original"]);
  return rows.map(normalizeOriginal);
}

export async function fetchListings(): Promise<Listing[]> {
  const data = await getJson("/ksedleni/bookstore/marketplace/listing?pagination.limit=100");
  const rows = pickArray<Record<string, unknown>>(data, ["Listing", "listing"]);
  return rows.map(normalizeListing);
}

export async function fetchCatalog(): Promise<{ originals: Original[]; listings: Listing[] }> {
  const [originals, listings] = await Promise.all([fetchOriginals(), fetchListings()]);
  return { originals, listings };
}

export function shortAddr(addr: string, head = 10, tail = 6): string {
  if (!addr || addr.length <= head + tail + 1) return addr || "—";
  return `${addr.slice(0, head)}…${addr.slice(-tail)}`;
}

export function shortHash(hash: string): string {
  if (!hash || hash.length < 18) return hash || "—";
  return `${hash.slice(0, 10)}…${hash.slice(-8)}`;
}
