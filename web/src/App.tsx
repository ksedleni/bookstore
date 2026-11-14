import { useCallback, useEffect, useState, type FormEvent } from "react";
import { Catalog } from "./components/Catalog";
import { ShelfArt } from "./components/ShelfArt";
import { fetchCatalog, lcdBase, rpcUrl, setLcdBase } from "./api";
import type { CatalogState } from "./types";

const LCD_STORAGE_KEY = "bookstore.lcdUrl";

function initialLcd(): string {
  try {
    return localStorage.getItem(LCD_STORAGE_KEY) || lcdBase();
  } catch {
    return lcdBase();
  }
}

export default function App() {
  const [lcdDraft, setLcdDraft] = useState(initialLcd);
  const [activeLcd, setActiveLcd] = useState(initialLcd);
  const [catalog, setCatalog] = useState<CatalogState>({ status: "idle" });

  const load = useCallback(async () => {
    setCatalog({ status: "loading" });
    setLcdBase(activeLcd);
    try {
      const data = await fetchCatalog();
      setCatalog({ status: "ok", ...data });
    } catch (err) {
      setCatalog({
        status: "error",
        message: err instanceof Error ? err.message : String(err),
      });
    }
  }, [activeLcd]);

  useEffect(() => {
    void load();
  }, [load]);

  function applyEndpoint(event: FormEvent) {
    event.preventDefault();
    const next = lcdDraft.trim().replace(/\/$/, "") || "/lcd";
    try {
      localStorage.setItem(LCD_STORAGE_KEY, next);
    } catch {
      /* ignore */
    }
    setActiveLcd(next);
  }

  const statusClass =
    catalog.status === "ok"
      ? "status-pill status-pill--ok"
      : catalog.status === "error"
        ? "status-pill status-pill--err"
        : "status-pill status-pill--busy";

  const statusLabel =
    catalog.status === "ok"
      ? `${catalog.originals.length} original${catalog.originals.length === 1 ? "" : "s"} · ${catalog.listings.filter((l) => l.active).length} active listing${catalog.listings.filter((l) => l.active).length === 1 ? "" : "s"}`
      : catalog.status === "error"
        ? "LCD unreachable"
        : "Connecting…";

  return (
    <>
      <div className="atmosphere" aria-hidden="true">
        <ShelfArt />
      </div>

      <header className="shell topbar rise">
        <a className="topbar__mark" href="#top">
          BookStore
        </a>
        <nav className="topbar__nav" aria-label="Primary">
          <a href="#catalog">Catalog</a>
          <a href="#keys">Keys</a>
          <a
            href="https://github.com/ksedleni/bookstore"
            target="_blank"
            rel="noreferrer"
          >
            Protocol
          </a>
        </nav>
      </header>

      <main id="top">
        <section className="shell hero">
          <div className="hero__copy">
            <h1 className="hero__brand rise">BookStore</h1>
            <p className="hero__headline rise rise-delay-1">
              A storefront. Not the store.
            </p>
            <p className="hero__lede rise rise-delay-2">
              Browse originals on your local ledger. Sales work without this UI —
              peer-to-peer, with a 0% protocol fee and a fixed 0.1% creator royalty.
            </p>
            <div className="cta-row rise rise-delay-3">
              <a className="btn btn--primary" href="#catalog">
                Open catalog
              </a>
              <a
                className="btn btn--ghost"
                href="https://github.com/ksedleni/bookstore#walkthrough--mint--fund-buyer--p2p-sale"
                target="_blank"
                rel="noreferrer"
              >
                Use CLI / P2P instead
              </a>
            </div>
          </div>
        </section>

        <section className="shell section" aria-labelledby="principles-title">
          <p className="section__eyebrow">Protocol first</p>
          <h2 id="principles-title" className="section__title">
            Ownership lives on the ledger
          </h2>
          <p className="section__lede">
            This page is a thin, forkable reference client. Anyone can build another
            bookstore — or none at all.
          </p>
          <div className="principles">
            <div className="principle">
              <h3>0% protocol fee</h3>
              <p>BookStore is not a fee-taking store. Settlement does not skim.</p>
            </div>
            <div className="principle">
              <h3>0.1% royalty</h3>
              <p>Fixed creator cut on paid ledger sales — not seller-configurable.</p>
            </div>
            <div className="principle">
              <h3>Bookstore-less trade</h3>
              <p>
                <code style={{ color: "var(--lamp)", fontSize: "0.9em" }}>
                  bookstored tx bookstore p2p-sale
                </code>{" "}
                is enough. No storefront required.
              </p>
            </div>
          </div>
        </section>

        <section className="shell section" id="catalog" aria-labelledby="catalog-title">
          <p className="section__eyebrow">Local ledger</p>
          <h2 id="catalog-title" className="section__title">
            Catalog
          </h2>
          <p className="section__lede">
            Originals and optional listings from the BookStore LCD. Content bytes
            stay off-chain; only <code style={{ color: "var(--lamp)" }}>contentHash</code>{" "}
            and metadata pointers are shown here.
          </p>

          <form className="toolbar" onSubmit={applyEndpoint}>
            <div className="field">
              <label htmlFor="lcd">LCD endpoint</label>
              <input
                id="lcd"
                name="lcd"
                value={lcdDraft}
                onChange={(e) => setLcdDraft(e.target.value)}
                placeholder="/lcd or http://127.0.0.1:1317"
                spellCheck={false}
                autoComplete="off"
              />
            </div>
            <button type="submit" className="btn btn--ghost btn--small">
              Use endpoint
            </button>
            <button type="button" className="btn btn--ghost btn--small" onClick={() => void load()}>
              Refresh
            </button>
            <span className={statusClass} role="status">
              {statusLabel}
            </span>
          </form>

          <Catalog state={catalog} onRefresh={() => void load()} />
        </section>

        <section className="shell section" id="keys" aria-labelledby="keys-title">
          <p className="section__eyebrow">Non-custodial</p>
          <h2 id="keys-title" className="section__title">
            Keys stay yours
          </h2>
          <div className="keys-panel">
            <p>
              This reference UI does not custody keys or invent a wallet. For localnet,
              mint, sell, and fetch content with the existing{" "}
              <code>bookstored</code> keyring (e.g. <code>--keyring-backend test</code>).
            </p>
            <p>
              Browser wallet connect (Keplr / Leap and similar) can be wired later against
              the same chain id and LCD — still non-custodial, still optional.
            </p>
            <p>
              CometBFT RPC default: <code>{rpcUrl()}</code>
            </p>
          </div>
        </section>
      </main>

      <footer className="shell footer">
        <span>Reference storefront · fork freely · protocol remains primary</span>
        <span>
          <a href="https://github.com/ksedleni/bookstore">ksedleni/bookstore</a>
        </span>
      </footer>
    </>
  );
}
