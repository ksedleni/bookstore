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

  let label = "Connecting…";
  if (ledgerOk && helperOk) label = "Connected";
  else if (ledgerOk && !helperOk) label = "Connected — helper offline";
  else if (!ledgerOk) label = "Not connected";

  const ok = ledgerOk && helperOk;

  return (
    <div className="connection" role="status">
      <span className={ok ? "dot dot--ok" : ledgerOk ? "dot" : "dot dot--err"} />
      <span>{label}</span>
      <button type="button" className="btn btn--ghost btn--small" onClick={onRefresh}>
        Refresh
      </button>
    </div>
  );
}
