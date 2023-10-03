export const ENDPOINT = import.meta.env.PROD ? window.location.origin : 'http://localhost';

export const HTTP_ENDPOINT = ENDPOINT + '/api/v1';