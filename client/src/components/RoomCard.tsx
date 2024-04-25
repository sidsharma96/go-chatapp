import { JSXElement, Show, createResource } from 'solid-js';

export default function RoomCard({ id, imgSrc, name }): JSXElement {
  return (
    <div class='min-w-32 min-h-52 flex flex-col items-center border bg-white px-6 py-6 shadow-md rounded-3xl'>
      <img src={imgSrc} alt='avatar' />
      <p class='self-center py-4 text-2xl font-semibold whitespace-nowrap dark:text-white'>
        {name}
      </p>
      <button class='border justify-center p-2 items-center content-center bg-indigo-700 text-white rounded-[4px]'>
        Join Room
      </button>
    </div>
  );
}
