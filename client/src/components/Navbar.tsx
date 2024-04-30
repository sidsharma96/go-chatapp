import { A, useNavigate } from '@solidjs/router';
import { JSXElement, useContext } from 'solid-js';
import { LogoutUrl } from '../constants/UrlConstants';
import { WebsocketContext } from '../context/WebsocketContextProvider';

const fetchLogout = async () => {
  try {
    const response = await fetch(LogoutUrl, {
      method: 'POST',
      headers: {
        'content-type': 'application/json',
      },
      credentials: 'include',
    });
    await response.json();
  } catch (err) {
    console.log(err);
  }
};

export default function Navbar(): JSXElement {
  const { conn, setConn } = useContext(WebsocketContext);
  const navigate = useNavigate();

  const clickHandler = () => {
    if (conn() != null) {
      conn().close();
      setConn(null);
    }
  };

  const logoutHandler = async (e) => {
    e.preventDefault();
    await fetchLogout();
    clickHandler();
    navigate('/login', { replace: true });
  };

  return (
    <nav class='bg-white w-screen border-gray-200 dark:bg-gray-900'>
      <div class='flex flex-wrap items-center justify-between mx-4 p-4'>
        <div class='flex items-center space-x-3 rtl:space-x-reverse'>
          <img
            src='https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQaH_aMEasvJijjfBdIjr9LJLidVqXtGyGhsCVR5p_W6w&s'
            class='h-8'
            alt='GopherChat logo'
          />
          <span class='self-center font-[Montserrat] text-blue-800 font-semibold text-2xl whitespace-nowrap dark:text-white'>
            GopherChat
          </span>
        </div>
        <div class='block px-2' id='navbar-default'>
          <ul class='font-medium flex p-0 mt-4 rounded-lg bg-white space-x-8 dark:bg-gray-800 dark:border-gray-700'>
            <li>
              <A
                href='/'
                onClick={clickHandler}
                class='block rounded text-black dark:text-white md:dark:text-blue-500'
              >
                Home
              </A>
            </li>
            <li>
              <A
                href='/profile'
                onClick={clickHandler}
                class='block rounded text-black dark:text-white md:dark:text-blue-500'
              >
                My Profile
              </A>
            </li>
            <li>
              <A
                href=''
                class='block rounded text-black dark:text-white md:dark:text-blue-500'
                onClick={logoutHandler}
              >
                Logout
              </A>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}
