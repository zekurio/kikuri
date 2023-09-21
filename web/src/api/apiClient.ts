import { HttpClient } from './httpClient';
import { Event } from './models';

export class APIClient extends HttpClient {
  private _onWsEvent: (e: Event<any>) => void = () => {};

  constructor() {
    super();
    this.connectWS(e => {
      this._onWsEvent(e);
    });
  }

  set onWsEvent(handler: (e: Event<any>) => void) {
    this._onWsEvent = handler;
  }

  loginCapabilities(): Promise<string[]> {
    return this.req('GET', 'auth/logincapabilities');
  }

  loginUrl(): string {
    return this.basePath(`auth/oauth2/discord/login`);
  }

  logoutUrl(): string {
    return this.basePath('auth/logout');
  }
}
