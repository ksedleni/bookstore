import { useState, type FormEvent } from "react";
import { publishBook, type PublishResult } from "../helperApi";

type Props = {
  helperReady: boolean;
  denom: string;
  onDone: () => void;
};

export function PublishWizard({ helperReady, denom, onDone }: Props) {
  const [title, setTitle] = useState("");
  const [editions, setEditions] = useState(2);
  const [listForSale, setListForSale] = useState(true);
  const [price, setPrice] = useState("1000000");
  const [file, setFile] = useState<File | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<PublishResult | null>(null);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setResult(null);
    if (!file) {
      setError("Choose a file to publish.");
      return;
    }
    if (!helperReady) {
      setError("Start the local helper first (make publisher).");
      return;
    }
    setBusy(true);
    try {
      const out = await publishBook({
        file,
        title: title.trim() || file.name,
        editions,
        price: listForSale ? price : "",
        list: listForSale,
      });
      setResult(out);
      onDone();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="panel" id="publish" aria-labelledby="publish-title">
      <p className="section__eyebrow">Publisher</p>
      <h2 id="publish-title" className="section__title">
        Publish a book
      </h2>
      <p className="section__lede">
        Pick a file, optionally mint numbered editions (each is a unique original),
        then list for sale. Protocol fee stays 0%; paid sales take a fixed 0.1% royalty.
      </p>

      <form className="wizard" onSubmit={onSubmit}>
        <label className="field">
          <span>Book file</span>
          <input
            type="file"
            onChange={(e) => setFile(e.target.files?.[0] ?? null)}
            required
          />
        </label>

        <label className="field">
          <span>Title</span>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder={file?.name || "My paperback"}
          />
        </label>

        <label className="field">
          <span>Numbered editions</span>
          <input
            type="number"
            min={1}
            max={20}
            value={editions}
            onChange={(e) => setEditions(Number(e.target.value) || 1)}
          />
          <span className="field__hint">
            Each edition gets its own file fingerprint and ownership record (not
            duplicate remints of the same bytes).
          </span>
        </label>

        <label className="check">
          <input
            type="checkbox"
            checked={listForSale}
            onChange={(e) => setListForSale(e.target.checked)}
          />
          <span>List for sale after mint</span>
        </label>

        {listForSale && (
          <label className="field">
            <span>Price ({denom})</span>
            <input
              type="text"
              inputMode="numeric"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              required={listForSale}
            />
          </label>
        )}

        <button className="btn btn--primary" type="submit" disabled={busy}>
          {busy ? "Publishing…" : listForSale ? "Mint & list" : "Mint only"}
        </button>
      </form>

      {error && <p className="banner banner--err">{error}</p>}

      {result && (
        <div className="banner banner--ok">
          <p>
            Published <strong>{result.title}</strong> — {result.editions} edition
            {result.editions === 1 ? "" : "s"}.
          </p>
          <ul className="result-list">
            {result.results.map((r) => (
              <li key={r.originalId}>
                Edition {r.edition}: original #{r.originalId}
                {r.listed ? ` · listed #${r.listingId} @ ${r.price}` : " · not listed"}
              </li>
            ))}
          </ul>
          <p className="muted">{result.feeNote}</p>
        </div>
      )}
    </section>
  );
}
