import {
  GuildSettings,
  AccessTokenModel,
  CodeResponse,
  Guild,
  ListResponse,
  User,
  Member,
  PermissionResponse,
} from "./models";

import { Client } from "./client";

import { SubClient } from "./subclient";

export class MiscClient extends SubClient {
  constructor(client: Client) {
    super(client, "");
  }

  me(): Promise<User> {
    return this.req("GET", "me");
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

  pushCode(code: string): Promise<CodeResponse> {
    return this.req("POST", "pushcode", { code });
  }
}
