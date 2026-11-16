import type { CatalogState, Listing, Original } from "../types";
import { shortAddr, shortHash } from "../api";

type Props = {
  state: CatalogState;
  onRefresh: () => void;
};

function listingFor(originalId: string, listings: Listing[]): Listing | undefined {
  return listings.find((l) => l.active && l.originalId === originalId);
}

function OriginalRow({
  original,
  listing,
  index,
}: {
  original: Original;
  listing?: Listing;
  index: number;
}) {
  return (
    <article className="row" style={{ animationDelay: `${0.04 * index}s` }}>
      <div className="row__id">#{original.id}</div>
      <div>
        <span className="row__label">contentHash</span>
        <div className="row__mono" title={original.contentHash}>
          {shortHash(original.contentHash)}
        </div>
        {original.metadataUri ? (
          <>
            <span className="row__label" style={{ marginTop: "0.55rem" }}>
              metadataUri
            </span>
            <div className="row__mono" title={original.metadataUri}>
              {original.metadataUri.length > 48
                ? `${original.metadataUri.slice(0, 48)}…`
                : original.metadataUri}
            </div>
          </>
        ) : (
          <div className="row__meta" style={{ marginTop: "0.45rem" }}>
            No metadata URI on-chain
          </div>
        )}
      </div>
      <div>
        <span className="row__label">owner</span>
        <div className="row__mono" title={original.owner}>
          {shortAddr(original.owner)}
        </div>
        <span className="row__label" style={{ marginTop: "0.55rem" }}>
          creator
        </span>
        <div className="row__mono" title={original.creator}>
          {shortAddr(original.creator)}
        </div>
      </div>
      <div>
        {listing ? (
          <>
            <span className="row__label">listed</span>
            <div className="row__meta">
              {listing.priceAmount} {listing.denom}
            </div>
            <span className="badge">listing #{listing.id}</span>
          </>
        ) : (
          <span className="row__meta">Not listed — sell via P2P CLI</span>
        )}
      </div>
    </article>
  );
}

export function Catalog({ state, onRefresh }: Props) {
  if (state.status === "loading" || state.status === "idle") {
    return (
      <p className="empty">
        Reading ledger via LCD…{" "}
        <button type="button" className="btn btn--ghost btn--small" onClick={onRefresh}>
          Retry
        </button>
      </p>
    );
  }

  if (state.status === "error") {
    return (
      <div className="empty">
        <p>Ledger API unreachable.</p>
        <p className="muted">{state.message}</p>
        <button type="button" className="btn btn--ghost btn--small" onClick={onRefresh}>
          Retry
        </button>
      </div>
    );
  }

  if (state.status !== "ok") {
    return null;
  }

  if (state.originals.length === 0) {
    return (
      <div className="empty">
        <p>No originals on this node yet.</p>
        <button type="button" className="btn btn--ghost btn--small" onClick={onRefresh}>
          Refresh
        </button>
      </div>
    );
  }

  return (
    <div className="catalog">
      {state.originals.map((o, i) => (
        <OriginalRow
          key={o.id}
          original={o}
          listing={listingFor(o.id, state.listings)}
          index={i}
        />
      ))}
    </div>
  );
}
