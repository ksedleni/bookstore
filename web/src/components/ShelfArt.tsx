/** Decorative full-bleed shelf / open-book silhouette — visual anchor, not a card. */
export function ShelfArt() {
  return (
    <svg
      className="atmosphere__shelf"
      viewBox="0 0 640 720"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      aria-hidden="true"
    >
      <defs>
        <linearGradient id="pageGlow" x1="120" y1="80" x2="520" y2="640" gradientUnits="userSpaceOnUse">
          <stop stopColor="#d4b57a" stopOpacity="0.55" />
          <stop offset="0.55" stopColor="#ebe6d8" stopOpacity="0.18" />
          <stop offset="1" stopColor="#2f4a3c" stopOpacity="0.05" />
        </linearGradient>
        <linearGradient id="spine" x1="0" y1="0" x2="1" y2="1">
          <stop stopColor="#2f4a3c" />
          <stop offset="1" stopColor="#1c3228" />
        </linearGradient>
      </defs>

      {/* Open book plane */}
      <path
        d="M70 160C140 120 250 110 320 150C390 110 500 120 570 160V560C500 520 390 510 320 545C250 510 140 520 70 560V160Z"
        fill="url(#pageGlow)"
        opacity="0.9"
      />
      <path
        d="M320 150V545"
        stroke="#ebe6d8"
        strokeOpacity="0.35"
        strokeWidth="2"
      />
      <path
        d="M110 210H250M110 250H230M110 290H245M110 330H210"
        stroke="#ebe6d8"
        strokeOpacity="0.28"
        strokeWidth="3"
        strokeLinecap="round"
      />
      <path
        d="M390 210H530M410 250H530M395 290H530M430 330H530"
        stroke="#ebe6d8"
        strokeOpacity="0.22"
        strokeWidth="3"
        strokeLinecap="round"
      />

      {/* Shelf ledge */}
      <rect x="40" y="600" width="560" height="14" rx="2" fill="#d4b57a" fillOpacity="0.35" />
      <rect x="48" y="614" width="544" height="6" fill="#12161c" fillOpacity="0.35" />

      {/* Spines on the ledge */}
      <g transform="translate(90 470)">
        <rect width="28" height="130" rx="2" fill="url(#spine)" />
        <rect x="36" width="22" height="130" rx="2" fill="#3a5646" />
        <rect x="66" width="34" height="130" rx="2" fill="#24362d" />
        <rect x="108" width="18" height="130" rx="2" fill="#4a6354" />
        <rect x="134" width="30" height="130" rx="2" fill="#1a2b22" />
        <rect x="172" width="24" height="130" rx="2" fill="#355043" />
      </g>
    </svg>
  );
}
