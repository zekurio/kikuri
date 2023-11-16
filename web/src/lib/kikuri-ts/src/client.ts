import { HttpClient } from "./httpclient";
import { AuthClient, GuildsClient, MiscClient, UsersClient } from "./bindings";

export class Client extends HttpClient {
  auth = new AuthClient(this);
  guilds = new GuildsClient(this);
  misc = new MiscClient(this);
  users = new UsersClient(this);

  constructor(endpoint: string = "/api") {
    super(endpoint);
  }

  public get clientEndpoint(): string {
    return this.endpoint;
  }
}
