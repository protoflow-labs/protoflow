import { useRoutes } from 'react-router-dom';
import Home from "@/pages";

export const AppRoutes = () => {
  const commonRoutes = [{ path: '/', element: <Home /> }];

  const element = useRoutes([...commonRoutes]);

  return <>{element}</>;
};
