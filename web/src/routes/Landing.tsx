import { Trans, useTranslation } from "react-i18next";
import styled from "styled-components";
import React from "react";
import { NavbarLanding } from "../components/Navbar";
import HeaderBackground from "../assets/header-background.jpg";
import KIIcon from "../assets/ki-icon.png";

type Props = NonNullable<unknown>;

const LandingContainer = styled.div`
  display: flex;
  flex-direction: column;
  min-height: 100vh;
`;

const Header = styled.div`
    display: flex;
    align-items: center;
    justify-content: space-around;
    background-image: url(${HeaderBackground});
    background-size: cover;
    background-repeat: no-repeat;
    animation: move 30s linear infinite;
    height: 400px;

    @keyframes move {
        0% { background-position: 0 0; }
        100% { background-position: 100% 0; }
    }
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

export const LandingRoute: React.FC<Props> = () => {
  const { t } = useTranslation("routes.landing");

  return (
    <LandingContainer>
      <NavbarLanding />
      <Header>
        <Brand>
          <img src={KIIcon} alt="kikuri icon" />
        </Brand>
      </Header>
      <Footer>
        <div>
          <span>KIKURI - きくり</span>
          <span>© {new Date().getFullYear()} Michael Schwieger</span>
          <a
            href="https://github.com/zekurio/kikuri/blob/master/LICENCE"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.license")}
          </a>
          <a
            href="https://github.com/zekurio/kikuri"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.repo")}
          </a>
        </div>
        <div>
          <a href="https://kikuri.xyz/invite" target="_blank" rel="noreferrer">
            {t("footer.invitestable")}
          </a>
          <a
            href="https://canary.kikuri.xyz/invite"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.invitecanary")}
          </a>
        </div>
        <div>
          <a
            href="https://github.com/zekurio/kikuri/wiki"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.wiki")}
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/selfhost"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.selfhost")}
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/commands"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.commands")}
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/permissions"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.permissions")}
          </a>
          <a
            href="https://github.com/zekurio/kikuri/wiki/restapi"
            target="_blank"
            rel="noreferrer"
          >
            {t("footer.restapi")}
          </a>
        </div>
      </Footer>
    </LandingContainer>
  );
};
