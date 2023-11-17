import { useEffect, useState } from 'react';

import { GuildMemberClient } from '../lib/kikuri-ts/src/bindings';
import { Member } from '../lib/kikuri-ts/src';
import { useApi } from './useApi';

type MemberRequester = <T>(req: (c: GuildMemberClient) => Promise<T>) => Promise<T>;

export const useMember = (
  guildid?: string,
  memberid?: string,
): [Member | undefined, MemberRequester] => {
  const fetch = useApi();
  const [member, setMember] = useState<Member>();

  const memberAction: MemberRequester = (req) =>
    fetch((c) => req(c.guilds.member(guildid!, memberid!)));

  useEffect(() => {
    if (!guildid || !memberid) return;
    memberAction((c) => c.get())
      .then((res) => setMember(res))
      .catch();
  }, [guildid, memberid]);

  return [member, memberAction];
};