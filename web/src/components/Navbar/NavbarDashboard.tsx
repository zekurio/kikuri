import styled from "styled-components";
import {Navbar} from "./Navbar";
import {NavContent} from "./Content";
import React from "react";
import {useSelfUser} from "../../hooks/useSelfUser";
import { IconBrandDiscordFilled} from "@tabler/icons-react";

const StyledNav = styled(Navbar)``;

const StyledContent = styled(NavContent)`
  display: flex;
  justify-content: flex-end;
  align-items: center;
    margin: 0 10px 0 0;
`;

const UserBox = styled.div`
    display: flex;
  align-items: center;
    
    > a {
        display: flex;
        align-items: center;
        gap: 10px;
        color: ${(p) => p.theme.text};
        text-decoration: none;
    }
    
    > a > img {
        height: 30px;
        width: 30px;
        border-radius: 50%;
    }
`;

export const NavbarDashboard: React.FC = () => {
  const self = useSelfUser();
  const avatarUrl = self?.avatar_url || IconBrandDiscordFilled;

  return (
    <StyledNav>
      <StyledContent>
        <UserBox>
          <a href="">
            <img src={avatarUrl} alt="user icon"/>
            {self?.username}
          </a>
        </UserBox>
      </StyledContent>
    </StyledNav>
  );
}