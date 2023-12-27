import React from "react";
import styled from "styled-components";
import { Navbar } from "./Navbar";
import { Button } from "../Button";
import { useTranslation } from "react-i18next";
import { NavContent } from "./Content";
import { useSelfUser } from "../../hooks/useSelfUser";
import {useNavigate} from "react-router";
import {useApi} from "../../hooks/useApi";

const StyledNav = styled(Navbar)``;

const NavLinks = styled.div`
  display: flex;
  align-items: center;
`;

const NavLink = styled.a`
  color: ${(p) => p.theme.text};
  text-decoration: none;
    margin: 0 5px;
`;

const RightNav = styled.div`
  display: flex;
  align-items: center;
`;

const DashboardButton = styled(Button)`
    background: transparent;
    border: 1px solid ${(p) => p.theme.accent};
    margin: 0 10px 0 10px;
    &:hover {
        background: ${(p) => p.theme.accent}};
}
`;

const LogoutButton = styled(Button)`
    background: transparent;
    border: 1px solid ${(p) => p.theme.red};
    margin: 0 10px 0 0;
    &:hover {
        background: ${(p) => p.theme.red}};
}
`;

export const NavbarLanding: React.FC = () => {
  const self = useSelfUser();
  const fetch = useApi();
  const { t } = useTranslation("components.navbar.landing");
  const nav = useNavigate();

  const _logout = () => {
    fetch((c) => c.auth.logout())
      .then(() => nav('/'))
      .catch();
  };

  return (
    <StyledNav>
      <NavContent>
        <NavLinks>
          <NavLink href="#Features">{t("links.features")}</NavLink>
          <NavLink href="#About">{t("links.about")}</NavLink>
          <NavLink href="/status">{t("links.status")}</NavLink>
        </NavLinks>
        <RightNav>
          <DashboardButton onClick={() => nav('/dashboard')}>{t("dashboard")}</DashboardButton>
          {self && <LogoutButton onClick={_logout}>{t("logout")}</LogoutButton>}
        </RightNav>
      </NavContent>
    </StyledNav>
  );
};
