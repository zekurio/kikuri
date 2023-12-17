import { Member } from "../lib/kikuri-api/src";

export function memberName(member: Member): string {
  return member.nick ? member.nick : member.user.username;
}
