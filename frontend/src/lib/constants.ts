export const HTTPS_ENABLED = process.env.PUBLIC_APP_URL?.startsWith('https://') ?? false;
export const ACCESS_TOKEN_COOKIE_NAME = HTTPS_ENABLED ? '__Host-access_token' : 'access_token';
