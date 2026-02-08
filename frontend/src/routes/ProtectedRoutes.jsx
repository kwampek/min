import { Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';

export default function ProtectedRoute({ children }) {
  const valid = useSelector((state) => state.user.valid);
  
  if (!valid) {
    return <Navigate to="/login" replace />;
  }

  return children
}
