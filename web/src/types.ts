export type Original = {
  id: string;
  contentHash: string;
  metadataHash: string;
  metadataUri: string;
  owner: string;
  royaltyPayee: string;
  creator: string;
};

export type Listing = {
  id: string;
  originalId: string;
  seller: string;
  priceAmount: string;
  denom: string;
  active: boolean;
};

export type CatalogState =
  | { status: "idle" | "loading" }
  | { status: "error"; message: string }
  | { status: "ok"; originals: Original[]; listings: Listing[] };
