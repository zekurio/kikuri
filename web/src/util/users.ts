import { Member } from "../lib/kikuri-ts/src";

export function memberName(member: Member): string {
  return member.nick ? member.nick : member.user.username;
}
