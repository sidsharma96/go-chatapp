enum RoomError {
  limitError = 'Limit of 6 rooms',
  emptyError = 'Room name cannot be empty',
  lengthError = 'Room name cannot be more than 10 letters',
  numericError = 'Room name can only have letters',
}

export default function ValidateRoomName(roomName: string, roomList): string {
  var letters = /^[A-Za-z]+$/;
  if (roomList.length >= 6) {
    return RoomError.limitError;
  } else if (roomName.length === 0) {
    return RoomError.emptyError;
  } else if (!roomName.match(letters)) {
    return RoomError.numericError;
  } else if (roomName.length > 10) {
    return RoomError.lengthError;
  }
  return '';
}
