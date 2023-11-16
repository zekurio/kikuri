import React from "react";
import styled from "styled-components";
import Color from "color";

import { IconLogin } from "@tabler/icons-react";
import { loginRoute } from "../services/api";

import { Button } from "../components/Button";

class Props {}

const StartContainer = styled.div``;

const LoginButton = styled(Button)`
  position: fixed;
  top: 1.5em;
  right: 1.5em;

  width: 3em;
  height: 3em;
  padding: 0 0.6em;
  display: flex;
  justify-content: flex-start;
  gap: 1em;
  overflow: hidden;
  background: ${(p) => p.theme.background3};
  opacity: 0.5;

  transition: all 0.25s ease;
  transform: none !important;

  > svg {
    min-height: 2em;
    min-width: 2em;
  }

  &:hover {
    width: 8em;
    background: ${(p) => p.theme.accent};
    opacity: 1;
    color: ${(p) => p.theme.white};
  }
`;

const Header = styled.header`
  display: flex;
  flex-direction: column;
  gap: 3em;
  align-items: center;
  padding-top: 10vh;

  > span {
    font-family: "Cantarell", serif;
    font-size: 1.1rem;
    font-weight: lighter;
    text-align: center;
    max-width: 20em;
    opacity: 0.9;
  }
`;

const HeaderButtons = styled.div`
  display: flex;
  gap: 2em;

  ${Button} {
    transition: all 0.25s ease;
    padding: 0.8em 2em;
    box-shadow: 0 0 2em 0 ${(p) => Color(p.theme.accent).alpha(0.2).hexa()};
    &:hover {
      box-shadow: 0 0 2em 0 ${(p) => Color(p.theme.accent).alpha(0.4).hexa()};
    }
  }
`;

const Footer = styled.footer`
  position: fixed;
  bottom: 0;
  width: 100%;
  display: flex;
  gap: 5em;
  padding: 2em;
  justify-content: center;
  color: ${(p) => p.theme.text};
  background-color: ${(p) => Color(p.theme.background2).alpha(0.5).hexa()};
  backdrop-filter: blur(5em);

  a {
    color: inherit;
    text-decoration: underline;
  }

  > div {
    > span,
    a {
      display: block;
      line-height: 1.8rem;
    }
  }
`;

export const StartRoute: React.FC<Props> = () => {
  const _loginRoute = loginRoute("/dashboard");

  return (
    <StartContainer>
      <LoginButton onClick={() => (window.location.href = _loginRoute)}>
        <IconLogin />
        Login
      </LoginButton>
      <Header>
        <HeaderButtons>
          <a href="/invite">
            <Button>Invite Kikuri</Button>
          </a>
        </HeaderButtons>
      </Header>
      <main></main>
      <Footer>
        <div>
          <span>kikuri</span>
          <span>© {new Date().getFullYear()} Michael Schwieger</span>
          <a
            href="https://github.com/zekurio/kikuri/blob/main/LICENSE"
            target="_blank"
            rel="noreferrer"
          >
            Covered by the MIT Licence.
          </a>
          <a
            href="https://github.com/zekurio/kikuri"
            target="_blank"
            rel="noreferrer"
          >
            GitHub Repository
          </a>
        </div>
        <div>
          <a href="https://kikuri.xyz/invite" target="_blank" rel="noreferrer">
            Invite Kikuri Stable
          </a>
          <a
            href="https://canary.kikuri.xyz/invite"
            target="_blank"
            rel="noreferrer"
          >
            Invite Kikuri Canary
          </a>
        </div>
      </Footer>
    </StartContainer>
  );
};
