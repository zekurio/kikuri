import { IconBrandDiscord } from "@tabler/icons";
import styled from "styled-components";

type ImgProps = {
  round?: boolean;
};

type Props = React.ImgHTMLAttributes<any> & ImgProps;

const StyledImg = styled.img<ImgProps>`
  border-radius: ${(p) => (p.round ? "100%" : "8px")};
`;

export const DiscordImage: React.FC<Props> = ({ src, ...props }) => {
  return <StyledImg src={!!src ? src : IconBrandDiscord} {...props} />;
};
