import { useRoutes } from 'react-router-dom';
import Home from "@/pages";
import ChatPage from "@/components/chat";

export const AppRoutes = () => {
  const commonRoutes = [{
    path: '/studio',
    element: <Home />
  }];

  const element = useRoutes([...commonRoutes]);

  return <>{element}</>;
};
