import { useState } from "react";
import type { Listing, Original } from "../types";
import { shortAddr } from "../api";
import { buyListing, type BuyResult } from "../helperApi";

type Props = {
  originals: Original[];
  listings: Listing[];
  helperReady: boolean;
  onDone: () => void;
};

function titleFor(o: Original): string {
  if (o.metadataUri?.startsWith("inline:edition-")) {
    return `Edition ${o.metadataUri.replace("inline:edition-", "")}`;
  }
  if (o.metadataUri) return o.metadataUri;
  return `Original #${o.id}`;
}

export function Shop({ originals, listings, helperReady, onDone }: Props) {
  const active = listings.filter((l) => l.active);
  const [busyId, setBusyId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [lastBuy, setLastBuy] = useState<BuyResult | null>(null);

  async function onBuy(listingId: string) {
    setError(null);
    setLastBuy(null);
    if (!helperReady) {
      setError("Start the local helper (make publisher) to buy on localnet.");
      return;
    }
    setBusyId(listingId);
    try {
      const out = await buyListing(listingId);
      setLastBuy(out);
      onDone();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setBusyId(null);
    }
  }

  return (
    <section className="panel" id="shop" aria-labelledby="shop-title">
      <p className="section__eyebrow">Shop</p>
      <h2 id="shop-title" className="section__title">
        Books for sale
      </h2>
      <p className="section__lede">
        Browse listings on your local ledger. Buying settles peer-to-peer on-chain
        with a 0% protocol fee and 0.1% creator royalty — this page is optional.
      </p>

      {active.length === 0 ? (
        <p className="empty">Nothing listed yet. Publish a book above, or sell via P2P without this UI.</p>
      ) : (
        <div className="shop-list">
          {active.map((l) => {
            const original = originals.find((o) => o.id === l.originalId);
            return (
              <article key={l.id} className="shop-row">
                <div>
                  <h3>{original ? titleFor(original) : `Original #${l.originalId}`}</h3>
                  <p className="muted">
                    Seller {shortAddr(l.seller)} · listing #{l.id}
                  </p>
                </div>
                <div className="shop-row__price">
                  <span>
                    {l.priceAmount} {l.denom}
                  </span>
                  <button
                    type="button"
                    className="btn btn--primary btn--small"
                    disabled={busyId === l.id}
                    onClick={() => void onBuy(l.id)}
                  >
                    {busyId === l.id ? "Buying…" : "Buy"}
                  </button>
                </div>
              </article>
            );
          })}
        </div>
      )}

      {error && <p className="banner banner--err">{error}</p>}
      {lastBuy && (
        <div className="banner banner--ok">
          <p>Purchase settled{lastBuy.buyTx ? ` (tx ${lastBuy.buyTx.slice(0, 10)}…)` : ""}.</p>
          {lastBuy.rewrapOk ? (
            <p>
              Content unlocked for the buyer key
              {lastBuy.fetchPath ? ` → ${lastBuy.fetchPath}` : ""}.
            </p>
          ) : (
            <p>
              {lastBuy.note ||
                "Ask the seller to hand off the encrypted file (rewrap), then download as owner."}
            </p>
          )}
          <p className="muted">{lastBuy.feeNote}</p>
        </div>
      )}
    </section>
  );
}
