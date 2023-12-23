import React from "react";
import { NavbarDashboard } from "../components/Navbar";

type Props = NonNullable<unknown>;

export const DashboardRoute: React.FC<Props> = () => {
  return (
    <div>
      <NavbarDashboard />
    </div>
  );
};
