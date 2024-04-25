import { JSXElement, createSignal, For } from 'solid-js';
import RoomCard from '../components/RoomCard';
import Navbar from '../components/Navbar';
import FetchRandomAvatarImage from '../utils/RandomAvatar';
import CreateRoomCard from '../components/CreateRoomCard';
import ValidateRoomName from '../utils/RoomnameValidator';

type Room = {
  id: string;
  imgSrc: string;
  name: string;
};

export default function Home(): JSXElement {
  const [roomName, setRoomName] = createSignal<string>('');
  const [roomCount, setRoomCount] = createSignal<Room[]>([]);
  const [roomErr, setRoomErr] = createSignal<string>('');

  const createRoomChangeEvent = (
    e: Event & {
      currentTarget: HTMLInputElement;
      target: HTMLInputElement;
    }
  ) => {
    e.preventDefault();
    setRoomName(e.target.value);
  };

  const createRoomHandleClick = () => {
    const roomNameValidation: string = ValidateRoomName(
      roomName(),
      roomCount()
    );
    if (roomNameValidation.length > 0) {
      setRoomErr(roomNameValidation);
    } else if (roomCount().length < 6) {
      var newRoom: Room = {
        id: crypto.randomUUID(),
        imgSrc: FetchRandomAvatarImage(),
        name: roomName(),
      };
      setRoomCount((roomCount) => [...roomCount, newRoom]);
      setRoomErr('');
    }
  };

  return (
    <div>
      <Navbar />
      <CreateRoomCard
        onChangeHandler={createRoomChangeEvent}
        onClickHandler={createRoomHandleClick}
        errorMsg={roomErr}
      />
      <div class='grid grid-cols-3 gap-8 px-8 py-8'>
        <For each={roomCount()}>
          {(room, index) => (
            <RoomCard id={room.id} imgSrc={room.imgSrc} name={room.name} />
          )}
        </For>
      </div>
    </div>
  );
}
