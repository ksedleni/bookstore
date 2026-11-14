import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, ".", "");
  const lcdTarget = env.VITE_LCD_PROXY_TARGET || "http://127.0.0.1:1317";

  return {
    plugins: [react()],
    server: {
      port: 5173,
      proxy: {
        // Browser → Vite → local bookstored LCD (avoids CORS during localnet).
        "/lcd": {
          target: lcdTarget,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/lcd/, ""),
        },
      },
    },
  };
});
