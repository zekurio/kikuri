import { useCallback, useState } from 'react';
import { useNavigate, useParams } from 'react-router';

import { Button } from '../../../components/Button';
import { Loader } from '../../../components/Loader';
import { Member } from '../../../lib/kikuri-ts/src';
import { MemberTileLarge } from '../../../components/MemberTileLarge';
import { MemberTile } from '../../../components/MemberTile';
import { SearchBar } from '../../../components/SearchBar';
import debounce from 'debounce';
import styled from 'styled-components';
import { useGuild } from '../../../hooks/useGuild';
import { useMembers } from '../../../hooks/useMembers';
import { useSelfMember } from '../../../hooks/useSelfMember';

type Props = {};

const MembersSection = styled.div`
  margin-top: 1em;
`;

const MemberTiles = styled.div`
  margin-top: 1em;
  display: flex;
  flex-wrap: wrap;
  gap: 1em;
`;

const LoadMoreButton = styled(Button)`
  margin: 1em auto 0 auto;
`;

const GuildMembersRoute: React.FC<Props> = () => {
  const { guildid } = useParams();
  const nav = useNavigate();
  const selfMember = useSelfMember(guildid);
  const guild = useGuild(guildid);
  const [search, setSearch] = useState('');
  const [members, loadMoreMembers] = useMembers(guildid, 100, search);

  const _onSearchInput = useCallback(debounce(setSearch, 500), []);

  const _navToMember = (member: Member) => {
    nav(member.user.id);
  };

  return (
    <>
      {(selfMember && <MemberTileLarge member={selfMember} guild={guild} onClick={_navToMember} />) || (
        <Loader width="100%" height="6em" />
      )}
      {(members && selfMember && (
        <MembersSection>
          <SearchBar onValueChange={_onSearchInput} placeholder="Search..." />
          <MemberTiles>
            {members
              .filter((m) => m.user.id !== selfMember.user.id)
              .map((m) => (
                <MemberTile key={`memb-${m.user.id}`} member={m} onClick={_navToMember} />
              ))}
          </MemberTiles>
          {members.length > 0 && !search && guild?.member_count! > members.length && (
            <LoadMoreButton onClick={() => loadMoreMembers()}>Laod more ...</LoadMoreButton>
          )}
        </MembersSection>
      )) || (
        <>
          <Loader width="100%" height="2em" margin="1em 0 0 0" />
          <Loader width="100%" height="6em" margin="1em 0 0 0" />
        </>
      )}
    </>
  );
};

export default GuildMembersRoute;