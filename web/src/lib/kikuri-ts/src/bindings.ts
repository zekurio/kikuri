import { GuildSettings, AccessTokenModel, CodeResponse } from "./models";

import { Client } from "./client";

import { SubClient } from "./subclient";

export class GuildsClient extends SubClient {
  constructor(private _client: Client) {
    super(_client, "guilds");
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
