import { IconBrandDiscord } from "@tabler/icons-react";
import styled from "styled-components";

type ImgProps = {
  round?: boolean;
};

type Props = React.ImgHTMLAttributes<any> & ImgProps;

const StyledImg = styled.img<ImgProps>`
  border-radius: ${(p) => (p.round ? "100%" : "8px")};
`;

export const DiscordImage: React.FC<Props> = ({ src, ...props }) => {
  if (src) {
    return <StyledImg src={src} {...props} />;
  } else {
    return <IconBrandDiscord {...props} />;
  }
};
