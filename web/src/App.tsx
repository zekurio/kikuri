import React, { useEffect } from "react";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import styled, { ThemeProvider, createGlobalStyle } from "styled-components";

import { stripSuffix } from "./util/utils";
import { useStoredTheme } from "./hooks/useStoredTheme";
import { LandingRoute } from "./routes/Landing";

const GlobalStyle = createGlobalStyle`
  body {
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
  }

  &::-webkit-scrollbar {
    width: 10px;
  }

  &::-webkit-scrollbar-track {
    background: ${(p) => p.theme.background};
  }

  &::-webkit-scrollbar-thumb {
    background: ${(p) => p.theme.background3};
  }

  * {
    box-sizing: border-box;
  }

  a {
    color: ${(p) => p.theme.accent};
  }

  h1, h2, h3, h4, h5 {
    font-family: 'Nunito Sans';
    font-weight: 400;
  }
`;

const AppContainer = styled.div`
  width: 100%;
  height: 100vh;
`;

export const App: React.FC = () => {
  const { theme } = useStoredTheme();

  useEffect(() => {
    if (
      import.meta.env.BASE_URL.length > 0 &&
      window.location.pathname === stripSuffix(import.meta.env.BASE_URL, "/")
    ) {
      window.location.assign(import.meta.env.BASE_URL);
    }
  }, []);

  return (
    <ThemeProvider theme={theme}>
      <AppContainer>
        <Router basename={import.meta.env.BASE_URL}>
          <Routes>
            <Route path="/" element={<LandingRoute />} />
          </Routes>
        </Router>
      </AppContainer>
      <GlobalStyle />
    </ThemeProvider>
  );
};

export default App;
