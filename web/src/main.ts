import './style.css';
import { APIClient } from './api/apiclient';

const apiClient = new APIClient();

apiClient.getSysinfo().then(systemInfo => {
  const textField = document.createElement('textarea');
  textField.value = JSON.stringify(systemInfo, null, 2);
  textField.readOnly = true;
  textField.style.width = '100%';
  textField.style.height = '400px';
  document.querySelector<HTMLDivElement>('#app')!.appendChild(textField);
}).catch(error => {
  console.error(error);
});