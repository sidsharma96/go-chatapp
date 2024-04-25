import { JSXElement } from 'solid-js';

export default function Navbar(): JSXElement {
  return (
    <nav class='bg-white border-gray-200 dark:bg-gray-900'>
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
              <a
                href='/'
                class='block rounded text-black dark:text-white md:dark:text-blue-500'
              >
                Home
              </a>
            </li>
            <li>
              <a
                href='#'
                class='block rounded text-black dark:text-white md:dark:text-blue-500'
              >
                My Profile
              </a>
            </li>
            <li>
              <a
                href='#'
                class='block rounded text-black dark:text-white md:dark:text-blue-500'
              >
                Logout
              </a>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}
