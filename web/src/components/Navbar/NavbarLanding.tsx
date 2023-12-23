import React from "react";
import { Navbar } from "./Navbar";
import { Button } from "../Button";
import { useTranslation } from "react-i18next";
import { useSelfUser } from "../../hooks/useSelfUser";
import { useApi } from "../../hooks/useApi";
import {useNavigate} from "react-router";

export const NavbarLanding: React.FC = () => {
    const { t } = useTranslation("components.navbar.landing");
    const selfUser = useSelfUser();
    const fetch = useApi();
    const nav = useNavigate();

    const handleLogout = () => {
        fetch((c) => c.auth.logout())
            .then((r) => {
                if (r.code === 200) {
                    window.location.reload();
                }
            })
            .catch();
    };

    return (
        <Navbar
            rightContent={
                <>
                    <Button
                        onClick={() => {
                            nav("/controlcenter");
                        }}
                    >
                        {t("controlcenter")}
                    </Button>
                    {selfUser && <Button onClick={handleLogout}>{t("logout")}</Button>}
                </>
            }
        >
            <a href="#about">{t("links.about")}</a>
            <a href="#features">{t("links.features")}</a>
        </Navbar>
    );
};
