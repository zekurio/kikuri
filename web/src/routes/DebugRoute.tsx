import React, { useState, useEffect } from 'react';
import { APIClient } from '../api/apiclient';
import { SystemInfo } from '../api/models';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import styled from 'styled-components';

type Props = {};

const DebugContainer = styled.div`
  background-color: ${(p) => p.theme.background};
  color: ${(p) => p.theme.text};
  padding: 16px;
`;

const LoginButton = styled.button`
  background-color: #7289da;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
`;

export const DebugRoute: React.FC<Props> = ({}) => {
  const [sysinfo, setSysinfo] = useState<SystemInfo | null>(null);

  useEffect(() => {
    const client = new APIClient();

    const fetchSysinfo = async () => {
      const info = await client.getSysinfo();
      setSysinfo(info);
    };

    fetchSysinfo();
  }, []);

  const handleLogin = async () => {
    const client = new APIClient();
    const response = await client.loginToDiscord();
    window.location.href = response.redirectUrl;
  };

  return (
    <DebugContainer>
      <h1>Debug</h1>
      <h2>API Endpoints</h2>
      <ul>
        <li>
          <strong>/others/sysinfo</strong>:
          {sysinfo ? (
            <SyntaxHighlighter language="json">
              {JSON.stringify(sysinfo, null, 2)}
            </SyntaxHighlighter>
          ) : (
            'Loading...'
          )}
        </li>
      </ul>
      <LoginButton onClick={handleLogin}>Login to Discord</LoginButton>
    </DebugContainer>
  );
};
