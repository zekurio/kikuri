import styled from "styled-components";
import { AppTheme } from "../theme/theme";
import { useStoredTheme } from "../hooks/useStoredTheme";
import { useStore } from "../services/store";
import { Button } from "../components/Button";

type Props = {};

const DebugContainer = styled.div`
  padding: 2em;
  > section {
    margin-bottom: 2em;
    > * {
      width: 100%;
      display: block;
      margin: 0 0.5em 0.5em 0.5em;
    }
  }
`;

interface ColorBlockProps {
  color: string;
}

const ColorBlock = styled.div<ColorBlockProps>`
  background-color: ${(p) => p.color};
  color: ${(p) => p.theme.textColor};
  padding: 1em;
  margin: 1em 0;
  border: 1px solid #ccc;
`;

export const DebugRoute: React.FC<Props> = () => {
  const currentTheme = useStoredTheme();
  const [scheme, setScheme] = useStore((s) => [
    s.theme,
    s.setTheme,
    s.accentColor,
    s.setAccentColor,
  ]);

  // Function to switch theme
  const switchTheme = () => {
    setScheme(scheme === AppTheme.LIGHT ? AppTheme.DARK : AppTheme.LIGHT);
  };

  return (
    <DebugContainer>
      <Button onClick={switchTheme}>Switch Theme</Button>{" "}
      {/* Button for theme switching */}
      <section>
        {Object.keys(currentTheme.theme).map((key) => (
          <ColorBlock
            key={key}
            color={currentTheme.theme[key]}
          >
            {key}
          </ColorBlock>
        ))}
      </section>
    </DebugContainer>
  );
};
