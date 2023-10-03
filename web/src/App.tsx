import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import styled, { createGlobalStyle, ThemeProvider } from 'styled-components';
import { DarkTheme } from './theme/theme';
import { DebugRoute } from './routes/DebugRoute';
import { MainRoute } from './routes/MainRoute';

const GlobalStyle = createGlobalStyle`
  body {
    font-family: 'Rubik', sans-serif;
    background-color: ${(p) => p.theme.background};
    color: ${(p) => p.theme.text};
    padding: 0;
    margin: 0;
  }

  * {
    box-sizing: border-box;
  }

  h1, h2, h3, h4, h5, h6 {
    margin-top: 0;
  }
`;

const Outlet = styled.div`
  height: 100vh;
`;

const App: React.FC = () => {
    return (
        <ThemeProvider theme={DarkTheme}>
          <Outlet>
            <BrowserRouter>
              <Routes>
                <Route path="/" element={ <MainRoute/> } />
                <Route path="/debug" element={ <DebugRoute/> } />
              </Routes>
            </BrowserRouter>
          </Outlet>
          <GlobalStyle />
        </ThemeProvider>
    );
};

export default App;
