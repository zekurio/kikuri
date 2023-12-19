import { Outlet, useLocation, useNavigate, useParams } from "react-router";

import LocalStorageUtil from "../util/localstorage";
import { SidebarDashboard } from "../components/Sidebar";
import styled from "styled-components";
import { useEffect } from "react";
import { useGuilds } from "../hooks/useGuilds";

class Props {}

const RouteContainer = styled.div`
  display: flex;
  height: 100%;
`;

const RouterOutlet = styled.main`
  padding: 1em;
  width: 100%;
  height: 100%;
  overflow-y: auto;
`;

export const DashboardRoute: React.FC<Props> = () => {
  const guilds = useGuilds();
  const { guildid } = useParams();
  const loc = useLocation();
  const nav = useNavigate();

  useEffect(() => {
    if (!guilds) {
      /* empty */
    } else if (guilds.length === 0) {
      nav("/welcome");
    } else if (loc.pathname.replaceAll("/", "") === "db" && !guildid) {
      const guild =
        guilds.find(
          (g) => g.id === LocalStorageUtil.get<string>("kikuri.selectedguild"),
        ) ?? guilds[0];
      nav(`guilds/${guild.id}/members`);
    }
  }, [guilds]);

  return (
    <RouteContainer>
      <SidebarDashboard />
      <RouterOutlet>
        <Outlet />
      </RouterOutlet>
    </RouteContainer>
  );
};
