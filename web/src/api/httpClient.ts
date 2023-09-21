import { APIError } from './errors';
import { AccessToken, Status, Event } from './models';
import { HTTP_ENDPOINT } from './static';
import { WSClient } from './wsClient';

export type HttpMethod =
  | 'GET'
  | 'PUT'
  | 'POST'
  | 'DELETE'
  | 'PATCH'
  | 'OPTIONS';

export type HttpHeadersMap = { [key: string]: string };

export class HttpClient {
  private _accessToken: AccessToken | undefined;
  private accessTokenRequest: Promise<AccessToken> | undefined;
  private ws: WSClient | undefined;

  async req<T>(
    method: HttpMethod,
    path: string,
    body?: object,
    headers: HttpHeadersMap = {},
  ): Promise<T> {
    const _headers = new Headers();
    _headers.set('Accept', 'application/json');
    Object.keys(headers).forEach(k => _headers.set(k, headers[k]));

    if (this._accessToken) {
      if (this.isAccessTokenExpired()) {
        this._accessToken = undefined;
        return await this.getAndSetAccessToken(() =>
          this.req(method, path, body, headers),
        );
      }
      _headers.set('Authorization', `bearer ${this._accessToken.access_token}`);
    }

    let _body = null;
    if (!!body) {
      if (body instanceof File) {
        const formData = new FormData();
        formData.append('file', body);
        _body = formData;
      } else {
        _headers.set('Content-Type', 'application/json');
        _body = JSON.stringify(body);
      }
    }

    const fullPath = replaceDoublePath(`${HTTP_ENDPOINT}/${path}`);
    const res = await window.fetch(fullPath, {
      method,
      headers: _headers,
      body: _body,
      credentials: 'include',
    });

    if (res.status === 204) {
      return {} as T;
    }

    let data = {};
    try {
      data = await res.json();
    } catch {}

    if (
      res.status === 401 &&
      (data as Status).message === 'invalid access token'
    ) {
      return await this.getAndSetAccessToken(() =>
        this.req(method, path, body, headers),
      );
    }

    if (res.status >= 400) throw new APIError(res, data as Status);

    return data as T;
  }

  async connectWS(onEvent: (e: Event<any>) => void) {
    const authTokenGetter = async (): Promise<string> => {
      if (this.isAccessTokenExpired()) {
        await this.getAndSetAccessToken();
      }
      return this._accessToken?.access_token!;
    };
    this.ws = new WSClient(authTokenGetter, onEvent);
    this.ws.connect();
  }

  protected basePath(path?: string): string {
    return replaceDoublePath(`${HTTP_ENDPOINT}/${path}`);
  }

  protected get accessToken(): string | undefined {
    return this._accessToken?.access_token;
  }

  private isAccessTokenExpired(): boolean {
    return (
      !this._accessToken ||
      Date.now() - this._accessToken.deadlineDate.getTime() > 0
    );
  }

  private async getAccessToken(): Promise<AccessToken> {
    if (!this.accessTokenRequest)
      this.accessTokenRequest = this.req('GET', 'auth/refresh');
    return this.accessTokenRequest;
  }

  private async getAndSetAccessToken<T>(replay?: () => Promise<T>): Promise<T> {
    const token = await this.getAccessToken();
    this.accessTokenRequest = undefined;
    this._accessToken = token as AccessToken;
    this._accessToken.deadlineDate = new Date(token.deadline);
    if (!!replay) return await replay();
    return Promise.resolve({} as T);
  }
}

function replaceDoublePath(url: string): string {
  const split = url.split('://');
  split[split.length - 1] = split[split.length - 1].replace(/\/\//g, '/');
  return split.join('://');
}
