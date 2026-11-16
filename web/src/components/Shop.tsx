import { useState } from "react";
import type { Listing, Original } from "../types";
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
  if (o.metadataUri && !o.metadataUri.startsWith("inline:")) return o.metadataUri;
  return `Edition #${o.id}`;
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
      setError("Helper offline — start it with make publisher.");
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
      <h2 id="shop-title" className="section__title">
        Shop
      </h2>

      {active.length === 0 ? (
        <p className="empty">No books listed yet.</p>
      ) : (
        <div className="shop-list">
          {active.map((l) => {
            const original = originals.find((o) => o.id === l.originalId);
            return (
              <article key={l.id} className="shop-row">
                <div>
                  <h3>{original ? titleFor(original) : `Book #${l.originalId}`}</h3>
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
          <p>
            {lastBuy.rewrapOk
              ? "Purchase complete — your copy is ready."
              : "Purchase complete. Content handoff from the seller may still be needed."}
          </p>
        </div>
      )}
    </section>
  );
}
