import { IconRobot } from "@tabler/icons-react";
import { Clickable } from "../styleParts";
import { Container } from "../Container";
import { DiscordImage } from "../DiscordImage";
import { Member } from "../../lib/kikuri-api/src";
import { memberName } from "../../util/users";
import styled from "styled-components";

type Props = {
  member: Member;
  onClick?: (member: Member) => void;
};

const StyledContainer = styled(Container)`
  ${Clickable()}

  display: flex;
  padding: 0.5em;

  > img {
    width: 3em;
    height: 3em;
  }
`;

const Details = styled.div`
  margin-left: 0.5em;

  > h4 {
    margin: 0 0 0.5em 0;
    font-weight: 600;
  }

  > span {
    font-size: 0.8rem;
  }
`;

const StyledBotIcon = styled(IconRobot)`
  color: ${(p) => p.theme.blurple};
  stroke-width: 2;
`;

export const MemberTile: React.FC<Props> = ({ member, onClick = () => {} }) => {
  return (
    <StyledContainer onClick={() => onClick(member)}>
      <DiscordImage src={member.avatar_url} />
      <Details>
        <h4>
          {memberName(member)} {member.user.bot && <StyledBotIcon />}
        </h4>
        <span>{member.user.username}</span>
      </Details>
    </StyledContainer>
  );
};
