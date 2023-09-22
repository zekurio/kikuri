export const ENDPOINT = import.meta.env.PROD ? window.location.origin : 'http://localhost:8080';

export const HTTP_ENDPOINT = ENDPOINT + '/api/v1';