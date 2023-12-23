import React, { useState } from "react";
import { Navbar } from "./Navbar"; // Adjust path as needed
import { Button } from "../Button";
import { useSelfUser } from "../../hooks/useSelfUser";
import { useStore } from "../../services/store";
import styled from "styled-components";
import { loginRoute } from "../../services/api";
import { DropdownMenu } from "../Dropdown/Dropdown";
import { DiscordImage } from "../DiscordImage";
import { DropdownEntry } from "../Dropdown";
import { IconLogout } from "@tabler/icons-react";
import { useNavigate } from "react-router";
import { useApi } from "../../hooks/useApi";

const Logo = styled.div`
  // Styling for the logo
`;

const Dropdown = styled.div`
  cursor: pointer;
  // Add styling for dropdown triggers
`;

const UserSection = styled.div`
  display: flex;
  align-items: center;
  > img {
    width: 2em;
    height: 2em;
    margin-right: 0.5em;
  }
  // Additional styling for user section
`;

export const NavbarDashboard: React.FC = () => {
  const self = useSelfUser();
  const selectedGuild = useStore((s) => s.selectedGuild);
  const [showDropdown, setShowDropdown] = useState(false);
  const nav = useNavigate();
  const fetch = useApi();

  const toggleDropdown = () => {
    setShowDropdown(!showDropdown);
  };

  const _logout = () => {
    fetch((c) => c.auth.logout())
      .then(() => nav("/"))
      .catch();
  };

  return (
    <Navbar
      rightContent={
        <>
          {self ? (
            <>
              <Dropdown>
                {selectedGuild ? (
                  <span>{selectedGuild.name}</span>
                ) : (
                  <span>No Guild Selected</span>
                )}
                {/* Guild Dropdown Logic Here */}
              </Dropdown>
              <UserSection onClick={toggleDropdown}>
                {/* ... existing code */}
                {showDropdown && (
                  <DropdownMenu show={showDropdown}>
                    <DropdownEntry>
                      <IconLogout />
                      <Button
                        onClick={() => {
                          _logout();
                        }}
                      >
                        Logout
                      </Button>
                    </DropdownEntry>
                  </DropdownMenu>
                )}
              </UserSection>
            </>
          ) : (
            <a href={loginRoute("/dashboard")}>Login</a>
          )}
        </>
      }
    >
      <Logo>Bot Logo</Logo>
    </Navbar>
  );
};
