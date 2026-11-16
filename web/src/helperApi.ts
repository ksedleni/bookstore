const DEFAULT_HELPER = "http://127.0.0.1:8787";

export function helperBase(): string {
  return (import.meta.env.VITE_HELPER_URL?.trim() || DEFAULT_HELPER).replace(/\/$/, "");
}

export type HelperStatus = {
  ok: boolean;
  helper: boolean;
  lcdOk: boolean;
  lcd: string;
  chainId: string;
  denom: string;
  keys: {
    publisher: { name: string; address: string } | null;
    buyer: { name: string; address: string } | null;
  };
  note?: string;
  error?: string;
};

export type PublishResult = {
  ok: boolean;
  title: string;
  editions: number;
  results: Array<{
    edition: number;
    filename: string;
    contentHash: string;
    originalId: string;
    mintTx: string;
    listed: boolean;
    price: string | null;
    listingId: string | null;
  }>;
  feeNote?: string;
  error?: string;
};

export type BuyResult = {
  ok: boolean;
  buyTx?: string;
  originalId?: string;
  contentHash?: string;
  rewrapOk?: boolean;
  fetchPath?: string | null;
  rewrapError?: string;
  note?: string;
  feeNote?: string;
  error?: string;
};

export async function fetchHelperStatus(): Promise<HelperStatus> {
  try {
    const res = await fetch(`${helperBase()}/api/health`);
    if (!res.ok) {
      return {
        ok: false,
        helper: false,
        lcdOk: false,
        lcd: "",
        chainId: "",
        denom: "",
        keys: { publisher: null, buyer: null },
        error: `${res.status} from helper`,
      };
    }
    return (await res.json()) as HelperStatus;
  } catch (err) {
    return {
      ok: false,
      helper: false,
      lcdOk: false,
      lcd: "",
      chainId: "",
      denom: "",
      keys: { publisher: null, buyer: null },
      error: err instanceof Error ? err.message : String(err),
    };
  }
}

export async function publishBook(input: {
  file: File;
  title: string;
  editions: number;
  price: string;
  list: boolean;
  from?: string;
}): Promise<PublishResult> {
  const fd = new FormData();
  fd.set("file", input.file);
  fd.set("title", input.title);
  fd.set("editions", String(input.editions));
  fd.set("list", input.list ? "1" : "0");
  if (input.price) fd.set("price", input.price);
  fd.set("from", input.from || "alice");

  const res = await fetch(`${helperBase()}/api/publish`, { method: "POST", body: fd });
  const data = (await res.json()) as PublishResult;
  if (!res.ok || !data.ok) {
    throw new Error(data.error || `publish failed (${res.status})`);
  }
  return data;
}

export async function buyListing(listingId: string, from = "bob"): Promise<BuyResult> {
  const res = await fetch(`${helperBase()}/api/buy`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ listingId, from, seller: "alice" }),
  });
  const data = (await res.json()) as BuyResult;
  if (!res.ok || !data.ok) {
    throw new Error(data.error || `buy failed (${res.status})`);
  }
  return data;
}
