/** @format */

export interface ListResponse<T> {
  n: number;
  data: T[];
}

export interface AccessTokenModel {
  token: string;
  expires: string;
}

export interface FlatUser {
  id: string;
  username: string;
  bot: boolean;
  avatar_url: string;
}

export interface User extends FlatUser {
  avatar: string;
  locale: string;
  verified: boolean;
  created_at?: string;
  bot_owner?: boolean;
  captcha_verified: boolean;
}

export interface Role {
  id: string;
  name: string;
  managed: boolean;
  mentionable: boolean;
  hoist: boolean;
  color: number;
  position: number;
  permission: number;

  _deleted: boolean;
}

export interface Member {
  user: User;
  guild_id: string;
  guild_name: string;
  joined_at: string;
  nick: string;
  avatar_url?: string;
  roles: string[];
  created_at?: string;
  dominance?: number;
  karma: number;
  karma_total: number;
  chat_muted: boolean;
}

export enum ChannelType {
  GUILD_TEXT = 0,
  DM = 1,
  GUILD_VOICE = 2,
  GROUP_DM = 3,
  GUILD_CATEGORY = 4,
  GUILD_NEWS = 5,
  GUILD_STORE = 6,
  GUILD_NEWS_THREAD = 10,
  GUILD_PUBLIC_THREAD = 11,
  GUILD_PRIVATE_THREAD = 12,
  GUILD_STAGE_VOICE = 13,
}

export interface Channel {
  id: string;
  guild_id: string;
  name: string;
  topic: string;
  type: ChannelType;
  nsfw: boolean;
  position: number;
  user_limit: number;
  parent_id: string;
}

export interface Guild {
  id: string;
  name: string;
  icon: string;
  icon_url: string;
  region: string;
  owner_id: string;
  joined_at: string;
  member_count: number;

  backups_enabled: boolean;
  latest_backup_entry: Date;
  invite_block_enabled: boolean;

  self_member?: Member;

  roles?: Role[];
  members?: Member[];
  channels?: Channel[];
}

export interface GuildSettings {
  autoroles: string[];
  autovoices: string[];
  perms: Map<string, string[]>;
}

export interface PermissionResponse {
  permissions: string[];
}

export interface SearchResult {
  guilds: Guild[];
  members: Member[];
}

export interface CodeResponse {
  code: number;
}

export interface ErrorResponse extends CodeResponse {
  error: string;
}
