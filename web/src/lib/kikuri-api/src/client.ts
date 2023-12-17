import { AxiosClient } from "./axios";
import {
  AuthClient,
  GuildsClient,
  MiscClient,
  SearchClient,
  UsersClient,
} from "./bindings";

export class Client extends AxiosClient {
  misc = new MiscClient(this);
  auth = new AuthClient(this);
  guilds = new GuildsClient(this);
  users = new UsersClient(this);
  search = new SearchClient(this);

  constructor(endpoint: string = "/api") {
    super(endpoint);
  }
}
