/** @format */

export interface AccessTokenModel {
  token: string;
  expires: string;
}

export interface GuildSettings {
  autoroles: string[];
  autovoices: string[];
  perms: Map<string, string[]>;
}

export interface CodeResponse {
  code: number;
}

export interface ErrorResponse extends CodeResponse {
  error: string;
}
