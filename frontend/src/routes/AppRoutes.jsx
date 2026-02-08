import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Main from '../pages/Main';
import Login from '../pages/Login';
import SetSafety from  '../pages/SetSafety'; 
import SetProfile from '../pages/SetProfile';
import SetGosUslugi from '../pages/SetGosUslugi';
import ProtectedRoute from './ProtectedRoutes';

const AppRoutes = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<ProtectedRoute>  <Main />  </ProtectedRoute>} />
      <Route path="/settings/safety"  element={<ProtectedRoute>  <SetSafety />  </ProtectedRoute>} />
      <Route path="/settings/profile"  element={<ProtectedRoute>  <SetProfile />  </ProtectedRoute>} />
      <Route path="/settings/gosuslugi" element={<ProtectedRoute>  <SetGosUslugi />  </ProtectedRoute>} />
    </Routes>
  );
};

export default AppRoutes;