import { PropsWithStyle } from "../props";
import styled from "styled-components";
import { PropsWithChildren } from "react";
import { EntryContainer } from "./EntryContainer";
import { MAX_WIDTH } from "../MaxWidthContainer";
import { Heading } from "../Heading";

type Props = PropsWithChildren & PropsWithStyle & NonNullable<unknown>;

export const SIDEBAR_WIDTH = "15rem";

const Brand = styled.div``;

const StyledBar = styled.div`
  display: flex;
  flex-direction: column;
  background-color: ${(p) => p.theme.background2};
  margin: 1rem 0 1rem 1rem;
  padding: 1rem;
  border-radius: 12px;
  width: 30vw;
  max-width: 15rem;

  @media (orientation: portrait) and (max-width: calc(${SIDEBAR_WIDTH} + ${MAX_WIDTH})) {
    width: fit-content;

    ${Brand} > svg {
      display: none;
    }

    ${Heading} {
      display: none;
    }
  }

  > ${EntryContainer} {
    &::-webkit-scrollbar {
      width: 5px;
    }
  }
`;

export const Sidebar: React.FC<Props> = ({ children, ...props }) => {
  return (
    <StyledBar {...props}>
      <Brand></Brand>
      {children}
    </StyledBar>
  );
};
