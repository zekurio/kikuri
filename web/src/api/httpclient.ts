import { HTTP_ENDPOINT } from "./static";

export type HttpMethod = 'GET' | 'PUT' | 'POST' | 'DELETE' | 'PATCH' | 'OPTIONS';

export type HttpHeadersMap = { [key: string]: string };

export class HTTPClient {

  async req<T>(
    method: HttpMethod,
    path: string,
    body?: object,
    headers: HttpHeadersMap = {},
  ): Promise<T> {
    const res = await fetch(HTTP_ENDPOINT + path, {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...headers,
      },
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      throw new Error(res.statusText);
    }

    return res.json();
  }

}