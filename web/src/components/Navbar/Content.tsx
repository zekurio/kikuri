import styled from "styled-components";
import React, { PropsWithChildren } from "react";
import { PropsWithStyle } from "../props";

type Props = PropsWithChildren & PropsWithStyle;

const Content = styled.div`
  display: flex;
  flex-grow: 1;
`;

export const NavContent: React.FC<Props> = ({ children, ...props }) => {
  return <Content {...props}>{children}</Content>;
};
