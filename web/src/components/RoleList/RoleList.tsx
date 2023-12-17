import { Flex } from "../Flex";
import { Role } from "../../lib/kikuri-api/src";
import { Tag } from "../Tag";
import styled from "styled-components";

type Props = {
  roleids: string[];
  guildroles: Role[];
};

const RolesContainer = styled(Flex)`
  flex-wrap: wrap;
  gap: 0.4em;
`;

export const RoleList: React.FC<Props> = ({ roleids, guildroles }) => {
  const roles = roleids
    .map((rid) => guildroles.find((r) => r.id === rid))
    .filter((r) => !!r)
    .sort((ra, rb) => rb!.position - ra!.position)
    .map((r) => (
      <Tag key={r!.id} colors={r!.color} borderRadius="8px">
        {r!.name}
      </Tag>
    ));
  return <RolesContainer>{roles}</RolesContainer>;
};
