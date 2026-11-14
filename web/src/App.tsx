import { useCallback, useEffect, useState } from "react";
import { Catalog } from "./components/Catalog";
import { ConnectionBar } from "./components/ConnectionBar";
import { PublishWizard } from "./components/PublishWizard";
import { ShelfArt } from "./components/ShelfArt";
import { Shop } from "./components/Shop";
import { fetchCatalog, lcdBase, setLcdBase } from "./api";
import { fetchHelperStatus, type HelperStatus } from "./helperApi";
import type { CatalogState } from "./types";

export default function App() {
  const [catalog, setCatalog] = useState<CatalogState>({ status: "idle" });
  const [helper, setHelper] = useState<HelperStatus | null>(null);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const load = useCallback(async () => {
    setLcdBase(lcdBase());
    setCatalog({ status: "loading" });
    const [helperStatus, cat] = await Promise.all([
      fetchHelperStatus(),
      fetchCatalog()
        .then((data) => ({ status: "ok" as const, ...data }))
        .catch((err: unknown) => ({
          status: "error" as const,
          message: err instanceof Error ? err.message : String(err),
        })),
    ]);
    setHelper(helperStatus);
    setCatalog(cat);
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  const denom = helper?.denom || "stake";
  const helperReady = Boolean(helper?.ok && helper.helper);
  const originals = catalog.status === "ok" ? catalog.originals : [];
  const listings = catalog.status === "ok" ? catalog.listings : [];

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
          <a href="#publish">Publish</a>
          <a href="#shop">Shop</a>
          <a href="#about">About</a>
        </nav>
      </header>

      <main id="top">
        <section className="shell hero">
          <div className="hero__copy">
            <h1 className="hero__brand rise">BookStore</h1>
            <p className="hero__headline rise rise-delay-1">
              Publish like a paperback.
            </p>
            <p className="hero__lede rise rise-delay-2">
              Mint numbered editions, list them for sale, or skip this page entirely —
              peer-to-peer trade still works with keys alone. 0% protocol fee. 0.1%
              creator royalty.
            </p>
            <div className="cta-row rise rise-delay-3">
              <a className="btn btn--primary" href="#publish">
                Publish a book
              </a>
              <a className="btn btn--ghost" href="#shop">
                Browse sales
              </a>
            </div>
          </div>
        </section>

        <div className="shell">
          <ConnectionBar helper={helper} catalog={catalog} onRefresh={() => void load()} />
        </div>

        <div className="shell">
          <PublishWizard
            helperReady={helperReady}
            denom={denom}
            onDone={() => void load()}
          />
        </div>

        <div className="shell">
          <Shop
            originals={originals}
            listings={listings}
            helperReady={helperReady}
            onDone={() => void load()}
          />
        </div>

        <section className="shell panel" id="about" aria-labelledby="about-title">
          <p className="section__eyebrow">Protocol first</p>
          <h2 id="about-title" className="section__title">
            This is a storefront, not the store
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
              <p>CLI P2P sale still works without any UI.</p>
            </div>
          </div>
          <p className="section__lede" style={{ marginTop: "1.5rem" }}>
            Localnet signing uses a small helper over your test keyring (alice /
            bob). Browser wallets like Keplr can replace that later — still
            non-custodial.
          </p>
        </section>

        <section className="shell panel">
          <button
            type="button"
            className="btn btn--ghost btn--small"
            onClick={() => setShowAdvanced((v) => !v)}
            aria-expanded={showAdvanced}
          >
            {showAdvanced ? "Hide advanced" : "Advanced / ledger details"}
          </button>
          {showAdvanced && (
            <div className="advanced">
              <p className="section__lede">
                Power-user view of originals, content fingerprints, and raw listings.
              </p>
              <Catalog state={catalog} onRefresh={() => void load()} />
            </div>
          )}
        </section>
      </main>

      <footer className="shell footer">
        <span>Reference UI · fork freely · protocol remains primary</span>
        <a href="https://github.com/ksedleni/bookstore">ksedleni/bookstore</a>
      </footer>
    </>
  );
}
