import { useTranslation } from "react-i18next";
import styled from "styled-components";
import React, {useEffect, useState} from "react";
import { NavbarLanding } from "../components/Navbar";
import KIIcon from "../assets/ki-icon.png";

type Props = NonNullable<unknown>;

const LandingContainer = styled.div``;

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

const Header = styled.header`
    display: flex;
    flex-direction: column;
    gap: 3em;
    align-items: center;
    padding-top: 20vh;

    > span {
        font-family: 'Montserrat', sans-serif;
        font-size: 1.1rem;
        font-weight: lighter;
        text-align: center;
        max-width: 20em;
    }
`;

const Main = styled.main`
  display: flex;
  flex-direction: column;
  gap: 2em;
  align-items: center;
  padding: 2em;
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

  const [backgroundPosition, setBackgroundPosition] = useState('center');

  const handleMouseMove = (event: MouseEvent) => {
    const { clientX, clientY } = event;
    const screenWidth = window.innerWidth;
    const screenHeight = window.innerHeight;

    const backgroundX = (clientX / screenWidth) * 100;
    const backgroundY = (clientY / screenHeight) * 100;

    setBackgroundPosition(`${backgroundX}% ${backgroundY}%`);
  };

  useEffect(() => {
    window.addEventListener('mousemove', handleMouseMove);

    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
    };
  }, []);

  return (
    <LandingContainer>
      <NavbarLanding />
      <Header style={{ backgroundPosition }}>
        <Brand>
          <img src={KIIcon} alt="kikuri icon" />
        </Brand>
      </Header>
      <Main>
        <h1>{t("main.wip")}</h1>
      </Main>
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
