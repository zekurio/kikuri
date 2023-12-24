import styled from "styled-components";
import React, {PropsWithChildren} from "react";
import {PropsWithStyle} from "../props";

type Props = PropsWithChildren & PropsWithStyle;

const Content = styled.div`
    display: flex;
    justify-content: space-between;
    margin: 10px 0;
    flex-grow: 1;
`;

export const NavContent: React.FC<Props> = ({ children, ...props }) => {
  return (
    <Content {...props}>
      {children}
    </Content>
  );
};