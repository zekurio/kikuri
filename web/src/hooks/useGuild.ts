import { useEffect, useState } from "react";

import { Guild } from "../lib/kikuri-api/src";
import { useApi } from "./useApi";

export const useGuild = (id?: string) => {
  const fetch = useApi();
  const [guild, setGuild] = useState<Guild>();

  useEffect(() => {
    if (!id) return;

    fetch((c) => c.guilds.guild(id))
      .then((res) => setGuild(res))
      .catch();
  }, [id]);

  return guild;
};
