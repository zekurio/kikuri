import React from "react";
import styled from "styled-components";
import { IconLogout } from "@tabler/icons-react"; // Import the logout icon

type DropdownMenuProps = {
  show: boolean;
};

const StyledDropdownMenu = styled.ul<DropdownMenuProps>`
  display: ${(props) => (props.show ? "block" : "none")};
  position: absolute;
  list-style: none;
  // Add more styles as needed
`;

type DropdownMenuComponentProps = {
  show: boolean;
  children: React.ReactNode;
};

export const DropdownMenu: React.FC<DropdownMenuComponentProps> = ({
  show,
  children,
}) => {
  return <StyledDropdownMenu show={show}>{children}</StyledDropdownMenu>;
};
