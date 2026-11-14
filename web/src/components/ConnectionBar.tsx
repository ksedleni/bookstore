import type { HelperStatus } from "../helperApi";
import type { CatalogState } from "../types";

type Props = {
  helper: HelperStatus | null;
  catalog: CatalogState;
  onRefresh: () => void;
};

export function ConnectionBar({ helper, catalog, onRefresh }: Props) {
  const ledgerOk = catalog.status === "ok";
  const helperOk = Boolean(helper?.ok && helper.helper);
  const connected = ledgerOk;

  return (
    <div className="connection" role="status">
      <span className={connected ? "dot dot--ok" : "dot dot--err"} />
      <span>
        {connected
          ? `Connected · ${catalog.status === "ok" ? catalog.originals.length : 0} books on ledger`
          : "Not connected to ledger"}
      </span>
      <span className="connection__sep">·</span>
      <span className={helperOk ? "muted" : "warn"}>
        {helperOk ? "Publisher helper ready" : "Helper offline — run make publisher"}
      </span>
      <button type="button" className="btn btn--ghost btn--small" onClick={onRefresh}>
        Refresh
      </button>
    </div>
  );
}
