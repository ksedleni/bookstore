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
      setError("Helper offline — start it with make publisher.");
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
      <h2 id="publish-title" className="section__title">
        Publish a book
      </h2>

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
            placeholder={file?.name || "Title"}
          />
        </label>

        <label className="field">
          <span>Editions</span>
          <input
            type="number"
            min={1}
            max={20}
            value={editions}
            onChange={(e) => setEditions(Number(e.target.value) || 1)}
          />
        </label>

        <label className="check">
          <input
            type="checkbox"
            checked={listForSale}
            onChange={(e) => setListForSale(e.target.checked)}
          />
          <span>List for sale</span>
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
          {busy ? "Publishing…" : listForSale ? "Mint & list" : "Mint"}
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
                Edition {r.edition}
                {r.listed ? ` · listed @ ${r.price}` : ""}
              </li>
            ))}
          </ul>
        </div>
      )}
    </section>
  );
}
