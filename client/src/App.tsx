import { lazy, children } from 'solid-js';
import { Route, RouteSectionProps, Router } from '@solidjs/router';
import type { Component, JSXElement } from 'solid-js';
import Home from './pages/Home';
import { WebSocketContextProvider } from './context/WebsocketContextProvider';
import { AuthContextProvider } from './context/AuthContextProvider';

const Login = lazy(() => import('./pages/Login'));
const Signup = lazy(() => import('./pages/Signup'));
const Profile = lazy(() => import('./pages/Profile'));
const Chat = lazy(() => import('./pages/Chat'));

const AppContext = (props: RouteSectionProps): JSXElement => {
  const safeChildren = children(() => props.children);
  return <>{safeChildren()}</>;
};

const App: Component = () => {
  return (
    <AuthContextProvider>
      <WebSocketContextProvider>
        <Router root={AppContext}>
          <Route path='/login' component={Login}></Route>
          <Route path='/signup' component={Signup}></Route>
          <Route path='/profile' component={Profile}></Route>
          <Route path='/chat/:name' component={Chat}></Route>
          <Route path='/' component={Home}></Route>
          <Route path={'*'} component={Login} />
        </Router>
      </WebSocketContextProvider>
    </AuthContextProvider>
  );
};

export default App;
