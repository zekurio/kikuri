import React from "react";
import styled from "styled-components";
import {NavbarDashboard} from "../components/Navbar/NavbarDashboard";

type Props = NonNullable<unknown>;

const DashboardContainer = styled.div``;

export const DashboardRoute: React.FC<Props> = () => {
  return (
    <DashboardContainer>
      <NavbarDashboard/>
    </DashboardContainer>
  );
};
