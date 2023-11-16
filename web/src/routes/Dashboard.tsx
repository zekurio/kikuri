import { Outlet, useLocation, useNavigate, useParams } from "react-router";
import styled from "styled-components";

type Props = {};

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

export const DashboardRoute: React.FC<Props> = () => {};
