import styled from "styled-components";

export type ButtonVariant =
  | "default"
  | "red"
  | "green"
  | "blue"
  | "yellow"
  | "orange"
  | "gray"
  | "pink";

export type ButtonProps = {
  variant?: ButtonVariant;
  nvp?: boolean;
  margin?: string;
};

export const Button = styled.button<ButtonProps>`
  font-size: 1rem;
  font-family: "Montserrat", sans-serif;
  color: ${(p) => p.theme.text};
  background: transparent; // Background initially transparent
  border: 2px solid ${(p) => p.theme.accent};
  padding: ${(p) => (p.nvp ? "0" : "0.8em")} 1em;
  border-radius: 3px;
  display: flex;
  gap: 0.8em;
  align-items: center;
  cursor: pointer;
  transition: background-color 0.2s ease;
  justify-content: center;
  margin: ${(p) => p.margin};

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &:enabled:hover {
    background-color: ${(p) => p.theme.accent}; // Fill background on hover
    color: ${(p) => p.theme.textAlt};
  }
`;
