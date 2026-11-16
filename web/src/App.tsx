import { useCallback, useEffect, useState } from "react";
import { Background } from "./components/Background";
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
        <section className="shell hero hero--compact">
          <div className="hero__copy">
            <h1 className="hero__brand rise">BookStore</h1>
            <p className="hero__headline rise rise-delay-1">Publish a book.</p>
            <div className="cta-row rise rise-delay-2">
              <a className="btn btn--primary" href="#publish">
                Publish
              </a>
              <a className="btn btn--ghost" href="#shop">
                Shop
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

        <Background
          catalog={catalog}
          helper={helper}
          onRefresh={() => void load()}
        />
      </main>

      <footer className="shell footer">
        <a href="#about">About</a>
        <a href="https://github.com/ksedleni/bookstore">Source</a>
      </footer>
    </>
  );
}
