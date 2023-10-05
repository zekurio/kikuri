export interface ListResponse<T> {
  n: number;
  data: T[];
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

export interface PermissionResponse {
  permissions: string[];
}

export interface GuildSettings {
  prefix: string;
  perms: Map<string, string[]>;
  autoroles: string[];
  modlogchannel: string;
  modnotchannel: string;
  voicelogchannel: string;
  joinmessagechannel: string;
  joinmessagetext: string;
  leavemessagechannel: string;
  leavemessagetext: string;
}

export interface PermissionsUpdate {
  perm: string;
  role_ids: string[];
  override?: boolean;
}

export type Status = 'online' | 'dnd' | 'idle' | 'invisible';

export interface Presence {
  game: string;
  status: Status;
}

export interface InviteSettingsRequest {
  guild_id: string;
  message: string;
  invite_code?: string;
}

export interface InviteSettingsResponse {
  guild: Guild;
  invite_url: string;
  message: string;
}

export interface Count {
  count: number;
}

export interface SystemInfo {
  version: string;
  commit_hash: string;
  build_date: Date;
  go_version: string;
  uptime: number;
  uptime_str: string;
  os: string;
  arch: string;
  cpus: number;
  go_routines: number;
  stack_use: number;
  stack_use_str: string;
  heap_use: number;
  heap_use_str: string;
  bot_user_id: string;
  bot_invite: string;
  guilds: number;
}

export interface Contact {
  title: string;
  value: string;
  url?: string;
}

export interface PrivacyInfo {
  noticeurl: string;
  contact: Contact[];
}

export interface APIToken {
  created: Date;
  expires: Date;
  last_access: Date;
  hits: number;
  token?: string;
}

export interface SubPermission {
  term: string;
  explicit: boolean;
  description: string;
}

export enum CommandOptionType {
  SUBCOMMAND = 1,
  SUBCOMMANDGROUP = 2,
  STRING = 3,
  INTEGER = 4,
  BOOLEAN = 5,
  USER = 6,
  CHANNEL = 7,
  ROLE = 8,
  MENTIONABLE = 9,
}

export interface CommandOptionChoise {
  name: string;
  value: string;
}

export interface CommandOption {
  type: CommandOptionType;
  name: string;
  description: string;
  required: boolean;
  choices: CommandOptionChoise[];
  options: CommandOption[];
}

export interface CommandInfo {
  name: string;
  description: string;
  version: string;
  perms: string;
  dm_capable: boolean;
  subperms: SubPermission[];
  options: CommandOption[];
  group: string;
}

export interface LandingPageInfo {
  localinvite: string;
  publicmaininvite: string;
  publiccaranyinvite: string;
}

export interface UserSettingsOTA {
  enabled: boolean;
}

export interface AccessTokenModel {
  token: string;
  expires: string;
}

export interface State {
  state: boolean;
}

export interface GuildSettingsApi {
  enabled: boolean;
  allowed_origins: string;
  token: string;
  reset_token: boolean;
  protected: boolean;
}

export interface MessageEmbedField {
  inline: boolean;
  name: string;
  value: string;
}

export interface MessageEmbedFooter {
  icon_url: string;
  proxy_icon_url: string;
  text: string;
}

export interface MessageEmbedImage {
  url: string;
  proxy_url: string;
  width: number;
  height: number;
}

export interface MessageEmbedThumbnail extends MessageEmbedImage {}

export interface MessageEmbedVideo extends MessageEmbedImage {}

export interface MessageEmbed {
  color: number;
  title: string;
  url: string;
  description: string;
  fields: MessageEmbedField[];
  footer: MessageEmbedFooter;
  image: MessageEmbedImage;
  thumbnail: MessageEmbedThumbnail;
  video: MessageEmbedVideo;

  color_hex: string;
}

export interface ChannelWithPermissions extends Channel {
  can_read: boolean;
  can_write: boolean;
}

export interface VerificationSiteKey {
  sitekey: string;
}

export interface GuildSettingsVerification {
  enabled: boolean;
}

export type PermissionsMap = { [key: string]: string[] };