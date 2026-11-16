import { useState } from "react";
import { Catalog } from "./Catalog";
import type { CatalogState } from "../types";
import type { HelperStatus } from "../helperApi";
import { lcdBase, rpcUrl } from "../api";
import { helperBase } from "../helperApi";

type Props = {
  catalog: CatalogState;
  helper: HelperStatus | null;
  onRefresh: () => void;
};

/** All protocol / technical / explainer copy lives here — not in Publish or Shop. */
export function Background({ catalog, helper, onRefresh }: Props) {
  const [showLedger, setShowLedger] = useState(false);

  return (
    <section className="shell panel background" id="about" aria-labelledby="about-title">
      <p className="section__eyebrow">Background</p>
      <h2 id="about-title" className="section__title">
        About BookStore
      </h2>

      <div className="principles">
        <div className="principle">
          <h3>0% protocol fee</h3>
          <p>BookStore does not take a cut of settlement.</p>
        </div>
        <div className="principle">
          <h3>0.1% royalty</h3>
          <p>Fixed creator share on paid sales — not seller-configurable.</p>
        </div>
        <div className="principle">
          <h3>Bookstore-less</h3>
          <p>
            This UI is optional. Peer-to-peer sale still works with keys and the
            CLI alone.
          </p>
        </div>
      </div>

      <div className="background__block">
        <h3 className="background__h">How publishing works</h3>
        <p>
          Each edition is a distinct file fingerprint and ownership record. When
          you mint more than one edition, a small edition marker is appended so
          hashes differ — the ledger rejects duplicate remints of the same bytes.
        </p>
        <p>
          Paid sales apply the fixed 0.1% creator royalty. Content bytes stay
          off-chain; the ledger stores only the digest.
        </p>
      </div>

      <div className="background__block">
        <h3 className="background__h">Localnet keys</h3>
        <p>
          Signing uses a small helper over your test keyring (alice publishes,
          bob buys). Browser wallets like Keplr can replace that later — still
          non-custodial.
        </p>
        {helper?.keys?.publisher && (
          <p className="muted mono-line">
            Publisher ({helper.keys.publisher.name}): {helper.keys.publisher.address}
          </p>
        )}
        {helper?.keys?.buyer && (
          <p className="muted mono-line">
            Buyer ({helper.keys.buyer.name}): {helper.keys.buyer.address}
          </p>
        )}
      </div>

      <div className="background__block">
        <h3 className="background__h">Endpoints</h3>
        <p className="muted mono-line">Ledger API: {lcdBase()}</p>
        <p className="muted mono-line">Helper: {helperBase()}</p>
        <p className="muted mono-line">RPC: {rpcUrl()}</p>
        <p className="muted mono-line">
          Chain: {helper?.chainId || "bookstore-1"} · denom {helper?.denom || "stake"}
        </p>
      </div>

      <button
        type="button"
        className="btn btn--ghost btn--small"
        onClick={() => setShowLedger((v) => !v)}
        aria-expanded={showLedger}
      >
        {showLedger ? "Hide ledger details" : "Show ledger details"}
      </button>
      {showLedger && (
        <div className="advanced">
          <p className="section__lede">
            Originals, content hashes, metadata URIs, and raw listings.
          </p>
          <Catalog state={catalog} onRefresh={onRefresh} />
        </div>
      )}
    </section>
  );
}
