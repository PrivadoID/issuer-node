import { z } from "zod";

import { getStrictParser } from "src/adapters/parsers";
import { Env } from "src/domain";
import { IMAGE_PLACEHOLDER_PATH } from "src/utils/constants";

export type EnvInput = {
  BASE_URL: string;
  VITE_BASE_URL?: string;
  VITE_BUILD_TAG?: string;
  VITE_DISPLAY_METHOD_BUILDER_URL: string;
  VITE_IPFS_GATEWAY_URL: string;
  VITE_ISSUER_LOGO?: string;
  VITE_ISSUER_NAME: string;
  VITE_PAYMENT_PAGES_ENABLED?: string;
  VITE_SCHEMA_EXPLORER_AND_BUILDER_URL?: string;
  VITE_WARNING_MESSAGE?: string;
};

export const envParser = getStrictParser<EnvInput, Env>()(
  z
    .object({
      BASE_URL: z.string(),
      VITE_BASE_URL: z.string().optional(),
      VITE_BUILD_TAG: z.string().optional(),
      VITE_DISPLAY_METHOD_BUILDER_URL: z.string(),
      VITE_IPFS_GATEWAY_URL: z.string().url(),
      VITE_ISSUER_LOGO: z.string().optional(),
      VITE_ISSUER_NAME: z.string().min(1),
      VITE_PAYMENT_PAGES_ENABLED: z.string().optional(),
      VITE_SCHEMA_EXPLORER_AND_BUILDER_URL: z
        .union([z.string().url(), z.literal("")])
        .transform((value) => value || undefined)
        .optional(),
      VITE_WARNING_MESSAGE: z.string().optional(),
    })
    .transform(
      ({
        BASE_URL,
        VITE_BASE_URL,
        VITE_BUILD_TAG,
        VITE_DISPLAY_METHOD_BUILDER_URL,
        VITE_IPFS_GATEWAY_URL,
        VITE_ISSUER_LOGO,
        VITE_ISSUER_NAME,
        VITE_PAYMENT_PAGES_ENABLED,
        VITE_SCHEMA_EXPLORER_AND_BUILDER_URL,
        VITE_WARNING_MESSAGE,
      }): Env => ({
        api: {
          // The UI's own server (nginx, or the Vite dev server) proxies /api to the issuer API and adds the API
          // credentials there, so they never reach the browser.
          url: `${BASE_URL.replace(/\/+$/, "")}/api`,
        },
        baseUrl: VITE_BASE_URL,
        buildTag: VITE_BUILD_TAG,
        displayMethodBuilderUrl: VITE_DISPLAY_METHOD_BUILDER_URL,
        ipfsGatewayUrl: VITE_IPFS_GATEWAY_URL,
        issuer: {
          logo: VITE_ISSUER_LOGO || IMAGE_PLACEHOLDER_PATH,
          name: VITE_ISSUER_NAME,
        },
        paymentPagesEnabled: VITE_PAYMENT_PAGES_ENABLED === "true",
        schemaExplorerAndBuilderUrl: VITE_SCHEMA_EXPLORER_AND_BUILDER_URL,
        warningMessage: VITE_WARNING_MESSAGE,
      })
    )
);
