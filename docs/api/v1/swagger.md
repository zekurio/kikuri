# Kikuri API
Kikuri API.

## Version: 1.0

---
## Authorization

### /auth/accesstoken

#### POST
##### Summary

Access Token Exchange

##### Description

Exchanges a refresh token for an access token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.AccessTokenResponse](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsaccesstokenresponse) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

### /auth/check

#### GET
##### Summary

Check

##### Description

Check if provided authorization token is valid.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Status](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsstatus) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

### /auth/logout

#### POST
##### Summary

Logout

##### Description

Reovkes the currently used access token and clears the refresh token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Status](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsstatus) |

---
## Guilds

### /guilds

#### GET
##### Summary

Get Guilds

##### Description

Returns all guilds the bot and the user have in common.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ [github_com_zekurio_kikuri_internal_services_webserver_v1_models.GuildReduced](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsguildreduced) ] |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

### /guilds/{id}

#### GET
##### Summary

Get Guild

##### Description

Returns a single guild object by its ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Guild](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsguild) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

---
## default

### /guilds/{id}/members

#### GET
##### Summary

Get Guild Member List

##### Description

Returns a list of guild members.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| after | query | Request members after the given member ID. | No | string |
| limit | query | The amount of results returned. | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wraped in models.ListResponse | [ [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Member](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsmember) ] |
| 400 | Bad Request | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}

#### GET
##### Summary

Get Guild Member

##### Description

Returns a single guild member by ID.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Member](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsmember) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/permissions

#### GET
##### Summary

Get Guild Member Permissions

##### Description

Returns the permission array of the given user.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.PermissionsResponse](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelspermissionsresponse) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

### /guilds/{id}/{memberid}/permissions/allowed

#### GET
##### Summary

Get Guild Member Allowed Permissions

##### Description

Returns all detailed permission DNS which the member is alloed to perform.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| memberid | path | The ID of the member. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Wrapped in models.ListResponse | [ string ] |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

---
## Guild Settings

### /guilds/{id}/settings

#### GET
##### Summary

Get Guild Settings

##### Description

Returns the specified general guild settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.GuildSettings](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsguildsettings) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

Get Guild Settings

##### Description

Returns the specified general guild settings.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| id | path | The ID of the guild. | Yes | string |
| payload | body | Modified guild settings payload. | Yes | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.GuildSettings](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsguildsettings) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Status](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsstatus) |
| 400 | Bad Request | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 403 | Forbidden | Object |
| 404 | Not Found | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

---
## Misc

### /me

#### GET
##### Summary

Me

##### Description

Returns the user object of the currently authenticated user.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.User](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsuser) |

---
## default

### /search

#### GET
##### Summary

Global Search

##### Description

Search through guilds and members by ID, name or displayname.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| query | query | The search query (either ID, name or displayname). | Yes | string |
| limit | query | The maximum amount of result items (per group). | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.SearchResult](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelssearchresult) |
| 400 | Bad Request | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

---
## default

### /token

#### GET
##### Summary

API Token Info

##### Description

Returns API Token metadata, not the token itself.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.APITokenResponse](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsapitokenresponse) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |
| 404 | Is returned when no token was generated before. | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

#### POST
##### Summary

API Token Generation

##### Description

Generates an API Token and returns it and its metadata.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.APITokenResponse](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsapitokenresponse) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

#### DELETE
##### Summary

API Token Deletion

##### Description

Deletes the users API token.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Status](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsstatus) |
| 401 | Unauthorized | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelserror) |

---
## Users

### /users/{id}

#### GET
##### Summary

User

##### Description

Returns a user by their id

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.User](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsuser) |

---
### Models

#### discordgo.Channel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| application_id | string | ApplicationID of the DM creator Zeroed if guild channel or not a bot user | No |
| applied_tags | [ string ] | The IDs of the set of tags that have been applied to a thread in a forum channel. | No |
| available_tags | [ [discordgo.ForumTag](#discordgoforumtag) ] | The set of tags that can be used in a forum channel. | No |
| bitrate | integer | The bitrate of the channel, if it is a voice channel. | No |
| default_forum_layout | [discordgo.ForumLayout](#discordgoforumlayout) | The default forum layout view used to display posts in forum channels. Defaults to ForumLayoutNotSet, which indicates a layout view has not been set by a channel admin. | No |
| default_reaction_emoji | [discordgo.ForumDefaultReaction](#discordgoforumdefaultreaction) | Emoji to use as the default reaction to a forum post. | No |
| default_sort_order | [discordgo.ForumSortOrderType](#discordgoforumsortordertype) | The default sort order type used to order posts in forum channels. Defaults to null, which indicates a preferred sort order hasn't been set by a channel admin. | No |
| default_thread_rate_limit_per_user | integer | The initial RateLimitPerUser to set on newly created threads in a channel. This field is copied to the thread at creation time and does not live update. | No |
| flags | [discordgo.ChannelFlags](#discordgochannelflags) | Channel flags. | No |
| guild_id | string | The ID of the guild to which the channel belongs, if it is in a guild. Else, this ID is empty (e.g. DM channels). | No |
| icon | string | Icon of the group DM channel. | No |
| id | string | The ID of the channel. | No |
| last_message_id | string | The ID of the last message sent in the channel. This is not guaranteed to be an ID of a valid message. | No |
| last_pin_timestamp | string | The timestamp of the last pinned message in the channel. nil if the channel has no pinned messages. | No |
| member_count | integer | An approximate count of users in a thread, stops counting at 50 | No |
| message_count | integer | An approximate count of messages in a thread, stops counting at 50 | No |
| name | string | The name of the channel. | No |
| nsfw | boolean | Whether the channel is marked as NSFW. | No |
| owner_id | string | ID of the creator of the group DM or thread | No |
| parent_id | string | The ID of the parent channel, if the channel is under a category. For threads - id of the channel thread was created in. | No |
| permission_overwrites | [ [discordgo.PermissionOverwrite](#discordgopermissionoverwrite) ] | A list of permission overwrites present for the channel. | No |
| position | integer | The position of the channel, used for sorting in client. | No |
| rate_limit_per_user | integer | Amount of seconds a user has to wait before sending another message or creating another thread (0-21600) bots, as well as users with the permission manage_messages or manage_channel, are unaffected | No |
| recipients | [ [discordgo.User](#discordgouser) ] | The recipients of the channel. This is only populated in DM channels. | No |
| thread_member | [discordgo.ThreadMember](#discordgothreadmember) | Thread member object for the current user, if they have joined the thread, only included on certain API endpoints | No |
| thread_metadata | [discordgo.ThreadMetadata](#discordgothreadmetadata) | Thread-specific fields not needed by other channels | No |
| topic | string | The topic of the channel. | No |
| type | [discordgo.ChannelType](#discordgochanneltype) | The type of the channel. | No |
| user_limit | integer | The user limit of the voice channel. | No |

#### discordgo.ChannelFlags

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ChannelFlags | integer |  |  |

#### discordgo.ChannelType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ChannelType | integer |  |  |

#### discordgo.ForumDefaultReaction

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| emoji_id | string | The id of a guild's custom emoji. | No |
| emoji_name | string | The unicode character of the emoji. | No |

#### discordgo.ForumLayout

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ForumLayout | integer |  |  |

#### discordgo.ForumSortOrderType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.ForumSortOrderType | integer |  |  |

#### discordgo.ForumTag

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| emoji_id | string |  | No |
| emoji_name | string |  | No |
| id | string |  | No |
| moderated | boolean |  | No |
| name | string |  | No |

#### discordgo.MfaLevel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.MfaLevel | integer |  |  |

#### discordgo.PermissionOverwrite

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| allow | string | *Example:* `"0"` | No |
| deny | string | *Example:* `"0"` | No |
| id | string |  | No |
| type | [discordgo.PermissionOverwriteType](#discordgopermissionoverwritetype) |  | No |

#### discordgo.PermissionOverwriteType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.PermissionOverwriteType | integer |  |  |

#### discordgo.PremiumTier

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.PremiumTier | integer |  |  |

#### discordgo.Role

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| color | integer | The hex color of this role. | No |
| hoist | boolean | Whether this role is hoisted (shows up separately in member list). | No |
| id | string | The ID of the role. | No |
| managed | boolean | Whether this role is managed by an integration, and thus cannot be manually added to, or taken from, members. | No |
| mentionable | boolean | Whether this role is mentionable. | No |
| name | string | The name of the role. | No |
| permissions | string | The permissions of the role on the guild (doesn't include channel overrides). This is a combination of bit masks; the presence of a certain permission can be checked by performing a bitwise AND between this int and the permission.<br>*Example:* `"0"` | No |
| position | integer | The position of this role in the guild's role hierarchy. | No |

#### discordgo.ThreadMember

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| flags | integer | Any user-thread settings, currently only used for notifications | No |
| id | string | The id of the thread | No |
| join_timestamp | string | The time the current user last joined the thread | No |
| user_id | string | The id of the user | No |

#### discordgo.ThreadMetadata

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| archive_timestamp | string | Timestamp when the thread's archive status was last changed, used for calculating recent activity | No |
| archived | boolean | Whether the thread is archived | No |
| auto_archive_duration | integer | Duration in minutes to automatically archive the thread after recent activity, can be set to: 60, 1440, 4320, 10080 | No |
| invitable | boolean | Whether non-moderators can add other non-moderators to a thread; only available on private threads | No |
| locked | boolean | Whether the thread is locked; when a thread is locked, only users with MANAGE_THREADS can unarchive it | No |

#### discordgo.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accent_color | integer | User's banner color, encoded as an integer representation of hexadecimal color code | No |
| avatar | string | The hash of the user's avatar. Use Session.UserAvatar to retrieve the avatar itself. | No |
| banner | string | The hash of the user's banner image. | No |
| bot | boolean | Whether the user is a bot. | No |
| discriminator | string | The discriminator of the user (4 numbers after name). | No |
| email | string | The email of the user. This is only present when the application possesses the email scope for the user. | No |
| flags | integer | The flags on a user's account. Only available when the request is authorized via a Bearer token. | No |
| id | string | The ID of the user. | No |
| locale | string | The user's chosen language option. | No |
| mfa_enabled | boolean | Whether the user has multi-factor authentication enabled. | No |
| premium_type | integer | The type of Nitro subscription on a user's account. Only available when the request is authorized via a Bearer token. | No |
| public_flags | [discordgo.UserFlags](#discordgouserflags) | The public flags on a user's account. This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag. | No |
| system | boolean | Whether the user is an Official Discord System user (part of the urgent message system). | No |
| token | string | The token of the user. This is only present for the user represented by the current session. | No |
| username | string | The user's username. | No |
| verified | boolean | Whether the user's email is verified. | No |

#### discordgo.UserFlags

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.UserFlags | integer |  |  |

#### discordgo.VerificationLevel

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| discordgo.VerificationLevel | integer |  |  |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.APITokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created | string |  | No |
| expires | string |  | No |
| hits | integer |  | No |
| last_access | string |  | No |
| token | string |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.AccessTokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| expires | string |  | No |
| token | string |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.Error

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| context | string |  | No |
| error | string |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.Guild

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| afk_channel_id | string |  | No |
| banner | string |  | No |
| channels | [ [discordgo.Channel](#discordgochannel) ] |  | No |
| description | string |  | No |
| icon | string |  | No |
| icon_url | string |  | No |
| id | string |  | No |
| joined_at | string |  | No |
| large | boolean |  | No |
| member_count | integer |  | No |
| mfa_level | [discordgo.MfaLevel](#discordgomfalevel) |  | No |
| name | string |  | No |
| owner_id | string |  | No |
| premium_subscription_count | integer |  | No |
| premium_tier | [discordgo.PremiumTier](#discordgopremiumtier) |  | No |
| region | string |  | No |
| roles | [ [discordgo.Role](#discordgorole) ] |  | No |
| self_member | [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Member](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsmember) |  | No |
| splash | string |  | No |
| unavailable | boolean |  | No |
| verification_level | [discordgo.VerificationLevel](#discordgoverificationlevel) |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.GuildReduced

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| icon | string |  | No |
| icon_url | string |  | No |
| id | string |  | No |
| joined_at | string |  | No |
| member_count | integer |  | No |
| name | string |  | No |
| online_member_count | integer |  | No |
| owner_id | string |  | No |
| region | string |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.GuildSettings

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| auto_roles | [ string ] |  | No |
| auto_voice | [ string ] |  | No |
| perms | object |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.Member

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | The hash of the avatar for the guild member, if any. | No |
| avatar_url | string |  | No |
| communication_disabled_until | string | The time at which the member's timeout will expire. Time in the past or nil if the user is not timed out. | No |
| created_at | string |  | No |
| deaf | boolean | Whether the member is deafened at a guild level. | No |
| guild_id | string | The guild ID on which the member exists. | No |
| guild_name | string |  | No |
| joined_at | string | The time at which the member joined the guild. | No |
| mute | boolean | Whether the member is muted at a guild level. | No |
| nick | string | The nickname of the member, if they have one. | No |
| pending | boolean | Is true while the member hasn't accepted the membership screen. | No |
| permissions | string | Total permissions of the member in the channel, including overrides, returned when in the interaction object.<br>*Example:* `"0"` | No |
| premium_since | string | When the user used their Nitro boost on the server | No |
| privilege | integer |  | No |
| roles | [ string ] | A list of IDs of the roles which are possessed by the member. | No |
| user | [discordgo.User](#discordgouser) | The underlying user on which the member is based. | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.PermissionsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| permissions | [ string ] |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.SearchResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| guilds | [ [github_com_zekurio_kikuri_internal_services_webserver_v1_models.GuildReduced](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsguildreduced) ] |  | No |
| members | [ [github_com_zekurio_kikuri_internal_services_webserver_v1_models.Member](#github_com_zekurio_kikuri_internal_services_webserver_v1_modelsmember) ] |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.Status

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | integer |  | No |
| message | string |  | No |

#### github_com_zekurio_kikuri_internal_services_webserver_v1_models.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accent_color | integer | User's banner color, encoded as an integer representation of hexadecimal color code | No |
| avatar | string | The hash of the user's avatar. Use Session.UserAvatar to retrieve the avatar itself. | No |
| avatar_url | string |  | No |
| banner | string | The hash of the user's banner image. | No |
| bot | boolean | Whether the user is a bot. | No |
| bot_owner | boolean |  | No |
| created_at | string |  | No |
| discriminator | string | The discriminator of the user (4 numbers after name). | No |
| email | string | The email of the user. This is only present when the application possesses the email scope for the user. | No |
| flags | integer | The flags on a user's account. Only available when the request is authorized via a Bearer token. | No |
| id | string | The ID of the user. | No |
| locale | string | The user's chosen language option. | No |
| mfa_enabled | boolean | Whether the user has multi-factor authentication enabled. | No |
| premium_type | integer | The type of Nitro subscription on a user's account. Only available when the request is authorized via a Bearer token. | No |
| public_flags | [discordgo.UserFlags](#discordgouserflags) | The public flags on a user's account. This is a combination of bit masks; the presence of a certain flag can be checked by performing a bitwise AND between this int and the flag. | No |
| system | boolean | Whether the user is an Official Discord System user (part of the urgent message system). | No |
| token | string | The token of the user. This is only present for the user represented by the current session. | No |
| username | string | The user's username. | No |
| verified | boolean | Whether the user's email is verified. | No |
