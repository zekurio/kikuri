import React from "react";
import styled from "styled-components";
import { Navbar } from "./Navbar";
import { Button } from "../Button";
import { useTranslation } from "react-i18next";
import {NavContent} from "./Content";

const StyledNav = styled(Navbar)``;

const NavLinks = styled.div`
  display: flex;
  align-items: center;
  margin-left: 15px;
`;

const NavLink = styled.a`
  color: ${(p) => p.theme.text};
  text-decoration: none;
  margin-right: 20px;
`;

const RightNav = styled.div`
  display: flex;
  align-items: center;
  padding: 0 15px;
`;

const DashboardButton = styled(Button)`
    background: transparent;
    border: 1px solid ${(p) => p.theme.accent};
    margin-right: 10px; // Adds space between the buttons

    &:hover {
        background: ${(p) => p.theme.accent}};
}
`;

const LogoutButton = styled(Button)`
    background: transparent;
    border: 1px solid ${(p) => p.theme.red};

    &:hover {
        background: ${(p) => p.theme.red}};
}
`;

export const NavbarLanding: React.FC = () => {
  const { t } = useTranslation("components.navbar.landing");

  return (
    <StyledNav>
      <NavContent>
        <NavLinks>
          <NavLink href="#Features">Features</NavLink>
          <NavLink href="#About">About</NavLink>
        </NavLinks>
        <RightNav>
          <DashboardButton>{t("Dashboard")}</DashboardButton>
          <LogoutButton>{t("logout")}</LogoutButton>
        </RightNav>
      </NavContent>
    </StyledNav>
  );
};
