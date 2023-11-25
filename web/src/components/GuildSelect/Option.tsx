import { DiscordImage } from "../DiscordImage";
import { Guild } from "../../lib/kikuri-ts/src";
import styled from "styled-components";

type Props = {
  guild: Guild;
};

const StyledDiv = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5em;

  > span {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  > img,
  svg {
    height: 1.2em;
    aspect-ratio: 1;
  }
`;

export const Option: React.FC<Props> = ({ guild }) => {
  return (
    <StyledDiv>
      <DiscordImage src={guild.icon_url} round />
      <span>{guild.name}</span>
    </StyledDiv>
  );
};
