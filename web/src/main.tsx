import { StyleProvider } from '@ant-design/cssinjs';
import { ConfigProvider } from 'antd';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from './App';
import FileList from './FileList';

import 'remixicon/fonts/remixicon.css';
import './style.scss';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
  },
  {
    path: `/browser/:model`,
    element: <FileList />,
  },
]);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ConfigProvider
      theme={{
        token: {
          colorPrimary: '#EB5424',
          colorText: '#313638',
          colorSuccess: '#4BAB4E',
          colorError: '#EB5424',
          colorInfo: '#2454BB',
          borderRadius: 4,
        },
      }}
    >
      <StyleProvider hashPriority="high">
        <React.StrictMode>
          <div className="p-0">
            <div className="p-4">
              <RouterProvider router={router} />
            </div>
          </div>
        </React.StrictMode>
      </StyleProvider>
    </ConfigProvider>
  </React.StrictMode>
);
