import { HttpClient } from "./httpclient";
import { AuthClient, GuildsClient } from "./bindings";

export class Client extends HttpClient {
  auth = new AuthClient(this);
  guilds = new GuildsClient(this);

  constructor(endpoint: string = "/api") {
    super(endpoint);
  }

  public get clientEndpoint(): string {
    return this.endpoint;
  }
}
