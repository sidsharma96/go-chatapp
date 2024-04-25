import { lazy } from 'solid-js';
import { Route, RouteSectionProps, Router } from '@solidjs/router';
import type { Component, JSXElement } from 'solid-js';
import Home from './pages/Home';

const Login = lazy(() => import('./pages/Login'));
const Signup = lazy(() => import('./pages/Signup'));

const AppContext = (props: RouteSectionProps): JSXElement => {
  return <>{props.children}</>;
};

const App: Component = () => {
  return (
    <Router root={AppContext}>
      <Route path='/login' component={Login}></Route>
      <Route path='/signup' component={Signup}></Route>
      <Route path='/' component={Home}></Route>
    </Router>
  );
};

export default App;
