export type Env = {
  api: {
    url: string;
  };
  baseUrl?: string;
  buildTag?: string;
  displayMethodBuilderUrl: string;
  ipfsGatewayUrl: string;
  issuer: {
    logo: string;
    name: string;
  };
  paymentPagesEnabled: boolean;
  schemaExplorerAndBuilderUrl?: string;
  warningMessage?: string;
};
