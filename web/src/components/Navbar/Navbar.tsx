import styled from "styled-components";
import React, { PropsWithChildren } from "react";
import KIIcon from "../../assets/ki-icon.png";
import { PropsWithStyle } from "../props";

type Props = PropsWithChildren & PropsWithStyle;

const Brand = styled.div`
  display: flex;
  table-layout: fixed;
  margin: 10px 5px 0 10px;
  > img {
    width: 40px;
    height: 40px;
  }
`;

const StyledNav = styled.nav`
  display: flex;
  background-color: ${(p) => p.theme.background3};
  color: ${(p) => p.theme.text};
  height: 60px;
  width: 100%;
`;

export const Navbar: React.FC<Props> = ({ children, ...props }) => {
  return (
    <StyledNav {...props}>
      <Brand>
        <img src={KIIcon} alt="kikuri logo" />
      </Brand>
      {children}
    </StyledNav>
  );
};
