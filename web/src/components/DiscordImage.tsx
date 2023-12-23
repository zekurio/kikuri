import React from "react";
import styled from "styled-components";
import { IconBrandDiscord } from "@tabler/icons-react";

type ImgProps = {
  round?: boolean;
};

type Props = React.ImgHTMLAttributes<HTMLImageElement> & ImgProps;

const StyledImg = styled.img<ImgProps>`
  border-radius: ${(p) => (p.round ? "100%" : "8px")};
`;

export const DiscordImage: React.FC<Props> = ({ src, ...props }) => {
  return src ? <StyledImg src={src} {...props} /> : <IconBrandDiscord />;
};
