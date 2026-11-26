import react from "@vitejs/plugin-react";
import { ProxyOptions, defineConfig, loadEnv } from "vite";
import { checker } from "vite-plugin-checker";
import svgr from "vite-plugin-svgr";
import tsconfigPaths from "vite-tsconfig-paths";

// Every VITE_* value is compiled into the public JS bundle. The UI used to send the API credentials from these two;
// it now reaches the API through the /api proxy, which adds them on the server. Refuse to build while they are still
// set (e.g. in a leftover ui/.env), so the credentials cannot end up in the bundle again.
const REMOVED_SECRET_VARS = ["VITE_API_USERNAME", "VITE_API_PASSWORD"];

// `npm start` has no nginx in front of it, so the dev server does nginx's job: it forwards <base>/api/* to the issuer
// API and adds the API credentials from ui/.env. These variables have no VITE_ prefix, so they stay out of the bundle.
function devApiProxy(base: string, env: Record<string, string>): Record<string, ProxyOptions> {
  const target = env.ISSUER_UI_API_PROXY_URL;
  const user = env.ISSUER_API_AUTH_USER;
  const password = env.ISSUER_API_AUTH_PASSWORD;
  if (!target || !user || !password) {
    throw new Error(
      "Set ISSUER_UI_API_PROXY_URL, ISSUER_API_AUTH_USER and ISSUER_API_AUTH_PASSWORD in ui/.env (see ui/.env.sample)."
    );
  }
  const prefix = `${base.replace(/\/+$/, "")}/api`;
  return {
    [`${prefix}/`]: {
      changeOrigin: true,
      headers: { Authorization: `Basic ${Buffer.from(`${user}:${password}`).toString("base64")}` },
      rewrite: (path) => path.slice(prefix.length),
      target,
    },
  };
}

// eslint-disable-next-line import/no-default-export
export default defineConfig(({ command, mode }) => {
  // All variables, not only VITE_*: the dev proxy needs the ISSUER_* ones. Only VITE_* reach the bundle.
  const env = loadEnv(mode, process.cwd(), "");

  const leftover = REMOVED_SECRET_VARS.filter((name) => env[name]);
  if (leftover.length > 0) {
    throw new Error(
      `${leftover.join(" and ")} must not be set: VITE_* values are compiled into the public JS bundle. ` +
        "The UI reaches the API through the /api proxy, which adds the credentials on the server. Remove them."
    );
  }

  const base = env.VITE_BASE_URL || "/";

  return {
    base,
    plugins: [
      checker({
        eslint: {
          lintCommand: 'eslint "./src/**/*.{ts,tsx}"',
        },
        overlay: false,
        typescript: true,
      }),
      react(),
      svgr(),
      tsconfigPaths(),
    ],
    server: command === "serve" ? { proxy: devApiProxy(base, env) } : undefined,
  };
});
