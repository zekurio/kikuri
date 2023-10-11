import { OthersClient } from './bindings';

import { HttpClient } from './httpclient';

export class Client extends HttpClient {

  other = new OthersClient(this);

  constructor(endpoint: string = '/api') {
    super(endpoint);
  }

  public get clientEndpoint(): string {
    return this.endpoint;
  }
}
