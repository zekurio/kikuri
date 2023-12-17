import styled from "styled-components";
import { BottomContainer, SelfContainer } from "./BottomContainer";
import { SIDEBAR_WIDTH, Sidebar } from "./Sidebar";
import { useNavigate, useParams } from "react-router";
import { Entry } from "./Entry";
import { GuildSelect } from "../GuildSelect/GuildSelect";
import { useGuilds } from "../../hooks/useGuilds";
import { useStore } from "../../services/store";
import { MAX_WIDTH } from "../MaxWidthContainer";
import { Guild } from "../../lib/kikuri-api/src";
import { useEffect } from "react";

class Props {}

const StyledEntry = styled(Entry)``;

const StyledGuildSelect = styled(GuildSelect)`
  margin-top: 1em;
`;

const StyledSidebar = styled(Sidebar)`
  @media (orientation: portrait) and (max-width: calc(${SIDEBAR_WIDTH} + ${MAX_WIDTH})) {
    ${StyledEntry}, ${SelfContainer} {
      justify-content: center;
      span {
        display: none;
      }
    }

    ${SelfContainer} > svg {
      display: none;
    }

    ${StyledGuildSelect} > div > div {
      justify-content: center;
      > span {
        display: none;
      }
    }
  }
`;

export const SidebarDashboard: React.FC<Props> = () => {
  const nav = useNavigate();
  const { guildid } = useParams();
  const guilds = useGuilds();
  const [selectedGuild, setSelectedGuild] = useStore((s) => [
    s.selectedGuild,
    s.setSelectedGuild,
  ]);

  useEffect(() => {
    if (!!guilds && !!guildid)
      setSelectedGuild(guilds.find((g) => g.id === guildid) ?? guilds[0]);
  }, [guildid, guilds]);

  const _onGuildSelect = (g: Guild) => {
    setSelectedGuild(g);
    nav(`guilds/${g.id}/members`);
  };

  return (
    <StyledSidebar>
      <StyledGuildSelect
        guilds={guilds ?? []}
        value={selectedGuild}
        onElementSelect={_onGuildSelect}
      />
      <BottomContainer />
    </StyledSidebar>
  );
};
