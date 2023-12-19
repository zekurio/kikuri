import { BottomContainer, SelfContainer } from './BottomContainer';
import { SIDEBAR_WIDTH, Sidebar } from './Sidebar';
import { useNavigate, useParams } from 'react-router';
import {MAX_WIDTH} from "../MaxWidthContainer";
import styled from "styled-components";
import {Entry} from "./Entry";
import {EntryContainer} from "./EntryContainer";
import {Section} from "./Section";
import {useTranslation} from "react-i18next";
import React, {useEffect} from "react";
import {useStore} from "../../services/store";
import {useGuilds} from "../../hooks/useGuilds";
import {Guild} from "../../lib/kikuri-api/src";
import {IconUser} from "@tabler/icons-react";

type Props = NonNullable<unknown>;

const StyledEntry = styled(Entry)``;

const StyledBar = styled(Sidebar)`
    @media (orientation: portrait) and (max-width: calc(${SIDEBAR_WIDTH} + ${MAX_WIDTH})) {
    ${StyledEntry}, ${SelfContainer} {
      justify-content: center;
      span {
        display: none;
      }
    }

    ${SelfContainer} > svg {
      display: none;
    }
  }
`;

export const SidebarDashboard: React.FC<Props> = () => {
    const { t } = useTranslation('components', { keyPrefix: 'sidebar' });
    const nav = useNavigate();
    const { guildid } = useParams();
    const guilds = useGuilds();
    const setSelectedGuild = useStore((s) => s.setSelectedGuild);
    const selectedGuild = useStore((s) => s.selectedGuild);

    useEffect(() => {
        if (!!guilds && !!guildid) setSelectedGuild(guilds.find((g) => g.id === guildid) ?? guilds[0]);
    }, [guildid, guilds, setSelectedGuild]);

    const _onGuildSelect = (g: Guild) => {
        setSelectedGuild(g);
        nav(`guilds/${g.id}/members`);
    };

    return (
        <StyledBar>
            <EntryContainer>
                <Section title={t('section.guilds.title')}>
                    <StyledEntry path={`/dashboard/guilds/${selectedGuild?.id}/members`}>
                        <IconUser />
                        <span>{t('section.guilds.members')}</span>
                    </StyledEntry>
                </Section>
            </EntryContainer>
        </StyledBar>
    );
};