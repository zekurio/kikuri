import { Trans, useTranslation } from "react-i18next";
import styled, { useTheme } from "styled-components";

import { Button } from "../components/Button";
import Color from "color";
import { LinearGradient } from "../components/styleParts";
import { IconLogin } from "@tabler/icons-react";
import KikuriIcon from "../assets/ki-icon.png";

import MockupVotingLight from "../assets/mockups/light/voting.png";
import MockupVotingDark from "../assets/mockups/dark/voting.png";
import { loginRoute } from "../services/api";
import React from "react";

type Props = NonNullable<unknown>;

const StartContainer = styled.div`
  display: flex;
  flex-direction: column;
  min-height: 100vh;
`;

const Header = styled.header`
  display: flex;
  flex-direction: column;
  gap: 3em;
  align-items: center;

  padding-top: 10vh;

  > span {
    font-family: "Nunito Sans", sans-serif;
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

const GlowLink = styled.a`
  ${(p) => LinearGradient(p.theme.accent)}
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  display: inline-block;
  text-decoration: none;
  text-shadow: 0 0 0.8em ${(p) => Color(p.theme.accent).alpha(0.8).hexa()};
`;

const Brand = styled.div`
  display: flex;
  gap: 1em;
  align-items: center;

  width: 80vw;
  height: 15vw;

  max-height: 6rem;
  max-width: 30rem;

  > img {
    height: 100%;
  }

  > svg {
    width: 100%;
    height: 100%;
  }
`;

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
  color: ${(p) => p.theme.text};

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
    color: ${(p) => p.theme.textAlt};
  }
`;

const Features = styled.div`
  display: grid;
  gap: 2em;
  align-items: start;
  margin: 8em 4em 4em 4em;

  grid-template-columns: repeat(auto-fit, minmax(300px, calc(100% / 3 - 2em)));

  @media (max-width: 50em) {
    margin: 8em 1em 4em 1em;
  }

  h1 {
    text-transform: uppercase;
    opacity: 0.8;
  }
`;

const Card = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 2em;
  width: 100%;
  padding: 2em;
  border-radius: 12px;
  background-color: ${(p) => Color(p.theme.background2).alpha(0.8).hexa()};
  backdrop-filter: blur(5em);

  > img {
    width: 30%;
    height: auto;
    border-radius: 8px;
    box-shadow: 0 1em 2em 0 rgba(0 0 0 / 25%);
    order: 2;
  }

  > div {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    order: 1;
    h1 {
      font-size: 1.6rem;
      margin-bottom: 0.5em;
    }

    span {
      font-size: 1.2rem;
      line-height: 1.4em;
    }
  }
`;

const MainContent = styled.main`
  flex-grow: 1;
`;

const Footer = styled.footer`
  display: flex;
  gap: 5em;
  padding: 2em;
  justify-content: center;
  color: ${(p) => p.theme.text};
  background-color: ${(p) => p.theme.background2};
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

type Feature = {
  heading: string;
  description?: string;
  features?: string[];
};

export const StartRoute: React.FC<Props> = () => {
  const _loginRoute = loginRoute();
  const { t } = useTranslation("routes.start");
  const theme = useTheme();
  const otherFeatures: Feature = t("other", { returnObjects: true }) as Feature;
  const upcomingFeatures: Feature = t("upcoming", {
    returnObjects: true,
  }) as Feature;

  return (
    <StartContainer>
      <LoginButton onClick={() => (window.location.href = _loginRoute)}>
        <IconLogin />
        {t("login")}
      </LoginButton>
      <Header>
        <Brand>
          <img src={KikuriIcon} alt="kikuri icon" />
        </Brand>
        <span>
          <Trans
            ns="routes.start"
            i18nKey="header.under"
            components={{
              "1": (
                <GlowLink
                  href="https://github.com/zekurio/kikuri"
                  target="_blank"
                  rel="noreferrer"
                >
                  _
                </GlowLink>
              ),
            }}
          />
        </span>
        <HeaderButtons>
          <a href="/invite">
            <Button>{t("header.invite")}</Button>
          </a>
          <a href="https://github.com/zekurio/kikuri/wiki/selfhost">
            <Button>{t("header.selfhost")}</Button>
          </a>
          <a href="https://discord.gg/ay5YXMv5nT">
            <Button>{t("header.support")}</Button>
          </a>
        </HeaderButtons>
      </Header>
      <MainContent>
        <Features>
          <Card>
            <img
              src={theme._isDark ? MockupVotingDark : MockupVotingLight}
              alt=""
            />
            <div>
              <h1>{t("features.votes.heading")}</h1>
              <span>{t("features.votes.description")}</span>
            </div>
          </Card>
          <Card>
            <div>
              <h1>{otherFeatures.heading}</h1>
              <ul>
                {otherFeatures.features?.map((feature, index) => (
                  <li key={index}>{feature}</li>
                ))}
              </ul>
            </div>
          </Card>
          <Card>
            <div>
              <h1>{upcomingFeatures.heading}</h1>
              <ul>
                {upcomingFeatures.features?.map((feature, index) => (
                  <li key={index}>{feature}</li>
                ))}
              </ul>
            </div>
          </Card>
        </Features>
      </MainContent>
      <Footer>
        <div>
          <span>KIKURI - きくり</span>
          <span>© {new Date().getFullYear()} Michael Schwieger</span>
          <a
            href="https://github.com/zekurio/kikuri/blob/master/LICENCE"
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
            Invite Stable
          </a>
          <a
            href="https://canary.kikuri.xyz/invite"
            target="_blank"
            rel="noreferrer"
          >
            Invite Canary
          </a>
        </div>
        <div>
          <a
            href="https://github.com/zekurio/kikuri/wiki"
            target="_blank"
            rel="noreferrer"
          >
            Wiki
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/selfhost"
            target="_blank"
            rel="noreferrer"
          >
            Self Host
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/commands"
            target="_blank"
            rel="noreferrer"
          >
            Commands
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/permissions"
            target="_blank"
            rel="noreferrer"
          >
            Permissions Guide
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/restapi"
            target="_blank"
            rel="noreferrer"
          >
            REST API
          </a>
        </div>
      </Footer>
    </StartContainer>
  );
};
