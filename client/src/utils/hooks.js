import React from 'react';
import { useCookies } from 'react-cookie';
import { logout } from '../api/authentication';

function useClientInfos() {
  const [cookies] = useCookies(['client-infos']);
  
  let clientInfos = null;
  
  try {
    // Try to parse the cookie into a JavaScript object
    clientInfos = cookies['client-infos'].split(',');
  } catch (error) {
    console.error('Error parsing client-infos cookie:', error);
    logout();
  }
  
  return {
    nickname: clientInfos[0],
    role: clientInfos[1]
  };
}

export {
  useClientInfos
};