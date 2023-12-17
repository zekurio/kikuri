import { Client } from "./client";
import { SubClient } from "./subclient";
import {
  APIToken,
  AccessTokenModel,
  CodeResponse,
  Guild,
  GuildSettings,
  ListResponse,
  Member,
  PermissionResponse,
  SearchResult,
  User,
} from "./models";

export class MiscClient extends SubClient {
  constructor(client: Client) {
    super(client, "");
  }

  me(): Promise<User> {
    return this.req("GET", "me");
  }
}

export class AuthClient extends SubClient {
  constructor(client: Client) {
    super(client, "auth");
  }

  accesstoken(): Promise<AccessTokenModel> {
    return this.req("POST", "accesstoken");
  }

  check(): Promise<CodeResponse> {
    return this.req("GET", "check");
  }

  logout(): Promise<CodeResponse> {
    return this.req("POST", "logout");
  }
}

export class SearchClient extends SubClient {
  constructor(client: Client) {
    super(client, "search");
  }

  query(query: string, limit: number = 50): Promise<SearchResult> {
    return this.req("GET", `/?query=${query}&limit=${limit}`);
  }
}

export class TokensClient extends SubClient {
  constructor(client: Client) {
    super(client, "token");
  }

  delete(): Promise<CodeResponse> {
    return this.req("DELETE", `/`);
  }

  info(): Promise<APIToken> {
    return this.req("GET", `/`);
  }

  generate(): Promise<APIToken> {
    return this.req("POST", `/`);
  }
}

export class GuildsClient extends SubClient {
  constructor(private _client: Client) {
    super(_client, "guilds");
  }

  list(): Promise<ListResponse<Guild>> {
    return this.req("GET", "/");
  }

  guild(id: string): Promise<Guild> {
    return this.req("GET", id);
  }

  settings(id: string): GuildSettingsClient {
    return new GuildSettingsClient(this._client, id);
  }

  members(
    id: string,
    limit = 50,
    after = "",
    filter = "",
  ): Promise<ListResponse<Member>> {
    return this.req(
      "GET",
      `${id}/members?limit=${limit}&after=${after}&filter=${filter}`,
    );
  }

  member(id: string, memberID: string): GuildMemberClient {
    return new GuildMemberClient(this._client, id, memberID);
  }
}

export class UsersClient extends SubClient {
  constructor(client: Client) {
    super(client, "users");
  }

  get(id: string): Promise<User> {
    return this.req("GET", id);
  }
}

export class GuildSettingsClient extends SubClient {
  constructor(client: Client, id: string) {
    super(client, `guilds/${id}/settings`);
  }

  settings(): Promise<GuildSettings> {
    return this.req("GET", "/");
  }

  setSettings(settings: GuildSettings): Promise<GuildSettings> {
    return this.req("POST", "/", settings);
  }
}

export class GuildMemberClient extends SubClient {
  constructor(client: Client, guildID: string, memberID: string) {
    super(client, `guilds/${guildID}/${memberID}`);
  }

  get(): Promise<Member> {
    return this.req("GET", "/");
  }

  permissions(): Promise<PermissionResponse> {
    return this.req("GET", "permissions");
  }

  permissionsAllowed(): Promise<ListResponse<string>> {
    return this.req("GET", "permissions/allowed");
  }
}
