import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import styled from 'styled-components';

type Props = {};

const MainContainer = styled.div`
  height: 100%;
  display: flex;

  > main {
    margin-left: 4em;
    width: 100%;

    @media screen and (orientation: portrait) {
      margin-left: 1.5em;
    }
  }
`;

export const MainRoute: React.FC<Props> = ({}) => {
  return (
    <MainContainer>
      <main>
        <Outlet />
      </main>
    </MainContainer>
  );
};
