import React from "react";
import styled from "styled-components";

// Styled components
const StyledNavbar = styled.nav`
  display: flex;
  justify-content: space-between;
  padding: 1em;
  background-color: ${(p) => p.theme.background};
`;

const NavSection = styled.div`
  display: flex;
  align-items: center;
  gap: 1em;
`;

export type NavbarProps = {
  children: React.ReactNode;
  rightContent?: React.ReactNode; // additional content on the right
};

export const Navbar: React.FC<NavbarProps> = ({ children, rightContent }) => {
  return (
    <StyledNavbar>
      <NavSection>{children}</NavSection>
      {rightContent && <NavSection>{rightContent}</NavSection>}
    </StyledNavbar>
  );
};
