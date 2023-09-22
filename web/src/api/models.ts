export interface SystemInfo {
    version: string;
    commit_hash: string;
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