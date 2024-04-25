import { createSignal, type JSXElement } from 'solid-js';
import LoginCard from '../components/LoginCard';

export default function Login(): JSXElement {
  return <LoginCard isLogin={true} />;
}
