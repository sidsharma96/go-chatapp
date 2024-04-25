import { createSignal, type JSXElement } from 'solid-js';
import LoginCard from '../components/LoginCard';

export default function Signup(): JSXElement {
  return <LoginCard isLogin={false} />;
}
